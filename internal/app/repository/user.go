package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wallet-user-svc/db"
	"wallet-user-svc/internal/app/errs"
	"wallet-user-svc/internal/app/model/domain"
	"wallet-user-svc/pkg/utils/cx"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// User domain model
type User struct {
	ID           string  `db:"id"`
	Email        *domain.Email  `db:"email"`
	Username     string  `db:"username"`
	CountryCode  *domain.CountryCode `db:"country_code"`
	Phone        *domain.PhoneNumber `db:"phone"`
	PasswordHash string  `db:"password_hash"`
	CreatedAt    int64   `db:"created_at"`
	UpdatedAt    int64   `db:"updated_at"`
}

func (u *User) ToDomain() *domain.User {
	username, err := domain.NewUsername(u.Username)
	if err != nil {
		// This should not happen in normal operation since we store validated usernames
		// But we need to handle it for backward compatibility
		username = domain.Username(u.Username)
	}

	id, err := uuid.Parse(u.ID)
	if err != nil {
		id = uuid.Nil
	}

	
	return &domain.User{
		ID:           id,
		Email:        u.Email,
		Username:     username,
		CountryCode:  u.CountryCode,
		Phone:        u.Phone,
		PasswordHash: domain.PasswordHash(u.PasswordHash),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

type UserRepository struct {
	db db.Store
}

func NewUserRepository(db db.Store) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, username, country_code, phone, password_hash, created_at, updated_at)
		VALUES (:id, :email, :username, :country_code, :phone, :password_hash, :created_at, :updated_at)
	`

	// Convert domain user to repository user
	repoUser := &User{
		ID:           user.ID.String(),
		Email:        user.Email,
		Username:     user.Username.String(),
		CountryCode:  user.CountryCode,
		Phone:        user.Phone,	
		PasswordHash: user.PasswordHash.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	// Check if we're in a transaction
	if tx, ok := ctx.Value(cx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		_, err := tx.NamedExecContext(ctx, query, repoUser)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	}

	// Use main database connection
	_, err := r.db.NamedExecContext(ctx, query, repoUser)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, username, country_code, phone, password_hash, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user User

	// Check if we're in a transaction
	if tx, ok := ctx.Value(cx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		err := tx.GetContext(ctx, &user, query, id.String())
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errs.ErrUserNotFound
			}
			return nil, fmt.Errorf("failed to get user by ID: %w", err)
		}
		return user.ToDomain(), nil
	}

	// Use main database connection
	err := r.db.GetContext(ctx, &user, query, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, country_code, phone, password_hash, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user User

	if tx, ok := ctx.Value(cx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		err := tx.GetContext(ctx, &user, query, email)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errs.ErrUserNotFound
			}
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}

		return user.ToDomain(), nil
	}

	// Use main database connection
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, countryCode, phone string) (*domain.User, error) {
	query := `
		SELECT id, email, username, country_code, phone, password_hash, created_at, updated_at
		FROM users 
		WHERE country_code = $1 AND phone = $2
	`

	var user User

	if tx, ok := ctx.Value(cx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		err := tx.GetContext(ctx, &user, query, countryCode, phone)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errs.ErrUserNotFound
			}
			return nil, fmt.Errorf("failed to get user by phone: %w", err)
		}

		return user.ToDomain(), nil
	}

	// Use main database connection
	err := r.db.GetContext(ctx, &user, query, countryCode, phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return user.ToDomain(), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	var result sql.Result
	var err error

	// Check if we're in a transaction
	if tx, ok := ctx.Value(cx.TransactionContextKey).(*sqlx.Tx); ok {
		// Use transaction
		result, err = tx.ExecContext(ctx, query, id.String())
	} else {
		// Use main database connection
		result, err = r.db.ExecContext(ctx, query, id.String())
	}

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
