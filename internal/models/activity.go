package models

import "time"

type Activity struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ResourceType string    `gorm:"not null" json:"resource_type"` // task, project, goal, habit
	ResourceID   uint      `gorm:"not null" json:"resource_id"`
	Action       string    `gorm:"not null" json:"action"` // created, updated, deleted, completed
	Description  string    `json:"description"`
	Metadata     string    `gorm:"type:text" json:"metadata"` // JSON data
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}
