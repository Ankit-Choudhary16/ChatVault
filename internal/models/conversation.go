package models

import (
	"time"
)

type Conversation struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Title     string    `gorm:"type:varchar(255);default:'New Chat'"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	User     User      `gorm:"foreignKey:UserID"`
	Messages []Message `gorm:"foreignKey:ConversationID"`
}

func (Conversation) TableName() string {
	return "conversations"
}
