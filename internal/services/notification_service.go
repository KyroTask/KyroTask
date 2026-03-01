package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
)

type NotificationService struct {
	telegramService *TelegramService
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		telegramService: NewTelegramService(),
	}
}

// SendTaskNotification sends a notification for a task
func (s *NotificationService) SendTaskNotification(user *models.User, task *models.Task) error {
	if user.TelegramID == nil {
		return nil
	}

	text := fmt.Sprintf("🔔 *Task Reminder*\n━━━━━━━━━━━━━━━━━━\n📝 *Task:* %s\n📅 *Due:* %s\n🏷️ *Project:* %s\n\n━━━━━━━━━━━━━━━━━━",
		s.telegramService.EscapeMarkdownV2(task.Title),
		task.DueDate.Format("2006-01-02 15:04"),
		s.telegramService.EscapeMarkdownV2(getProjectName(task)),
	)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "✅ Done & Remove", "callback_data": fmt.Sprintf("task_done:%d", task.ID)},
			},
		},
	}

	msgID, err := s.telegramService.SendMessage(*user.TelegramID, text, replyMarkup)
	if err != nil {
		return err
	}

	notification := models.Notification{
		UserID:            user.ID,
		Type:              "task_due",
		RelatedID:         task.ID,
		TelegramMessageID: msgID,
		SentAt:            time.Now(),
	}

	return database.DB.Create(&notification).Error
}

// SendHabitNotification sends a daily habit reminder
func (s *NotificationService) SendHabitNotification(user *models.User, habit *models.Habit) error {
	if user.TelegramID == nil {
		return nil
	}

	text := fmt.Sprintf("✨ *Daily Habit Check-in*\n━━━━━━━━━━━━━━━━━━\nTime to log your daily habits!\n\n✅ *%s* - [Streak: %d days 🔥]\n\n━━━━━━━━━━━━━━━━━━",
		s.telegramService.EscapeMarkdownV2(habit.Name),
		habit.CurrentStreak,
	)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "✅ Done & Remove", "callback_data": fmt.Sprintf("habit_done:%d", habit.ID)},
			},
		},
	}

	msgID, err := s.telegramService.SendMessage(*user.TelegramID, text, replyMarkup)
	if err != nil {
		return err
	}

	notification := models.Notification{
		UserID:            user.ID,
		Type:              "habit_reminder",
		RelatedID:         habit.ID,
		TelegramMessageID: msgID,
		SentAt:            time.Now(),
	}

	return database.DB.Create(&notification).Error
}

// SendWeeklyHabitNotification sends a weekly habit review with tasks
func (s *NotificationService) SendWeeklyHabitNotification(user *models.User, habit *models.Habit, tasks []models.Task) error {
	if user.TelegramID == nil {
		return nil
	}

	var taskList strings.Builder
	for _, t := range tasks {
		taskList.WriteString(fmt.Sprintf("• %s\n", s.telegramService.EscapeMarkdownV2(t.Title)))
	}

	text := fmt.Sprintf("🗓️ *Weekly Habit Review*\n━━━━━━━━━━━━━━━━━━\n*Habit:* %s\n\n*Tasks for this week:*\n%s\n━━━━━━━━━━━━━━━━━━",
		s.telegramService.EscapeMarkdownV2(habit.Name),
		taskList.String(),
	)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "✅ Weekly Done & Remove", "callback_data": fmt.Sprintf("habit_weekly_done:%d", habit.ID)},
			},
		},
	}

	msgID, err := s.telegramService.SendMessage(*user.TelegramID, text, replyMarkup)
	if err != nil {
		return err
	}

	notification := models.Notification{
		UserID:            user.ID,
		Type:              "habit_weekly",
		RelatedID:         habit.ID,
		TelegramMessageID: msgID,
		SentAt:            time.Now(),
	}

	return database.DB.Create(&notification).Error
}

// HandleCallback processes Telegram callback queries
func (s *NotificationService) HandleCallback(telegramID int64, messageID int, callbackData string) error {
	parts := strings.Split(callbackData, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid callback data")
	}

	action := parts[0]
	idStr := parts[1]

	switch action {
	case "task_done":
		var task models.Task
		if err := database.DB.First(&task, idStr).Error; err != nil {
			return err
		}
		now := time.Now()
		task.Status = "completed"
		task.CompletedAt = &now
		database.DB.Save(&task)

	case "habit_done", "habit_weekly_done":
		var habit models.Habit
		if err := database.DB.First(&habit, idStr).Error; err != nil {
			return err
		}
		// Log habit completion
		log := models.HabitLog{
			HabitID:   habit.ID,
			UserID:    habit.UserID,
			LogDate:   time.Now(),
			Completed: true,
		}
		database.DB.Create(&log)

		// Update streak (simplified)
		habit.CurrentStreak++
		if habit.CurrentStreak > habit.BestStreak {
			habit.BestStreak = habit.CurrentStreak
		}
		database.DB.Save(&habit)

	case "goal_dismiss":
		// Just delete the message, which happens below

	}

	// Always delete the message after successful action
	return s.telegramService.DeleteMessage(telegramID, messageID)
}

// SyncTaskCompletion deletes any active notifications for a task that was completed elsewhere
func (s *NotificationService) SyncTaskCompletion(taskID uint) error {
	var notifications []models.Notification
	database.DB.Where("type = ? AND related_id = ? AND status = ?", "task_due", taskID, "sent").Find(&notifications)

	for _, n := range notifications {
		var user models.User
		if err := database.DB.First(&user, n.UserID).Error; err == nil {
			s.telegramService.DeleteMessage(*user.TelegramID, n.TelegramMessageID)
		}
		n.Status = "deleted"
		database.DB.Save(&n)
	}
	return nil
}

// SyncHabitCompletion deletes any active notifications for a habit that was logged elsewhere
func (s *NotificationService) SyncHabitCompletion(habitID uint) error {
	var notifications []models.Notification
	database.DB.Where("type IN (?) AND related_id = ? AND status = ?", []string{"habit_reminder", "habit_weekly"}, habitID, "sent").Find(&notifications)

	for _, n := range notifications {
		var user models.User
		if err := database.DB.First(&user, n.UserID).Error; err == nil {
			s.telegramService.DeleteMessage(*user.TelegramID, n.TelegramMessageID)
		}
		n.Status = "deleted"
		database.DB.Save(&n)
	}
	return nil
}

func getProjectName(task *models.Task) string {
	if task.Project != nil {
		return task.Project.Name
	}
	return "No Project"
}

// SendGoalProgressNotification sends a notification for goal progress
func (s *NotificationService) SendGoalProgressNotification(user *models.User, goal *models.Goal) error {
	if user.TelegramID == nil {
		return nil
	}

	progressBar := generateProgressBar(float64(goal.Progress))

	text := fmt.Sprintf("🎉 *Goal Progress Update!*\n━━━━━━━━━━━━━━━━━━\nYour goal *%s* is now:\n\n%s *%d%% Complete*\n\n━━━━━━━━━━━━━━━━━━",
		s.telegramService.EscapeMarkdownV2(goal.Title),
		progressBar,
		goal.Progress,
	)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🗑️ Dismiss", "callback_data": fmt.Sprintf("goal_dismiss:%d", goal.ID)},
			},
		},
	}

	msgID, err := s.telegramService.SendMessage(*user.TelegramID, text, replyMarkup)
	if err != nil {
		return err
	}

	notification := models.Notification{
		UserID:            user.ID,
		Type:              "goal_progress",
		RelatedID:         goal.ID,
		TelegramMessageID: msgID,
		SentAt:            time.Now(),
	}

	return database.DB.Create(&notification).Error
}

// SendGoalTargetDateNotification sends a reminder for goal deadline
func (s *NotificationService) SendGoalTargetDateNotification(user *models.User, goal *models.Goal) error {
	if user.TelegramID == nil {
		return nil
	}

	text := fmt.Sprintf("🎯 *Goal Target Date Alert*\n━━━━━━━━━━━━━━━━━━\nYour goal *%s* is due on *%s*.\n\nCurrent Progress: *%d%%*\n\n━━━━━━━━━━━━━━━━━━",
		s.telegramService.EscapeMarkdownV2(goal.Title),
		goal.TargetDate.Format("2006-01-02"),
		goal.Progress,
	)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🗑️ Dismiss", "callback_data": fmt.Sprintf("goal_dismiss:%d", goal.ID)},
			},
		},
	}

	msgID, err := s.telegramService.SendMessage(*user.TelegramID, text, replyMarkup)
	if err != nil {
		return err
	}

	notification := models.Notification{
		UserID:            user.ID,
		Type:              "goal_deadline",
		RelatedID:         goal.ID,
		TelegramMessageID: msgID,
		SentAt:            time.Now(),
	}

	return database.DB.Create(&notification).Error
}

// SendGoalCompletedNotification sends a celebration for goal completion
func (s *NotificationService) SendGoalCompletedNotification(user *models.User, goal *models.Goal) error {
	if user.TelegramID == nil {
		return nil
	}

	text := fmt.Sprintf("🏆 *Goal Completed!*\n━━━━━━━━━━━━━━━━━━\nCongratulations! You've achieved your goal:\n\n*%s*\n\nGreat work! 🌟\n━━━━━━━━━━━━━━━━━━",
		s.telegramService.EscapeMarkdownV2(goal.Title),
	)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "🗑️ Dismiss", "callback_data": fmt.Sprintf("goal_dismiss:%d", goal.ID)},
			},
		},
	}

	msgID, err := s.telegramService.SendMessage(*user.TelegramID, text, replyMarkup)
	if err != nil {
		return err
	}

	notification := models.Notification{
		UserID:            user.ID,
		Type:              "goal_completed",
		RelatedID:         goal.ID,
		TelegramMessageID: msgID,
		SentAt:            time.Now(),
	}

	return database.DB.Create(&notification).Error
}

func generateProgressBar(progress float64) string {
	totalBars := 10
	filledBars := int(progress / 10)
	if filledBars > totalBars {
		filledBars = totalBars
	}

	var sb strings.Builder
	for i := 0; i < filledBars; i++ {
		sb.WriteString("▓")
	}
	for i := filledBars; i < totalBars; i++ {
		sb.WriteString("░")
	}
	return sb.String()
}
