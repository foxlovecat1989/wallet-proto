package grpc

import (
	"context"
	"runtime/debug"
	"time"

	logutils "wallet-user-svc/pkg/utils/log"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PanicRecoveryInterceptor is a gRPC interceptor that recovers from panics
func PanicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Get logger from context, fallback to default if not available
		logger := logutils.GetLoggerOrDefault(ctx)

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

// StreamPanicRecoveryInterceptor is a gRPC stream interceptor that recovers from panics
func StreamPanicRecoveryInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Get logger from context, fallback to default if not available
		logger := logutils.GetLoggerOrDefault(stream.Context())

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
