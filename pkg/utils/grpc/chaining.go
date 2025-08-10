package grpc

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GetUnaryInterceptors returns a single chained unary interceptor as server option
func GetUnaryInterceptors(logger *logrus.Logger) []grpc.ServerOption {
	// Chain the interceptors in the desired order
	// ContextLoggerInterceptor should be first to ensure logger is available in context
	chainedInterceptor := grpc.ChainUnaryInterceptor(
		ContextLoggerInterceptor(logger),
		PanicRecoveryInterceptor(),
		LoggingInterceptor(),
		ErrorHandlingInterceptor(),
	)

	return []grpc.ServerOption{chainedInterceptor}
}

// GetStreamInterceptors returns a single chained stream interceptor as server option
func GetStreamInterceptors(logger *logrus.Logger) []grpc.ServerOption {
	// Chain the stream interceptors in the desired order
	// ContextLoggerStreamInterceptor should be first to ensure logger is available in context
	chainedInterceptor := grpc.ChainStreamInterceptor(
		ContextLoggerStreamInterceptor(logger),
		StreamPanicRecoveryInterceptor(),
		StreamLoggingInterceptor(),
	)

	return []grpc.ServerOption{chainedInterceptor}
}
