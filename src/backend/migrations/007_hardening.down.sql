BEGIN;

DROP INDEX IF EXISTS idx_auth_sessions_expires;
DROP TABLE IF EXISTS auth_sessions;

DROP INDEX IF EXISTS idx_ideas_status_deadline;
DROP INDEX IF EXISTS idx_contributions_author_id;
DROP INDEX IF EXISTS idx_votes_voter_id;

DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP TABLE IF EXISTS audit_logs;

DROP TRIGGER IF EXISTS trg_no_self_vote ON votes;
DROP FUNCTION IF EXISTS check_no_self_vote();

COMMIT;
