package handlers

import (
	"net/http"
	"strconv"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

// ListByTask returns all comments for a task
func (h *CommentHandler) ListByTask(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	taskID, _ := strconv.Atoi(c.Param("id"))

	// Verify task belongs to user
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var comments []models.Comment
	if err := database.DB.Preload("User").Where("task_id = ?", taskID).Order("created_at ASC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Create adds a comment to a task
func (h *CommentHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	taskID, _ := strconv.Atoi(c.Param("id"))

	// Verify task belongs to user
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.TaskID = uint(taskID)
	comment.UserID = userID

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	LogActivity(userID, "Commented", "Task", uint(taskID), comment.Content)
	c.JSON(http.StatusCreated, comment)
}

// Delete removes a comment
func (h *CommentHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("comment_id"))

	var comment models.Comment
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
