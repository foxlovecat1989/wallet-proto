package errs

import (
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorWrapper_Error(t *testing.T) {
	err := NewError(codes.InvalidArgument, "test error")
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

func TestErrorWrapper_WithDetail(t *testing.T) {
	err := NewError(codes.InvalidArgument, "validation failed").
		WithDetail("field", "email").
		WithDetail("value", "invalid")

	details := err.GetDetails()
	if len(details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(details))
	}

	if details["field"] != "email" {
		t.Errorf("Expected field 'email', got '%v'", details["field"])
	}

	if details["value"] != "invalid" {
		t.Errorf("Expected value 'invalid', got '%v'", details["value"])
	}
}

func TestErrorWrapper_WithContext(t *testing.T) {
	err := NewError(codes.NotFound, "user not found").
		WithRequestID("req-123").
		WithUserID("user-456").
		WithOperation("GetUser")

	if err.RequestID != "req-123" {
		t.Errorf("Expected RequestID 'req-123', got '%s'", err.RequestID)
	}

	if err.UserID != "user-456" {
		t.Errorf("Expected UserID 'user-456', got '%s'", err.UserID)
	}

	if err.Operation != "GetUser" {
		t.Errorf("Expected Operation 'GetUser', got '%s'", err.Operation)
	}
}

func TestErrorWrapper_GRPCStatus(t *testing.T) {
	err := NewError(codes.Unauthenticated, "token expired")
	st := err.GRPCStatus()

	if st.Code() != codes.Unauthenticated {
		t.Errorf("Expected code %v, got %v", codes.Unauthenticated, st.Code())
	}

	if st.Message() != "token expired" {
		t.Errorf("Expected message 'token expired', got '%s'", st.Message())
	}
}

func TestWrapError(t *testing.T) {
	originalErr := status.Error(codes.Internal, "database error")
	wrappedErr := WrapError(originalErr, codes.InvalidArgument, "validation failed")

	if wrappedErr.Code != codes.InvalidArgument {
		t.Errorf("Expected code %v, got %v", codes.InvalidArgument, wrappedErr.Code)
	}

	if wrappedErr.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Unwrap() != originalErr {
		t.Errorf("Expected original error, got %v", wrappedErr.Unwrap())
	}
}

func TestErrorWrapper_Timestamp(t *testing.T) {
	err := NewError(codes.InvalidArgument, "test")

	// Check that timestamp is set (should be close to now)
	now := time.Now()
	diff := now.Sub(err.Timestamp)

	if diff < -time.Second || diff > time.Second {
		t.Errorf("Timestamp should be close to now, got diff: %v", diff)
	}
}

func TestErrorWrapper_GetDetail(t *testing.T) {
	err := NewError(codes.InvalidArgument, "test").
		WithDetail("key1", "value1").
		WithDetail("key2", 123)

	// Test existing key
	if value, exists := err.GetDetail("key1"); !exists || value != "value1" {
		t.Errorf("Expected 'value1', got %v (exists: %t)", value, exists)
	}

	// Test existing key with different type
	if value, exists := err.GetDetail("key2"); !exists || value != 123 {
		t.Errorf("Expected 123, got %v (exists: %t)", value, exists)
	}

	// Test non-existing key
	if value, exists := err.GetDetail("nonexistent"); exists {
		t.Errorf("Expected no value, got %v (exists: %t)", value, exists)
	}
}

func TestErrorWrapper_Chaining(t *testing.T) {
	err := NewError(codes.InvalidArgument, "validation failed").
		WithDetail("field", "email").
		WithRequestID("req-123").
		WithUserID("user-456").
		WithOperation("ValidateEmail").
		WithDetail("value", "invalid")

	// Verify all chained operations worked
	if err.Code != codes.InvalidArgument {
		t.Errorf("Expected code %v, got %v", codes.InvalidArgument, err.Code)
	}

	if err.RequestID != "req-123" {
		t.Errorf("Expected RequestID 'req-123', got '%s'", err.RequestID)
	}

	if err.UserID != "user-456" {
		t.Errorf("Expected UserID 'user-456', got '%s'", err.UserID)
	}

	if err.Operation != "ValidateEmail" {
		t.Errorf("Expected Operation 'ValidateEmail', got '%s'", err.Operation)
	}

	details := err.GetDetails()
	if len(details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(details))
	}
}

func TestToGRPCError_ErrorWrapper(t *testing.T) {
	err := NewError(codes.NotFound, "user not found")
	grpcErr := ToGRPCError(err)

	st, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.NotFound {
		t.Errorf("Expected code %v, got %v", codes.NotFound, st.Code())
	}
}

func TestToGRPCError_LegacyErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{"InvalidEmail", ErrInvalidEmailLegacy, codes.InvalidArgument},
		{"UserNotFound", ErrUserNotFoundLegacy, codes.NotFound},
		{"UserExists", ErrUserExistsLegacy, codes.AlreadyExists},
		{"TokenExpired", ErrTokenExpiredLegacy, codes.Unauthenticated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcErr := ToGRPCError(tt.err)
			st, ok := status.FromError(grpcErr)
			if !ok {
				t.Fatal("Expected gRPC status error")
			}

			if st.Code() != tt.expected {
				t.Errorf("Expected code %v, got %v", tt.expected, st.Code())
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{"ErrInvalidEmail", ErrInvalidEmail, codes.InvalidArgument},
		{"ErrInvalidUsername", ErrInvalidUsername, codes.InvalidArgument},
		{"ErrInvalidPassword", ErrInvalidPassword, codes.InvalidArgument},
		{"ErrUserNotFound", ErrUserNotFound, codes.NotFound},
		{"ErrUserExists", ErrUserExists, codes.AlreadyExists},
		{"ErrInvalidToken", ErrInvalidToken, codes.InvalidArgument},
		{"ErrTokenExpired", ErrTokenExpired, codes.Unauthenticated},
		{"ErrTokenRevoked", ErrTokenRevoked, codes.Unauthenticated},
		{"ErrTokenNotFound", ErrTokenNotFound, codes.NotFound},
		{"ErrTokenIsRequired", ErrTokenIsRequired, codes.InvalidArgument},
		{"ErrInvalidCredentials", ErrInvalidCredentials, codes.Unauthenticated},
		{"ErrEmailIsRequired", ErrEmailIsRequired, codes.InvalidArgument},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper, ok := tt.err.(*ErrorWrapper)
			if !ok {
				t.Fatalf("Expected ErrorWrapper, got %T", tt.err)
			}

			if wrapper.Code != tt.expected {
				t.Errorf("Expected code %v, got %v", tt.expected, wrapper.Code)
			}
		})
	}
}
