package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/claway/server/internal/config"
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
// Kept in sync with migrations/001_init.up.sql.
const migrationSQL = `
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    openclaw_id     TEXT NOT NULL UNIQUE,
    username        TEXT NOT NULL,
    agent_api_key   TEXT,
    credits_balance NUMERIC(12, 4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ideas (
    id                    BIGSERIAL PRIMARY KEY,
    title                 TEXT NOT NULL,
    description           TEXT NOT NULL DEFAULT '',
    target_user_hint      TEXT NOT NULL DEFAULT '',
    problem_definition    TEXT NOT NULL DEFAULT '',
    initiator_id          BIGINT NOT NULL REFERENCES users(id),
    initiator_cut_percent NUMERIC(5, 2) NOT NULL DEFAULT 0,
    package_type          TEXT NOT NULL DEFAULT 'standard' CHECK (package_type IN ('light', 'standard')),
    status                TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'completed', 'cancelled')),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deadline              TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS tasks (
    id                   BIGSERIAL PRIMARY KEY,
    idea_id              BIGINT NOT NULL REFERENCES ideas(id),
    type                 TEXT NOT NULL CHECK (type IN ('D1','D2','D3','D4','D5','D6','D7','D8','D9')),
    title                TEXT NOT NULL,
    description          TEXT NOT NULL DEFAULT '',
    acceptance_criteria  TEXT NOT NULL DEFAULT '',
    dependencies         TEXT NOT NULL DEFAULT '',
    token_limit_hint     INT NOT NULL DEFAULT 0,
    status               TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','claimed','submitted','approved','rejected')),
    claimed_by           BIGINT REFERENCES users(id),
    claimed_at           TIMESTAMPTZ,
    submitted_at         TIMESTAMPTZ,
    approved_at          TIMESTAMPTZ,
    output_content       TEXT,
    output_note          TEXT,
    quality_score        NUMERIC(5, 2),
    reject_reason        TEXT,
    cost_usd_accumulated NUMERIC(12, 6) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS documents (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT NOT NULL UNIQUE REFERENCES tasks(id),
    content         TEXT NOT NULL DEFAULT '',
    current_version INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS document_versions (
    id                 BIGSERIAL PRIMARY KEY,
    document_id        BIGINT NOT NULL REFERENCES documents(id),
    version            INT NOT NULL,
    content            TEXT NOT NULL DEFAULT '',
    diff_from_previous TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by         BIGINT NOT NULL REFERENCES users(id),
    UNIQUE (document_id, version)
);

CREATE TABLE IF NOT EXISTS token_usage_logs (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id),
    task_id    BIGINT NOT NULL REFERENCES tasks(id),
    model      TEXT NOT NULL,
    tokens_in  INT NOT NULL DEFAULT 0,
    tokens_out INT NOT NULL DEFAULT 0,
    cost_usd   NUMERIC(12, 6) NOT NULL DEFAULT 0,
    timestamp  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contributions (
    id             BIGSERIAL PRIMARY KEY,
    idea_id        BIGINT NOT NULL REFERENCES ideas(id),
    task_id        BIGINT NOT NULL REFERENCES tasks(id),
    user_id        BIGINT NOT NULL REFERENCES users(id),
    cost_usd       NUMERIC(12, 6) NOT NULL DEFAULT 0,
    quality_score  NUMERIC(5, 2) NOT NULL DEFAULT 0,
    weighted_score NUMERIC(12, 6) NOT NULL DEFAULT 0,
    weight_percent NUMERIC(7, 4) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS credit_transactions (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id),
    type           TEXT NOT NULL,
    amount         NUMERIC(12, 4) NOT NULL,
    reference_type TEXT NOT NULL DEFAULT '',
    reference_id   BIGINT NOT NULL DEFAULT 0,
    description    TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS prds (
    id            BIGSERIAL PRIMARY KEY,
    idea_id       BIGINT NOT NULL UNIQUE REFERENCES ideas(id),
    content       TEXT NOT NULL DEFAULT '',
    published_at  TIMESTAMPTZ,
    price_credits NUMERIC(12, 4) NOT NULL DEFAULT 0,
    read_count    INT NOT NULL DEFAULT 0
);
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

	// Truncate in reverse dependency order using CASCADE.
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE
			credit_transactions,
			contributions,
			token_usage_logs,
			document_versions,
			documents,
			prds,
			tasks,
			ideas,
			users
		CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

// CreateTestService creates a Service instance configured for testing.
func CreateTestService(t *testing.T, pool *pgxpool.Pool) *service.Service {
	t.Helper()
	cfg := &config.Config{
		DatabaseURL: testDatabaseURL(),
		JWTSecret:   TestJWTSecret,
		Port:        "8080",
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

// SimulateTokenUsage inserts a token_usage_log and updates task cost_usd_accumulated.
// This simulates LLM proxy usage during task work.
func SimulateTokenUsage(t *testing.T, pool *pgxpool.Pool, userID, taskID int64, costUSD float64) {
	t.Helper()
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO token_usage_logs (user_id, task_id, model, tokens_in, tokens_out, cost_usd)
		 VALUES ($1, $2, 'gpt-4o', 1000, 500, $3)`,
		userID, taskID, costUSD,
	)
	if err != nil {
		t.Fatalf("failed to insert token_usage_log: %v", err)
	}

	_, err = pool.Exec(ctx,
		`UPDATE tasks SET cost_usd_accumulated = cost_usd_accumulated + $1 WHERE id = $2`,
		costUSD, taskID,
	)
	if err != nil {
		t.Fatalf("failed to update task cost: %v", err)
	}
}

// GiveCredits directly updates a user's credits balance via SQL.
func GiveCredits(t *testing.T, pool *pgxpool.Pool, userID int64, amount float64) {
	t.Helper()
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`UPDATE users SET credits_balance = credits_balance + $1, updated_at = NOW() WHERE id = $2`,
		amount, userID,
	)
	if err != nil {
		t.Fatalf("failed to give credits to user %d: %v", userID, err)
	}
}

// GetUserBalance reads the current credits_balance for a user.
func GetUserBalance(t *testing.T, pool *pgxpool.Pool, userID int64) float64 {
	t.Helper()
	ctx := context.Background()

	var balance float64
	err := pool.QueryRow(ctx, `SELECT credits_balance FROM users WHERE id = $1`, userID).Scan(&balance)
	if err != nil {
		t.Fatalf("failed to get user balance for user %d: %v", userID, err)
	}
	return balance
}

// TaskIDByType finds a task ID by its type within an idea.
func TaskIDByType(t *testing.T, pool *pgxpool.Pool, ideaID int64, taskType string) int64 {
	t.Helper()
	ctx := context.Background()

	var id int64
	err := pool.QueryRow(ctx,
		`SELECT id FROM tasks WHERE idea_id = $1 AND type = $2`, ideaID, taskType,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to find task %s for idea %d: %v", taskType, ideaID, err)
	}
	return id
}

// MustFormat is a helper that wraps fmt.Sprintf for cleaner test messages.
func MustFormat(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
