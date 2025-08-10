-- Drop all tables and functions created in the up migration

-- Drop triggers first
DROP TRIGGER IF EXISTS update_notification_event_logs_updated_at ON notification_event_logs;
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON refresh_tokens;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop tables
DROP TABLE IF EXISTS notification_event_logs;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
