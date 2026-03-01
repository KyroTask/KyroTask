package handlers

import (
	"net/http"
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

type DashboardResponse struct {
	TaskStats         TaskStats         `json:"task_stats"`
	HabitStats        HabitStats        `json:"habit_stats"`
	ActiveGoals       []models.Goal     `json:"active_goals"`
	RecentActivity    []models.Activity `json:"recent_activity"`
	ProductivityScore int               `json:"productivity_score"`
}

type TaskStats struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Completed int64 `json:"completed"`
	Overdue   int64 `json:"overdue"`
}

type HabitStats struct {
	Total        int64 `json:"total"`
	DoneToday    int64 `json:"done_today"`
	PendingToday int64 `json:"pending_today"`
	TopStreak    int   `json:"top_streak"`
}

func (h *DashboardHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))

	var resp DashboardResponse

	// Task stats
	database.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&resp.TaskStats.Total)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&resp.TaskStats.Pending)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, "completed").Count(&resp.TaskStats.Completed)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status != ? AND due_date < ?", userID, "completed", now).Count(&resp.TaskStats.Overdue)

	// Habit stats
	var habits []models.Habit
	database.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&habits)
	resp.HabitStats.Total = int64(len(habits))

	for _, habit := range habits {
		var count int64
		database.DB.Model(&models.HabitLog{}).Where("habit_id = ? AND log_date >= ?", habit.ID, startOfDay).Count(&count)
		if count > 0 {
			resp.HabitStats.DoneToday++
		} else {
			resp.HabitStats.PendingToday++
		}
		if habit.BestStreak > resp.HabitStats.TopStreak {
			resp.HabitStats.TopStreak = habit.BestStreak
		}
	}

	// Active goals (top 5)
	database.DB.Where("user_id = ? AND status = ?", userID, "active").
		Order("updated_at DESC").Limit(5).Find(&resp.ActiveGoals)

	// Recent activity (last 5)
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(5).Find(&resp.RecentActivity)

	// Productivity score: % of tasks completed this week
	var weekTotal int64
	var weekCompleted int64
	database.DB.Model(&models.Task{}).Where("user_id = ? AND created_at >= ?", userID, startOfWeek).Count(&weekTotal)
	database.DB.Model(&models.Task{}).Where("user_id = ? AND status = ? AND completed_at >= ?", userID, "completed", startOfWeek).Count(&weekCompleted)
	if weekTotal > 0 {
		resp.ProductivityScore = int(float64(weekCompleted) / float64(weekTotal) * 100)
	}

	c.JSON(http.StatusOK, resp)
}
