package domain

import (
	"regexp"
	"wallet-user-svc/internal/app/errs"
)

type PhoneNumber string

func NewPhoneNumber(phone string) (PhoneNumber, error) {
	p := PhoneNumber(phone)
	if err := p.Validate(); err != nil {
		return "", err
	}

	return p, nil
}	

func (p PhoneNumber) Validate() error {
	if len(p) < 10 || len(p) > 15 {
		return errs.ErrInvalidPhoneNumber
	}

	if err := p.validateFormat(p.String()); err != nil {
		return err
	}

	return nil
}

func (p PhoneNumber) validateFormat(phone string) error {
	if !regexp.MustCompile(`^\+[1-9]\d{1,14}$`).MatchString(phone) {
		return errs.ErrInvalidPhoneNumber
	}

	return nil
}

func (p PhoneNumber) String() string {
	return string(p)
}

func (p PhoneNumber) IsSet() bool {
	return p != ""
}

func NewPhoneNumberPtr(phone *string) (*PhoneNumber, error) {
	if phone == nil {
		return nil, errs.ErrInvalidPhoneNumber
	}

	phoneNumber, err := NewPhoneNumber(*phone)
	if err != nil {
		return nil, err
	}

	return &phoneNumber, nil
}

func (p PhoneNumber) ToPtrString() *string {
	s := p.String()
	return &s
}

type CountryCode string

func NewCountryCode(code string) (CountryCode, error) {
	c := CountryCode(code)
	if err := c.Validate(); err != nil {
		return "", err
	}
	return c, nil
}

func NewCountryCodePtr(code *string) (*CountryCode, error) {
	if code == nil {
		return nil, errs.ErrInvalidCountryCode
	}

	countryCode, err := NewCountryCode(*code)
	if err != nil {
		return nil, err
	}

	return &countryCode, nil
}

func (c CountryCode) Validate() error {
	if len(c) != 2 {
		return errs.ErrInvalidCountryCode
	}

	return nil
}

func (c CountryCode) IsSet() bool {
	return c != ""
}

func (c CountryCode) ToPtrString() *string {
	s := string(c)
	return &s
}