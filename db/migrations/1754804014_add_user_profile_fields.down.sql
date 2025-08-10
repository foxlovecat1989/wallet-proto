-- Remove profile fields from users table
DROP INDEX IF EXISTS idx_users_phone;

ALTER TABLE users 
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS date_of_birth,
DROP COLUMN IF EXISTS profile_picture_url;
