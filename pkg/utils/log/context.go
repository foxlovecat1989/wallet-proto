package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Context key for logger
type loggerContextKey struct{}

var LoggerContextKey = loggerContextKey{}

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}

// GetLoggerFromContext retrieves a logger from the context
func GetLoggerFromContext(ctx context.Context) (*logrus.Entry, bool) {
	logger, ok := ctx.Value(LoggerContextKey).(*logrus.Entry)
	return logger, ok
}

// GetLoggerOrDefault retrieves a logger from context or returns the default logger
func GetLoggerOrDefault(ctx context.Context) *logrus.Entry {
	if logger, ok := GetLoggerFromContext(ctx); ok {
		return logger
	}
	return GetLogger().WithField("context", "default")
}

// WithRequestID adds a request ID to the logger and context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	logger := GetLoggerOrDefault(ctx).WithField("request_id", requestID)
	return WithLogger(ctx, logger)
}

// WithUserID adds a user ID to the logger and context
func WithUserID(ctx context.Context, userID string) context.Context {
	logger := GetLoggerOrDefault(ctx).WithField("user_id", userID)
	return WithLogger(ctx, logger)
}

// WithContextFields adds multiple fields to the logger and context
func WithContextFields(ctx context.Context, fields logrus.Fields) context.Context {
	logger := GetLoggerOrDefault(ctx).WithFields(fields)
	return WithLogger(ctx, logger)
}
