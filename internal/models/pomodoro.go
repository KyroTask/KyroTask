package models

import "time"

// UserPomodoroProgress tracks the global state of the user in the Pomodoro journey.
type UserPomodoroProgress struct {
	UserID               uint       `json:"user_id" gorm:"primaryKey"`
	CurrentLevel         int        `json:"current_level" gorm:"default:1"`
	LastLevelCompletedAt *time.Time `json:"last_level_completed_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// PomodoroSession tracks individual work sessions for a specific day/level.
type PomodoroSession struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	UserID          uint   `json:"user_id" gorm:"index:idx_pomo_user_status"`
	Level           int    `json:"level"`
	ProjectID       *uint  `json:"project_id"`
	TargetCycles    int    `json:"target_cycles"`
	CompletedCycles int    `json:"completed_cycles" gorm:"default:0"`
	Status          string `json:"status" gorm:"default:'in_progress';index:idx_pomo_user_status"` // in_progress, completed, abandoned
	WorkDuration    int    `json:"work_duration"`                                                  // in minutes
	BreakDuration   int    `json:"break_duration"`                                                 // in minutes

	// Timer state — persisted so we can resume across refreshes/devices
	CycleStartedAt *time.Time `json:"cycle_started_at"` // when the current work/break timer started
	IsOnBreak      bool       `json:"is_on_break" gorm:"default:false"`
	IsPaused       bool       `json:"is_paused" gorm:"default:false"`
	PausedAt       *time.Time `json:"paused_at"`      // when pause was triggered
	PausedElapsed  int        `json:"paused_elapsed"` // total seconds spent paused during current cycle

	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

// SecondsRemaining computes how many seconds are left in the current cycle.
func (s *PomodoroSession) SecondsRemaining() int {
	if s.CycleStartedAt == nil {
		if s.IsOnBreak {
			return s.BreakDuration * 60
		}
		return s.WorkDuration * 60
	}

	// Total duration for this cycle
	totalSecs := s.WorkDuration * 60
	if s.IsOnBreak {
		totalSecs = s.BreakDuration * 60
	}

	// Elapsed = time since cycle started - paused time
	var elapsed int
	if s.IsPaused && s.PausedAt != nil {
		// Timer is paused: count elapsed up to the pause moment
		elapsed = int(s.PausedAt.Sub(*s.CycleStartedAt).Seconds()) - s.PausedElapsed
	} else {
		elapsed = int(time.Since(*s.CycleStartedAt).Seconds()) - s.PausedElapsed
	}

	remaining := totalSecs - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetLevelConfig returns the required Work/Break durations and Target Cycles based on the Level phase.
func GetLevelConfig(level int) (workDuration int, breakDuration int, targetCycles int, requiresProject bool) {
	if level >= 1 && level <= 9 {
		// Normal Phase
		workDuration = 25
		breakDuration = 5
		targetCycles = level
		requiresProject = false
	} else if level >= 20 && level <= 29 {
		// Focus Phase
		workDuration = 50
		breakDuration = 15
		targetCycles = level - 19
		requiresProject = false
	} else if level >= 30 && level <= 39 {
		// Deep Work Phase
		workDuration = 180
		breakDuration = 20
		targetCycles = level - 29
		requiresProject = false
	} else if level >= 40 && level <= 49 {
		// Hard Work Phase
		workDuration = 180
		breakDuration = 20
		targetCycles = level - 39
		requiresProject = true
	} else {
		// Default / fallback
		workDuration = 25
		breakDuration = 5
		targetCycles = 1
		requiresProject = false
	}

	return
}
