package domain

import (
	"time"

	"wallet-user-svc/internal/app/errs"

	"github.com/google/uuid"
)

// User represents a user in the authentication system
type User struct {
	ID           uuid.UUID    `json:"id" `
	Email        *Email       `json:"email" `
	Username     Username     `json:"username" `
	CountryCode  *CountryCode `json:"country_code,omitempty" `
	Phone        *PhoneNumber `json:"phone,omitempty" `
	PasswordHash PasswordHash `json:"-" `
	CreatedAt    int64        `json:"created_at" `
	UpdatedAt    int64        `json:"updated_at" `
}

// NewUser creates a new user with generated ID and timestamps
func NewUser(email, passwordHash, username string, countryCode, phone *string) (*User, error) {
	if err := validateUserInput(email, countryCode, phone); err != nil {
		return nil, err
	}

	emailObj, countryCodeObj, phoneObj, err := createContactInfo(email, countryCode, phone)
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
		CountryCode:  countryCodeObj,
		Phone:        phoneObj,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func validateUserInput(email string, countryCode, phone *string) error {
	hasEmail := email != ""
	hasCountryCode := countryCode != nil && *countryCode != ""
	hasPhone := phone != nil && *phone != ""

	if !hasEmail && !(hasCountryCode && hasPhone) {
		return errs.ErrEmailOrPhoneRequired
	}
	return nil
}

func createContactInfo(email string, countryCode, phone *string) (*Email, *CountryCode, *PhoneNumber, error) {
	var emailObj *Email
	var countryCodeObj *CountryCode
	var phoneObj *PhoneNumber

	if email != "" {
		obj, err := NewEmail(email)
		if err != nil {
			return nil, nil, nil, err
		}
		emailObj = &obj
	}

	if countryCode != nil && *countryCode != "" {
		obj, err := NewCountryCode(*countryCode)
		if err != nil {
			return nil, nil, nil, err
		}
		countryCodeObj = &obj
	}

	if phone != nil && *phone != "" {
		obj, err := NewPhoneNumber(*phone)
		if err != nil {
			return nil, nil, nil, err
		}
		phoneObj = &obj
	}

	return emailObj, countryCodeObj, phoneObj, nil
}

// NewUserWithPassword creates a new user with password validation
func NewUserWithPassword(
	email *string,
	password, username string,
	countryCode, phone *string,
) (*User, error) {
	// Check if either email OR both country code and phone are provided
	hasEmail := email != nil && *email != ""
	hasCountryCode := countryCode != nil && *countryCode != ""
	hasPhone := phone != nil && *phone != ""

	if !hasEmail && !hasCountryCode && !hasPhone {
		return nil, errs.ErrEmailOrPhoneRequired
	}

	usernameObj, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	pwd, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	// Hash the password
	passwordHash, err := NewPasswordHashFromPlain(string(pwd))
	if err != nil {
		return nil, err
	}

	emailObj, err := NewEmailPtr(email)
	if err != nil {
		return nil, err
	}

	countryCodeObj, err := NewCountryCodePtr(countryCode)
	if err != nil {
		return nil, err
	}

	phoneObj, err := NewPhoneNumberPtr(phone)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()

	return &User{
		ID:           uuid.New(),
		Email:        emailObj,
		PasswordHash: passwordHash,
		Username:     usernameObj,
		CountryCode:  countryCodeObj,
		Phone:        phoneObj,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsValid checks if the user data is valid
func (u *User) IsValid() error {
	// Check if either email OR both country code and phone are provided
	hasEmail := u.Email.IsSet()
	hasCountryCode := u.CountryCode != nil && *u.CountryCode != ""
	hasPhone := u.Phone != nil && *u.Phone != ""

	if !hasEmail && !(hasCountryCode && hasPhone) {
		return errs.ErrEmailOrPhoneRequired
	}

	// If email is provided, validate it
	if hasEmail {
		if err := u.Email.Validate(); err != nil {
			return err
		}
	}

	if err := u.Username.Validate(); err != nil {
		return err
	}
	if err := u.PasswordHash.Validate(); err != nil {
		return err
	}

	return nil
}
