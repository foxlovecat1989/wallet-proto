package service

import (
	"context"
	"time"

	"user-svc/internal/app/domains/dto"
	"user-svc/internal/app/domains/errs"
	"user-svc/internal/app/domains/models"
	"user-svc/pkg/utils/crypt/token"
	"user-svc/pkg/utils/tx"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *models.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*models.RefreshToken, error)
}

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(*tx.TxWrapper) error) error
}

// UserService handles business logic for user operations
type UserService struct {
	userRepo         UserRepository
	refreshTokenRepo RefreshTokenRepository
	txManager        TxManager
	tokenMaker       token.TokenMaker
	tokenDuration    time.Duration
	refreshDuration  time.Duration
}

// NewUserService creates a new UserService instance
func NewUserService(
	userRepo UserRepository,
	refreshTokenRepo RefreshTokenRepository,
	txManager TxManager,
	tokenMaker token.TokenMaker,
	tokenDuration time.Duration,
	refreshDuration time.Duration,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		txManager:        txManager,
		tokenMaker:       tokenMaker,
		tokenDuration:    tokenDuration,
		refreshDuration:  refreshDuration,
	}
}

// Register handles user registration
func (s *UserService) Register(ctx context.Context, req dto.RegisterReq) (*dto.RegisterResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := models.NewUserWithPassword(req.Email, req.Password, req.Username)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.tokenMaker.CreateTokenPair(
		user.ID.String(),
		user.Username.String(),
		int64(s.tokenDuration),
	)
	if err != nil {
		return nil, err
	}

	err = s.txManager.WithTransaction(ctx, func(tx *tx.TxWrapper) error {
		if err := s.userRepo.Create(ctx, user); err != nil {
			return err
		}

		refreshTokenModel, err := models.NewRefreshToken(
			user.ID,
			refreshToken,
			time.Now().Add(s.refreshDuration).UnixMilli(),
		)
		if err != nil {
			return err
		}

		if err := s.refreshTokenRepo.Create(ctx, refreshTokenModel); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResp{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login handles user login
func (s *UserService) Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !user.PasswordHash.VerifyPassword(req.Password) {
		return nil, errs.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.tokenMaker.CreateTokenPair(
		user.ID.String(),
		user.Username.String(),
		int64(s.tokenDuration),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel, err := models.NewRefreshToken(
		user.ID,
		refreshToken,
		time.Now().Add(s.refreshDuration).UnixMilli(),
	)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &dto.LoginResp{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req dto.RefreshTokenReq) (*dto.RefreshTokenResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		if err == errs.ErrTokenNotFound {
			return nil, errs.ErrTokenNotFound
		}

		return nil, err
	}

	if refreshToken.IsRevoked {
		return nil, errs.ErrTokenRevoked
	}

	if refreshToken.ExpiresAt < time.Now().UnixMilli() {
		return nil, errs.ErrTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenMaker.CreateAccessToken(
		user.ID.String(),
		user.Username.String(),
		int64(s.tokenDuration),
	)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResp{
		AccessToken: accessToken,
	}, nil
}
