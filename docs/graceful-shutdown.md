# Graceful Shutdown Mechanism

This document explains the graceful shutdown mechanism implemented in the user-svc application.

## Overview

The application implements a robust graceful shutdown mechanism that ensures all components are properly stopped when the application receives a shutdown signal or encounters an error.

## Components

### 1. Main Application Context

The application uses a main application context (`appCtx`) that is shared across all components:

```go
appCtx, appCancel := context.WithCancel(context.Background())
defer appCancel()
```

### 2. Signal Handling

The application listens for OS signals (SIGINT, SIGTERM) to initiate graceful shutdown:

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
```

### 3. Server Error Handling

The gRPC server runs in a goroutine and any errors are captured to trigger shutdown:

```go
serverErrChan := make(chan error, 1)
go func() {
    if err := grpcServer.Serve(lis); err != nil {
        serverErrChan <- err
    }
}()
```

## Shutdown Flow

### 1. Shutdown Trigger

Shutdown can be triggered by:
- OS signal (SIGINT, SIGTERM)
- Server error
- Manual cancellation

### 2. Context Cancellation

When shutdown is triggered:
1. The main application context is cancelled (`appCancel()`)
2. This signals all components to stop gracefully
3. A shutdown timeout context is created (30 seconds default)

### 3. Component Shutdown

Components are shut down in the following order:

#### Notification Worker
- Receives context cancellation signal
- Processes any remaining events before stopping
- Uses single-threaded sequential processing for predictable behavior
- Processes events immediately on startup, then follows configured intervals

#### gRPC Server
- Calls `GracefulStop()` to stop accepting new connections
- Allows existing requests to complete
- Waits for all active connections to close

### 4. Timeout Handling

If shutdown takes longer than the timeout:
- A warning is logged
- The server is force-stopped using `Stop()`
- The application exits

## Key Features

### Context-Aware Processing

The notification worker checks for context cancellation before processing each event:

```go
select {
case <-ctx.Done():
    s.logger.Info("Context cancelled, stopping event processing")
    return
default:
}
```

### Single-Threaded Processing

Event processing uses sequential processing in a single thread for predictable behavior:

```go
// Process events sequentially in a single thread
for _, event := range events {
    // Check for context cancellation before processing each event
    select {
    case <-ctx.Done():
        s.logger.Info("Context cancelled, stopping event processing")
        return
    default:
    }

    if err := s.processEvent(ctx, event); err != nil {
        s.logger.WithError(err).WithField("eventID", event.ID).Error("Failed to process event")
    }
}
```

### Immediate Processing

The worker processes events immediately on startup and then follows the configured interval.

### Graceful Event Processing

When shutting down, the worker processes any remaining events before stopping:

```go
case <-ctx.Done():
    s.logger.Info("Stopping notification worker (context cancelled)")
    // Process any remaining events before stopping
    s.processPendingEvents(context.Background())
    return
```

## Configuration

The shutdown timeout can be configured by modifying the `shutdownTimeout` variable in `main.go`:

```go
shutdownTimeout := 30 * time.Second
```

## Testing

The graceful shutdown mechanism is tested with various scenarios:

- `TestGracefulShutdown`: Tests basic context cancellation
- `TestShutdownTimeout`: Tests timeout scenarios
- `TestContextCancellation`: Tests immediate cancellation
- `TestGracefulShutdownWithTimeout`: Tests successful timeout handling
- `TestGracefulShutdownTimeoutExceeded`: Tests timeout exceeded scenarios

Run tests with:
```bash
go test ./cmd/api -v
```

## Best Practices

1. **Always check context cancellation** before starting long-running operations
2. **Use timeouts** to prevent indefinite waiting
3. **Log shutdown progress** for debugging
4. **Handle cleanup** in defer statements
5. **Use WaitGroups** to coordinate component shutdown
6. **Implement force shutdown** as a fallback for timeout scenarios

## Monitoring

The application logs shutdown progress at each stage:

- Shutdown signal received
- Worker shutdown initiated
- Server shutdown initiated
- Shutdown completion or timeout
- Force shutdown if necessary

This provides visibility into the shutdown process and helps with debugging. 