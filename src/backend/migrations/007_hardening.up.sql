-- 007: Hardening — vote constraints, audit logs, missing indexes, auth sessions
-- From CEO review: 6A, T2, T3, 2A

BEGIN;

-- 6A: Prevent voting for own contributions via trigger
-- (CHECK constraints can't reference other tables, so we use a trigger)
CREATE OR REPLACE FUNCTION check_no_self_vote() RETURNS TRIGGER AS $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM contributions
    WHERE id = NEW.contribution_id AND author_id = NEW.voter_id
  ) THEN
    RAISE EXCEPTION 'cannot vote for own contribution';
  END IF;

  -- Also verify contribution is submitted (not draft)
  IF NOT EXISTS (
    SELECT 1 FROM contributions
    WHERE id = NEW.contribution_id AND status = 'submitted'
  ) THEN
    RAISE EXCEPTION 'can only vote for submitted contributions';
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_no_self_vote
  BEFORE INSERT ON votes
  FOR EACH ROW
  EXECUTE FUNCTION check_no_self_vote();

-- T2: Audit log table for tracking key operations
CREATE TABLE IF NOT EXISTS audit_logs (
  id          BIGSERIAL PRIMARY KEY,
  user_id     BIGINT NOT NULL REFERENCES users(id),
  action      TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id BIGINT NOT NULL,
  metadata    JSONB DEFAULT '{}',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action, created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- T3: Missing indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_votes_voter_id ON votes(voter_id);
CREATE INDEX IF NOT EXISTS idx_contributions_author_id ON contributions(author_id);
CREATE INDEX IF NOT EXISTS idx_ideas_status_deadline ON ideas(status, deadline);

-- 2A: Auth sessions table (migrate from in-memory sync.Map)
CREATE TABLE IF NOT EXISTS auth_sessions (
  id          TEXT PRIMARY KEY,
  user_id     BIGINT REFERENCES users(id),
  token       TEXT NOT NULL,
  metadata    JSONB DEFAULT '{}',
  expires_at  TIMESTAMPTZ NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_sessions_expires ON auth_sessions(expires_at);

COMMIT;
