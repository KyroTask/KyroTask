package handlers

import (
	"net/http"
	"strconv"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/bif12/kyrotask/internal/services"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	notificationService *services.NotificationService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		notificationService: services.NewNotificationService(),
	}
}

func (h *TaskHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var tasks []models.Task

	query := database.DB.Select("id, title, status, priority, due_date, project_id, goal_id, user_id, slug, milestone_id").Preload("Project").Preload("Goal").Where("user_id = ?", userID)

	// Filter by project if provided
	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	// Filter by goal if provided
	if goalID := c.Query("goal_id"); goalID != "" {
		query = query.Where("goal_id = ?", goalID)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UserID = userID
	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	LogActivity(userID, "Created", "Task", task.ID, task.Title)

	if task.GoalID != nil {
		go services.UpdateGoalProgress(*task.GoalID)
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	param := c.Param("id")

	var task models.Task
	query := database.DB.Preload("Project").Preload("Subtasks").Where("user_id = ?", userID)

	if id, err := strconv.Atoi(param); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", param)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	param := c.Param("id")

	var task models.Task
	query := database.DB.Where("user_id = ?", userID)

	if id, err := strconv.Atoi(param); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", param)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	oldStatus := task.Status
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&task)

	if oldStatus != "completed" && task.Status == "completed" {
		LogActivity(userID, "Completed", "Task", task.ID, task.Title)
		go h.notificationService.SyncTaskCompletion(task.ID)
	} else {
		LogActivity(userID, "Updated", "Task", task.ID, task.Title)
	}

	if task.GoalID != nil {
		go services.UpdateGoalProgress(*task.GoalID)
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	param := c.Param("id")

	var task models.Task
	query := database.DB.Where("user_id = ?", userID)

	if id, err := strconv.Atoi(param); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", param)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	LogActivity(userID, "Deleted", "Task", task.ID, task.Title)
	go h.notificationService.SyncTaskCompletion(task.ID)

	if task.GoalID != nil {
		go services.UpdateGoalProgress(*task.GoalID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
