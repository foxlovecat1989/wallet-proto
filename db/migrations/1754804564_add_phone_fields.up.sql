-- Add phone fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS country_code VARCHAR(5);
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(15);

-- Create composite index for phone lookup
CREATE INDEX IF NOT EXISTS idx_users_country_code_phone ON users(country_code, phone);
