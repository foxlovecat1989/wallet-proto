package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher provides password hashing and verification functionality
type Hasher struct {
	cost int
}

// NewHasher creates a new password hasher with the specified cost
func NewHasher(cost int) *Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &Hasher{cost: cost}
}

// HashPassword hashes a plain text password using bcrypt
func (h *Hasher) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against a hashed password
func (h *Hasher) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// DefaultHasher returns a hasher with default bcrypt cost
func DefaultHasher() *Hasher {
	return NewHasher(bcrypt.DefaultCost)
}
