# Scripts Directory

This directory contains test scripts for the User Service gRPC API.

## Available Scripts

### Test Scripts

#### `test-all.sh`
Comprehensive test script that tests all gRPC endpoints in one file:
- **Register method**: 3 test scenarios (basic, different email, admin user)
- **Login method**: 4 test scenarios (existing users, non-existent user)
- **RefreshToken method**: 5 test scenarios (valid, different, expired, empty, long tokens)

**Usage:**
```bash
make test-all
# or
./scripts/test-all.sh
```

**Features:**
- ✅ No duplication
- ✅ Comprehensive coverage
- ✅ Clear section organization
- ✅ Detailed test summary
- ✅ Single source of truth

## Prerequisites

Before running any test script:

1. **Start the server:**
   ```bash
   make server
   # or
   make server
   ```

2. **Ensure grpcurl is installed:**
   ```bash
   # macOS
   brew install grpcurl
   
   # Linux
   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
   ```

## Server Requirements

- Server must be running on `localhost:50051`
- gRPC reflection must be enabled
- UserService must be registered

## Test Data

The comprehensive test script uses the following test data:

### Register Tests
- `test@example.com` / `testuser` / `password123`
- `alice@example.com` / `alice` / `securepass456`
- `admin@example.com` / `admin` / `admin123`

### Login Tests
- All users from Register tests
- `nonexistent@example.com` / `wrongpassword` (for error testing)

### RefreshToken Tests
- `mock_refresh_token`
- `another_refresh_token`
- `expired_token_123`
- Empty token
- Long token

## Output Format

Each script provides:
- Clear request/response display
- Formatted JSON output
- Test summary with counts
- Error handling for server issues
- Section organization for easy reading

## Script History

Previously, there were separate scripts for each method:
- `test-register.sh` (removed)
- `test-login.sh` (removed)
- `test-refresh-token.sh` (removed)
- `test-grpc.sh` (removed)

These have been consolidated into `test-all.sh` to eliminate duplication and provide a single comprehensive test suite.

## Adding New Tests

When adding new test scenarios:

1. Add them to `test-all.sh` in the appropriate section
2. Update the test summary at the end
3. Ensure the test data is documented above
4. Test the script to ensure it works correctly 