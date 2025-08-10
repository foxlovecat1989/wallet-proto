package handler

import (
	"context"

	pb "wallet-user-svc/api/proto"
	"wallet-user-svc/internal/app/model/dto"
	logutils "wallet-user-svc/pkg/utils/log"

	"github.com/sirupsen/logrus"
)

// UserHandler handles gRPC requests for user operations
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService UserService
}

// UserServiceInterface defines the methods that the user service should implement
type UserService interface {
	Register(ctx context.Context, req dto.RegisterReq) (*dto.RegisterResp, error)
	Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenReq) (*dto.RefreshTokenResp, error)
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Get logger from context
	logger := logutils.GetLoggerOrDefault(ctx)

	logger.WithFields(logrus.Fields{
		"username":     req.Username,
		"email":        req.Email,
		"country_code": req.CountryCode,
		"phone":        req.Phone,
	}).Info("User registration request received")

	// Create RegisterReq with proper handling of optional fields
	registerReq := dto.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	}

	// Handle email (can be empty if using phone)
	if req.Email != "" {
		registerReq.Email = &req.Email
	}

	// Handle country code and phone (can be empty if using email)
	if req.CountryCode != "" {
		registerReq.CountryCode = &req.CountryCode
	}
	if req.Phone != "" {
		registerReq.Phone = &req.Phone
	}

	resp, err := h.userService.Register(ctx, registerReq)
	if err != nil {
		logger.WithError(err).Error("User registration failed")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"user_id":  resp.User.ID.String(),
		"username": resp.User.Username.String(),
	}).Info("User registration successful")

	user := &pb.User{
		Id:       resp.User.ID.String(),
		Username: resp.User.Username.String(),
	}

	// Handle optional email
	if resp.User.Email != nil {
		user.Email = resp.User.Email.ToPtrString()
	}

	// Handle optional country code
	if resp.User.CountryCode != nil {
		user.CountryCode = resp.User.CountryCode.ToPtrString()
	}

	// Handle optional phone
	if resp.User.Phone != nil {
		user.Phone = resp.User.Phone.ToPtrString()
	}

	return &pb.RegisterResponse{
		User:         user,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// Login handles user login
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Get logger from context
	logger := logutils.GetLoggerOrDefault(ctx)

	resp, err := h.userService.Login(ctx, dto.LoginReq{
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		logger.WithError(err).Error("User login failed")
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	resp, err := h.userService.RefreshToken(ctx, dto.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
