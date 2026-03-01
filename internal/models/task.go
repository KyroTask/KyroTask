package models

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   *uint      `gorm:"index" json:"project_id"`
	Project     *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	GoalID      *uint      `gorm:"index" json:"goal_id"`
	Goal        *Goal      `gorm:"foreignKey:GoalID" json:"goal,omitempty"`
	MilestoneID *uint      `gorm:"index" json:"milestone_id"`
	Milestone   *Milestone `gorm:"foreignKey:MilestoneID" json:"milestone,omitempty"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID    *uint      `gorm:"index" json:"parent_id"`
	Title       string     `gorm:"not null" json:"title"`
	Slug        string     `gorm:"uniqueIndex" json:"slug"`
	Description string     `json:"description"`
	Status      string     `gorm:"default:'pending'" json:"status"`  // pending, in_progress, completed
	Priority    string     `gorm:"default:'medium'" json:"priority"` // low, medium, high, urgent
	DueDate     *time.Time `json:"due_date"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	IsDaily     bool       `gorm:"default:false" json:"is_daily"`
	CompletedAt *time.Time `json:"completed_at"`
	Position    int        `gorm:"default:0" json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Subtasks []Task    `gorm:"foreignKey:ParentID" json:"subtasks,omitempty"`
	Comments []Comment `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Tags     []Tag     `gorm:"many2many:task_tags;" json:"tags,omitempty"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Slug == "" {
		baseSlug := slug.Make(t.Title)
		t.Slug = baseSlug

		// Ensure uniqueness
		var count int64
		tx.Model(&Task{}).Where("slug = ?", t.Slug).Count(&count)
		if count > 0 {
			t.Slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix())
		}
	}
	return
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"not null;index" json:"task_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}
