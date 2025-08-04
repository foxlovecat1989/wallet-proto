package models

import (
	"time"

	"user-svc/internal/app/domains/errs"

	"github.com/google/uuid"
)

// User represents a user in the authentication system
type User struct {
	ID           uuid.UUID    `json:"id" `
	Email        Email        `json:"email" `
	Username     Username     `json:"username" `
	PasswordHash PasswordHash `json:"-" `
	CreatedAt    int64        `json:"created_at" `
	UpdatedAt    int64        `json:"updated_at" `
}

// NewUser creates a new user with generated ID and timestamps
func NewUser(email, passwordHash, username string) (*User, error) {
	if email == "" {
		return nil, errs.ErrEmailIsRequired
	}

	emailObj, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	usernameObj, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	passwordHashObj, err := NewPasswordHash(passwordHash)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	id := uuid.New()

	return &User{
		ID:           id,
		Email:        emailObj,
		PasswordHash: passwordHashObj,
		Username:     usernameObj,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// NewUserWithPassword creates a new user with password validation
func NewUserWithPassword(email, password, username string) (*User, error) {
	if email == "" {
		return nil, errs.ErrEmailIsRequired
	}

	emailObj, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Validate username
	usernameObj, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	// Validate password
	pwd, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	passwordHash, err := NewPasswordHash(string(pwd))
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()

	return &User{
		ID:           uuid.New(),
		Email:        emailObj,
		PasswordHash: passwordHash,
		Username:     usernameObj,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsValid checks if the user data is valid
func (u *User) IsValid() error {
	if u.Email == "" {
		return errs.ErrEmailIsRequired
	}
	if err := u.Email.Validate(); err != nil {
		return err
	}
	if err := u.Username.Validate(); err != nil {
		return err
	}
	if err := u.PasswordHash.Validate(); err != nil {
		return err
	}

	return nil
}
