package services

import (
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
)

// CalculateAndUpdateStreak calculates the current streak for a habit and updates it in the DB
func CalculateAndUpdateStreak(habitID uint) error {
	var habit models.Habit
	if err := database.DB.First(&habit, habitID).Error; err != nil {
		return err
	}

	// Get all logs for this habit, ordered by date descending
	var logs []models.HabitLog
	if err := database.DB.Where("habit_id = ? AND completed = ?", habitID, true).
		Order("log_date DESC").Find(&logs).Error; err != nil {
		return err
	}

	if len(logs) == 0 {
		habit.CurrentStreak = 0
		return database.DB.Model(&habit).Updates(map[string]interface{}{
			"current_streak": 0,
		}).Error
	}

	// Calculate current streak
	streak := 0
	today := time.Now().Truncate(24 * time.Hour)

	for i, log := range logs {
		logDate := log.LogDate.Truncate(24 * time.Hour)
		expectedDate := today.AddDate(0, 0, -i)

		if logDate.Equal(expectedDate) {
			streak++
		} else if i == 0 && logDate.Equal(today.AddDate(0, 0, -1)) {
			// Allow checking from yesterday if not logged today yet
			streak++
			today = today.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	habit.CurrentStreak = streak
	if streak > habit.BestStreak {
		habit.BestStreak = streak
	}

	return database.DB.Model(&habit).Updates(map[string]interface{}{
		"current_streak": habit.CurrentStreak,
		"best_streak":    habit.BestStreak,
	}).Error
}
