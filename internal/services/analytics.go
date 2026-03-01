package services

import (
	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
)

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

func (s *AnalyticsService) GetUserAnalytics(userID uint) (models.AnalyticsDashboardData, error) {
	var data models.AnalyticsDashboardData
	db := database.DB

	var completedTasks int64
	db.Model(&models.Task{}).Where("user_id = ? AND status = 'completed'", userID).Count(&completedTasks)
	data.Tasks.TotalCompleted = int(completedTasks)

	var pendingTasks int64
	db.Model(&models.Task{}).Where("user_id = ? AND status != 'completed'", userID).Count(&pendingTasks)
	data.Tasks.TotalPending = int(pendingTasks)

	// Habit Stats
	var activeHabits int64
	db.Model(&models.Habit{}).Where("user_id = ? AND is_active = true", userID).Count(&activeHabits)
	data.Habits.ActiveHabits = int(activeHabits)

	var totalLogs int64
	db.Model(&models.HabitLog{}).Where("user_id = ? AND completed = true", userID).Count(&totalLogs)
	data.Habits.TotalLogs = int(totalLogs)

	type Result struct {
		MaxStreak int
	}
	var res Result
	db.Model(&models.Habit{}).Where("user_id = ?", userID).Select("COALESCE(MAX(best_streak), 0) as max_streak").Scan(&res)
	data.Habits.HighestStreak = res.MaxStreak

	// Pomodoro Stats
	var completedSessions int64
	db.Model(&models.PomodoroSession{}).Where("user_id = ? AND status = 'completed'", userID).Count(&completedSessions)
	data.Pomodoros.TotalSessionsCompleted = int(completedSessions)

	var focusMinutes struct {
		TotalMinutes int
	}
	db.Model(&models.PomodoroSession{}).Where("user_id = ? AND status = 'completed'", userID).Select("COALESCE(SUM(work_duration * completed_cycles), 0) as total_minutes").Scan(&focusMinutes)
	data.Pomodoros.TotalFocusMinutes = focusMinutes.TotalMinutes

	var progress models.UserPomodoroProgress
	db.Where("user_id = ?", userID).Attrs(models.UserPomodoroProgress{
		CurrentLevel: 1,
	}).FirstOrCreate(&progress)

	data.Pomodoros.CurrentLevel = progress.CurrentLevel

	if progress.CurrentLevel >= 1 && progress.CurrentLevel <= 9 {
		data.Pomodoros.CurrentPhase = "Normal"
	} else if progress.CurrentLevel >= 20 && progress.CurrentLevel <= 29 {
		data.Pomodoros.CurrentPhase = "Focus"
	} else if progress.CurrentLevel >= 30 && progress.CurrentLevel <= 39 {
		data.Pomodoros.CurrentPhase = "Deep Work"
	} else if progress.CurrentLevel >= 40 && progress.CurrentLevel <= 49 {
		data.Pomodoros.CurrentPhase = "Hard Work"
	} else {
		data.Pomodoros.CurrentPhase = "Unknown"
	}

	// Goal Stats
	var activeGoals int64
	db.Model(&models.Goal{}).Where("user_id = ? AND status = 'active'", userID).Count(&activeGoals)
	data.Goals.ActiveGoals = int(activeGoals)

	var completedGoals int64
	db.Model(&models.Goal{}).Where("user_id = ? AND status = 'completed'", userID).Count(&completedGoals)
	data.Goals.CompletedGoals = int(completedGoals)

	return data, nil
}
