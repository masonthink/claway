-- 002_oauth_accounts.down.sql
DROP TABLE IF EXISTS user_oauth_accounts;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
ALTER TABLE users DROP COLUMN IF EXISTS display_name;
ALTER TABLE users ALTER COLUMN openclaw_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN openclaw_id DROP DEFAULT;
