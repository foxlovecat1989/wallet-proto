# Protocol Buffers

This directory contains the Protocol Buffer (protobuf) definitions for the user service.

## Files

- `user-svc.proto` - User service definitions (authentication, registration, etc.)
- `order-svc.proto` - Order service definitions (order management)

## Usage

### Generate Go code

From the project root:

```bash
make proto
```

This will generate the Go code from the `.proto` files into the `api/proto/` directory.

### Format and validate

From the `proto/` directory:

```bash
make format    # Format proto files
make validate  # Validate proto files
make all       # Format and validate
```

## Dependencies

- `protobuf` - Protocol Buffer compiler
- `buf` - Modern protobuf tooling (optional)

Install with:

```bash
make install
```

## Generated Files

After running `make proto`, the following files will be generated in `api/proto/`:

- `user-svc.pb.go` - User service Go code
- `user-svc_grpc.pb.go` - User service gRPC Go code
- `order-svc.pb.go` - Order service Go code
- `order-svc_grpc.pb.go` - Order service gRPC Go code
