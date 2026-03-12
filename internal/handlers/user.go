package handlers

import (
	"net/http"

	"github.com/Ankit-Choudhary16/ChatVault/internal/models"
	"github.com/Ankit-Choudhary16/ChatVault/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
			return
		}

		var user *models.User
		user, err := database.GetUserByID(db, userID.(uint))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		})
	}
}
