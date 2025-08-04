package grpc

import (
	"context"
	"runtime/debug"
	"time"

	"user-svc/internal/app/domains/errs"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PanicRecoveryInterceptor is a gRPC interceptor that recovers from panics
func PanicRecoveryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				logger.WithFields(logrus.Fields{
					"method":      info.FullMethod,
					"panic":       r,
					"stack_trace": string(debug.Stack()),
					"timestamp":   time.Now().UTC(),
				}).Error("gRPC panic recovered")

				// Create a proper gRPC error response
				err = status.Error(codes.Internal, "Internal server error occurred")
			}
		}()

		// Call the actual handler
		return handler(ctx, req)
	}
}

// ErrorHandlingInterceptor is a gRPC interceptor that handles errors and converts them to proper gRPC status codes
func ErrorHandlingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Call the handler
		resp, err = handler(ctx, req)

		// If there's an error, handle it
		if err != nil {
			// Log the error
			logger.WithFields(logrus.Fields{
				"method":    info.FullMethod,
				"error":     err.Error(),
				"timestamp": time.Now().UTC(),
			}).Error("gRPC error occurred")

			// Convert to gRPC error if it's not already
			if _, ok := status.FromError(err); !ok {
				err = errs.ToGRPCError(err)
			}
		}

		return resp, err
	}
}

// LoggingInterceptor is a gRPC interceptor that logs request/response information
func LoggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		// Log the incoming request
		logger.WithFields(logrus.Fields{
			"method":    info.FullMethod,
			"timestamp": start.UTC(),
		}).Info("gRPC request started")

		// Call the handler
		resp, err = handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log the response
		if err != nil {
			logger.WithFields(logrus.Fields{
				"method":    info.FullMethod,
				"duration":  duration,
				"error":     err.Error(),
				"timestamp": time.Now().UTC(),
			}).Error("gRPC request failed")
		} else {
			logger.WithFields(logrus.Fields{
				"method":    info.FullMethod,
				"duration":  duration,
				"timestamp": time.Now().UTC(),
			}).Info("gRPC request completed")
		}

		return resp, err
	}
}

// GetUnaryInterceptors returns a single chained unary interceptor as server option
func GetUnaryInterceptors(logger *logrus.Logger) []grpc.ServerOption {
	// Chain the interceptors in the desired order
	chainedInterceptor := grpc.ChainUnaryInterceptor(
		PanicRecoveryInterceptor(logger),
		LoggingInterceptor(logger),
		ErrorHandlingInterceptor(logger),
	)

	return []grpc.ServerOption{chainedInterceptor}
}

// CustomErrorHandler provides custom error handling for gRPC streams
func CustomErrorHandler(logger *logrus.Logger) func(error) {
	return func(err error) {
		logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"timestamp": time.Now().UTC(),
		}).Error("gRPC stream error")
	}
}

// StreamPanicRecoveryInterceptor is a gRPC stream interceptor that recovers from panics
func StreamPanicRecoveryInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				logger.WithFields(logrus.Fields{
					"method":      info.FullMethod,
					"panic":       r,
					"stack_trace": string(debug.Stack()),
					"timestamp":   time.Now().UTC(),
				}).Error("gRPC stream panic recovered")

				// Create a proper gRPC error response
				err = status.Error(codes.Internal, "Internal server error occurred")
			}
		}()

		// Call the actual handler
		return handler(srv, stream)
	}
}

// StreamLoggingInterceptor is a gRPC stream interceptor that logs stream information
func StreamLoggingInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()

		// Log the incoming stream
		logger.WithFields(logrus.Fields{
			"method":           info.FullMethod,
			"is_client_stream": info.IsClientStream,
			"is_server_stream": info.IsServerStream,
			"timestamp":        start.UTC(),
		}).Info("gRPC stream started")

		// Call the handler
		err = handler(srv, stream)

		// Calculate duration
		duration := time.Since(start)

		// Log the stream completion
		if err != nil {
			logger.WithFields(logrus.Fields{
				"method":    info.FullMethod,
				"duration":  duration,
				"error":     err.Error(),
				"timestamp": time.Now().UTC(),
			}).Error("gRPC stream failed")
		} else {
			logger.WithFields(logrus.Fields{
				"method":    info.FullMethod,
				"duration":  duration,
				"timestamp": time.Now().UTC(),
			}).Info("gRPC stream completed")
		}

		return err
	}
}

// GetStreamInterceptors returns a single chained stream interceptor as server option
func GetStreamInterceptors(logger *logrus.Logger) []grpc.ServerOption {
	// Chain the stream interceptors in the desired order
	chainedInterceptor := grpc.ChainStreamInterceptor(
		StreamPanicRecoveryInterceptor(logger),
		StreamLoggingInterceptor(logger),
	)

	return []grpc.ServerOption{chainedInterceptor}
}
