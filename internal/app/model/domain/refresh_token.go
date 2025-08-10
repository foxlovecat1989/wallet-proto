package domain

import (
	"time"

	"wallet-user-svc/internal/app/errs"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token domain model
type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt int64     `json:"expiresAt"`
	IsRevoked bool      `json:"isRevoked"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

// NewRefreshToken creates a new RefreshToken
func NewRefreshToken(userID uuid.UUID, tokenHash string, expiresAt int64) (*RefreshToken, error) {
	if userID == uuid.Nil {
		return nil, errs.ErrInvalidToken
	}

	if tokenHash == "" {
		return nil, errs.ErrInvalidToken
	}

	if expiresAt <= time.Now().UnixMilli() {
		return nil, errs.ErrTokenExpired
	}

	return &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tokenHash,
		ExpiresAt: expiresAt,
		IsRevoked: false,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}, nil
}

// IsValid checks if the refresh token is valid
func (rt *RefreshToken) IsValid() error {
	if rt.ID == uuid.Nil {
		return errs.ErrInvalidToken
	}

	if rt.UserID == uuid.Nil {
		return errs.ErrInvalidToken
	}

	if rt.Token == "" {
		return errs.ErrInvalidToken
	}

	if rt.IsRevoked {
		return errs.ErrTokenRevoked
	}

	if rt.ExpiresAt <= time.Now().UnixMilli() {
		return errs.ErrTokenExpired
	}

	return nil
}
