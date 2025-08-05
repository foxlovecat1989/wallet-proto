# User Service

A gRPC-based user authentication service with clean architecture, domain-driven design, and modern Go practices.

## 🚀 Features

- **User Authentication**: Registration and login with email/password
- **Token Management**: Real JWT token support with access and refresh tokens
- **Database Persistence**: PostgreSQL database with full CRUD operations
- **Domain Models**: Clean domain models with comprehensive validation
- **Repository Pattern**: Real data access layer with transaction support
- **Service Layer**: Business logic separation with transaction management
- **Transaction Management**: Clean transaction handling with configurable isolation levels
- **gRPC API**: Protocol buffer definitions and gRPC server setup
- **Clean Architecture**: Separation of concerns with internal packages
- **Graceful Shutdown**: Robust shutdown mechanism with context cancellation and timeout handling
- **Exception Handling**: Comprehensive panic recovery and error handling system
- **Code Quality**: Clean, maintainable code with no unused functions
- **Error Handling**: Comprehensive customized error wrapper system with rich metadata
- **Configuration Management**: Flexible configuration with environment variables and YAML
- **Event-Driven Architecture**: Asynchronous notification system with event logging
- **Background Workers**: Notification worker with graceful shutdown and concurrency control
- **Redis Integration**: Asynq-based task queue for asynchronous processing
- **Context Management**: Proper context propagation and cancellation throughout the application

## 🔄 Graceful Shutdown

The service implements a production-ready graceful shutdown mechanism that ensures data integrity and proper resource cleanup:

### Key Capabilities

- **Context Coordination**: Single application context coordinates shutdown across all components
- **Signal Handling**: Responds to OS signals (SIGINT, SIGTERM) and server errors
- **Timeout Protection**: 30-second graceful shutdown with force shutdown fallback
- **Worker Management**: Notification worker processes remaining events before stopping
- **Server Graceful Stop**: gRPC server stops accepting new connections gracefully
- **Comprehensive Logging**: Detailed shutdown progress for monitoring and debugging

### Shutdown Process

1. **Trigger**: OS signal or server error initiates shutdown
2. **Coordination**: Main context cancellation signals all components
3. **Worker Cleanup**: Notification worker processes pending events
4. **Server Stop**: gRPC server stops gracefully
5. **Timeout Handling**: Force shutdown if graceful shutdown times out

See [`docs/graceful-shutdown.md`](docs/graceful-shutdown.md) for detailed documentation.

## 🔧 Implementation Status

The service currently uses **REAL implementations** for all major components:

- ✅ **Business Logic**: Fully implemented with proper validation
- ✅ **Domain Models**: Complete with validation rules
- ✅ **Service Layer**: Real service with transaction support
- ✅ **Transaction Management**: Clean transaction handling with configurable isolation levels
- ✅ **Error Handling**: Comprehensive customized error wrapper system with gRPC status codes
- ✅ **Exception Handling**: Panic recovery and error handling interceptors
- ✅ **Repositories**: REAL implementations with PostgreSQL database operations
- ✅ **Database**: REAL PostgreSQL connection with full transaction support
- ✅ **Token Management**: REAL JWT implementation with access and refresh tokens
- ✅ **Event Logging**: Notification event log system with status tracking
- ✅ **Background Workers**: Notification worker with graceful shutdown and concurrency control
- ✅ **Task Queue**: Redis-based Asynq integration for async processing
- ✅ **Graceful Shutdown**: Robust shutdown mechanism with context cancellation and timeout handling

**Note**: The service is now fully functional with real database persistence, JWT token generation, comprehensive error handling, panic recovery, and an event-driven notification system. All components are production-ready implementations.

## 📋 Prerequisites

- **Go 1.24.4** or later
- **Protocol Buffers** compiler (protoc)
- **PostgreSQL** database (for data persistence)
- **Redis** database (for task queue and async processing)
- **Make** (for build automation)

## 🛠️ Installation & Setup

### 1. Clone and Setup

```bash
git clone <repository-url>
cd user-svc
```

### 2. Install Dependencies

```bash
# Install protoc (macOS)
brew install protobuf

# Install Go protobuf plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. Setup Protocol Buffers

```bash
# Generate protobuf files from proto/ to api/proto/
make proto
```

### 4. Setup Database

The service requires a PostgreSQL database and Redis for task queue processing. You can use Docker for quick setup:

```bash
# Start PostgreSQL and Redis with Docker
docker run --name postgres-user-svc \
  -e POSTGRES_DB=users \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15

docker run --name redis-user-svc \
  -p 6379:6379 \
  -d redis:7-alpine

# Or use the provided docker-compose
make docker-up
```

The database schema will be automatically initialized when the service starts.

## ⚙️ Configuration

The service uses a comprehensive configuration system built with [Viper](https://github.com/spf13/viper) that supports multiple formats and sources.

### Configuration Options

- **YAML Format**: Support for YAML configuration files
- **Environment Variables**: Automatic binding with dot-to-underscore conversion
- **Default Values**: Sensible defaults for all settings
- **Validation**: Built-in configuration validation
- **Flexible Loading**: Multiple ways to load configuration

### Quick Start

```bash
# Use default configuration
make server

# Override with environment variables
export SERVER_PORT=50052
export DATABASE_HOST=localhost
make server

# Use configuration file
./user-svc-api -config=config.yaml
```

### Configuration Files

A sample configuration file is provided:
- `config.yaml` - YAML format

### Environment Variables

All configuration can be set via environment variables:

```bash
# Server settings
export SERVER_PORT=50051
export SERVER_HOST=0.0.0.0

# Database settings
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=password
export DATABASE_DB_NAME=user_svc

# Redis settings
export REDIS_HOST=localhost
export REDIS_PORT=6379

# JWT settings
export JWT_SECRET_KEY=your-secret-key
export JWT_ACCESS_TOKEN_DURATION=15m
export JWT_REFRESH_TOKEN_DURATION=168h
```

For detailed configuration documentation, see [`internal/app/config/README.md`](internal/app/config/README.md).

## 🏃‍♂️ Running the Service

### Development

```bash
# Build and run
make run

# Or build and run separately
make build
./user-svc-api
```

The gRPC server will start on `0.0.0.0:50051`.

## 📚 API Documentation

### User Service

#### Register User

```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse)
```

**Request:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username"
  },
  "access_token": "jwt_token_here",
  "refresh_token": "refresh_token_here"
}
```

#### Login User

```protobuf
rpc Login(LoginRequest) returns (LoginResponse)
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username"
  },
  "access_token": "jwt_token_here",
  "refresh_token": "refresh_token_here"
}
```

#### Refresh Token

```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse)
```

**Request:**
```json
{
  "refresh_token": "refresh_token_here"
}
```

**Response:**
```json
{
  "access_token": "new_jwt_token_here"
}
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific test suites
go test ./cmd/api -v                    # Graceful shutdown tests
go test ./internal/app/domains/errs -v  # Error handling tests
go test ./pkg/utils/crypt/password -v   # Password hashing tests
```

### Graceful Shutdown Testing

The graceful shutdown mechanism includes comprehensive tests:

- **Context Cancellation**: Tests proper context propagation and cancellation
- **Timeout Handling**: Tests shutdown timeout scenarios
- **Worker Coordination**: Tests notification worker shutdown behavior
- **Error Scenarios**: Tests shutdown behavior during errors
- **Force Shutdown**: Tests fallback to force shutdown when timeout exceeded

## 🧹 Code Quality & Cleanup

### Recent Improvements

The codebase has been cleaned up to remove unused code and improve maintainability:

#### ✅ **Removed Unused Code**
- **Unused TxWrapper Methods**: Removed 7 unused transaction wrapper methods that were defined but never called
- **Hardcoded Strings**: Replaced hardcoded transaction context keys with proper constants
- **Build Artifacts**: Cleaned up build artifacts and binaries

#### ✅ **Code Quality Improvements**
- **Proper Constants**: Updated repositories to use `tx.TransactionContextKey` instead of hardcoded `"tx"` strings
- **Clean Imports**: Added proper imports and removed unused dependencies
- **Consistent Patterns**: Standardized transaction context usage across all repositories

#### ✅ **Maintained Functionality**
- **Transaction Support**: All transaction functionality preserved and improved
- **API Compatibility**: No breaking changes to the public API
- **Test Coverage**: All existing tests continue to pass

### Code Quality Standards

- **No Unused Functions**: All functions are actively used or removed
- **Proper Error Handling**: Comprehensive error handling throughout
- **Clean Architecture**: Clear separation of concerns
- **Type Safety**: Strong typing with proper validation

### gRPC API Testing

Test the gRPC endpoints using the provided test script:

```bash
# Test all gRPC endpoints
make test-all

# Or run script directly
./scripts/test-all.sh
```

**Prerequisites:**
- Server must be running (`make server`)
- grpcurl must be installed (`brew install grpcurl`)

See `scripts/README.md` for detailed documentation of test scripts.

## 🔧 Development

### Using gRPC Tools

The server has gRPC reflection enabled for development:

```bash
# List services
grpcurl -plaintext localhost:50051 list

# List methods
grpcurl -plaintext localhost:50051 list user.UserService

# Call register method
grpcurl -plaintext -d '{
  "email": "test@example.com", 
  "username": "testuser", 
  "password": "password123"
}' localhost:50051 user.UserService/Register
```

### Protocol Buffer Development

```bash
# Generate protobuf files
make proto
```

### Transaction Management

The service uses a clean transaction management system with configurable isolation levels:

- **TransactionManager**: Handles database transaction lifecycle
- **TxWrapper**: Wraps database transactions with helper methods
- **Context Integration**: Transactions are passed through context
- **Automatic Rollback**: Failed transactions are automatically rolled back
- **Proper Cleanup**: All transactions are properly committed or rolled back
- **Configurable Isolation**: Support for different transaction isolation levels

#### Available Transaction Methods

```go
// Default transaction (Read Committed)
err = s.txManager.WithTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
    txCtx := context.WithValue(ctx, tx.TransactionContextKey, txWrapper.GetTx())
    // Use txCtx for database operations
    return nil
})

// Custom isolation level
err = s.txManager.WithTransactionIsolation(ctx, func(txWrapper *tx.TxWrapper) error {
    // Use serializable isolation
    return nil
}, sql.LevelSerializable)

// Read-only transaction
err = s.txManager.WithReadOnlyTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
    // Read-only operations only
    return nil
})

// Convenience methods for common isolation levels
err = s.txManager.WithSerializableTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
    // Serializable isolation
    return nil
})

err = s.txManager.WithRepeatableReadTransaction(ctx, func(txWrapper *tx.TxWrapper) error {
    // Repeatable read isolation
    return nil
})

// Custom transaction options
opts := &sql.TxOptions{
    Isolation: sql.LevelSerializable,
    ReadOnly:  true,
}
err = s.txManager.WithTransactionOptions(ctx, func(txWrapper *tx.TxWrapper) error {
    // Custom options
    return nil
}, opts)
```

#### Isolation Levels

- **Read Committed** (default): Prevents dirty reads
- **Read Uncommitted**: Lowest isolation, allows dirty reads
- **Repeatable Read**: Prevents non-repeatable reads
- **Serializable**: Highest isolation, prevents phantom reads

## 📁 Project Structure

```
user-svc/
├── api/
│   └── proto/              # Generated protobuf files
├── cmd/
│   └── api/
│       ├── main.go         # Application entry point with graceful shutdown
│       └── main_test.go    # Graceful shutdown tests
├── deployments/            # Deployment configurations
│   ├── Dockerfile
│   └── k8s.yaml
├── docs/                   # Documentation
│   └── graceful-shutdown.md # Graceful shutdown documentation
├── internal/               # Private application code
│   ├── app/               # Application layer
│   │   ├── config/        # Configuration system
│   │   ├── domains/       # Domain models and business rules
│   │   │   ├── dto/       # Data transfer objects
│   │   │   ├── errs/      # Domain errors
│   │   │   ├── events/    # Event definitions and types
│   │   │   └── models/    # Domain models
│   │   ├── handler/       # gRPC handlers
│   │   ├── repository/    # Data access layer
│   │   └── service/       # Business logic layer
│   └── db/                # Database layer
│       ├── init.sql       # Database initialization
│       └── store.go       # Database store
├── pkg/                   # Public utilities
│   └── utils/             # Utility functions
│       ├── crypt/         # Cryptography utilities
│       │   └── token/     # Token management
│       ├── grpc/          # gRPC interceptors and utilities
│       ├── log/           # Logging utilities
│       └── tx/            # Transaction management utilities
├── workers/               # Background workers
│   └── notificaiton.go    # Notification worker with graceful shutdown
├── scripts/               # Test and utility scripts
│   ├── test-all.sh        # Comprehensive gRPC tests (all methods)
│   └── README.md          # Scripts documentation
├── proto/                 # Protocol buffer definitions
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── Makefile              # Build automation
└── README.md             # Project documentation
```

## 🐳 Docker

### Build Docker Image

```bash
# Using Makefile
make docker-build

# Or manually
docker build -f deployments/Dockerfile -t user-svc .
```

### Run Docker Container

```bash
# Using Makefile
make docker-run

# Or manually
docker run -p 50051:50051 --name user-svc-container user-svc
```

### Docker Image Details

- **Base Image**: Alpine Linux (lightweight)
- **Go Version**: 1.24
- **Multi-stage Build**: Optimized for size
- **Security**: Runs as non-root user
- **Health Check**: Built-in health monitoring
- **Port**: 50051 (gRPC)

## 🐳 Docker Compose

### Production Environment

Start all services (PostgreSQL + User Service):

```bash
# Start all services
make docker-up

# Or manually
docker-compose up -d
```



### Available Services

#### Production (`docker-compose.yml`)
- **user-svc**: User service on port 50051
- **postgres**: PostgreSQL database on port 5432
- **redis**: Redis cache and task queue on port 6379
- **pgadmin**: Database management (optional) on port 8080



### Docker Commands

```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# Clean up docker resources
make docker-clean

# View logs
docker-compose logs -f user-svc

# Access database
docker-compose exec postgres psql -U user -d users

# Access pgAdmin
# Open http://localhost:8080
# Email: admin@example.com
# Password: admin
```

### Environment Variables

The docker-compose files include the following environment variables:

```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=users
DB_SSL_MODE=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
JWT_TOKEN_DURATION=15m
JWT_REFRESH_DURATION=7d

# Server
GRPC_PORT=50051
GRPC_HOST=0.0.0.0

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 📋 Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make clean         # Clean build artifacts
make run           # Build and run the application
make test          # Run all tests
make proto         # Generate protobuf files
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make docker-up     # Start all services
make docker-down   # Stop all services
make docker-clean  # Clean up docker resources
```

## 🔒 Security Features

- **Password Hashing**: Bcrypt with configurable cost
- **Token Security**: JWT token support with refresh tokens
- **Input Validation**: Comprehensive validation for all inputs
- **Error Handling**: Secure error responses without information leakage

## 🛡️ Exception Handling

The service implements a comprehensive exception handling system that prevents server crashes and provides proper error responses:

### Panic Recovery

- **Automatic Panic Recovery**: All panics are caught and converted to gRPC Internal errors
- **Server Stability**: Server continues running even after unexpected panics
- **Detailed Logging**: Panic details are logged with stack traces for debugging
- **Structured Error Responses**: Clients receive proper gRPC status codes instead of connection failures

### Error Handling Interceptors

The system uses gRPC interceptors to handle exceptions at the middleware level:

- **PanicRecoveryInterceptor**: Catches panics and prevents server crashes
- **ErrorHandlingInterceptor**: Converts errors to proper gRPC status codes
- **LoggingInterceptor**: Provides comprehensive request/response logging

### Implementation

```go
// Automatically configured in main.go
unaryInterceptors := grpcutils.GetUnaryInterceptors(logger)
streamInterceptors := grpcutils.GetStreamInterceptors(logger)
serverOptions := append(unaryInterceptors, streamInterceptors...)
grpcServer := grpc.NewServer(serverOptions...)
```

## 🚨 Error Handling

The service uses a comprehensive customized error wrapper system with rich metadata and gRPC status codes:

### Error Wrapper Features

- **Rich Metadata**: Request ID, User ID, Operation name, Timestamp
- **Custom Details**: Key-value pairs for additional context
- **Stack Traces**: Optional stack trace information
- **gRPC Integration**: Automatic conversion to gRPC status errors
- **Method Chaining**: Fluent API for building complex errors
- **Error Wrapping**: Wrap existing errors with additional context

### Standard gRPC Status Codes

- `INVALID_ARGUMENT`: Missing required fields or invalid input
- `NOT_FOUND`: User or token not found
- `ALREADY_EXISTS`: User already exists (registration)
- `UNAUTHENTICATED`: Invalid credentials, expired/revoked tokens
- `INTERNAL`: Server errors
- `PERMISSION_DENIED`: Insufficient permissions
- `RESOURCE_EXHAUSTED`: Rate limiting, quota exceeded

### Usage Examples

#### Basic Error Wrapper

```go
// Create a simple error
err := errs.NewError(codes.InvalidArgument, "validation failed")

// Add details and context
err = errs.NewError(codes.NotFound, "user not found").
    WithDetail("user_id", "123").
    WithRequestID("req-456").
    WithUserID("user-789").
    WithOperation("GetUser")
```

#### Method Chaining

```go
// Chain multiple operations
err := errs.NewError(codes.InvalidArgument, "validation failed").
    WithDetail("field", "email").
    WithDetail("value", "invalid-email").
    WithRequestID("req-123").
    WithUserID("user-456").
    WithOperation("ValidateEmail").
    WithStackTrace(getStackTrace())
```

#### Error Wrapping

```go
// Wrap existing errors with context
dbErr := fmt.Errorf("database connection failed")
err := errs.WrapError(dbErr, codes.Internal, "failed to save user").
    WithDetail("database", "postgres").
    WithDetail("table", "users").
    WithRequestID("req-123")
```

#### Service Layer Usage

```go
func (s *Service) ValidateUser(ctx context.Context, email string) error {
    if email == "" {
        return errs.ErrEmailIsRequired.
            WithDetail("operation", "user registration").
            WithRequestID(getRequestID(ctx))
    }
    
    if !isValidEmail(email) {
        return errs.ErrInvalidEmail.
            WithDetail("provided_email", email).
            WithDetail("expected_format", "user@domain.com")
    }
    
    return nil
}
```

#### Handler Layer (Direct Return)

```go
func (h *Handler) SomeMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    resp, err := h.service.SomeOperation(ctx, req)
    if err != nil {
        return nil, err // ErrorWrapper automatically converts to gRPC status
    }
    return resp, nil
}
```

#### Error Recovery

```go
func handleError(err error) {
    if wrapper, ok := err.(*errs.ErrorWrapper); ok {
        fmt.Printf("Error Code: %s\n", wrapper.Code)
        fmt.Printf("Request ID: %s\n", wrapper.RequestID)
        fmt.Printf("User ID: %s\n", wrapper.UserID)
        fmt.Printf("Operation: %s\n", wrapper.Operation)
        fmt.Printf("Details: %v\n", wrapper.GetDetails())
        
        // Extract specific details
        if field, exists := wrapper.GetDetail("field"); exists {
            fmt.Printf("Failed Field: %v\n", field)
        }
    }
}
```

## 📊 Monitoring & Logging

### Logging Configuration

The service uses structured logging with JSON format by default.

### Graceful Shutdown

The service implements a robust graceful shutdown mechanism that ensures all components are properly stopped when the application receives a shutdown signal or encounters an error.

#### Key Features

- **Context Cancellation**: Main application context coordinates shutdown across all components
- **Signal Handling**: Responds to OS signals (SIGINT, SIGTERM) and server errors
- **Timeout Protection**: 30-second graceful shutdown timeout with force shutdown fallback
- **Worker Coordination**: Notification worker stops gracefully with remaining event processing
- **Server Graceful Stop**: gRPC server stops accepting new connections and waits for active requests
- **Comprehensive Logging**: Detailed shutdown progress logging for debugging

#### Shutdown Flow

1. **Shutdown Trigger**: OS signal or server error triggers shutdown
2. **Context Cancellation**: Main application context is cancelled, signaling all components
3. **Worker Shutdown**: Notification worker processes remaining events and stops
4. **Server Shutdown**: gRPC server stops gracefully
5. **Timeout Handling**: Force shutdown if graceful shutdown exceeds timeout

For detailed documentation, see [`docs/graceful-shutdown.md`](docs/graceful-shutdown.md).

## 🔄 Event-Driven Architecture

The service implements an event-driven architecture for asynchronous processing:

### Notification Event Logging

- **Event Persistence**: All notification events are logged to the database
- **Status Tracking**: Events have pending, success, and failed states
- **Retry Mechanism**: Failed events can be retried automatically
- **Event Types**: Support for different event types (login notifications, etc.)

### Background Worker

- **Notification Worker**: Processes pending notification events sequentially in a single thread
- **Graceful Shutdown**: Worker stops cleanly when service shuts down, processing remaining events
- **Context Cancellation**: Uses context for proper shutdown coordination and cancellation checks
- **Single-Threaded Processing**: Events are processed sequentially for predictable behavior and easier debugging
- **WaitGroup Integration**: Coordinates with main service for graceful shutdown
- **Error Handling**: Comprehensive error handling and logging with event-level error tracking
- **Immediate Processing**: Processes events immediately on startup, then follows configured intervals

### Task Queue Integration

- **Asynq**: Redis-based task queue for asynchronous processing
- **Task Serialization**: Events are serialized and queued for processing
- **Queue Management**: Proper queue management with error handling
- **Redis Integration**: Uses Redis for task persistence and delivery

### Event Flow

```
User Login → Event Logged → Worker Processes → Task Queued → Notification Sent
```

1. **User Login**: When a user logs in, a notification event is logged to the database
2. **Worker Processing**: Background worker polls for pending events
3. **Task Creation**: Worker creates Asynq tasks for notification processing
4. **Queue Processing**: Tasks are queued in Redis for async processing
5. **Status Update**: Event status is updated to success/failed based on processing result

### Worker Configuration

The notification worker runs with configurable intervals and graceful shutdown:

```go
// Worker configuration for single-threaded processing
notificationWorker := workers.NewNotificationWorker(
    logger,
    asyncQClient,
    notificationEventLogRepo,
    &wg,
    time.Second*10, // Polling interval
    3,              // Max retries
    100,            // Batch size
)

// Start worker with application context for coordinated shutdown
go func() {
    notificationWorker.Start(appCtx)
}()
```

### Testing the Graceful Shutdown

The graceful shutdown mechanism is thoroughly tested:

```bash
# Run shutdown tests
go test ./cmd/api -v

# Test scenarios covered:
# - Basic context cancellation
# - Shutdown timeout handling
# - Immediate context cancellation
# - Graceful shutdown with timeout
# - Timeout exceeded scenarios
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Ensure no unused code is introduced
7. Submit a pull request

### Code Quality Guidelines

- **Remove Unused Code**: Don't leave unused functions, variables, or imports
- **Use Constants**: Avoid hardcoded strings, use proper constants
- **Follow Patterns**: Maintain consistency with existing code patterns
- **Test Coverage**: Ensure new functionality is properly tested
- **Clean Architecture**: Maintain separation of concerns
- **Context Management**: Always propagate and check context cancellation
- **Graceful Shutdown**: Ensure all components support graceful shutdown
- **Error Handling**: Use the error wrapper system for consistent error responses

## 📄 License

[Add your license information here]

## 🆘 Support

For issues and questions:
- Create an issue in the repository
- Check the documentation
- Review the test files for usage examples 