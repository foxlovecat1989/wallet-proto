package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"wallet-user-svc/internal/app/config"
	"wallet-user-svc/internal/app/errs"
	"wallet-user-svc/internal/app/model/domain"
	"wallet-user-svc/internal/app/model/dto"
	"wallet-user-svc/internal/app/model/events"
	"wallet-user-svc/internal/app/repository"
	"wallet-user-svc/pkg/utils/crypt/token"
	"wallet-user-svc/pkg/utils/cx"
	logutils "wallet-user-svc/pkg/utils/log"
	"wallet-user-svc/pkg/utils/tx"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, countryCode, phone string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *domain.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
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

type NotificationEventLogRepository interface {
	Create(ctx context.Context, event *repository.NotificationEventLog) error
}

// UserService handles business logic for user operations
type UserService struct {
	config                   *config.Config
	userRepo                 UserRepository
	refreshTokenRepo         RefreshTokenRepository
	txManager                TxManager
	tokenMaker               token.TokenMaker
	notificationEventLogRepo NotificationEventLogRepository
}

// NewUserService creates a new UserService instance
func NewUserService(
	config *config.Config,
	userRepo UserRepository,
	refreshTokenRepo RefreshTokenRepository,
	txManager TxManager,
	tokenMaker token.TokenMaker,
	notificationEventLogRepo NotificationEventLogRepository,
) *UserService {
	logutils.Info("Initializing UserService")

	service := &UserService{
		config:                   config,
		userRepo:                 userRepo,
		refreshTokenRepo:         refreshTokenRepo,
		txManager:                txManager,
		tokenMaker:               tokenMaker,
		notificationEventLogRepo: notificationEventLogRepo,
	}

	logutils.WithFields(logrus.Fields{
		"access_token_duration":  config.JWT.AccessTokenDuration.String(),
		"refresh_token_duration": config.JWT.RefreshTokenDuration.String(),
	}).Info("UserService initialized successfully")

	return service
}

// Register handles user registration
func (s *UserService) Register(ctx context.Context, req dto.RegisterReq) (*dto.RegisterResp, error) {
	// Get logger from context
	logger := logutils.GetLoggerOrDefault(ctx)

	if err := req.Validate(); err != nil {
		logger.WithError(err).Error("Request validation failed")
		return nil, err
	}

	user, err := domain.NewUserWithPassword(
		req.Email,
		req.Password,
		req.Username,
		req.CountryCode,
		req.Phone,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to create user with password")
		return nil, err
	}

	accessToken, refreshToken, err := s.tokenMaker.CreateTokenPair(
		user.ID.String(),
		user.Username.String(),
		int64(s.config.JWT.AccessTokenDuration),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to create token pair")
		return nil, err
	}

	err = s.txManager.WithTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
		txCtx := context.WithValue(ctx, cx.TransactionContextKey, txWrapper.GetTx())

		if err := s.userRepo.Create(txCtx, user); err != nil {
			logger.WithError(err).Error("Failed to create user in database")
			return err
		}

		logger.Debug("Creating refresh token model")
		refreshTokenModel, err := domain.NewRefreshToken(
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
	// Get logger from context
	logger := logutils.GetLoggerOrDefault(ctx)

	logger.Info("Starting user login")

	// Validate email is provided
	if req.Email == "" {
		logger.Error("Email is required for login")
		return nil, errs.ErrEmailIsRequired
	}

	user, err := s.authenticateUser(ctx, req, logger)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.createTokenPair(user, logger)
	if err != nil {
		return nil, err
	}

	if err := s.storeRefreshToken(ctx, user, refreshToken, logger); err != nil {
		return nil, err
	}

	s.logLoginSuccess(user, logger)

	if err := s.createLoginNotification(ctx, user, logger); err != nil {
		return nil, err
	}

	return &dto.LoginResp{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) authenticateUser(ctx context.Context, req dto.LoginReq, logger *logrus.Entry) (*domain.User, error) {
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

	return user, nil
}

func (s *UserService) createTokenPair(user *domain.User, logger *logrus.Entry) (string, string, error) {
	logger.WithField("user_id", user.ID.String()).Debug("Creating token pair")
	accessToken, refreshToken, err := s.tokenMaker.CreateTokenPair(
		user.ID.String(),
		user.Username.String(),
		int64(s.config.JWT.AccessTokenDuration),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to create token pair")
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *UserService) storeRefreshToken(ctx context.Context, user *domain.User, refreshToken string, logger *logrus.Entry) error {
	logger.Debug("Starting database transaction")
	return s.txManager.WithTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
		txCtx := context.WithValue(ctx, cx.TransactionContextKey, txWrapper.GetTx())

		logger.Debug("Creating refresh token model")
		refreshTokenModel, err := domain.NewRefreshToken(
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
}

func (s *UserService) logLoginSuccess(user *domain.User, logger *logrus.Entry) {
	logFields := logrus.Fields{
		"user_id":  user.ID.String(),
		"username": user.Username.String(),
	}
	if user.Email != nil {
		logFields["email"] = user.Email.String()
	}
	if user.CountryCode != nil && user.Phone != nil {
		logFields["country_code"] = *user.CountryCode
		logFields["phone"] = *user.Phone
	}
	logger.WithFields(logFields).Info("User login completed successfully")
}

func (s *UserService) createLoginNotification(ctx context.Context, user *domain.User, logger *logrus.Entry) error {
	notificationParams := dto.SendLoginNotificationParams{
		UserID:   user.ID.String(),
		Username: user.Username.String(),
		LoginAt:  time.Now(),
	}
	if user.Email != nil {
		email := user.Email.String()
		notificationParams.Email = &email
	}
	payload, err := json.Marshal(notificationParams)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal notification payload")
		return err
	}

	if err := s.notificationEventLogRepo.Create(ctx, &repository.NotificationEventLog{
		ID:        uuid.New().String(),
		EventName: string(events.LoginEventType),
		Payload:   payload,
		Status:    repository.NotificationEventLogStatusPending,
	}); err != nil {
		logger.WithError(err).Error("Failed to create notification event log")
		return err
	}

	return nil
}

func (s *UserService) RefreshToken(ctx context.Context, req dto.RefreshTokenReq) (*dto.RefreshTokenResp, error) {
	// Get logger from context
	logger := logutils.GetLoggerOrDefault(ctx)

	logger.Info("Starting token refresh")

	// Validate refresh token is provided
	if req.RefreshToken == "" {
		logger.Error("Refresh token is required")
		return nil, errs.ErrTokenIsRequired
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
