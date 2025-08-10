package domain

import "wallet-user-svc/internal/app/errs"

// Username represents a validated username
type Username string

// NewUsername creates a new Username and validates it
func NewUsername(username string) (Username, error) {
	u := Username(username)
	if err := u.Validate(); err != nil {
		return "", err
	}
	return u, nil
}

// Validate checks if the username meets requirements
func (u Username) Validate() error {
	username := string(u)

	if err := u.validateLength(username); err != nil {
		return err
	}
	if err := u.validateCharacters(username); err != nil {
		return err
	}
	if err := u.validateBoundaries(username); err != nil {
		return err
	}
	if err := u.validateConsecutive(username); err != nil {
		return err
	}

	return nil
}

func (u Username) validateLength(username string) error {
	if username == "" || len(username) < 3 || len(username) > 30 {
		return errs.ErrInvalidUsername
	}
	return nil
}

func (u Username) validateCharacters(username string) error {
	for _, char := range username {
		if !u.isValidChar(char) {
			return errs.ErrInvalidUsername
		}
	}
	return nil
}

func (u Username) isValidChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_' || char == '-'
}

func (u Username) validateBoundaries(username string) error {
	if len(username) > 0 && (username[0] == '_' || username[0] == '-' ||
		username[len(username)-1] == '_' || username[len(username)-1] == '-') {
		return errs.ErrInvalidUsername
	}
	return nil
}

func (u Username) validateConsecutive(username string) error {
	for i := 0; i < len(username)-1; i++ {
		if (username[i] == '_' && username[i+1] == '_') ||
			(username[i] == '-' && username[i+1] == '-') {
			return errs.ErrInvalidUsername
		}
	}
	return nil
}

// String returns the username as a string
func (u Username) String() string {
	return string(u)
}
