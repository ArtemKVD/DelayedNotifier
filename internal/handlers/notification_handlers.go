package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"DelayedNotifier/internal/models"
	"DelayedNotifier/internal/repository"

	"github.com/wb-go/wbf/zlog"
)

type NotificationHandler struct {
	repo repository.Repository
}

func NewNotificationHandler(repo repository.Repository) *NotificationHandler {
	return &NotificationHandler{
		repo: repo,
	}
}

type CreateNotificationRequest struct {
	UserID      string    `json:"user_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Message     string    `json:"message" binding:"required"`
	ChatID      string    `json:"chat_id" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("Invalid request body")
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

	err := h.repo.Create(notification)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Create notification error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Create notification error",
		})
		return
	}

	zlog.Logger.Info().Str("id", notification.ID).Time("scheduled_at", notification.ScheduledAt).Msg("Notification created")
	c.JSON(http.StatusCreated, notification)
}

func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.repo.GetByID(id)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("id", id).Msg("Get notification error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if notification == nil {
		zlog.Logger.Warn().Str("id", id).Msg("Notification not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notification not found",
		})
		return
	}

	c.JSON(http.StatusOK, notification)
}

func (h *NotificationHandler) CancelNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.repo.GetByID(id)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("id", id).Msg("Failed to get notification")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if notification == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notification not found",
		})
		return
	}

	err = h.repo.UpdateStatus(id, "cancelled")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("id", id).Msg("Cancel notification error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cancel notification error",
		})
		return
	}

	zlog.Logger.Info().Str("id", id).Msg("Notification cancelled")
	c.Status(http.StatusNoContent)
}
