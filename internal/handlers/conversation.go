package handlers

import (
	"net/http"
	"strconv"

	"github.com/Ankit-Choudhary16/ChatVault/internal/models"
	"github.com/Ankit-Choudhary16/ChatVault/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateConversationRequest struct {
	Title string `json:"title"`
}

type CreateConversationResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}
const DateTimeFormat = "02 Jan 2006, 03:04 PM"

func CreateConversation(db *gorm.DB) gin.HandlerFunc {
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

		var req CreateConversationRequest
		_ = c.ShouldBindJSON(&req) // title is optional, ignore bind errors

		title := req.Title
		if title == "" {
			title = "New Chat"
		}

		conversation := models.Conversation{
			UserID: uid,
			Title:  title,
		}
		err := database.CreateConversation(db, &conversation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create conversation"})
			return
		}
		
		c.JSON(http.StatusCreated, CreateConversationResponse{
			ID:        conversation.ID,
			UserID:    conversation.UserID,
			Title:     conversation.Title,
			CreatedAt: conversation.CreatedAt.Format(DateTimeFormat),
		})
	}
}

type ListConversationsResponse struct {
	Conversations []ConversationItem `json:"conversations"`
	Total         int64              `json:"total"`
	Page          int                `json:"page"`
}

type ConversationItem struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SingleConversation struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	MessageCount int64  `json:"message_count"`
}

func ListConversations(db *gorm.DB) gin.HandlerFunc {
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

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}

		offset := (page - 1) * limit

		var total int64
		
		total ,err:= database.CountConversations(db, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count conversations"})
			return
		}
		
		if total == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "no conversations found"})
			return
		}


		var conversations []models.Conversation
		
		conversations, err = database.GetAllConversations(db, uid, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list conversations"})
			return
		}

		items := make([]ConversationItem, len(conversations))
		for i, conv := range conversations {
			items[i] = ConversationItem{
				ID:        conv.ID,
				UserID:    conv.UserID,
				Title:     conv.Title,
				CreatedAt: conv.CreatedAt.Format(DateTimeFormat),
				UpdatedAt: conv.UpdatedAt.Format(DateTimeFormat),
			}
		}

		c.JSON(http.StatusOK, ListConversationsResponse{
			Conversations: items,
			Total:         total,
			Page:          page,
		})
	}
}

func GetConversation(db *gorm.DB) gin.HandlerFunc {
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

		var conversation *models.Conversation
		conversation, err = database.GetConversationByID(db, uint(conversationID), uid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch conversation"})
			return
		}


		var messageCount int64
		messageCount, err = database.CountMessages(db, uint(conversationID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count messages"})
			return
		}

		c.JSON(http.StatusOK, SingleConversation{
			ID:           conversation.ID,
			UserID:       conversation.UserID,
			Title:        conversation.Title,
			CreatedAt:    conversation.CreatedAt.Format(DateTimeFormat),
			UpdatedAt:    conversation.UpdatedAt.Format(DateTimeFormat),
			MessageCount: messageCount,
		})
	}
}

func DeleteConversation(db *gorm.DB) gin.HandlerFunc {
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

		error, rowsAffected := database.DeleteConversation(db, uint(conversationID), uid)
		if error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete conversation"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "conversation deleted"})
	}
}
