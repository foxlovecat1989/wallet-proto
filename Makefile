# Makefile for user-svc

.PHONY: all build test clean run proto help

# Default target
all: build

# Build the application
build:
	@echo "Building user-svc..."
	go build -o user-svc-api ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f user-svc-api

# Run the application
run: build
	@echo "Running user-svc..."
	./user-svc-api

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

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Build the application (default)"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application"
	@echo "  proto        - Update submodule and generate proto files"
	@echo "  help         - Show this help message"