package repository

import (
	"context"
	"encoding/json"
	"user-svc/internal/app/domains/models"
	"user-svc/internal/db"

	"github.com/samber/lo"
)

type NotificationEventLogStatus string

const (
	NotificationEventLogStatusPending NotificationEventLogStatus = "pending"
	NotificationEventLogStatusSuccess NotificationEventLogStatus = "success"
	NotificationEventLogStatusFailed  NotificationEventLogStatus = "failed"
)

type NotificationEventLog struct {
	ID        string                     `db:"id"`
	EventName string                     `db:"event_name"`
	Payload   json.RawMessage            `db:"payload"`
	Status    NotificationEventLogStatus `db:"status"`
	CreatedAt int64                      `db:"created_at"`
	UpdatedAt int64                      `db:"updated_at"`
}

func (e *NotificationEventLog) ToModel() *models.NotificationEventLog {
	return &models.NotificationEventLog{
		ID:        e.ID,
		EventName: e.EventName,
		Payload:   e.Payload,
		Status:    models.NotificationEventLogStatus(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type NotificationEventLogRepository struct {
	store db.Store
}

func NewNotificationEventLogRepository(store db.Store) *NotificationEventLogRepository {
	return &NotificationEventLogRepository{store: store}
}

func (r *NotificationEventLogRepository) Create(ctx context.Context, event *NotificationEventLog) error {
	_, err := r.store.ExecContext(
		ctx,
		`INSERT INTO notification_event_logs (id, event_name, payload, status) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		event.ID, event.EventName, event.Payload, event.Status,
	)

	return err
}

func (r *NotificationEventLogRepository) FindPendingEvents(
	ctx context.Context,
	eventName string,
	batchSize int,
) ([]*models.NotificationEventLog, error) {
	events := make([]*NotificationEventLog, 0)
	err := r.store.SelectContext(
		ctx,
		&events,
		`SELECT id, event_name, payload, status, created_at, updated_at 
		FROM notification_event_logs 
		WHERE event_name = $1 AND status = $2 
		ORDER BY created_at ASC 
		LIMIT $3`,
		eventName, NotificationEventLogStatusPending, batchSize,
	)

	return lo.Map(events, func(event *NotificationEventLog, _ int) *models.NotificationEventLog {
		return event.ToModel()
	}), err
}

func (r *NotificationEventLogRepository) UpdateStatusSuccess(ctx context.Context, id string) error {
	_, err := r.store.ExecContext(
		ctx,
		`UPDATE notification_event_logs SET status = $1 WHERE id = $2`,
		NotificationEventLogStatusSuccess, id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationEventLogRepository) UpdateStatusFailed(ctx context.Context, id string) error {
	_, err := r.store.ExecContext(
		ctx,
		`UPDATE notification_event_logs SET status = $1 WHERE id = $2`,
		NotificationEventLogStatusFailed, id,
	)

	return err
}
