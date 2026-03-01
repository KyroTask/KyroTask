package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type MilestoneHandler struct{}

func NewMilestoneHandler() *MilestoneHandler {
	return &MilestoneHandler{}
}

// verifyMilestoneOwnership checks the milestone's goal belongs to the user
func verifyMilestoneOwnership(milestoneID uint, userID uint) (*models.Milestone, error) {
	var milestone models.Milestone
	if err := database.DB.Preload("Goal").First(&milestone, milestoneID).Error; err != nil {
		return nil, err
	}
	if milestone.Goal.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return &milestone, nil
}

func (h *MilestoneHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	goalIDStr := c.Query("goal_id")

	if goalIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "goal_id query parameter is required"})
		return
	}

	goalID, _ := strconv.Atoi(goalIDStr)

	// Verify goal belongs to user
	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	var milestones []models.Milestone
	if err := database.DB.Preload("Tasks").Where("goal_id = ?", goalID).Find(&milestones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch milestones"})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

func (h *MilestoneHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var milestone models.Milestone

	if err := c.ShouldBindJSON(&milestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify goal belongs to user
	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", milestone.GoalID, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	if err := database.DB.Create(&milestone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create milestone"})
		return
	}

	LogActivity(userID, "Created", "Milestone", milestone.ID, milestone.Title)
	c.JSON(http.StatusCreated, milestone)
}

func (h *MilestoneHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	milestone, err := verifyMilestoneOwnership(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found or unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(milestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(milestone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update milestone"})
		return
	}

	LogActivity(userID, "Updated", "Milestone", milestone.ID, milestone.Title)
	c.JSON(http.StatusOK, milestone)
}

func (h *MilestoneHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	milestone, err := verifyMilestoneOwnership(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found or unauthorized"})
		return
	}

	if err := database.DB.Delete(milestone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete milestone"})
		return
	}

	LogActivity(userID, "Deleted", "Milestone", uint(id), milestone.Title)
	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted"})
}
