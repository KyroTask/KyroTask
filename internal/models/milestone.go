package models

import "time"

type Milestone struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	GoalID      uint       `gorm:"not null;index" json:"goal_id"`
	Goal        Goal       `gorm:"foreignKey:GoalID" json:"goal,omitempty"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Status      string     `gorm:"default:'pending'" json:"status"` // pending, completed
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Tasks       []Task     `gorm:"foreignKey:MilestoneID" json:"tasks,omitempty"`
}
