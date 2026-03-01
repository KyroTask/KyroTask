package handlers

import (
	"net/http"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct{}

func NewActivityHandler() *ActivityHandler {
	return &ActivityHandler{}
}

func (h *ActivityHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var activities []models.Activity

	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(50).Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// LogActivity is a helper to log activities from other handlers
func LogActivity(userID uint, action, resourceType string, resourceID uint, description string) {
	activity := models.Activity{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Description:  description,
	}
	database.DB.Create(&activity)
}
