package services

import (
	"math"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
)

// UpdateProjectProgress recalculates and updates the progress of a project based on its goals
func UpdateProjectProgress(projectID uint) error {
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		return err
	}

	var totalGoals int64
	var totalProgress float64

	// Get all goals for this project
	var goals []models.Goal
	if err := database.DB.Where("project_id = ?", projectID).Find(&goals).Error; err != nil {
		return err
	}

	totalGoals = int64(len(goals))

	if totalGoals == 0 {
		project.Progress = 0
	} else {
		for _, goal := range goals {
			totalProgress += float64(goal.Progress)
		}

		progress := totalProgress / float64(totalGoals)
		project.Progress = int(math.Round(progress))
	}

	// Update the project
	return database.DB.Model(&project).Update("progress", project.Progress).Error
}
