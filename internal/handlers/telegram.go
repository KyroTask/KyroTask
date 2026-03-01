package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bif12/kyrotask/internal/config"
	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-gonic/gin"
)

type TelegramHandler struct {
	telegramService     *services.TelegramService
	notificationService *services.NotificationService
}

func NewTelegramHandler() *TelegramHandler {
	return &TelegramHandler{
		telegramService:     services.NewTelegramService(),
		notificationService: services.NewNotificationService(),
	}
}

type VerifyWebAppRequest struct {
	InitData string `json:"initData" binding:"required"`
}

type VerifyWebAppResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// VerifyTelegramWebApp verifies the Telegram WebApp initData and returns a JWT
func (h *TelegramHandler) VerifyTelegramWebApp(c *gin.Context) {
	var req VerifyWebAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Verify initData
	telegramData, err := h.telegramService.VerifyInitData(req.InitData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram data: " + err.Error()})
		return
	}

	// Create or update user
	user, err := h.telegramService.CreateOrUpdateUser(telegramData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Generate JWT
	token, err := h.telegramService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, VerifyWebAppResponse{
		Token: token,
		User:  user,
	})
}

// VerifyTelegramWidget verifies the Telegram Login Widget data and returns a JWT
func (h *TelegramHandler) VerifyTelegramWidget(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Verify widget data
	if err := h.telegramService.VerifyWidgetData(req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram data: " + err.Error()})
		return
	}

	// Create or update user
	user, err := h.telegramService.CreateOrUpdateUserFromWidget(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Generate JWT
	token, err := h.telegramService.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, VerifyWebAppResponse{
		Token: token,
		User:  user,
	})
}

// GetCurrentUser returns the currently authenticated user
func (h *TelegramHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch user from database using userID
	// For now, just return the user ID
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Authenticated",
	})
}

// HandleWebhook handles incoming Telegram bot webhooks
func (h *TelegramHandler) HandleWebhook(c *gin.Context) {
	var update services.TelegramUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	h.ProcessUpdate(update)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ProcessUpdate handles a single Telegram update (from webhook or polling)
func (h *TelegramHandler) ProcessUpdate(update services.TelegramUpdate) {
	if update.CallbackQuery != nil {
		err := h.notificationService.HandleCallback(
			update.CallbackQuery.From.ID,
			update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.Data,
		)
		if err != nil {
			log.Printf("Error handling callback: %v", err)
		}
	}

	if update.Message != nil {
		text := update.Message.Text
		telegramID := update.Message.From.ID

		if strings.HasPrefix(text, "/start") {
			h.handleStart(telegramID)
		} else {
			// Check if authenticated
			user, err := h.telegramService.GetUserByTelegramID(telegramID)
			if err != nil {
				h.telegramService.SendMessage(telegramID, "⚠️ *Authentication Required*\n━━━━━━━━━━━━━━━━━━\nPlease authenticate first by opening the Mini App using the button below or typing /start\\.", map[string]interface{}{
					"inline_keyboard": [][]map[string]interface{}{
						{
							{"text": "🚀 Open Mini App", "url": fmt.Sprintf("https://t.me/%s/app", config.AppConfig.TelegramBotUsername)},
						},
					},
				})
				return
			}

			if strings.HasPrefix(text, "/help") {
				h.handleHelp(telegramID)
			} else if strings.HasPrefix(text, "/calendar") {
				h.handleCalendar(telegramID, user)
			} else {
				h.telegramService.SendMessage(telegramID, "❓ *Unknown Command*\n━━━━━━━━━━━━━━━━━━\nI don't recognize that command\\. Type /help for a list of available commands\\.", nil)
			}
		}
	}
}

// StartPolling starts the long polling loop for Telegram updates
func (h *TelegramHandler) StartPolling() {
	log.Println("🤖 Telegram bot starting in Long Polling mode...")

	// Delete webhook first to ensure polling works
	if err := h.telegramService.DeleteWebhook(); err != nil {
		log.Printf("Warning: Failed to delete webhook: %v", err)
	}

	offset := 0
	for {
		updates, err := h.telegramService.GetUpdates(offset)
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates {
			h.ProcessUpdate(update)
			offset = update.UpdateID + 1
		}
	}
}

func (h *TelegramHandler) handleStart(telegramID int64) {
	user, err := h.telegramService.GetUserByTelegramID(telegramID)
	if err == nil {
		// User is already authenticated
		text := fmt.Sprintf("👋 *Welcome back, %s\\!* \n━━━━━━━━━━━━━━━━━━\nYou are already authenticated and ready to go\\. You can manage your tasks, habits, and goals directly from the Mini App or use commands here\\.\n\nType /help to see what I can do for you\\.",
			h.telegramService.EscapeMarkdownV2(user.FirstName))

		replyMarkup := map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{"text": "🚀 Open Mini App", "url": fmt.Sprintf("https://t.me/%s/app", config.AppConfig.TelegramBotUsername)},
				},
			},
		}
		h.telegramService.SendMessage(telegramID, text, replyMarkup)
		return
	}

	// New user - Authentication required
	text := "👋 *Welcome to Task Manager Bot\\!*\n━━━━━━━━━━━━━━━━━━\nI am your personal productivity assistant\\. To get started, please authenticate your account by opening the Mini App below\\.\n\nOnce authenticated, you can receive notifications and manage your tasks directly from Telegram\\."

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🚀 Open Mini App", "url": fmt.Sprintf("https://t.me/%s/app", config.AppConfig.TelegramBotUsername)},
			},
		},
	}

	h.telegramService.SendMessage(telegramID, text, replyMarkup)
}

func (h *TelegramHandler) handleHelp(telegramID int64) {
	text := "📖 *Available Commands*\n━━━━━━━━━━━━━━━━━━\n/start \\- Welcome message & App link\n/calendar \\- View upcoming tasks\n/help \\- Show this help message\n\n*Mini App Features:*\n• Manage Projects & Goals\n• Track Habits & Streaks\n• Set Task Reminders\n\nNeed more help? Open the Mini App to access full settings\\."
	h.telegramService.SendMessage(telegramID, text, nil)
}

func (h *TelegramHandler) handleCalendar(telegramID int64, user *models.User) {
	var tasks []models.Task
	now := time.Now()
	// Fetch next 7 days of tasks
	nextWeek := now.AddDate(0, 0, 7)

	if err := database.DB.Where("user_id = ? AND due_date >= ? AND due_date <= ? AND status != ?", user.ID, now, nextWeek, "completed").Order("due_date asc").Limit(10).Find(&tasks).Error; err != nil {
		h.telegramService.SendMessage(telegramID, "❌ *Error*\n━━━━━━━━━━━━━━━━━━\nFailed to fetch your calendar\\. Please try again later\\.", nil)
		return
	}

	if len(tasks) == 0 {
		h.telegramService.SendMessage(telegramID, "📅 *Your Calendar*\n━━━━━━━━━━━━━━━━━━\nYou have no upcoming tasks for the next 7 days\\. Enjoy your free time\\! 🎉", nil)
		return
	}

	var sb strings.Builder
	sb.WriteString("📅 *Upcoming Tasks*\n━━━━━━━━━━━━━━━━━━\n")
	for _, t := range tasks {
		sb.WriteString(fmt.Sprintf("• *%s*\n  ⏰ %s\n\n",
			h.telegramService.EscapeMarkdownV2(t.Title),
			t.DueDate.Format("Mon, Jan 02 15:04"),
		))
	}
	sb.WriteString("━━━━━━━━━━━━━━━━━━\n_Showing next 10 tasks_")

	h.telegramService.SendMessage(telegramID, sb.String(), nil)
}
