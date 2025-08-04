package errs

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc/codes"
)

// ExampleUsage demonstrates how to use the customized error wrapper
func ExampleUsage() {
	fmt.Println("=== Customized Error Wrapper Examples ===")

	// Example 1: Basic error wrapper usage
	fmt.Println("1. Basic Error Wrapper:")
	basicErr := NewError(codes.InvalidArgument, "validation failed")
	fmt.Printf("   - Error: %s (Code: %s)\n", basicErr.Error(), basicErr.Code)

	// Example 2: Error wrapper with details
	fmt.Println("\n2. Error Wrapper with Details:")
	detailedErr := NewError(codes.InvalidArgument, "user validation failed").
		WithDetail("field", "email").
		WithDetail("value", "invalid-email").
		WithDetail("rule", "must be valid email format")
	fmt.Printf("   - Error: %s\n", detailedErr.Error())
	fmt.Printf("   - Details: %v\n", detailedErr.GetDetails())

	// Example 3: Error wrapper with context
	fmt.Println("\n3. Error Wrapper with Context:")
	contextErr := NewError(codes.NotFound, "user not found").
		WithRequestID("req-123").
		WithUserID("user-456").
		WithOperation("GetUserByEmail")
	fmt.Printf("   - Error: %s\n", contextErr.Error())
	fmt.Printf("   - RequestID: %s\n", contextErr.RequestID)
	fmt.Printf("   - UserID: %s\n", contextErr.UserID)
	fmt.Printf("   - Operation: %s\n", contextErr.Operation)

	// Example 4: Wrapping existing errors
	fmt.Println("\n4. Wrapping Existing Errors:")
	dbErr := fmt.Errorf("database connection failed")
	wrappedErr := WrapError(dbErr, codes.Internal, "failed to retrieve user").
		WithDetail("database", "postgres").
		WithDetail("table", "users")
	fmt.Printf("   - Wrapped Error: %s\n", wrappedErr.Error())
	fmt.Printf("   - Original Error: %v\n", wrappedErr.Unwrap())

	// Example 5: Error with stack trace
	fmt.Println("\n5. Error with Stack Trace:")
	stackErr := NewError(codes.Internal, "unexpected error").
		WithStackTrace(getStackTrace())
	fmt.Printf("   - Error: %s\n", stackErr.Error())
	fmt.Printf("   - Stack Trace: %s\n", stackErr.StackTrace)
}

// ExampleServiceMethod shows how to use error wrappers in service methods
func ExampleServiceMethod(ctx context.Context, email string, userID string) error {
	// Validation with detailed error
	if email == "" {
		return ErrEmailIsRequired.
			WithDetail("operation", "user registration").
			WithUserID(userID).
			WithRequestID(getRequestID(ctx))
	}

	if !isValidEmail(email) {
		return ErrInvalidEmail.
			WithDetail("provided_email", email).
			WithDetail("expected_format", "user@domain.com").
			WithUserID(userID).
			WithRequestID(getRequestID(ctx))
	}

	// Business logic with custom error
	if !hasPermission(ctx, userID) {
		return NewError(codes.PermissionDenied, "insufficient permissions").
			WithDetail("required_permission", "admin").
			WithDetail("user_permissions", []string{"read"}).
			WithUserID(userID).
			WithRequestID(getRequestID(ctx))
	}

	// Database operations with wrapped errors
	if err := saveToDatabase(email); err != nil {
		return WrapError(err, codes.Internal, "failed to save user").
			WithDetail("database_operation", "INSERT").
			WithDetail("table", "users").
			WithUserID(userID).
			WithRequestID(getRequestID(ctx)).
			WithStackTrace(getStackTrace())
	}

	return nil
}

// ExampleHandlerMethod shows how handlers can use error wrappers
func ExampleHandlerMethod(ctx context.Context, req interface{}) (interface{}, error) {
	// Call service method
	err := ExampleServiceMethod(ctx, "test@example.com", "user-123")
	if err != nil {
		// The error wrapper automatically converts to gRPC status
		return nil, err
	}

	return "success", nil
}

// ExampleErrorRecovery shows how to extract information from error wrappers
func ExampleErrorRecovery(err error) {
	fmt.Println("\n=== Error Recovery Example ===")

	if wrapper, ok := err.(*ErrorWrapper); ok {
		fmt.Printf("Error Code: %s\n", wrapper.Code)
		fmt.Printf("Error Message: %s\n", wrapper.Message)
		fmt.Printf("Timestamp: %s\n", wrapper.Timestamp)
		fmt.Printf("Request ID: %s\n", wrapper.RequestID)
		fmt.Printf("User ID: %s\n", wrapper.UserID)
		fmt.Printf("Operation: %s\n", wrapper.Operation)
		fmt.Printf("Details: %v\n", wrapper.GetDetails())
		fmt.Printf("Stack Trace: %s\n", wrapper.StackTrace)

		// Extract specific details
		if field, exists := wrapper.GetDetail("field"); exists {
			fmt.Printf("Failed Field: %v\n", field)
		}
	} else {
		fmt.Printf("Standard Error: %v\n", err)
	}
}

// ExampleChaining shows how to chain error wrapper methods
func ExampleChaining() {
	fmt.Println("\n=== Method Chaining Example ===")

	err := NewError(codes.InvalidArgument, "validation failed").
		WithDetail("field", "email").
		WithDetail("value", "invalid").
		WithRequestID("req-789").
		WithUserID("user-456").
		WithOperation("ValidateEmail").
		WithStackTrace(getStackTrace())

	fmt.Printf("Chained Error: %s\n", err.Error())
	fmt.Printf("All Details: %v\n", err.GetDetails())
}

// Helper functions
func isValidEmail(email string) bool {
	return len(email) > 0 && email != "invalid"
}

func hasPermission(ctx context.Context, userID string) bool {
	// Simulate permission check
	return false
}

func saveToDatabase(email string) error {
	// Simulate database error
	return fmt.Errorf("database connection failed")
}

func getRequestID(ctx context.Context) string {
	// Simulate getting request ID from context
	return "req-123"
}

func getStackTrace() string {
	// Get current stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// ExampleErrorTypes shows different ways to create error wrappers
func ExampleErrorTypes() {
	fmt.Println("\n=== Error Types Examples ===")

	// 1. Predefined domain errors
	fmt.Println("1. Predefined Domain Errors:")
	fmt.Printf("   - Invalid Email: %s\n", ErrInvalidEmail.Error())
	fmt.Printf("   - User Not Found: %s\n", ErrUserNotFound.Error())

	// 2. Custom errors with specific codes
	fmt.Println("\n2. Custom Errors:")
	rateLimitErr := NewError(codes.ResourceExhausted, "rate limit exceeded").
		WithDetail("limit", 100).
		WithDetail("current", 101).
		WithDetail("reset_time", "2024-01-01T00:00:00Z")
	fmt.Printf("   - Rate Limit: %s\n", rateLimitErr.Error())

	// 3. Wrapped database errors
	fmt.Println("\n3. Wrapped Database Errors:")
	dbErr := fmt.Errorf("connection timeout")
	wrappedDBErr := WrapError(dbErr, codes.Unavailable, "database unavailable").
		WithDetail("timeout", "30s").
		WithDetail("retry_count", 3)
	fmt.Printf("   - Database Error: %s\n", wrappedDBErr.Error())

	// 4. Authentication errors
	fmt.Println("\n4. Authentication Errors:")
	authErr := NewError(codes.Unauthenticated, "invalid credentials").
		WithDetail("attempt_count", 5).
		WithDetail("max_attempts", 3).
		WithDetail("lockout_duration", "15m")
	fmt.Printf("   - Auth Error: %s\n", authErr.Error())
}
