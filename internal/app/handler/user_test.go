package handler

import (
	"context"
	"testing"
	"time"

	pb "wallet-user-svc/api/proto"
	"wallet-user-svc/internal/app/model/domain"
	"wallet-user-svc/internal/app/model/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserService is a mock implementation of UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req dto.RegisterReq) (*dto.RegisterResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResp), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req dto.LoginReq) (*dto.LoginResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResp), args.Error(1)
}

func (m *MockUserService) RefreshToken(ctx context.Context, req dto.RefreshTokenReq) (*dto.RefreshTokenResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RefreshTokenResp), args.Error(1)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.RegisterRequest
		mockResponse   *dto.RegisterResp
		mockError      error
		expectedError  bool
		expectedFields map[string]interface{}
	}{
		{
			name: "successful registration with email",
			request: &pb.RegisterRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			mockResponse: &dto.RegisterResp{
				User: &domain.User{
					ID:       uuid.New(),
					Email:    func() *domain.Email { e, _ := domain.NewEmail("test@example.com"); return &e }(),
					Username: func() domain.Username { u, _ := domain.NewUsername("testuser"); return u }(),
				},
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_123",
			},
			mockError:     nil,
			expectedError: false,
			expectedFields: map[string]interface{}{
				"email":         "test@example.com",
				"username":      "testuser",
				"access_token":  "access_token_123",
				"refresh_token": "refresh_token_123",
			},
		},
		{
			name: "successful registration with phone",
			request: &pb.RegisterRequest{
				Username:    "testuser",
				Password:    "password123",
				CountryCode: "+1",
				Phone:       "1234567890",
			},
			mockResponse: &dto.RegisterResp{
				User: &domain.User{
					ID:          uuid.New(),
					Username:    func() domain.Username { u, _ := domain.NewUsername("testuser"); return u }(),
					CountryCode: func() *domain.CountryCode { c, _ := domain.NewCountryCode("+1"); return &c }(),
					Phone:       func() *domain.PhoneNumber { p, _ := domain.NewPhoneNumber("+11234567890"); return &p }(),
				},
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_123",
			},
			mockError:     nil,
			expectedError: false,
			expectedFields: map[string]interface{}{
				"username":      "testuser",
				"country_code":  "+1",
				"phone":         "+11234567890",
				"access_token":  "access_token_123",
				"refresh_token": "refresh_token_123",
			},
		},
		{
			name: "registration with service error",
			request: &pb.RegisterRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "password123",
			},
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations
			if tt.mockResponse != nil {
				mockService.On("Register", mock.Anything, mock.MatchedBy(func(req dto.RegisterReq) bool {
					return req.Username == tt.request.Username && req.Password == tt.request.Password
				})).Return(tt.mockResponse, tt.mockError)
			} else {
				mockService.On("Register", mock.Anything, mock.Anything).Return(nil, tt.mockError)
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Call the handler
			response, err := handler.Register(ctx, tt.request)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotNil(t, response.User)
				assert.Equal(t, tt.expectedFields["username"], response.User.Username)
				assert.Equal(t, tt.expectedFields["access_token"], response.AccessToken)
				assert.Equal(t, tt.expectedFields["refresh_token"], response.RefreshToken)

				// Check optional fields
				if email, ok := tt.expectedFields["email"]; ok {
					assert.NotNil(t, response.User.Email)
					assert.Equal(t, email, *response.User.Email)
				}
				if countryCode, ok := tt.expectedFields["country_code"]; ok {
					assert.NotNil(t, response.User.CountryCode)
					assert.Equal(t, countryCode, *response.User.CountryCode)
				}
				if phone, ok := tt.expectedFields["phone"]; ok {
					assert.NotNil(t, response.User.Phone)
					assert.Equal(t, phone, *response.User.Phone)
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.LoginRequest
		mockResponse   *dto.LoginResp
		mockError      error
		expectedError  bool
		expectedFields map[string]interface{}
	}{
		{
			name: "successful login with email",
			request: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockResponse: &dto.LoginResp{
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_123",
			},
			mockError:     nil,
			expectedError: false,
			expectedFields: map[string]interface{}{
				"access_token":  "access_token_123",
				"refresh_token": "refresh_token_123",
			},
		},
		{
			name: "successful login with phone",
			request: &pb.LoginRequest{
				CountryCode: "+1",
				Phone:       "1234567890",
				Password:    "password123",
			},
			mockResponse: &dto.LoginResp{
				AccessToken:  "access_token_123",
				RefreshToken: "refresh_token_123",
			},
			mockError:     nil,
			expectedError: false,
			expectedFields: map[string]interface{}{
				"access_token":  "access_token_123",
				"refresh_token": "refresh_token_123",
			},
		},
		{
			name: "login with service error",
			request: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations
			if tt.mockResponse != nil {
				mockService.On("Login", mock.Anything, mock.MatchedBy(func(req dto.LoginReq) bool {
					return req.Password == tt.request.Password
				})).Return(tt.mockResponse, tt.mockError)
			} else {
				mockService.On("Login", mock.Anything, mock.Anything).Return(nil, tt.mockError)
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Call the handler
			response, err := handler.Login(ctx, tt.request)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedFields["access_token"], response.AccessToken)
				assert.Equal(t, tt.expectedFields["refresh_token"], response.RefreshToken)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_RefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.RefreshTokenRequest
		mockResponse   *dto.RefreshTokenResp
		mockError      error
		expectedError  bool
		expectedFields map[string]interface{}
	}{
		{
			name: "successful token refresh",
			request: &pb.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			mockResponse: &dto.RefreshTokenResp{
				AccessToken: "new_access_token_123",
			},
			mockError:     nil,
			expectedError: false,
			expectedFields: map[string]interface{}{
				"access_token": "new_access_token_123",
			},
		},
		{
			name: "refresh with invalid token",
			request: &pb.RefreshTokenRequest{
				RefreshToken: "invalid_refresh_token",
			},
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations
			if tt.mockResponse != nil {
				mockService.On("RefreshToken", mock.Anything, mock.MatchedBy(func(req dto.RefreshTokenReq) bool {
					return req.RefreshToken == tt.request.RefreshToken
				})).Return(tt.mockResponse, tt.mockError)
			} else {
				mockService.On("RefreshToken", mock.Anything, mock.Anything).Return(nil, tt.mockError)
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Call the handler
			response, err := handler.RefreshToken(ctx, tt.request)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedFields["access_token"], response.AccessToken)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

// Integration test helper functions
func TestUserHandler_Integration(t *testing.T) {
	t.Skip("Integration test - requires running service and database")

	// This test would require:
	// 1. A running gRPC server
	// 2. A test database
	// 3. Real service implementation
	// 4. gRPC client to make actual calls

	t.Run("full registration and login flow", func(t *testing.T) {
		// TODO: Implement integration test
		// This would test the full flow from registration to login
		// using a real gRPC client and server
	})
}

// Benchmark tests
func BenchmarkUserHandler_Register(b *testing.B) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Set up mock response
	mockResponse := &dto.RegisterResp{
		User: &domain.User{
			ID:       uuid.New(),
			Email:    func() *domain.Email { e, _ := domain.NewEmail("test@example.com"); return &e }(),
			Username: func() domain.Username { u, _ := domain.NewUsername("testuser"); return u }(),
		},
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_123",
	}

	mockService.On("Register", mock.Anything, mock.Anything).Return(mockResponse, nil)

	request := &pb.RegisterRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handler.Register(ctx, request)
		require.NoError(b, err)
	}
}
