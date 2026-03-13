package store_test

import (
	"context"
	"testing"

	"github.com/claway/server/internal/testutil"
)

func TestIncrementRateLimit(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	// First increment should return 1
	count, err := s.IncrementRateLimit(ctx, userID, "vote")
	if err != nil {
		t.Fatalf("IncrementRateLimit failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Second increment should return 2
	count, err = s.IncrementRateLimit(ctx, userID, "vote")
	if err != nil {
		t.Fatalf("IncrementRateLimit (2nd) failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestIncrementRateLimit_DifferentActions(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	// Different actions should have independent counters
	s.IncrementRateLimit(ctx, userID, "vote")
	s.IncrementRateLimit(ctx, userID, "vote")
	s.IncrementRateLimit(ctx, userID, "post_idea")

	voteCount, _ := s.GetRateLimitCount(ctx, userID, "vote")
	ideaCount, _ := s.GetRateLimitCount(ctx, userID, "post_idea")

	if voteCount != 2 {
		t.Errorf("expected vote count 2, got %d", voteCount)
	}
	if ideaCount != 1 {
		t.Errorf("expected post_idea count 1, got %d", ideaCount)
	}
}

func TestIncrementRateLimit_DifferentUsers(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	alice := testutil.CreateTestUser(t, pool, "oc1", "alice")
	bob := testutil.CreateTestUser(t, pool, "oc2", "bob")

	s.IncrementRateLimit(ctx, alice, "vote")
	s.IncrementRateLimit(ctx, alice, "vote")
	s.IncrementRateLimit(ctx, bob, "vote")

	aliceCount, _ := s.GetRateLimitCount(ctx, alice, "vote")
	bobCount, _ := s.GetRateLimitCount(ctx, bob, "vote")

	if aliceCount != 2 {
		t.Errorf("expected alice vote count 2, got %d", aliceCount)
	}
	if bobCount != 1 {
		t.Errorf("expected bob vote count 1, got %d", bobCount)
	}
}

func TestGetRateLimitCount_NoRecord(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	// No record should return 0
	count, err := s.GetRateLimitCount(ctx, userID, "vote")
	if err != nil {
		t.Fatalf("GetRateLimitCount failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

func TestIncrementRateLimit_AtomicIncrement(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	// Simulate hitting the daily limit (e.g., 10 votes)
	for i := 1; i <= 11; i++ {
		count, err := s.IncrementRateLimit(ctx, userID, "vote")
		if err != nil {
			t.Fatalf("IncrementRateLimit iteration %d failed: %v", i, err)
		}
		if count != i {
			t.Errorf("iteration %d: expected count %d, got %d", i, i, count)
		}
	}

	// Verify final count
	count, _ := s.GetRateLimitCount(ctx, userID, "vote")
	if count != 11 {
		t.Errorf("expected final count 11, got %d", count)
	}
}
