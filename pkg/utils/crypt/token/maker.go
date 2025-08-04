package token

type TokenMaker interface {
	CreateTokenPair(userID string, username string, duration int64) (string, string, error)
	CreateAccessToken(userID string, username string, duration int64) (string, error)
	CreateRefreshToken(userID string, username string, duration int64) (string, error)
	VerifyAccessToken(token string) (*Payload, error)
	VerifyRefreshToken(token string) (*Payload, error)
}
