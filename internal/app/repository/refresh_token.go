package repository

import (
	"context"
	"database/sql"
	"fmt"

	"user-svc/internal/app/domains/errs"
	"user-svc/internal/app/domains/models"
	"user-svc/internal/db"
	"user-svc/pkg/utils/tx"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt int64     `db:"expires_at"`
	IsRevoked bool      `db:"is_revoked"`
	CreatedAt int64     `db:"created_at"`
	UpdatedAt int64     `db:"updated_at"`
}

func (rt *RefreshToken) ToDomain() *models.RefreshToken {
	return &models.RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
		IsRevoked: rt.IsRevoked,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
	}
}

type RefreshTokenRepository struct {
	db db.Store
}

func NewRefreshTokenRepository(db db.Store) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, is_revoked, created_at, updated_at)
		VALUES (:id, :user_id, :token, :expires_at, :is_revoked, :created_at, :updated_at)
	`

	repoRefreshToken := &RefreshToken{
		ID:        refreshToken.ID,
		UserID:    refreshToken.UserID,
		Token:     refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
		IsRevoked: refreshToken.IsRevoked,
		CreatedAt: refreshToken.CreatedAt,
		UpdatedAt: refreshToken.UpdatedAt,
	}

	// Check if we're in a transaction
	if tx, ok := ctx.Value(tx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		_, err := tx.NamedExecContext(ctx, query, repoRefreshToken)
		if err != nil {
			return fmt.Errorf("failed to create refresh token: %w", err)
		}
		return nil
	}

	// Use main database connection
	_, err := r.db.NamedExecContext(ctx, query, repoRefreshToken)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetByTokenHash retrieves a refresh token by token hash
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, is_revoked, created_at, updated_at
		FROM refresh_tokens 
		WHERE token = $1
	`

	var refreshToken RefreshToken

	// Check if we're in a transaction
	if tx, ok := ctx.Value(tx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		err := tx.QueryRowContext(ctx, query, tokenHash).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.IsRevoked, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errs.ErrTokenNotFound
			}
			return nil, fmt.Errorf("failed to get refresh token by token: %w", err)
		}
		return refreshToken.ToDomain(), nil
	}

	// Use main database connection
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.IsRevoked, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to get refresh token by token: %w", err)
	}

	return refreshToken.ToDomain(), nil
}
