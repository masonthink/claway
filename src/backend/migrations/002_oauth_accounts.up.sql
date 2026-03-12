-- 002_oauth_accounts.up.sql
-- Add multi-provider OAuth support

-- Make openclaw_id optional (legacy field)
ALTER TABLE users ALTER COLUMN openclaw_id DROP NOT NULL;
ALTER TABLE users ALTER COLUMN openclaw_id SET DEFAULT '';

-- Add profile fields to users
ALTER TABLE users ADD COLUMN display_name TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN avatar_url TEXT NOT NULL DEFAULT '';

-- OAuth linked accounts
CREATE TABLE IF NOT EXISTS user_oauth_accounts (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider          TEXT NOT NULL,  -- 'x', 'google', 'linkedin'
    provider_user_id  TEXT NOT NULL,
    provider_username TEXT NOT NULL DEFAULT '',
    provider_email    TEXT NOT NULL DEFAULT '',
    access_token      TEXT NOT NULL DEFAULT '',
    refresh_token     TEXT NOT NULL DEFAULT '',
    token_expires_at  TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_oauth_accounts_user_id ON user_oauth_accounts(user_id);
