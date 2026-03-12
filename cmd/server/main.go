package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/Ankit-Choudhary16/ChatVault/internal/config"
    "github.com/Ankit-Choudhary16/ChatVault/internal/database"
    "github.com/Ankit-Choudhary16/ChatVault/internal/handlers"
    "github.com/Ankit-Choudhary16/ChatVault/internal/middleware"
)

func main() {
    // Load config
    cfg := config.Load()
    
    // Connect to database
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }else {
		log.Println("Successfully connected to the database")
	}
	
    
    // Setup router
    r := gin.Default()
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API routes
    api := r.Group("/api/v1")
    {
        // Auth (no middleware)
        api.POST("/auth/register", handlers.Register(db))
        api.POST("/auth/login", handlers.Login(db))
        
        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.JWTAuth())
        {
            protected.GET("/users/me", handlers.GetProfile(db))
            protected.POST("/conversations", handlers.CreateConversation(db))
            protected.GET("/conversations", handlers.ListConversations(db))
            protected.GET("/conversations/:id", handlers.GetConversation(db))
            protected.DELETE("/conversations/:id", handlers.DeleteConversation(db))
            protected.POST("/conversations/:id/messages", handlers.AddMessage(db))
            protected.GET("/conversations/:id/messages", handlers.GetMessages(db))
        }
    }
    
    // // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}

