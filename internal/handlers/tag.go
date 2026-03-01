package handlers

import (
	"net/http"
	"strconv"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

func (h *TagHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var tags []models.Tag

	if err := database.DB.Where("user_id = ?", userID).Order("name ASC").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var tag models.Tag

	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.UserID = userID

	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Tag{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
}
