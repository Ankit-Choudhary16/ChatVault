package database

import (
	// "fmt"
	// "log"

	"github.com/Ankit-Choudhary16/ChatVault/internal/models"
	"gorm.io/gorm"
)

func EmailExist(db *gorm.DB, email string) bool {
	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return true
	}
	return false

}

func CreateUser(db *gorm.DB, user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateConversation(db *gorm.DB, conversation *models.Conversation) error {
	if err := db.Create(conversation).Error; err != nil {
		return err
	}
	return nil
}

func CountConversations(db *gorm.DB, userID uint) (int64, error) {
	var count int64
	if err := db.Model(&models.Conversation{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetAllConversations(db *gorm.DB, userID uint, limit int, offset int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	if err := db.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Offset(offset).Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func GetConversationByID(db *gorm.DB, conversationID uint, userID uint) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := db.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func CountMessages(db *gorm.DB, conversationID uint) (int64, error) {
	var count int64
	if err := db.Model(&models.Message{}).Where("conversation_id = ?", conversationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func DeleteConversation(db *gorm.DB, conversationID uint, userID uint) (error, int64) {
	result := db.Where("id = ? AND user_id = ?", conversationID, userID).Delete(&models.Conversation{})
	return result.Error, result.RowsAffected
}

func AddMessage(db *gorm.DB, message *models.Message) error {
	if err := db.Create(message).Error; err != nil {
		return err
	}
	return nil
}

func GetMessages(db *gorm.DB, conversationID uint) ([]models.Message, error) {
	var messages []models.Message
	if err := db.Where("conversation_id = ?", conversationID).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
