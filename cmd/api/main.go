package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "user-svc/api/proto"
	"user-svc/internal/app/config"
	"user-svc/internal/app/handler"
	"user-svc/internal/app/repository"
	"user-svc/internal/app/service"
	"user-svc/internal/db"
	"user-svc/pkg/utils/crypt/token"
	grpcutils "user-svc/pkg/utils/grpc"
	logutils "user-svc/pkg/utils/log"
	"user-svc/pkg/utils/tx"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	if err := logutils.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger := logutils.GetLogger()

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatalf("Configuration validation failed: %v", err)
	}

	// Get interceptors for exception handling
	unaryInterceptors := grpcutils.GetUnaryInterceptors(logger)
	streamInterceptors := grpcutils.GetStreamInterceptors(logger)

	// Create gRPC server with interceptors
	serverOptions := append(unaryInterceptors, streamInterceptors...)
	grpcServer := grpc.NewServer(serverOptions...)

	db, err := db.NewStore(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to create database store: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	txManager := tx.NewTransactionManager(db.DB())
	tokenMaker := token.NewJWTTokenMaker(cfg.JWT.SecretKey)

	userService := service.NewUserService(
		cfg,
		userRepo,
		refreshTokenRepo,
		txManager,
		tokenMaker,
	)
	userHandler := handler.NewUserHandler(userService)

	// Register services
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Enable reflection for development
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcAddr := cfg.Server.GetServerAddr()
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	logger.WithFields(logrus.Fields{
		"address":              grpcAddr,
		"port":                 cfg.Server.Port,
		"host":                 cfg.Server.Host,
		"db_host":              cfg.Database.Host,
		"db_port":              cfg.Database.Port,
		"jwt_access_duration":  cfg.JWT.AccessTokenDuration,
		"jwt_refresh_duration": cfg.JWT.RefreshTokenDuration,
		"log_level":            cfg.Log.Level,
		"reflection":           "enabled",
	}).Info("gRPC server starting")

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.WithError(err).Error("Failed to serve gRPC server")
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	logger.WithField("signal", sig).Info("Received shutdown signal, initiating graceful shutdown")

	// Create a context with timeout for graceful shutdown
	shutdownTimeout := 30 * time.Second
	_, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Gracefully stop the gRPC server
	logger.Info("Stopping gRPC server...")
	grpcServer.GracefulStop()
	logger.Info("gRPC server stopped")

	logger.Info("Graceful shutdown completed")
}
