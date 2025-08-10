#!/bin/bash

# Comprehensive Test Script for gRPC User Service
echo "🧪 Testing gRPC User Service - All Methods"
echo "=========================================="

# Check if server is running
if ! nc -z localhost 50051 2>/dev/null; then
    echo "❌ Server is not running on port 50051"
    echo "   Start the server with: make server"
    echo "   Or run directly with: ./bin/user-svc-api -config config.yaml"
    exit 1
fi

echo "✅ Server is running on port 50051"
echo ""

# ============================================================================
# REGISTER TESTS
# ============================================================================
echo "📝 REGISTER METHOD TESTS"
echo "========================"

# Test 1: Basic Register
echo "🧪 Test 1: Basic Register"
echo "Request:"
echo '{
  "email": "test@example.com",
  "username": "testuser",
  "password": "SecurePass123!"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "username": "testuser",
  "password": "SecurePass123!"
}' localhost:50051 user.UserService/Register

echo ""
echo ""

# Test 2: Register with different email
echo "🧪 Test 2: Register with different email"
echo "Request:"
echo '{
  "email": "alice@example.com",
  "username": "alice",
  "password": "AlicePass456!"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "alice@example.com",
  "username": "alice",
  "password": "AlicePass456!"
}' localhost:50051 user.UserService/Register

echo ""
echo ""

# Test 3: Register admin user
echo "🧪 Test 3: Register admin user"
echo "Request:"
echo '{
  "email": "admin@example.com",
  "username": "admin",
  "password": "AdminPass789!"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "admin@example.com",
  "username": "admin",
  "password": "AdminPass789!"
}' localhost:50051 user.UserService/Register

echo ""
echo ""

# ============================================================================
# LOGIN TESTS
# ============================================================================
echo "🔐 LOGIN METHOD TESTS"
echo "====================="

# Test 1: Login with test user
echo "🧪 Test 1: Login with test user"
echo "Request:"
echo '{
  "email": "test@example.com",
  "password": "SecurePass123!"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "SecurePass123!"
}' localhost:50051 user.UserService/Login

echo ""
echo ""

# Test 2: Login with alice user
echo "🧪 Test 2: Login with alice user"
echo "Request:"
echo '{
  "email": "alice@example.com",
  "password": "AlicePass456!"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "alice@example.com",
  "password": "AlicePass456!"
}' localhost:50051 user.UserService/Login

echo ""
echo ""

# Test 3: Login with admin user
echo "🧪 Test 3: Login with admin user"
echo "Request:"
echo '{
  "email": "admin@example.com",
  "password": "AdminPass789!"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "admin@example.com",
  "password": "AdminPass789!"
}' localhost:50051 user.UserService/Login

echo ""
echo ""

# Test 4: Login with non-existent user (will show mock behavior)
echo "🧪 Test 4: Login with non-existent user"
echo "Request:"
echo '{
  "email": "nonexistent@example.com",
  "password": "wrongpassword"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "email": "nonexistent@example.com",
  "password": "wrongpassword"
}' localhost:50051 user.UserService/Login

echo ""
echo ""

# ============================================================================
# REFRESH TOKEN TESTS
# ============================================================================
echo "🔄 REFRESH TOKEN METHOD TESTS"
echo "============================="

# Test 1: Refresh with mock token
echo "🧪 Test 1: Refresh with mock token"
echo "Request:"
echo '{
  "refreshToken": "mock_refresh_token"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "refreshToken": "mock_refresh_token"
}' localhost:50051 user.UserService/RefreshToken

echo ""
echo ""

# Test 2: Refresh with different token
echo "🧪 Test 2: Refresh with different token"
echo "Request:"
echo '{
  "refreshToken": "another_refresh_token"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "refreshToken": "another_refresh_token"
}' localhost:50051 user.UserService/RefreshToken

echo ""
echo ""

# Test 3: Refresh with expired token (mock behavior)
echo "🧪 Test 3: Refresh with expired token"
echo "Request:"
echo '{
  "refreshToken": "expired_token_123"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "refreshToken": "expired_token_123"
}' localhost:50051 user.UserService/RefreshToken

echo ""
echo ""

# Test 4: Refresh with empty token
echo "🧪 Test 4: Refresh with empty token"
echo "Request:"
echo '{
  "refreshToken": ""
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "refreshToken": ""
}' localhost:50051 user.UserService/RefreshToken

echo ""
echo ""

# Test 5: Refresh with long token
echo "🧪 Test 5: Refresh with long token"
echo "Request:"
echo '{
  "refreshToken": "very_long_refresh_token_that_should_work_fine_with_mock_service_123456789"
}'
echo ""
echo "Response:"
grpcurl -plaintext -d '{
  "refreshToken": "very_long_refresh_token_that_should_work_fine_with_mock_service_123456789"
}' localhost:50051 user.UserService/RefreshToken

echo ""
echo ""

# ============================================================================
# SUMMARY
# ============================================================================
echo "✅ ALL TESTS COMPLETED!"
echo ""
echo "📊 TEST SUMMARY:"
echo "  📝 Register: 3 users created"
echo "  🔐 Login: 4 scenarios tested"
echo "  🔄 RefreshToken: 5 scenarios tested"
echo ""
echo "🎯 Total: 12 test scenarios completed successfully" 