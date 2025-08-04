# Makefile for user-svc

.PHONY: all build test clean run proto help

# Default target
all: build

# Build the application
build:
	@echo "Building user-svc..."
	@mkdir -p bin
	go build -o bin/user-svc-api ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run the application
run: build
	@echo "Running user-svc..."
	./bin/user-svc-api

# Run server (alias for run)
server: run



# Test all gRPC endpoints
test-all:
	@echo "Testing all gRPC endpoints..."
	./scripts/test-all.sh

# Development: start database and server
dev:
	@echo "Starting development environment..."
	@echo "Starting PostgreSQL database..."
	docker compose up postgres -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Starting user service..."
	@echo "Note: You may need to run with config file: ./bin/user-svc-api -config config.yaml"
	$(MAKE) server

# Setup proto (update submodule and generate files)
proto:
	@echo "Cleaning up existing proto files..."
	rm -rf api/proto/*.pb.go
	@echo "Updating proto submodule..."
	git submodule update --remote proto
	@echo "Generating protobuf files from proto/ to api/proto/..."
	
	protoc --proto_path=proto \
		--go_out=api/proto --go_opt=paths=source_relative \
		--go-grpc_out=api/proto --go-grpc_opt=paths=source_relative \
		proto/*.proto
	@echo "Proto setup completed!"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -f deployments/Dockerfile -t user-svc .
	@echo "Docker image built successfully!"

docker-run:
	@echo "Running Docker container..."
	docker run -p 50051:50051 --name user-svc-container user-svc

docker-up:
	@echo "Starting services with docker compose..."
	docker compose up -d
	@echo "Services started!"

docker-down:
	@echo "Stopping services with docker compose..."
	docker compose down
	@echo "Services stopped!"

docker-clean:
	@echo "Cleaning up docker resources..."
	docker compose down -v --remove-orphans
	docker system prune -f
	@echo "Docker cleanup completed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application"
	@echo "  server       - Run server (alias for run)"
	@echo "  dev          - Start database and server for development"
	@echo "  test-all     - Test all gRPC endpoints"
	@echo "  proto        - Update submodule and generate proto files"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-up    - Start all services with docker-compose"
	@echo "  docker-down  - Stop all services with docker-compose"
	@echo "  docker-clean - Clean up docker volumes and containers"
	@echo "  help         - Show this help message"