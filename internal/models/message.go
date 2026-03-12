package models

import (
	"time"
)

type Message struct {
	ID             uint      `gorm:"primaryKey"`
	ConversationID uint      `gorm:"not null;index"`
	Role           string    `gorm:"type:varchar(20);not null"` // 'user' or 'assistant'
	Content        string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Conversation Conversation `gorm:"foreignKey:ConversationID"`
}

func (Message) TableName() string {
	return "messages"
}
