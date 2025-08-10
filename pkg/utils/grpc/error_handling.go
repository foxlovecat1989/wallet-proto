package grpc

import (
	"context"
	"time"

	"wallet-user-svc/internal/app/errs"
	logutils "wallet-user-svc/pkg/utils/log"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ErrorHandlingInterceptor is a gRPC interceptor that handles errors and converts them to proper gRPC status codes
func ErrorHandlingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Get logger from context, fallback to default if not available
		logger := logutils.GetLoggerOrDefault(ctx)

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

// CustomErrorHandler provides custom error handling for gRPC streams
func CustomErrorHandler(logger *logrus.Logger) func(error) {
	return func(err error) {
		logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"timestamp": time.Now().UTC(),
		}).Error("gRPC stream error")
	}
}
