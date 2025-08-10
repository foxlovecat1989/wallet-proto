package dto

import "time"

type SendLoginNotificationParams struct {
	UserID   string    `json:"userID"`
	Email    *string   `json:"email,omitempty"`
	Username string    `json:"username"`
	LoginAt  time.Time `json:"loginAt"`
}
