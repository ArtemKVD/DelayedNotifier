package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"DelayedNotifier/internal/models"
)

type NotificationHandler struct {
	notifications map[string]*models.Notification
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notifications: make(map[string]*models.Notification),
	}
}

type NotificationRequest struct {
	UserID      string    `json:"user_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Message     string    `json:"message" binding:"required"`
	ChatID      string    `json:"chat_id" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req NotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	notification := models.NewNotification(
		req.UserID,
		req.Title,
		req.Message,
		req.ChatID,
		req.ScheduledAt,
	)

	h.notifications[notification.ID] = notification

	c.JSON(http.StatusCreated, notification)
}

func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, exists := h.notifications[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notification not found",
		})
		return
	}

	c.JSON(http.StatusOK, notification)
}

func (h *NotificationHandler) CancelNotification(c *gin.Context) {
	id := c.Param("id")

	notification, exists := h.notifications[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notification not found",
		})
		return
	}
	notification.Status = "cancelled"

	c.Status(http.StatusNoContent)
}
