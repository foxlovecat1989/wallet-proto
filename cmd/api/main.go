package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// TODO: Load configuration
	// cfg, err := app.LoadConfig("configs/config.yaml")
	// if err != nil {
	// 	log.Fatalf("Failed to load config: %v", err)
	// }

	// TODO: Connect to database
	// dbConn, err := db.NewConnection(cfg.Database)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	// TODO: Initialize database schema
	// log.Println("Initializing database schema...")
	// if err := db.InitDatabase(dbConn.DB.DB); err != nil {
	// 	log.Fatalf("Failed to initialize database schema: %v", err)
	// }

	// Initialize repositories
	// userRepo := db.NewUserRepository(dbConn.DB)
	// refreshTokenRepo := db.NewRefreshTokenRepository(dbConn.DB)

	// Initialize token maker
	// tokenMaker := utils.NewJWTTokenMaker(cfg.Security.JWT.SecretKey)

	// txManager := utils.NewTransactionManager(dbConn.DB)

	// Initialize gRPC server with error handling
	// userServer := app.NewUserHandler(nil) // TODO: Implement proper service

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	// pb.RegisterUserServiceServer(grpcServer, userServer)

	// Enable reflection for development
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcAddr := ":50051" // TODO: Use config
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on %s", grpcAddr)
	log.Printf("Server configuration:")
	log.Printf("  - Database: TODO")
	log.Printf("  - Token Duration: TODO")
	log.Printf("  - Reflection: enabled")

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal %v, initiating graceful shutdown...", sig)

	// Create a context with timeout for graceful shutdown
	shutdownTimeout := 30 * time.Second // TODO: Use config
	_, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Gracefully stop the gRPC server
	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")

	// TODO: Close database connection
	// log.Println("Closing database connection...")
	// if err := dbConn.Close(); err != nil {
	// 	log.Printf("Error closing database connection: %v", err)
	// } else {
	// 	log.Println("Database connection closed")
	// }

	log.Println("Graceful shutdown completed")
}
