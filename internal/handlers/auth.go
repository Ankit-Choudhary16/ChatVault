package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"golang.org/x/crypto/bcrypt"
	"github.com/Ankit-Choudhary16/ChatVault/internal/models"
	"github.com/Ankit-Choudhary16/ChatVault/internal/services"
	"github.com/Ankit-Choudhary16/ChatVault/internal/database"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var emailExists = database.EmailExist(db, req.Email)
		if emailExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already registered",
			})
			return
		}
		

		hashedPassword, err := services.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to hash password",
			})
			return
		}

		user := models.User{
			Email:        req.Email,
			Name:         req.Name,
			PasswordHash: string(hashedPassword),
		}
		err = database.CreateUser(db, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user",
			})
			return
		}

		c.JSON(http.StatusCreated, RegisterResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var user *models.User
		user, err := database.GetUserByEmail(db, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email",
			})
			return
		}

		if !services.ComparePassword(req.Password, user.PasswordHash) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid password",
			})
			return
		}

		token, err := services.GenerateJWT(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate token",
			})
			return
		}
		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}
