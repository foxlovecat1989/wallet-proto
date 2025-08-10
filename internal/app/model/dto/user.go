package dto

import (
	"wallet-user-svc/internal/app/errs"
	"wallet-user-svc/internal/app/model/domain"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    *string `json:"email"`
	CountryCode *string `json:"countryCode"`
	Phone       *string `json:"phone"`
}

func (r *RegisterReq) Validate() error {
	if r.Email == nil && (r.CountryCode == nil || r.Phone == nil) {
		return errs.ErrEmailOrPhoneRequired
	}

	return nil
}

type RegisterResp struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}	

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}