-- 005_v3_refactor.up.sql
-- Refactor schema for v3: competitive bidding model with blind voting.

-- ============================================================
-- 1. Drop old tables (order matters for FK dependencies)
-- ============================================================
DROP TABLE IF EXISTS credit_transactions CASCADE;
DROP TABLE IF EXISTS contributions CASCADE;
DROP TABLE IF EXISTS prds CASCADE;
DROP TABLE IF EXISTS token_usage_logs CASCADE;
DROP TABLE IF EXISTS document_versions CASCADE;
DROP TABLE IF EXISTS documents CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;

-- ============================================================
-- 2. Modify users table: drop v2-only columns
-- ============================================================
ALTER TABLE users DROP COLUMN IF EXISTS agent_api_key;
ALTER TABLE users DROP COLUMN IF EXISTS credits_balance;

-- ============================================================
-- 3. Modify ideas table for v3
-- ============================================================
-- Drop old columns
ALTER TABLE ideas DROP COLUMN IF EXISTS target_user_hint;
ALTER TABLE ideas DROP COLUMN IF EXISTS problem_definition;
ALTER TABLE ideas DROP COLUMN IF EXISTS initiator_cut_percent;
ALTER TABLE ideas DROP COLUMN IF EXISTS package_type;

-- Add new columns
ALTER TABLE ideas ADD COLUMN target_user TEXT NOT NULL DEFAULT '';
ALTER TABLE ideas ADD COLUMN core_problem TEXT NOT NULL DEFAULT '';
ALTER TABLE ideas ADD COLUMN out_of_scope TEXT;
ALTER TABLE ideas ADD COLUMN revealed_at TIMESTAMPTZ;

-- Make deadline NOT NULL (v3 always sets a 7-day deadline)
ALTER TABLE ideas ALTER COLUMN deadline SET NOT NULL;
ALTER TABLE ideas ALTER COLUMN deadline SET DEFAULT (NOW() + INTERVAL '7 days');

-- Update status constraint for v3
ALTER TABLE ideas DROP CONSTRAINT IF EXISTS ideas_status_check;
ALTER TABLE ideas ADD CONSTRAINT ideas_status_check
    CHECK (status IN ('open', 'closed', 'cancelled'));

-- Migrate existing status values
UPDATE ideas SET status = 'open' WHERE status IN ('draft', 'active');
UPDATE ideas SET status = 'closed' WHERE status = 'completed';

-- ============================================================
-- 4. Create new tables
-- ============================================================

-- Contributions (each user can submit at most 1 per idea)
CREATE TABLE contributions (
    id           BIGSERIAL PRIMARY KEY,
    idea_id      BIGINT NOT NULL REFERENCES ideas(id),
    author_id    BIGINT NOT NULL REFERENCES users(id),
    content      TEXT NOT NULL DEFAULT '',
    decision_log JSONB NOT NULL DEFAULT '[]',
    status       TEXT NOT NULL DEFAULT 'draft'
                 CHECK (status IN ('draft', 'submitted')),
    view_count   INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submitted_at TIMESTAMPTZ,
    UNIQUE (idea_id, author_id)
);

-- Votes (blind voting, one vote per user per idea)
CREATE TABLE votes (
    id              BIGSERIAL PRIMARY KEY,
    idea_id         BIGINT NOT NULL REFERENCES ideas(id),
    voter_id        BIGINT NOT NULL REFERENCES users(id),
    contribution_id BIGINT NOT NULL REFERENCES contributions(id),
    voted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (idea_id, voter_id)
);

-- Rate limits (daily quotas, atomic upsert)
CREATE TABLE rate_limits (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    action      TEXT NOT NULL CHECK (action IN ('post_idea', 'vote')),
    action_date DATE NOT NULL DEFAULT CURRENT_DATE,
    count       INT NOT NULL DEFAULT 1,
    UNIQUE (user_id, action, action_date)
);

-- Reveal snapshots (frozen ranking data at reveal time)
CREATE TABLE reveal_snapshots (
    id             BIGSERIAL PRIMARY KEY,
    idea_id        BIGINT NOT NULL UNIQUE REFERENCES ideas(id),
    ranked_results JSONB NOT NULL,
    total_votes    INT NOT NULL,
    revealed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 5. Indexes
-- ============================================================
CREATE INDEX idx_ideas_deadline ON ideas(deadline);
CREATE INDEX idx_contributions_idea_id ON contributions(idea_id);
CREATE INDEX idx_contributions_author_id ON contributions(author_id);
CREATE INDEX idx_votes_idea_id ON votes(idea_id);
CREATE INDEX idx_votes_contribution_id ON votes(contribution_id);
CREATE INDEX idx_votes_voter_id ON votes(voter_id);
CREATE INDEX idx_rate_limits_user_action_date ON rate_limits(user_id, action, action_date);
