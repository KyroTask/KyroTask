package services

import (
	"math"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/models"
)

// UpdateGoalProgress recalculates and updates the progress of a goal based on its tasks
func UpdateGoalProgress(goalID uint) error {
	var goal models.Goal
	if err := database.DB.First(&goal, goalID).Error; err != nil {
		return err
	}

	oldProgress := goal.Progress

	var totalTasks int64
	var completedTasks int64

	// Count total tasks for this goal
	if err := database.DB.Model(&models.Task{}).Where("goal_id = ?", goalID).Count(&totalTasks).Error; err != nil {
		return err
	}

	if totalTasks == 0 {
		// No tasks, progress is 0 (or keep manual if we supported that, but let's say 0 for now)
		// Or maybe 100 if it was manually marked? Let's stick to task-based.
		goal.Progress = 0
	} else {
		// Count completed tasks
		if err := database.DB.Model(&models.Task{}).Where("goal_id = ? AND status = ?", goalID, "completed").Count(&completedTasks).Error; err != nil {
			return err
		}

		progress := float64(completedTasks) / float64(totalTasks) * 100
		goal.Progress = int(math.Round(progress))
	}

	// Update the goal
	if err := database.DB.Model(&goal).Update("progress", goal.Progress).Error; err != nil {
		return err
	}

	// Trigger Notifications
	if goal.Progress > oldProgress {
		ns := NewNotificationService()
		var user models.User
		if err := database.DB.First(&user, goal.UserID).Error; err == nil {
			if goal.Progress == 100 {
				ns.SendGoalCompletedNotification(&user, &goal)
			} else {
				ns.SendGoalProgressNotification(&user, &goal)
			}
		}
	}

	// If goal belongs to a project, update project progress
	if goal.ProjectID != nil {
		return UpdateProjectProgress(*goal.ProjectID)
	}

	return nil
}
