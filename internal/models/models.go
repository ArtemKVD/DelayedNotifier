package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	ChatID      string    `json:"chat_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	RetryCount  int       `json:"RetryCount"`
}

func NewNotification(userID, title, message, chatID string, scheduledAt time.Time) *Notification {
	return &Notification{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Message:     message,
		ChatID:      chatID,
		ScheduledAt: scheduledAt,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
}
