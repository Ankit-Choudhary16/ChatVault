package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ankit-Choudhary16/ChatVault/internal/database"
	"github.com/Ankit-Choudhary16/ChatVault/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddMessageRequest struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type AddMessageResponse struct {
	ID        uint   `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func AddMessage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
			return
		}

		conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
			return
		}

		var req AddMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Role != "user" && req.Role != "assistant" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role must be 'user' or 'assistant'"})
			return
		}

		var _ *models.Conversation
		_, err = database.GetConversationByID(db, uint(conversationID), uid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch conversation"})
			return
		}
		

		message := models.Message{
			ConversationID: uint(conversationID),
			Role:           req.Role,
			Content:        req.Content,
		}
		error := database.AddMessage(db, &message)
		if error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create message"})
			return
		}
		c.JSON(http.StatusCreated, AddMessageResponse{
			ID:        message.ID,
			Role:      message.Role,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Format(DateTimeFormat),
		})
	}
}

type MessageItem struct {
	ID        uint   `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type GetMessagesResponse struct {
	Messages []MessageItem `json:"messages"`
}

func GetMessages(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
			return
		}

		conversationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
			return
		}

		var _ *models.Conversation
		_, err = database.GetConversationByID(db, uint(conversationID), uid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch conversation"})
			return
		}

		var messages []models.Message
		messages, err = database.GetMessages(db, uint(conversationID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
			return
		}

		items := make([]MessageItem, len(messages))
		for i, msg := range messages {
			items[i] = MessageItem{
				ID:        msg.ID,
				Role:      msg.Role,
				Content:   msg.Content,
				CreatedAt: msg.CreatedAt.Format(DateTimeFormat),
			}
		}

		c.JSON(http.StatusOK, GetMessagesResponse{Messages: items})
	}
}
