package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TelegramID   *int64    `gorm:"uniqueIndex" json:"telegram_id,omitempty"`
	FirebaseUID  *string   `gorm:"uniqueIndex" json:"firebase_uid,omitempty"`
	Email        *string   `gorm:"uniqueIndex" json:"email,omitempty"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhotoURL     string    `json:"photo_url"`
	IsBot        bool      `json:"is_bot"`
	LanguageCode string    `json:"language_code"`
	IsPremium    bool      `json:"is_premium"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TelegramAccount struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ChatID    int64     `gorm:"uniqueIndex" json:"chat_id"`
	IsLinked  bool      `json:"is_linked"`
	CreatedAt time.Time `json:"created_at"`
}
