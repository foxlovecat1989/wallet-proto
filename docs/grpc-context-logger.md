# gRPC Context Logger Integration

This document explains how to use `context.WithLogger` in gRPC middleware for structured logging with request-specific context.

## Overview

The gRPC middleware system has been enhanced to automatically inject a logger into the request context, allowing all subsequent interceptors and handlers to access a logger with request-specific fields.

## How It Works

### 1. Context Logger Interceptor

The `ContextLoggerInterceptor` is the first interceptor in the chain that:

- Creates a logger entry with request-specific fields (method name, type)
- Injects this logger into the context using `logutils.WithLogger()`
- Passes the enhanced context to subsequent interceptors and handlers

### 2. Updated Interceptors

All existing interceptors have been updated to:

- Use `logutils.GetLoggerOrDefault(ctx)` to retrieve the logger from context
- Fall back to a default logger if no context logger is available
- Maintain the same functionality while using the context logger

### 3. Handler Integration

In your gRPC handlers, you can now access the context logger:

```go
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // Get logger from context
    logger := logutils.GetLoggerOrDefault(ctx)
    
    // Log with request-specific fields
    logger.WithFields(logrus.Fields{
        "username":    req.Username,
        "email":       req.Email,
        "country_code": req.CountryCode,
        "phone":       req.Phone,
    }).Info("User registration request received")

    // ... your business logic ...

    if err != nil {
        logger.WithError(err).Error("User registration failed")
        return nil, err
    }

    logger.WithFields(logrus.Fields{
        "user_id": resp.User.ID.String(),
        "username": resp.User.Username.String(),
    }).Info("User registration successful")

    return response, nil
}
```

## Interceptor Chain Order

The interceptors are chained in this order:

1. **ContextLoggerInterceptor** - Injects logger into context
2. **PanicRecoveryInterceptor** - Recovers from panics
3. **LoggingInterceptor** - Logs request/response information
4. **ErrorHandlingInterceptor** - Handles and converts errors

## Benefits

1. **Request Tracing**: Each request gets its own logger with method-specific fields
2. **Structured Logging**: All logs include consistent fields like `grpc_method`, `grpc_type`
3. **Context Propagation**: The logger is available throughout the entire request lifecycle
4. **Fallback Safety**: If no context logger is available, it falls back to a default logger
5. **Consistent Format**: All logs use the same structured format with timestamps

## Example Log Output

```json
{
  "level": "info",
  "msg": "User registration request received",
  "time": "2024-01-15T10:30:45.123Z",
  "grpc_method": "/user.UserService/Register",
  "grpc_type": "unary",
  "username": "john_doe",
  "email": "john@example.com",
  "country_code": "+1",
  "phone": "5551234567"
}
```

## Adding Custom Fields

You can add custom fields to the context logger at any point:

```go
// Add request ID
ctx = logutils.WithRequestID(ctx, "req-123")

// Add user ID
ctx = logutils.WithUserID(ctx, "user-456")

// Add multiple fields
ctx = logutils.WithContextFields(ctx, logrus.Fields{
    "correlation_id": "corr-789",
    "client_ip": "192.168.1.1",
})
```

## Stream Support

The same pattern works for gRPC streams with `ContextLoggerStreamInterceptor`:

```go
func (h *UserHandler) StreamUsers(stream pb.UserService_StreamUsersServer) error {
    logger := logutils.GetLoggerOrDefault(stream.Context())
    
    logger.Info("Stream started")
    
    for {
        req, err := stream.Recv()
        if err != nil {
            logger.WithError(err).Error("Stream receive error")
            return err
        }
        
        logger.WithField("user_id", req.UserId).Info("Processing stream request")
        // ... process request ...
    }
}
```

## Testing

When testing your handlers, you can create a context with a test logger:

```go
func TestUserHandler_Register(t *testing.T) {
    testLogger := logrus.New()
    testLogger.SetOutput(io.Discard) // Suppress output in tests
    
    ctx := logutils.WithLogger(context.Background(), testLogger.WithField("test", true))
    
    handler := NewUserHandler(mockService)
    resp, err := handler.Register(ctx, testRequest)
    
    // ... assertions ...
}
```

## Migration Guide

If you have existing handlers that don't use the context logger:

1. Add the import: `logutils "wallet-user-svc/pkg/utils/log"`
2. Get the logger: `logger := logutils.GetLoggerOrDefault(ctx)`
3. Replace any direct logger usage with the context logger
4. Add structured fields for better observability

The system is backward compatible - handlers that don't use the context logger will continue to work with the default logger.
