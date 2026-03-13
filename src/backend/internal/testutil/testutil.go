package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/claway/server/internal/config"
	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/service"
	"github.com/claway/server/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestJWTSecret is used for all tests.
const TestJWTSecret = "test-jwt-secret"

// defaultTestDatabaseURL is the fallback when TEST_DATABASE_URL is not set.
const defaultTestDatabaseURL = "postgres://mason@localhost:5432/claway_test?sslmode=disable"

// adminDatabaseURL connects to the default postgres database for administrative operations.
const adminDatabaseURL = "postgres://mason@localhost:5432/postgres?sslmode=disable"

// migrationSQL is the schema used to initialize the test database.
// Kept in sync with migrations/005_v3_refactor.up.sql.
const migrationSQL = `
CREATE TABLE IF NOT EXISTS users (
    id           BIGSERIAL PRIMARY KEY,
    openclaw_id  TEXT DEFAULT '',
    username     TEXT NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    avatar_url   TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ideas (
    id           BIGSERIAL PRIMARY KEY,
    initiator_id BIGINT NOT NULL REFERENCES users(id),
    title        TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    target_user  TEXT NOT NULL DEFAULT '',
    core_problem TEXT NOT NULL DEFAULT '',
    out_of_scope TEXT,
    status       TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed', 'cancelled')),
    deadline     TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
    revealed_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contributions (
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

CREATE TABLE IF NOT EXISTS votes (
    id              BIGSERIAL PRIMARY KEY,
    idea_id         BIGINT NOT NULL REFERENCES ideas(id),
    voter_id        BIGINT NOT NULL REFERENCES users(id),
    contribution_id BIGINT NOT NULL REFERENCES contributions(id),
    voted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (idea_id, voter_id)
);

CREATE TABLE IF NOT EXISTS rate_limits (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    action      TEXT NOT NULL CHECK (action IN ('post_idea', 'vote')),
    action_date DATE NOT NULL DEFAULT CURRENT_DATE,
    count       INT NOT NULL DEFAULT 1,
    UNIQUE (user_id, action, action_date)
);

CREATE TABLE IF NOT EXISTS reveal_snapshots (
    id             BIGSERIAL PRIMARY KEY,
    idea_id        BIGINT NOT NULL UNIQUE REFERENCES ideas(id),
    ranked_results JSONB NOT NULL,
    total_votes    INT NOT NULL,
    revealed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_oauth_accounts (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider          TEXT NOT NULL,
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

CREATE INDEX IF NOT EXISTS idx_ideas_deadline ON ideas(deadline);
CREATE INDEX IF NOT EXISTS idx_contributions_idea_id ON contributions(idea_id);
CREATE INDEX IF NOT EXISTS idx_contributions_author_id ON contributions(author_id);
CREATE INDEX IF NOT EXISTS idx_votes_idea_id ON votes(idea_id);
CREATE INDEX IF NOT EXISTS idx_votes_contribution_id ON votes(contribution_id);
CREATE INDEX IF NOT EXISTS idx_votes_voter_id ON votes(voter_id);
CREATE INDEX IF NOT EXISTS idx_rate_limits_user_action_date ON rate_limits(user_id, action, action_date);
`

// testDatabaseURL returns the test database connection string.
func testDatabaseURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return defaultTestDatabaseURL
}

// SetupTestDB connects to claway_test, runs migration, and returns a pool.
// It creates the database if it does not exist.
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	// Connect to the default postgres database to create claway_test if needed.
	adminPool, err := pgxpool.New(ctx, adminDatabaseURL)
	if err != nil {
		t.Fatalf("failed to connect to admin database: %v", err)
	}

	var exists bool
	err = adminPool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'claway_test')").Scan(&exists)
	if err != nil {
		adminPool.Close()
		t.Fatalf("failed to check if test database exists: %v", err)
	}
	if !exists {
		_, err = adminPool.Exec(ctx, "CREATE DATABASE claway_test")
		if err != nil {
			adminPool.Close()
			t.Fatalf("failed to create test database: %v", err)
		}
	}
	adminPool.Close()

	// Connect to the test database.
	pool, err := pgxpool.New(ctx, testDatabaseURL())
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Drop all tables first to ensure clean schema (test DB only).
	_, _ = pool.Exec(ctx, `
		DROP TABLE IF EXISTS user_oauth_accounts CASCADE;
		DROP TABLE IF EXISTS reveal_snapshots CASCADE;
		DROP TABLE IF EXISTS rate_limits CASCADE;
		DROP TABLE IF EXISTS votes CASCADE;
		DROP TABLE IF EXISTS contributions CASCADE;
		DROP TABLE IF EXISTS ideas CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
	`)

	// Run migration.
	_, err = pool.Exec(ctx, migrationSQL)
	if err != nil {
		pool.Close()
		t.Fatalf("failed to run migration: %v", err)
	}

	return pool
}

// CleanupDB truncates all tables in dependency-safe order.
func CleanupDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE
			user_oauth_accounts,
			reveal_snapshots,
			rate_limits,
			votes,
			contributions,
			ideas,
			users
		CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

// CreateTestStore creates a Store instance for testing.
func CreateTestStore(t *testing.T, pool *pgxpool.Pool) *store.Store {
	t.Helper()
	return store.New(pool)
}

// CreateTestService creates a Service instance configured for testing.
func CreateTestService(t *testing.T, pool *pgxpool.Pool) *service.Service {
	t.Helper()
	cfg := &config.Config{
		DatabaseURL: testDatabaseURL(),
		JWTSecret:   TestJWTSecret,
		Port:        "8080",
		FrontendURL: "http://localhost:3000",
	}
	s := store.New(pool)
	return service.New(s, cfg)
}

// CreateTestUser inserts a user and returns the user's internal ID.
func CreateTestUser(t *testing.T, pool *pgxpool.Pool, openclawID, username string) int64 {
	t.Helper()
	ctx := context.Background()
	var id int64
	err := pool.QueryRow(ctx,
		`INSERT INTO users (openclaw_id, username) VALUES ($1, $2) RETURNING id`,
		openclawID, username,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test user %s: %v", username, err)
	}
	return id
}

// CreateTestIdea inserts an open idea with a 7-day deadline and returns it.
func CreateTestIdea(t *testing.T, pool *pgxpool.Pool, initiatorID int64, title string) *model.Idea {
	t.Helper()
	ctx := context.Background()
	s := store.New(pool)
	idea, err := s.CreateIdea(ctx, &model.Idea{
		InitiatorID: initiatorID,
		Title:       title,
		Description: "Test description",
		TargetUser:  "developers",
		CoreProblem: "test problem",
		Status:      model.IdeaStatusOpen,
		Deadline:    time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("failed to create test idea %q: %v", title, err)
	}
	return idea
}

// CreateTestIdeaWithDeadline creates an idea with a specific deadline.
func CreateTestIdeaWithDeadline(t *testing.T, pool *pgxpool.Pool, initiatorID int64, title string, deadline time.Time) *model.Idea {
	t.Helper()
	ctx := context.Background()
	s := store.New(pool)
	idea, err := s.CreateIdea(ctx, &model.Idea{
		InitiatorID: initiatorID,
		Title:       title,
		Description: "Test description",
		TargetUser:  "developers",
		CoreProblem: "test problem",
		Status:      model.IdeaStatusOpen,
		Deadline:    deadline,
	})
	if err != nil {
		t.Fatalf("failed to create test idea %q: %v", title, err)
	}
	return idea
}

// CreateTestContribution inserts a draft contribution and returns it.
func CreateTestContribution(t *testing.T, pool *pgxpool.Pool, ideaID, authorID int64, content string) *model.Contribution {
	t.Helper()
	ctx := context.Background()
	s := store.New(pool)
	c, err := s.CreateContribution(ctx, &model.Contribution{
		IdeaID:      ideaID,
		AuthorID:    authorID,
		Content:     content,
		DecisionLog: json.RawMessage("[]"),
	})
	if err != nil {
		t.Fatalf("failed to create test contribution: %v", err)
	}
	return c
}

// SubmitTestContribution creates and submits a contribution in one step.
func SubmitTestContribution(t *testing.T, pool *pgxpool.Pool, ideaID, authorID int64, content string) *model.Contribution {
	t.Helper()
	ctx := context.Background()
	s := store.New(pool)

	c, err := s.CreateContribution(ctx, &model.Contribution{
		IdeaID:      ideaID,
		AuthorID:    authorID,
		Content:     content,
		DecisionLog: json.RawMessage("[]"),
	})
	if err != nil {
		t.Fatalf("failed to create test contribution: %v", err)
	}

	submitted, err := s.SubmitContribution(ctx, c.ID)
	if err != nil {
		t.Fatalf("failed to submit test contribution: %v", err)
	}
	return submitted
}

// MustFormat is a helper that wraps fmt.Sprintf for cleaner test messages.
func MustFormat(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
