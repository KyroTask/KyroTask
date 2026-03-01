package services

import (
	"log"
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
)

type Scheduler struct {
	notificationService *NotificationService
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		notificationService: NewNotificationService(),
	}
}

func (s *Scheduler) Start() {
	log.Println("Starting scheduler...")
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			s.checkTasksDueSoon()
			s.checkHabitReminders()
			s.checkGoalDeadlines()
		}
	}()
}

func (s *Scheduler) checkTasksDueSoon() {
	now := time.Now()
	oneHourLater := now.Add(1 * time.Hour)

	// Find tasks due between now and 1 hour + 5 mins (buffer) that haven't been notified
	var tasks []models.Task
	// This is a simplified query. In production, we'd check the Notification table to ensure we haven't sent this specific notification yet.
	// For now, let's assume we check if due_date is within the specific minute window to avoid duplicates,
	// or better, check if a notification exists.

	startWindow := oneHourLater.Truncate(time.Minute)
	endWindow := startWindow.Add(time.Minute)

	if err := database.DB.Preload("Project").Where("due_date >= ? AND due_date < ? AND status != ?", startWindow, endWindow, "completed").Find(&tasks).Error; err != nil {
		log.Printf("Error checking due tasks: %v", err)
		return
	}

	for _, task := range tasks {
		// Check if already notified
		var count int64
		database.DB.Model(&models.Notification{}).Where("type = ? AND related_id = ?", "task_due", task.ID).Count(&count)
		if count > 0 {
			continue
		}

		var user models.User
		if err := database.DB.First(&user, task.UserID).Error; err != nil {
			continue
		}

		if err := s.notificationService.SendTaskNotification(&user, &task); err != nil {
			log.Printf("Failed to send notification for task %d: %v", task.ID, err)
		}
	}
}

func (s *Scheduler) checkHabitReminders() {
	now := time.Now()
	// Default reminder at 9:00 PM
	if now.Hour() == 21 && now.Minute() == 0 {
		var habits []models.Habit
		// Find active habits
		if err := database.DB.Where("is_active = ?", true).Find(&habits).Error; err != nil {
			log.Printf("Error checking habits: %v", err)
			return
		}

		for _, habit := range habits {
			// Check if already logged today
			var count int64
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			database.DB.Model(&models.HabitLog{}).Where("habit_id = ? AND log_date >= ?", habit.ID, startOfDay).Count(&count)
			if count > 0 {
				continue // Already done
			}

			var user models.User
			if err := database.DB.First(&user, habit.UserID).Error; err != nil {
				continue
			}

			if err := s.notificationService.SendHabitNotification(&user, &habit); err != nil {
				log.Printf("Failed to send notification for habit %d: %v", habit.ID, err)
			}
		}
	}
}

func (s *Scheduler) checkGoalDeadlines() {
	now := time.Now()
	oneWeekLater := now.AddDate(0, 0, 7)

	startWindow := oneWeekLater.Truncate(24 * time.Hour)
	endWindow := startWindow.Add(24 * time.Hour)

	var goals []models.Goal
	if err := database.DB.Where("target_date >= ? AND target_date < ? AND status = ?", startWindow, endWindow, "active").Find(&goals).Error; err != nil {
		log.Printf("Error checking goal deadlines: %v", err)
		return
	}

	for _, goal := range goals {
		// Check if already notified
		var count int64
		database.DB.Model(&models.Notification{}).Where("type = ? AND related_id = ?", "goal_deadline", goal.ID).Count(&count)
		if count > 0 {
			continue
		}

		var user models.User
		if err := database.DB.First(&user, goal.UserID).Error; err != nil {
			continue
		}

		if err := s.notificationService.SendGoalTargetDateNotification(&user, &goal); err != nil {
			log.Printf("Failed to send notification for goal %d: %v", goal.ID, err)
		}
	}
}
