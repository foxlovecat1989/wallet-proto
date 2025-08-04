package app

import (
	"context"

	"user-svc/internal/app/domains/dto"
	pb "user-svc/api/proto"
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
	resp, err := h.userService.Register(ctx, dto.RegisterReq{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:       resp.User.ID.String(),
			Email:    resp.User.Email.String(),
			Username: resp.User.Username.String(),
		},
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// Login handles user login
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := h.userService.Login(ctx, dto.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Id:       resp.User.ID.String(),
			Email:    resp.User.Email.String(),
			Username: resp.User.Username.String(),
		},
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
