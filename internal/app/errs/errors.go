package errs

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Domain errors with gRPC status codes
var (
	ErrInvalidEmail         = NewError(codes.InvalidArgument, "invalid email")
	ErrInvalidUsername      = NewError(codes.InvalidArgument, "invalid username")
	ErrInvalidPassword      = NewError(codes.InvalidArgument, "invalid password")
	ErrUserNotFound         = NewError(codes.NotFound, "user not found")
	ErrUserExists           = NewError(codes.AlreadyExists, "user already exists")
	ErrInvalidToken         = NewError(codes.InvalidArgument, "invalid token")
	ErrTokenExpired         = NewError(codes.Unauthenticated, "token expired")
	ErrTokenRevoked         = NewError(codes.Unauthenticated, "token revoked")
	ErrTokenNotFound        = NewError(codes.NotFound, "token not found")
	ErrTokenIsRequired      = NewError(codes.InvalidArgument, "token is required")
	ErrInvalidCredentials   = NewError(codes.Unauthenticated, "invalid credentials")
	ErrEmailIsRequired      = NewError(codes.InvalidArgument, "email is required")
	ErrEmailOrPhoneRequired = NewError(codes.InvalidArgument, "either email or both country code and phone are required")
	ErrInvalidPhoneNumber   = NewError(codes.InvalidArgument, "invalid phone number")
	ErrInvalidCountryCode   = NewError(codes.InvalidArgument, "invalid country code")
)	

// ErrorWrapper is a customizable error wrapper with rich metadata
type ErrorWrapper struct {
	Code       codes.Code
	Message    string
	Details    map[string]interface{}
	Timestamp  time.Time
	RequestID  string
	UserID     string
	Operation  string
	Err        error
	StackTrace string
}

// Error implements the error interface
func (e *ErrorWrapper) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ErrorWrapper) Unwrap() error {
	return e.Err
}

// GRPCStatus returns the gRPC status
func (e *ErrorWrapper) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// WithDetail adds a key-value detail to the error
func (e *ErrorWrapper) WithDetail(key string, value interface{}) *ErrorWrapper {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithRequestID adds a request ID to the error
func (e *ErrorWrapper) WithRequestID(requestID string) *ErrorWrapper {
	e.RequestID = requestID
	return e
}

// WithUserID adds a user ID to the error
func (e *ErrorWrapper) WithUserID(userID string) *ErrorWrapper {
	e.UserID = userID
	return e
}

// WithOperation adds an operation name to the error
func (e *ErrorWrapper) WithOperation(operation string) *ErrorWrapper {
	e.Operation = operation
	return e
}

// WithStackTrace adds stack trace information
func (e *ErrorWrapper) WithStackTrace(stackTrace string) *ErrorWrapper {
	e.StackTrace = stackTrace
	return e
}

// GetDetail retrieves a detail value by key
func (e *ErrorWrapper) GetDetail(key string) (interface{}, bool) {
	if e.Details == nil {
		return nil, false
	}
	value, exists := e.Details[key]
	return value, exists
}

// GetDetails returns all details
func (e *ErrorWrapper) GetDetails() map[string]interface{} {
	return e.Details
}

// NewError creates a new error wrapper
func NewError(code codes.Code, message string) *ErrorWrapper {
	return &ErrorWrapper{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code codes.Code, message string) *ErrorWrapper {
	return &ErrorWrapper{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Err:       err,
		Details:   make(map[string]interface{}),
	}
}

// Legacy error variables for backward compatibility
var (
	ErrInvalidEmailLegacy         = errors.New("invalid email")
	ErrInvalidUsernameLegacy      = errors.New("invalid username")
	ErrInvalidPasswordLegacy      = errors.New("invalid password")
	ErrUserNotFoundLegacy         = errors.New("user not found")
	ErrUserExistsLegacy           = errors.New("user already exists")
	ErrInvalidTokenLegacy         = errors.New("invalid token")
	ErrTokenExpiredLegacy         = errors.New("token expired")
	ErrTokenRevokedLegacy         = errors.New("token revoked")
	ErrTokenNotFoundLegacy        = errors.New("token not found")
	ErrTokenIsRequiredLegacy      = errors.New("token is required")
	ErrInvalidCredentialsLegacy   = errors.New("invalid credentials")
	ErrEmailIsRequiredLegacy      = errors.New("email is required")
	ErrEmailOrPhoneRequiredLegacy = errors.New("either email or both country code and phone are required")
)

// ToGRPCError converts any error to a gRPC error
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already an error wrapper with gRPC status
	if wrapper, ok := err.(*ErrorWrapper); ok {
		return wrapper.GRPCStatus().Err()
	}

	// Check if it's already a gRPC status error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// Map common errors to appropriate gRPC status codes
	switch err {
	case ErrInvalidEmailLegacy, ErrInvalidUsernameLegacy, ErrInvalidPasswordLegacy,
		ErrInvalidTokenLegacy, ErrTokenIsRequiredLegacy, ErrEmailIsRequiredLegacy,
		ErrEmailOrPhoneRequiredLegacy:
		return status.Error(codes.InvalidArgument, err.Error())
	case ErrUserNotFoundLegacy, ErrTokenNotFoundLegacy:
		return status.Error(codes.NotFound, err.Error())
	case ErrUserExistsLegacy:
		return status.Error(codes.AlreadyExists, err.Error())
	case ErrTokenExpiredLegacy, ErrTokenRevokedLegacy, ErrInvalidCredentialsLegacy:
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		// For unknown errors, return internal error
		return status.Error(codes.Internal, err.Error())
	}
}

// FromGRPCError converts a gRPC error back to an error wrapper
func FromGRPCError(err error) *ErrorWrapper {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return NewError(codes.Internal, err.Error())
	}

	return NewError(st.Code(), st.Message())
}
