package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bif12/kyrotask/internal/database"
	"github.com/bif12/kyrotask/internal/middleware"
	"github.com/bif12/kyrotask/internal/models"
	"github.com/gin-gonic/gin"
)

type PomodoroHandler struct{}

func NewPomodoroHandler() *PomodoroHandler {
	return &PomodoroHandler{}
}

// GetStatus returns the user's current level, timer state, and seconds remaining.
func (h *PomodoroHandler) GetStatus(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var progress models.UserPomodoroProgress
	if err := database.DB.Where("user_id = ?", userID).Attrs(models.UserPomodoroProgress{
		CurrentLevel: 1,
	}).FirstOrCreate(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	canStartToday := true
	if progress.LastLevelCompletedAt != nil {
		today := time.Now().Truncate(24 * time.Hour)
		lastDate := progress.LastLevelCompletedAt.Truncate(24 * time.Hour)
		if lastDate.Equal(today) || lastDate.After(today) {
			canStartToday = false
		}
	}

	workDur, breakDur, targets, reqProj := models.GetLevelConfig(progress.CurrentLevel)

	// Check if there is an active session
	var activeSession *models.PomodoroSession
	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err == nil {
		// Auto-repair sessions from before CycleStartedAt was added
		if session.CycleStartedAt == nil {
			session.CycleStartedAt = &session.CreatedAt
			database.DB.Save(&session)
		}
		activeSession = &session
	}

	resp := gin.H{
		"current_level":           progress.CurrentLevel,
		"can_start_today":         canStartToday,
		"required_work_duration":  workDur,
		"required_break_duration": breakDur,
		"target_cycles":           targets,
		"requires_project":        reqProj,
		"active_session":          activeSession,
	}

	// If there's an active session, include computed timer info
	if activeSession != nil {
		resp["seconds_remaining"] = activeSession.SecondsRemaining()
		resp["is_on_break"] = activeSession.IsOnBreak
		resp["is_paused"] = activeSession.IsPaused
	}

	c.JSON(http.StatusOK, resp)
}

type StartSessionRequest struct {
	ProjectID *uint `json:"project_id"`
}

// StartSession begins a new Pomodoro session for the current level.
func (h *PomodoroHandler) StartSession(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var progress models.UserPomodoroProgress
	if err := database.DB.Where("user_id = ?", userID).First(&progress).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Progress not found"})
		return
	}

	// Verify they are allowed to start today
	if progress.LastLevelCompletedAt != nil {
		today := time.Now().Truncate(24 * time.Hour)
		lastDate := progress.LastLevelCompletedAt.Truncate(24 * time.Hour)
		if lastDate.Equal(today) || lastDate.After(today) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You have already completed a level today. Rest until tomorrow."})
			return
		}
	}

	// Check if a session is already in progress
	var count int64
	database.DB.Model(&models.PomodoroSession{}).Where("user_id = ? AND status = ?", userID, "in_progress").Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A session is already in progress"})
		return
	}

	workDur, breakDur, targets, reqProj := models.GetLevelConfig(progress.CurrentLevel)

	if reqProj && req.ProjectID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This level requires selecting a Project"})
		return
	}

	now := time.Now()
	session := models.PomodoroSession{
		UserID:          userID,
		Level:           progress.CurrentLevel,
		ProjectID:       req.ProjectID,
		TargetCycles:    targets,
		CompletedCycles: 0,
		Status:          "in_progress",
		WorkDuration:    workDur,
		BreakDuration:   breakDur,
		CycleStartedAt:  &now,
		IsOnBreak:       false,
		IsPaused:        false,
		PausedElapsed:   0,
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Log Activity
	actionDesc := "Started Focus Level " + strconv.Itoa(progress.CurrentLevel)
	LogActivity(userID, "Started", "Pomodoro", session.ID, actionDesc)

	go BroadcastPomodoroStatus(userID)
	c.JSON(http.StatusCreated, gin.H{
		"session":           session,
		"seconds_remaining": session.SecondsRemaining(),
	})

	// Broadcast to all user's tabs
	go BroadcastPomodoroStatus(userID)
}

// LogCycle increments the completed cycle count and starts a break.
func (h *PomodoroHandler) LogCycle(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	now := time.Now()
	session.CompletedCycles += 1
	session.IsOnBreak = true
	session.CycleStartedAt = &now
	session.PausedElapsed = 0
	session.IsPaused = false
	session.PausedAt = nil

	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cycle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session":           session,
		"seconds_remaining": session.SecondsRemaining(),
	})

	go BroadcastPomodoroStatus(userID)
}

// ResumeWork is called when a break ends and the user starts the next work cycle.
func (h *PomodoroHandler) ResumeWork(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	if !session.IsOnBreak {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not currently on a break"})
		return
	}

	now := time.Now()
	session.IsOnBreak = false
	session.CycleStartedAt = &now
	session.PausedElapsed = 0
	session.IsPaused = false
	session.PausedAt = nil

	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resume work"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session":           session,
		"seconds_remaining": session.SecondsRemaining(),
	})

	go BroadcastPomodoroStatus(userID)
}

// PauseSession pauses the timer.
func (h *PomodoroHandler) PauseSession(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	if session.IsPaused {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already paused"})
		return
	}

	now := time.Now()
	session.IsPaused = true
	session.PausedAt = &now
	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pause"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session":           session,
		"seconds_remaining": session.SecondsRemaining(),
	})

	go BroadcastPomodoroStatus(userID)
}

// ResumeSession unpauses the timer.
func (h *PomodoroHandler) ResumeSession(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	if !session.IsPaused {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not paused"})
		return
	}

	// Add the paused duration to PausedElapsed
	if session.PausedAt != nil {
		session.PausedElapsed += int(time.Since(*session.PausedAt).Seconds())
	}
	session.IsPaused = false
	session.PausedAt = nil

	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session":           session,
		"seconds_remaining": session.SecondsRemaining(),
	})

	go BroadcastPomodoroStatus(userID)
}

// CompleteLevel is called when the target iterations are finally met. It levels the user up.
func (h *PomodoroHandler) CompleteLevel(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	if session.CompletedCycles < session.TargetCycles {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot complete level before reaching target cycles"})
		return
	}

	// Update user progress (Level Up)
	var progress models.UserPomodoroProgress
	if err := database.DB.Where("user_id = ?", userID).First(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	now := time.Now()
	progress.LastLevelCompletedAt = &now

	// Level progression jump logic based on phase boundaries
	// The instruction snippet simplifies this, so we'll follow that.
	progress.CurrentLevel++
	switch progress.CurrentLevel {
	case 10:
		progress.CurrentLevel = 20
	case 30:
		progress.CurrentLevel = 30
	case 40:
		progress.CurrentLevel = 40
	case 50:
		progress.CurrentLevel = 50
	}

	if err := database.DB.Save(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
		return
	}

	// Update the session state to completed
	session.Status = "completed"
	session.FinishedAt = &now
	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark session as completed"})
		return
	}

	// Log Activity
	actionDesc := "Completed Focus Level " + strconv.Itoa(session.Level)
	LogActivity(userID, "Completed", "Pomodoro", session.ID, actionDesc)

	go BroadcastPomodoroStatus(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Level completed successfully"})
}

// AbandonSession cancels the current session and allows them to try again right away.
func (h *PomodoroHandler) AbandonSession(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var session models.PomodoroSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "in_progress").First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active session found"})
		return
	}

	session.Status = "abandoned"
	now := time.Now()
	session.FinishedAt = &now
	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to abandon session"})
		return
	}

	// Log Activity
	actionDesc := "Abandoned Focus Level " + strconv.Itoa(session.Level)
	LogActivity(userID, "Abandoned", "Pomodoro", session.ID, actionDesc)

	c.JSON(http.StatusOK, gin.H{"message": "Session abandoned"})

	go BroadcastPomodoroStatus(userID)
}
