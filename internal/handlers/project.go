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

type ProjectHandler struct{}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

func (h *ProjectHandler) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var projects []models.Project

	// Optimize by selecting only necessary fields. Progress is now stored in the DB.
	if err := database.DB.Select("id, name, slug, description, status, progress, user_id, color, icon, is_archived, created_at, updated_at").Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var project models.Project

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.UserID = userID
	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	LogActivity(userID, "Created", "Project", project.ID, project.Name)

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	param := c.Param("id")

	var project models.Project
	var err error

	if id, convErr := strconv.Atoi(param); convErr == nil {
		err = database.DB.Preload("Goals").Preload("Goals.Tasks").Where("id = ? AND user_id = ?", id, userID).First(&project).Error
	} else {
		err = database.DB.Preload("Goals").Preload("Goals.Tasks").Where("slug = ? AND user_id = ?", param, userID).First(&project).Error
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.Project
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&project)
	LogActivity(userID, "Updated", "Project", project.ID, project.Name)
	go services.UpdateProjectProgress(project.ID)
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.Project
	database.DB.Where("id = ? AND user_id = ?", id, userID).First(&project)

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Project{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	LogActivity(userID, "Deleted", "Project", uint(id), project.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}
