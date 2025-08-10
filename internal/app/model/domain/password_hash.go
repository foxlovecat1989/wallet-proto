package domain

import (
	"wallet-user-svc/internal/app/errs"
	"wallet-user-svc/pkg/utils/crypt/password"
)

// PasswordHash represents a hashed password
type PasswordHash string

// NewPasswordHash creates a new PasswordHash and validates it
func NewPasswordHash(hash string) (PasswordHash, error) {
	ph := PasswordHash(hash)
	if err := ph.Validate(); err != nil {
		return "", err
	}
	return ph, nil
}

// NewPasswordHashFromPlain creates a new PasswordHash from a plain text password
func NewPasswordHashFromPlain(plainPassword string) (PasswordHash, error) {
	hasher := password.DefaultHasher()
	hashedPassword, err := hasher.HashPassword(plainPassword)
	if err != nil {
		return "", err
	}
	return PasswordHash(hashedPassword), nil
}

// Validate checks if the password hash is valid (non-empty)
func (ph PasswordHash) Validate() error {
	if string(ph) == "" {
		return errs.ErrInvalidPassword
	}
	return nil
}

// String returns the password hash as a string
func (ph PasswordHash) String() string {
	return string(ph)
}

// VerifyPassword checks if the password hash matches the provided password
func (ph PasswordHash) VerifyPassword(plainPassword string) bool {
	hasher := password.DefaultHasher()
	return hasher.VerifyPassword(string(ph), plainPassword)
}
