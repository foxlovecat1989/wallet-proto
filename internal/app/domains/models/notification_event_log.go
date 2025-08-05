package models

import (
	"encoding/json"
)

type NotificationEventLogStatus string

const (
	NotificationEventLogStatusPending NotificationEventLogStatus = "pending"
	NotificationEventLogStatusSuccess NotificationEventLogStatus = "success"
	NotificationEventLogStatusFailed  NotificationEventLogStatus = "failed"
)

type NotificationEventLog struct {
	ID        string                     `db:"id" json:"id"`
	EventName string                     `db:"event_name" json:"eventName"`
	Payload   json.RawMessage            `db:"payload" json:"payload"`
	Status    NotificationEventLogStatus `db:"status" json:"status"`
	CreatedAt int64                      `db:"created_at" json:"createdAt"`
	UpdatedAt int64                      `db:"updated_at" json:"updatedAt"`
}
