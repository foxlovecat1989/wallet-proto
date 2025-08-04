package service

import (
	"context"
	"database/sql"
	"time"

	"user-svc/internal/app/config"
	"user-svc/internal/app/domains/dto"
	"user-svc/internal/app/domains/errs"
	"user-svc/internal/app/domains/models"
	"user-svc/pkg/utils/crypt/token"
	"user-svc/pkg/utils/log"
	"user-svc/pkg/utils/tx"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	WithTransactionOptions(ctx context.Context, fn func(*tx.TxWrapper) error, opts *sql.TxOptions) error
	WithTransactionIsolation(ctx context.Context, fn func(*tx.TxWrapper) error, isolation sql.IsolationLevel) error
	WithReadOnlyTransaction(ctx context.Context, fn func(*tx.TxWrapper) error) error
	WithSerializableTransaction(ctx context.Context, fn func(*tx.TxWrapper) error) error
	WithRepeatableReadTransaction(ctx context.Context, fn func(*tx.TxWrapper) error) error
	WithReadUncommittedTransaction(ctx context.Context, fn func(*tx.TxWrapper) error) error
}

// UserService handles business logic for user operations
type UserService struct {
	config           *config.Config
	userRepo         UserRepository
	refreshTokenRepo RefreshTokenRepository
	txManager        TxManager
	tokenMaker       token.TokenMaker
}

// NewUserService creates a new UserService instance
func NewUserService(
	config *config.Config,
	userRepo UserRepository,
	refreshTokenRepo RefreshTokenRepository,
	txManager TxManager,
	tokenMaker token.TokenMaker,
) *UserService {
	log.Info("Initializing UserService")

	service := &UserService{
		config:           config,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		txManager:        txManager,
		tokenMaker:       tokenMaker,
	}

	log.WithFields(logrus.Fields{
		"access_token_duration":  config.JWT.AccessTokenDuration.String(),
		"refresh_token_duration": config.JWT.RefreshTokenDuration.String(),
	}).Info("UserService initialized successfully")

	return service
}

// Register handles user registration
func (s *UserService) Register(ctx context.Context, req dto.RegisterReq) (*dto.RegisterResp, error) {
	logger := log.WithFields(logrus.Fields{
		"method":   "Register",
		"email":    req.Email,
		"username": req.Username,
	})

	logger.Info("Starting user registration")

	if err := req.Validate(); err != nil {
		logger.WithError(err).Error("Request validation failed")
		return nil, err
	}

	logger.Debug("Creating new user with password")
	user, err := models.NewUserWithPassword(req.Email, req.Password, req.Username)
	if err != nil {
		logger.WithError(err).Error("Failed to create user with password")
		return nil, err
	}

	logger.WithField("user_id", user.ID.String()).Debug("Creating token pair")
	accessToken, refreshToken, err := s.tokenMaker.CreateTokenPair(
		user.ID.String(),
		user.Username.String(),
		int64(s.config.JWT.AccessTokenDuration),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to create token pair")
		return nil, err
	}

	logger.Debug("Starting database transaction")
	err = s.txManager.WithTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, tx.TransactionContextKey, txWrapper.GetTx())

		logger.Debug("Creating user in database")
		if err := s.userRepo.Create(txCtx, user); err != nil {
			logger.WithError(err).Error("Failed to create user in database")
			return err
		}

		logger.Debug("Creating refresh token model")
		refreshTokenModel, err := models.NewRefreshToken(
			user.ID,
			refreshToken,
			time.Now().Add(s.config.JWT.RefreshTokenDuration).UnixMilli(),
		)
		if err != nil {
			logger.WithError(err).Error("Failed to create refresh token model")
			return err
		}

		logger.Debug("Storing refresh token in database")
		if err := s.refreshTokenRepo.Create(txCtx, refreshTokenModel); err != nil {
			logger.WithError(err).Error("Failed to store refresh token in database")
			return err
		}

		logger.Debug("Database transaction completed successfully")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("Database transaction failed")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"user_id":  user.ID.String(),
		"email":    user.Email.String(),
		"username": user.Username.String(),
	}).Info("User registration completed successfully")

	return &dto.RegisterResp{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login handles user login
func (s *UserService) Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error) {
	logger := log.WithFields(logrus.Fields{
		"method": "Login",
		"email":  req.Email,
	})

	logger.Info("Starting user login")

	if err := req.Validate(); err != nil {
		logger.WithError(err).Error("Request validation failed")
		return nil, err
	}

	logger.Debug("Retrieving user by email")
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve user by email")
		return nil, err
	}

	logger.WithField("user_id", user.ID.String()).Debug("Verifying password")
	if !user.PasswordHash.VerifyPassword(req.Password) {
		logger.WithFields(logrus.Fields{
			"user_id": user.ID.String(),
			"email":   user.Email.String(),
		}).Warn("Invalid password provided")
		return nil, errs.ErrInvalidCredentials
	}

	logger.WithField("user_id", user.ID.String()).Debug("Creating token pair")
	accessToken, refreshToken, err := s.tokenMaker.CreateTokenPair(
		user.ID.String(),
		user.Username.String(),
		int64(s.config.JWT.AccessTokenDuration),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to create token pair")
		return nil, err
	}

	logger.Debug("Starting database transaction")
	err = s.txManager.WithTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, tx.TransactionContextKey, txWrapper.GetTx())

		logger.Debug("Creating refresh token model")
		refreshTokenModel, err := models.NewRefreshToken(
			user.ID,
			refreshToken,
			time.Now().Add(s.config.JWT.RefreshTokenDuration).UnixMilli(),
		)
		if err != nil {
			logger.WithError(err).Error("Failed to create refresh token model")
			return err
		}

		logger.Debug("Storing refresh token in database")
		if err := s.refreshTokenRepo.Create(txCtx, refreshTokenModel); err != nil {
			logger.WithError(err).Error("Failed to store refresh token in database")
			return err
		}

		logger.Debug("Database transaction completed successfully")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("Database transaction failed")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"user_id":  user.ID.String(),
		"email":    user.Email.String(),
		"username": user.Username.String(),
	}).Info("User login completed successfully")

	return &dto.LoginResp{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req dto.RefreshTokenReq) (*dto.RefreshTokenResp, error) {
	logger := log.WithFields(logrus.Fields{
		"method":       "RefreshToken",
		"token_length": len(req.RefreshToken),
	})

	logger.Info("Starting token refresh")

	if err := req.Validate(); err != nil {
		logger.WithError(err).Error("Request validation failed")
		return nil, err
	}

	logger.Debug("Retrieving refresh token from database")
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		if err == errs.ErrTokenNotFound {
			logger.Warn("Refresh token not found in database")
			return nil, errs.ErrTokenNotFound
		}

		logger.WithError(err).Error("Failed to retrieve refresh token from database")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"token_id":   refreshToken.ID.String(),
		"user_id":    refreshToken.UserID.String(),
		"expires_at": refreshToken.ExpiresAt,
		"is_revoked": refreshToken.IsRevoked,
	}).Debug("Retrieved refresh token")

	if refreshToken.IsRevoked {
		logger.WithFields(logrus.Fields{
			"token_id": refreshToken.ID.String(),
			"user_id":  refreshToken.UserID.String(),
		}).Warn("Refresh token is revoked")
		return nil, errs.ErrTokenRevoked
	}

	if refreshToken.ExpiresAt < time.Now().UnixMilli() {
		logger.WithFields(logrus.Fields{
			"token_id":     refreshToken.ID.String(),
			"user_id":      refreshToken.UserID.String(),
			"expires_at":   refreshToken.ExpiresAt,
			"current_time": time.Now().UnixMilli(),
		}).Warn("Refresh token has expired")
		return nil, errs.ErrTokenExpired
	}

	logger.WithField("user_id", refreshToken.UserID.String()).Debug("Retrieving user by ID")
	user, err := s.userRepo.GetByID(ctx, refreshToken.UserID)
	if err != nil {
		logger.WithError(err).WithField("user_id", refreshToken.UserID.String()).Error("Failed to retrieve user by ID")
		return nil, err
	}

	logger.WithField("user_id", user.ID.String()).Debug("Creating new access token")
	accessToken, err := s.tokenMaker.CreateAccessToken(
		user.ID.String(),
		user.Username.String(),
		int64(s.config.JWT.AccessTokenDuration),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to create access token")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"user_id":  user.ID.String(),
		"email":    user.Email.String(),
		"username": user.Username.String(),
		"token_id": refreshToken.ID.String(),
	}).Info("Token refresh completed successfully")

	return &dto.RefreshTokenResp{
		AccessToken: accessToken,
	}, nil
}
