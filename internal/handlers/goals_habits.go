package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-gonic/gin"
)

type GoalHandler struct{}

func NewGoalHandler() *GoalHandler {
	return &GoalHandler{}
}

func (h *GoalHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var goals []models.Goal

	query := database.DB.Select("id, title, description, status, progress, project_id, user_id, slug, target_date").Where("user_id = ?", userID)

	// Filter by project if provided
	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&goals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	param := c.Param("id")

	var goal models.Goal
	query := database.DB.Preload("Milestones").
		Preload("Milestones.Tasks").
		Preload("Habits").
		Where("user_id = ?", userID)

	if id, err := strconv.Atoi(param); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", param)
	}

	if err := query.First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var goal models.Goal

	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal.UserID = userID
	if err := database.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	LogActivity(userID, "Created", "Goal", goal.ID, goal.Title)

	if goal.ProjectID != nil {
		go services.UpdateProjectProgress(*goal.ProjectID)
	}

	c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	if err := database.DB.Delete(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goal"})
		return
	}

	LogActivity(userID, "Deleted", "Goal", uint(id), goal.Title)

	if goal.ProjectID != nil {
		go services.UpdateProjectProgress(*goal.ProjectID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted"})
}

func (h *GoalHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goal"})
		return
	}

	LogActivity(userID, "Updated", "Goal", goal.ID, goal.Title)

	if goal.ProjectID != nil {
		go services.UpdateProjectProgress(*goal.ProjectID)
	}

	c.JSON(http.StatusOK, goal)
}

type HabitHandler struct {
	notificationService *services.NotificationService
}

func NewHabitHandler() *HabitHandler {
	return &HabitHandler{
		notificationService: services.NewNotificationService(),
	}
}

func (h *HabitHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var habits []models.Habit

	// Only preload logs from the last 30 days to improve performance
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := database.DB.Select("id, name, description, frequency, goal_id, user_id, current_streak, best_streak, is_active, scheduled_days").Preload("Logs", "log_date >= ?", thirtyDaysAgo).Where("user_id = ?", userID).Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch habits"})
		return
	}

	c.JSON(http.StatusOK, habits)
}

func (h *HabitHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var habit models.Habit

	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	habit.UserID = userID
	if err := database.DB.Create(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create habit"})
		return
	}

	LogActivity(userID, "Created", "Habit", habit.ID, habit.Name)

	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var habit models.Habit
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := database.DB.Preload("Logs", "log_date >= ?", thirtyDaysAgo).Where("id = ? AND user_id = ?", id, userID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Log(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var log models.HabitLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.LogDate = log.LogDate.Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)

	if !log.LogDate.Equal(today) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Habits can only be logged for today"})
		return
	}

	// Verify habit belongs to user
	var habit models.Habit
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// Prevent duplicate logs for the same day
	var existingLogCount int64
	database.DB.Model(&models.HabitLog{}).Where("habit_id = ? AND log_date = ?", id, log.LogDate).Count(&existingLogCount)
	if existingLogCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Habit is already checked in for today"})
		return
	}

	log.HabitID = uint(id)
	log.UserID = userID
	if err := database.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log habit"})
		return
	}

	LogActivity(userID, "Logged", "Habit", uint(id), habit.Name)
	go h.notificationService.SyncHabitCompletion(uint(id))
	go services.CalculateAndUpdateStreak(uint(id))

	c.JSON(http.StatusCreated, log)
}

func (h *HabitHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var habit models.Habit
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	if err := database.DB.Delete(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete habit"})
		return
	}

	LogActivity(userID, "Deleted", "Habit", uint(id), habit.Name)
	go h.notificationService.SyncHabitCompletion(uint(id))

	c.JSON(http.StatusOK, gin.H{"message": "Habit deleted"})
}

func (h *HabitHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var habit models.Habit
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update habit"})
		return
	}

	LogActivity(userID, "Updated", "Habit", habit.ID, habit.Name)

	c.JSON(http.StatusOK, habit)
}
