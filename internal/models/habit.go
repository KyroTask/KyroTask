package models

import "time"

type Habit struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"not null;index" json:"user_id"`
	User          User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GoalID        *uint      `gorm:"index" json:"goal_id"`
	Goal          *Goal      `gorm:"foreignKey:GoalID" json:"goal,omitempty"`
	Name          string     `gorm:"not null" json:"name"`
	Icon          string     `json:"icon"` // Custom emoji/icon for the habit
	Description   string     `json:"description"`
	Frequency     string     `gorm:"default:'daily'" json:"frequency"` // daily, weekly
	CurrentStreak int        `gorm:"default:0" json:"current_streak"`
	BestStreak    int        `gorm:"default:0" json:"best_streak"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	ScheduledDays string     `json:"scheduled_days"` // Comma-separated days (0-6, 0=Sunday)
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Logs          []HabitLog `gorm:"foreignKey:HabitID" json:"logs,omitempty"`
}

type HabitLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HabitID   uint      `gorm:"not null;index" json:"habit_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	LogDate   time.Time `gorm:"not null;index" json:"log_date"`
	Completed bool      `gorm:"default:true" json:"completed"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
