package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null"`
	Name         string    `gorm:"type:varchar(100);not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Conversations []Conversation `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
