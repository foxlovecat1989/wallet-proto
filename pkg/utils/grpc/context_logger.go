package grpc

import (
	"context"

	logutils "wallet-user-svc/pkg/utils/log"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// ContextLoggerInterceptor is a gRPC interceptor that injects a logger into the context
// This should be the first interceptor in the chain to ensure all subsequent interceptors
// can access the context logger
func ContextLoggerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Create a logger entry with request-specific fields
		requestLogger := logger.WithFields(logrus.Fields{
			"grpc_method": info.FullMethod,
			"grpc_type":   "unary",
		})

		// Inject the logger into the context
		ctxWithLogger := logutils.WithLogger(ctx, requestLogger)

		// Call the handler with the enhanced context
		return handler(ctxWithLogger, req)
	}
}

// ContextLoggerStreamInterceptor is a gRPC stream interceptor that injects a logger into the context
func ContextLoggerStreamInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// Create a logger entry with stream-specific fields
		requestLogger := logger.WithFields(logrus.Fields{
			"grpc_method":      info.FullMethod,
			"grpc_type":        "stream",
			"is_client_stream": info.IsClientStream,
			"is_server_stream": info.IsServerStream,
		})

		// Inject the logger into the context
		ctxWithLogger := logutils.WithLogger(stream.Context(), requestLogger)

		// Create a wrapped stream that provides the enhanced context
		wrappedStream := &contextStream{
			ServerStream: stream,
			ctx:          ctxWithLogger,
		}

		// Call the handler with the wrapped stream
		return handler(srv, wrappedStream)
	}
}

// contextStream wraps grpc.ServerStream to provide enhanced context
type contextStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *contextStream) Context() context.Context {
	return s.ctx
}
