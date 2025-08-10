package grpc

import (
	"context"
	"time"

	logutils "wallet-user-svc/pkg/utils/log"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// LoggingInterceptor is a gRPC interceptor that logs request/response information
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Get logger from context, fallback to default if not available
		logger := logutils.GetLoggerOrDefault(ctx)

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

// StreamLoggingInterceptor is a gRPC stream interceptor that logs stream information
func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Get logger from context, fallback to default if not available
		logger := logutils.GetLoggerOrDefault(stream.Context())

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
