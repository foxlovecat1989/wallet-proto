package events

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type LoginEvent struct {
	EventMetadata EventMetadata `json:"eventMetadata"`
	UserID        string        `json:"userId"`
	Email         string        `json:"email"`
	Username      string        `json:"username"`
	LoginAt       time.Time     `json:"loginAt"`
}

func (e *LoginEvent) ToTask() (*asynq.Task, error) {
	payload, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(string(LoginEventType), payload), nil
}
