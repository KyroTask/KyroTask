package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/bif12/kyrotask/internal/config"
	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt/v5"
)

type TelegramUpdate = tgbotapi.Update

type TelegramService struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramService() *TelegramService {
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.TelegramBotToken)
	if err != nil {
		log.Printf("Error initializing Telegram Bot: %v", err)
		return &TelegramService{}
	}
	return &TelegramService{Bot: bot}
}

// DeleteWebhook removes any existing webhook (required for long polling)
func (s *TelegramService) DeleteWebhook() error {
	if s.Bot == nil {
		return errors.New("bot not initialized")
	}
	_, err := s.Bot.Request(tgbotapi.DeleteWebhookConfig{})
	return err
}

// GetUpdates fetches new updates from Telegram using long polling
func (s *TelegramService) GetUpdates(offset int) ([]tgbotapi.Update, error) {
	if s.Bot == nil {
		return nil, errors.New("bot not initialized")
	}
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = 30
	return s.Bot.GetUpdates(u)
}

// VerifyInitData verifies Telegram WebApp initData HMAC signature
func (s *TelegramService) VerifyInitData(initData string) (map[string]string, error) {
	if initData == "" {
		return nil, errors.New("initData is empty")
	}

	// Parse the init data
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initData: %w", err)
	}

	// Extract hash
	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, errors.New("hash not found in initData")
	}

	// Remove hash from values for verification
	values.Del("hash")

	// Create data-check-string
	var dataCheckArray []string
	for key := range values {
		dataCheckArray = append(dataCheckArray, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sort.Strings(dataCheckArray)
	dataCheckString := strings.Join(dataCheckArray, "\n")

	// Compute HMAC
	botToken := config.AppConfig.TelegramBotToken
	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	// Compare hashes
	if computedHash != receivedHash {
		return nil, errors.New("invalid hash: initData verification failed")
	}

	// Convert values to map
	result := make(map[string]string)
	for key := range values {
		result[key] = values.Get(key)
	}

	return result, nil
}

// CreateOrUpdateUser creates or updates a user from Telegram data
func (s *TelegramService) CreateOrUpdateUser(telegramData map[string]string) (*models.User, error) {
	// Parse user data from initData
	userStr := telegramData["user"]
	if userStr == "" {
		return nil, errors.New("user data not found in initData")
	}

	// Define a struct to match Telegram's user object
	type TelegramUser struct {
		ID           int64  `json:"id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
		IsBot        bool   `json:"is_bot"`
		PhotoURL     string `json:"photo_url"`
	}

	var tgUser TelegramUser
	if err := json.Unmarshal([]byte(userStr), &tgUser); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	var user models.User

	// Check if user exists
	result := database.DB.Where("telegram_id = ?", tgUser.ID).First(&user)

	if result.Error == nil {
		// User exists, update
		user.FirstName = tgUser.FirstName
		user.LastName = tgUser.LastName
		user.Username = tgUser.Username
		user.LanguageCode = tgUser.LanguageCode
		user.PhotoURL = tgUser.PhotoURL
		database.DB.Save(&user)
		return &user, nil
	}

	// Create new user
	user = models.User{
		TelegramID:   &tgUser.ID,
		FirstName:    tgUser.FirstName,
		LastName:     tgUser.LastName,
		Username:     tgUser.Username,
		LanguageCode: tgUser.LanguageCode,
		IsBot:        tgUser.IsBot,
		PhotoURL:     tgUser.PhotoURL,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *TelegramService) GenerateJWT(user *models.User) (string, error) {
	var tID int64
	if user.TelegramID != nil {
		tID = *user.TelegramID
	}
	claims := middleware.Claims{
		UserID:     user.ID,
		TelegramID: tID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyWidgetData verifies Telegram Login Widget data
func (s *TelegramService) VerifyWidgetData(data map[string]interface{}) error {
	receivedHash, ok := data["hash"].(string)
	if !ok || receivedHash == "" {
		return errors.New("hash not found in widget data")
	}

	// Create data-check-string
	var dataCheckArray []string
	for key, value := range data {
		if key == "hash" {
			continue
		}
		// Convert value to string
		var valStr string
		switch v := value.(type) {
		case float64:
			valStr = fmt.Sprintf("%.0f", v)
		default:
			valStr = fmt.Sprintf("%v", v)
		}
		dataCheckArray = append(dataCheckArray, fmt.Sprintf("%s=%s", key, valStr))
	}
	sort.Strings(dataCheckArray)
	dataCheckString := strings.Join(dataCheckArray, "\n")

	// Compute HMAC
	botToken := config.AppConfig.TelegramBotToken
	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	// Compare hashes
	if computedHash != receivedHash {
		fmt.Printf("Widget Verification Failed:\n")
		fmt.Printf("Received Hash: %s\n", receivedHash)
		fmt.Printf("Computed Hash: %s\n", computedHash)
		fmt.Printf("Data Check String:\n%s\n", dataCheckString)
		fmt.Printf("Bot Token: %s...\n", botToken[:10]) // Log partial token for verification
		return errors.New("invalid hash: widget data verification failed")
	}

	// Check auth_date
	if authDate, ok := data["auth_date"].(float64); ok {
		if time.Now().Unix()-int64(authDate) > 86400 {
			return errors.New("auth_date is outdated")
		}
	}

	return nil
}

// CreateOrUpdateUserFromWidget creates or updates a user from Telegram Widget data
func (s *TelegramService) CreateOrUpdateUserFromWidget(data map[string]interface{}) (*models.User, error) {
	id, _ := data["id"].(float64)
	telegramID := int64(id)
	firstName, _ := data["first_name"].(string)
	lastName, _ := data["last_name"].(string)
	username, _ := data["username"].(string)
	photoURL, _ := data["photo_url"].(string)

	var user models.User

	// Check if user exists
	result := database.DB.Where("telegram_id = ?", telegramID).First(&user)

	if result.Error == nil {
		// User exists, update
		user.FirstName = firstName
		user.LastName = lastName
		user.Username = username
		user.PhotoURL = photoURL
		database.DB.Save(&user)
		return &user, nil
	}

	// Create new user
	user = models.User{
		TelegramID: &telegramID,
		FirstName:  firstName,
		LastName:   lastName,
		Username:   username,
		PhotoURL:   photoURL,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// SendMessage sends a message to a Telegram user
func (s *TelegramService) SendMessage(telegramID int64, text string, replyMarkup interface{}) (int, error) {
	if s.Bot == nil {
		return 0, errors.New("bot not initialized")
	}

	msg := tgbotapi.NewMessage(telegramID, text)
	msg.ParseMode = "MarkdownV2"

	if replyMarkup != nil {
		// Convert generic replyMarkup to tgbotapi.InlineKeyboardMarkup if possible
		// This is a bit tricky with the current generic map structure,
		// but let's try to handle the common case.
		if markupMap, ok := replyMarkup.(map[string]interface{}); ok {
			if keyboard, ok := markupMap["inline_keyboard"].([][]map[string]interface{}); ok {
				var inlineKeyboard [][]tgbotapi.InlineKeyboardButton
				for _, row := range keyboard {
					var inlineRow []tgbotapi.InlineKeyboardButton
					for _, btn := range row {
						text, _ := btn["text"].(string)
						callbackData, _ := btn["callback_data"].(string)
						url, _ := btn["url"].(string)

						button := tgbotapi.InlineKeyboardButton{Text: text}
						if callbackData != "" {
							button.CallbackData = &callbackData
						}
						if url != "" {
							button.URL = &url
						}
						inlineRow = append(inlineRow, button)
					}
					inlineKeyboard = append(inlineKeyboard, inlineRow)
				}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard...)
			}
		}
	}

	sentMsg, err := s.Bot.Send(msg)
	if err != nil {
		return 0, err
	}

	return sentMsg.MessageID, nil
}

// DeleteMessage deletes a message from a Telegram chat
func (s *TelegramService) DeleteMessage(telegramID int64, messageID int) error {
	if s.Bot == nil {
		return errors.New("bot not initialized")
	}
	deleteConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    telegramID,
		MessageID: messageID,
	}
	_, err := s.Bot.Request(deleteConfig)
	return err
}

// GetUserByTelegramID fetches a user by their Telegram ID
func (s *TelegramService) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// EscapeMarkdownV2 escapes special characters for Telegram MarkdownV2
func (s *TelegramService) EscapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(", "\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>", "#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|", "\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
	)
	return replacer.Replace(text)
}
