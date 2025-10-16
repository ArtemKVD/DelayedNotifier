package repository

import (
	"context"
	"database/sql"

	"DelayedNotifier/internal/models"

	"github.com/wb-go/wbf/zlog"

	"github.com/wb-go/wbf/dbpg"
)

type PostgresRepository struct {
	db *dbpg.DB
}

type Repository interface {
	Create(notification *models.Notification) error
	GetByID(id string) (*models.Notification, error)
	UpdateStatus(id, status string) error
	UpdateWithRetry(id string, retryCount int, lastError string) error
	ListByStatus(status string) ([]*models.Notification, error)
}

func NewPostgresRepository(db *dbpg.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(notification *models.Notification) error {
	query := `
        INSERT INTO notifications (id, user_id, title, message, chat_id, scheduled_at, status, created_at, retry_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := r.db.ExecContext(context.Background(), query,
		notification.ID,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.ChatID,
		notification.ScheduledAt,
		notification.Status,
		notification.CreatedAt,
		notification.RetryCount,
	)

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Create notification error")
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByID(id string) (*models.Notification, error) {
	query := `
        SELECT id, user_id, title, message, chat_id, scheduled_at, status, created_at, retry_count, last_error
        FROM notifications WHERE id = $1
    `

	var notification models.Notification
	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Message,
		&notification.ChatID,
		&notification.ScheduledAt,
		&notification.Status,
		&notification.CreatedAt,
		&notification.RetryCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		zlog.Logger.Error().Err(err).Str("id", id).Msg("Get notification error")
		return nil, err
	}

	return &notification, nil
}

func (r *PostgresRepository) UpdateStatus(id, status string) error {
	query := `UPDATE notifications SET status = $1 WHERE id = $2`

	_, err := r.db.ExecContext(context.Background(), query, status, id)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("id", id).Str("status", status).Msg("Update notification error")
		return err
	}

	return nil
}

func (r *PostgresRepository) UpdateWithRetry(id string, retryCount int, lastError string) error {
	query := `UPDATE notifications SET retry_count = $1, last_error = $2, status = $3 WHERE id = $4`

	status := "pending"
	if retryCount >= 5 {
		status = "failed"
	}

	_, err := r.db.ExecContext(context.Background(), query, retryCount, lastError, status, id)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("id", id).Msg("Failed to update notification retry info")
		return err
	}

	return nil
}

func (r *PostgresRepository) ListByStatus(status string) ([]*models.Notification, error) {
	query := `
        SELECT id, user_id, title, message, chat_id, scheduled_at, status, created_at, retry_count, last_error
        FROM notifications WHERE status = $1 ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(context.Background(), query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Message,
			&notification.ChatID,
			&notification.ScheduledAt,
			&notification.Status,
			&notification.CreatedAt,
			&notification.RetryCount,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}
