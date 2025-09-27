package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"DelayedNotifier/internal/config"
	"DelayedNotifier/internal/handlers"
)

func main() {
	cfg := config.Load()
	r := gin.Default()

	notificationHandler := handlers.NewNotificationHandler()

	api := r.Group("/api")
	{
		api.POST("/notify", notificationHandler.CreateNotification)
		api.GET("/notify/:id", notificationHandler.GetNotification)
		api.DELETE("/notify/:id", notificationHandler.CancelNotification)
	}

	log.Printf("Server starting on port %v", cfg.HTTPPort)
	err := r.Run(":" + cfg.HTTPPort)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
