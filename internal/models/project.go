package models

import (
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name        string    `gorm:"not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex" json:"slug"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	Progress    int       `gorm:"default:0" json:"progress"`      // 0-100
	Status      string    `gorm:"default:'active'" json:"status"` // active, completed, archived
	IsArchived  bool      `gorm:"default:false" json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Goals       []Goal    `gorm:"foreignKey:ProjectID" json:"goals,omitempty"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Slug == "" {
		baseSlug := slug.Make(p.Name)
		p.Slug = baseSlug

		// Ensure uniqueness
		var count int64
		tx.Model(&Project{}).Where("slug = ?", p.Slug).Count(&count)
		if count > 0 {
			p.Slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix())
		}
	}
	return
}
