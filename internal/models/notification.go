package models

import "time"

type Notification struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"not null;index" json:"user_id"`
	User              User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type              string    `gorm:"not null" json:"type"`    // task_due, habit_reminder, etc.
	RelatedID         uint      `gorm:"index" json:"related_id"` // ID of the task, habit, or goal
	TelegramMessageID int       `json:"telegram_message_id"`
	Status            string    `gorm:"default:'sent'" json:"status"` // sent, failed, deleted
	SentAt            time.Time `json:"sent_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
