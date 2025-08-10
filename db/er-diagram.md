# User Service Database ER Diagram

This document contains the Entity-Relationship diagram for the User Service database schema.

## Database Schema Overview

The User Service database consists of three main tables:
- **users**: Core user authentication and profile information
- **refresh_tokens**: Session management and token storage
- **notification_event_logs**: Event logging for notifications

## ER Diagram

```mermaid
erDiagram
    users {
        UUID id PK "Primary Key"
        VARCHAR(255) email UK "Unique, Not Null"
        VARCHAR(100) username "Not Null"
        VARCHAR(255) password_hash "Not Null"
        VARCHAR(100) first_name "Profile Field"
        VARCHAR(100) last_name "Profile Field"
        VARCHAR(5) country_code "Phone Country Code"
        VARCHAR(15) phone "Phone Number"
        DATE date_of_birth "Profile Field"
        VARCHAR(500) profile_picture_url "Profile Field"
        BIGINT created_at "Timestamp (epoch ms)"
        BIGINT updated_at "Timestamp (epoch ms)"
    }

    refresh_tokens {
        UUID id PK "Primary Key"
        UUID user_id FK "Foreign Key to users.id"
        VARCHAR(500) token "Not Null"
        BIGINT expires_at "Not Null"
        BOOLEAN is_revoked "Default: FALSE"
        BIGINT created_at "Timestamp (epoch ms)"
        BIGINT updated_at "Timestamp (epoch ms)"
    }

    notification_event_logs {
        UUID id PK "Primary Key"
        VARCHAR(255) event_name "Not Null"
        JSONB payload "Not Null"
        VARCHAR(50) status "Default: 'pending'"
        BIGINT created_at "Timestamp (epoch ms)"
        BIGINT updated_at "Timestamp (epoch ms)"
    }

    %% Relationships
    users ||--o{ refresh_tokens : "has many"
    users ||--o{ notification_event_logs : "generates"

    %% Indexes
    users {
        INDEX idx_users_email "email"
        INDEX idx_users_username "username"
        INDEX idx_users_country_code_phone "country_code, phone"
        INDEX idx_users_created_at "created_at"
    }

    refresh_tokens {
        INDEX idx_refresh_tokens_user_id "user_id"
        INDEX idx_refresh_tokens_token_hash "token"
        INDEX idx_refresh_tokens_expires_at "expires_at"
        INDEX idx_refresh_tokens_is_revoked "is_revoked"
        INDEX idx_refresh_tokens_created_at "created_at"
    }

    notification_event_logs {
        INDEX idx_notification_event_logs_event_name_status "event_name, status"
    }
```

## Table Descriptions

### users
The main user table containing authentication and profile information.

**Key Features:**
- UUID primary key with auto-generation
- Unique email constraint
- Password hash storage
- Profile fields (first_name, last_name, country_code, phone, date_of_birth, profile_picture_url)
- Automatic timestamp management (created_at, updated_at)
- Multiple indexes for performance optimization

### refresh_tokens
Manages user session tokens for authentication.

**Key Features:**
- UUID primary key with auto-generation
- Foreign key relationship to users table with CASCADE delete
- Token expiration management
- Revocation tracking
- Automatic timestamp management
- Comprehensive indexing for performance

### notification_event_logs
Stores notification events for processing and tracking.

**Key Features:**
- UUID primary key
- JSONB payload for flexible event data storage
- Status tracking for event processing
- Automatic timestamp management
- Composite index on event_name and status

## Database Features

### Triggers
- **update_updated_at_column()**: Automatically updates the `updated_at` timestamp on any table update

### Indexes
- **Performance Indexes**: Optimized for common query patterns
- **Composite Indexes**: Multi-column indexes for complex queries
- **Unique Constraints**: Ensures data integrity

### Data Types
- **UUID**: Primary keys for scalability and security
- **BIGINT**: Timestamps stored as epoch milliseconds
- **JSONB**: Flexible JSON storage for event payloads
- **VARCHAR**: String fields with appropriate length limits

## Relationships

1. **users → refresh_tokens**: One-to-many relationship
   - A user can have multiple refresh tokens
   - Tokens are automatically deleted when user is deleted (CASCADE)

2. **users → notification_event_logs**: One-to-many relationship
   - Users can generate multiple notification events
   - Events are tracked for audit and processing purposes

## Migration History

- **000001**: Initial schema creation
- **1754804014**: Added user profile fields (first_name, last_name, phone, date_of_birth, profile_picture_url)

## Usage Notes

- All timestamps are stored as BIGINT (epoch milliseconds) for consistency
- UUIDs are used for all primary keys to ensure scalability
- Foreign key relationships maintain referential integrity
- Automatic triggers ensure data consistency
- Comprehensive indexing optimizes query performance
