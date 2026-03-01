package models

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Goal struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null;index" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectID   *uint       `gorm:"index" json:"project_id"`
	Project     *Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Title       string      `gorm:"not null" json:"title"`
	Slug        string      `gorm:"uniqueIndex" json:"slug"`
	Description string      `json:"description"`
	Motivation  string      `json:"motivation"` // The "Why"
	Notes       string      `json:"notes"`      // Log/Notes section
	TargetDate  *time.Time  `json:"target_date"`
	Progress    int         `gorm:"default:0" json:"progress"`      // 0-100
	Status      string      `gorm:"default:'active'" json:"status"` // active, completed, abandoned
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Tasks       []Task      `gorm:"foreignKey:GoalID" json:"tasks,omitempty"`
	Milestones  []Milestone `gorm:"foreignKey:GoalID" json:"milestones,omitempty"`
	Habits      []Habit     `gorm:"foreignKey:GoalID" json:"habits,omitempty"`
}

func (g *Goal) BeforeCreate(tx *gorm.DB) (err error) {
	if g.Slug == "" {
		baseSlug := slug.Make(g.Title)
		g.Slug = baseSlug

		// Ensure uniqueness
		var count int64
		tx.Model(&Goal{}).Where("slug = ?", g.Slug).Count(&count)
		if count > 0 {
			g.Slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix())
		}
	}
	return
}
