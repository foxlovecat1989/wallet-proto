-- Remove phone fields from users table
DROP INDEX IF EXISTS idx_users_country_code_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users DROP COLUMN IF EXISTS country_code;
