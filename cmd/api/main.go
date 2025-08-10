package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "user-svc/api/proto"
	"user-svc/internal/app/config"
	"user-svc/internal/app/handler"
	"user-svc/internal/app/repository"
	"user-svc/internal/app/service"
	"user-svc/db"
	"user-svc/internal/workers"
	"user-svc/pkg/utils/crypt/token"
	grpcutils "user-svc/pkg/utils/grpc"
	logutils "user-svc/pkg/utils/log"
	"user-svc/pkg/utils/tx"

	"github.com/hibiken/asynq"
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
	notificationEventLogRepo := repository.NewNotificationEventLogRepository(db)

	userService := service.NewUserService(
		cfg,
		userRepo,
		refreshTokenRepo,
		txManager,
		tokenMaker,
		notificationEventLogRepo,
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

	// Create main application context with cancellation
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Start notification worker if enabled
	var notificationWorker *workers.NotificationWorker
	var wg sync.WaitGroup

	if cfg.Worker.Notification.Enabled {
		asyncQClient := asynq.NewClient(asynq.RedisClientOpt{
			Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		})
		defer asyncQClient.Close()

		notificationWorker = workers.NewNotificationWorker(
			logger,
			asyncQClient,
			notificationEventLogRepo,
			&wg,
			cfg.Worker.Notification.Interval,
			cfg.Worker.Notification.MaxRetries,
			cfg.Worker.Notification.BatchSize,
		)

		// Start worker with application context
		go func() {
			notificationWorker.Start(appCtx)
		}()

		logger.WithFields(logrus.Fields{
			"interval":    cfg.Worker.Notification.Interval,
			"max_retries": cfg.Worker.Notification.MaxRetries,
			"batch_size":  cfg.Worker.Notification.BatchSize,
		}).Info("Notification worker started")
	} else {
		logger.Info("Notification worker disabled")
	}

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.WithError(err).Error("gRPC server error")
			serverErrChan <- err
		}
	}()

	logger.Info("gRPC server is running and ready to accept connections")

	// Wait for either shutdown signal or server error
	select {
	case sig := <-sigChan:
		logger.WithField("signal", sig).Info("Received shutdown signal, initiating graceful shutdown")
	case err := <-serverErrChan:
		logger.WithError(err).Error("Server error occurred, initiating shutdown")
	}

	// Create a context with timeout for graceful shutdown
	shutdownTimeout := 30 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Cancel the main application context to signal all components to stop
	appCancel()

	// Wait for all components to finish with timeout
	shutdownDone := make(chan struct{})
	go func() {
		// Wait for notification worker to finish
		if cfg.Worker.Notification.Enabled {
			logger.Info("Waiting for notification worker to stop...")
			wg.Wait()
			logger.Info("Notification worker stopped")
		}

		// Gracefully stop the gRPC server
		logger.Info("Stopping gRPC server...")
		grpcServer.GracefulStop()
		logger.Info("gRPC server stopped")

		close(shutdownDone)
	}()

	// Wait for shutdown to complete or timeout
	select {
	case <-shutdownDone:
		logger.Info("Graceful shutdown completed successfully")
	case <-shutdownCtx.Done():
		logger.Warn("Shutdown timeout reached, forcing shutdown")
		// Force stop the server if graceful shutdown times out
		grpcServer.Stop()
		logger.Info("Forced shutdown completed")
	}
}
