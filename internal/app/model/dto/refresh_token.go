package dto

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}		