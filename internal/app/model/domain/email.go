package domain

import "wallet-user-svc/internal/app/errs"

// Email represents a validated email address
type Email string

// NewEmail creates a new Email and validates it
func NewEmail(email string) (Email, error) {
	e := Email(email)
	if err := e.Validate(); err != nil {
		return "", err
	}
	return e, nil
}

func NewEmailPtr(email *string) (*Email, error) {
	if email == nil {
		return nil, errs.ErrInvalidEmail
	}

	emailObj, err := NewEmail(*email)
	if err != nil {
		return nil, err
	}

	return &emailObj, nil
}

// Validate checks if the email format is valid
func (e Email) Validate() error {
	email := string(e)

	if err := e.validateLength(email); err != nil {
		return err
	}
	if err := e.validateAtSymbol(email); err != nil {
		return err
	}
	if err := e.validateDomain(email); err != nil {
		return err
	}

	return nil
}

func (e Email) validateLength(email string) error {
	if len(email) < 5 || len(email) > 254 {
		return errs.ErrInvalidEmail
	}
	return nil
}

func (e Email) validateAtSymbol(email string) error {
	atIndex := e.findAtSymbol(email)
	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return errs.ErrInvalidEmail
	}
	return nil
}

func (e Email) findAtSymbol(email string) int {
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return -1 // Multiple @ symbols
			}
			atIndex = i
		}
	}
	return atIndex
}

func (e Email) validateDomain(email string) error {
	atIndex := e.findAtSymbol(email)
	if atIndex == -1 {
		return errs.ErrInvalidEmail
	}

	domain := email[atIndex+1:]
	if len(domain) < 2 || len(domain) > 253 {
		return errs.ErrInvalidEmail
	}

	if !e.hasDotInDomain(domain) {
		return errs.ErrInvalidEmail
	}

	return nil
}

func (e Email) hasDotInDomain(domain string) bool {
	for _, char := range domain {
		if char == '.' {
			return true
		}
	}
	return false
}

// String returns the email as a string
func (e Email) String() string {
	return string(e)
}

func (e Email) IsSet() bool {
	return e != ""
}

func (e Email) ToPtrString() *string {
	s := e.String()
	return &s
}
