package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
	"github.com/claway/server/internal/testutil"
)

func TestCreateIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	idea, err := s.CreateIdea(ctx, &model.Idea{
		InitiatorID: userID,
		Title:       "Test Idea",
		Description: "A test idea description",
		TargetUser:  "developers",
		CoreProblem: "testing is hard",
		Status:      model.IdeaStatusOpen,
		Deadline:    time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("CreateIdea failed: %v", err)
	}

	if idea.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if idea.Title != "Test Idea" {
		t.Errorf("expected title 'Test Idea', got %q", idea.Title)
	}
	if idea.Status != model.IdeaStatusOpen {
		t.Errorf("expected status 'open', got %q", idea.Status)
	}
	if idea.InitiatorID != userID {
		t.Errorf("expected initiator_id %d, got %d", userID, idea.InitiatorID)
	}
	if idea.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
}

func TestGetIdeaByID(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	idea := testutil.CreateTestIdea(t, pool, userID, "Find Me")

	got, err := s.GetIdeaByID(ctx, idea.ID)
	if err != nil {
		t.Fatalf("GetIdeaByID failed: %v", err)
	}
	if got.Title != "Find Me" {
		t.Errorf("expected title 'Find Me', got %q", got.Title)
	}
}

func TestGetIdeaByID_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	_, err := s.GetIdeaByID(ctx, 99999)
	if err == nil {
		t.Fatal("expected error for non-existent idea")
	}
}

func TestListIdeas(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	testutil.CreateTestIdea(t, pool, userID, "Idea 1")
	testutil.CreateTestIdea(t, pool, userID, "Idea 2")

	// List all
	ideas, total, err := s.ListIdeas(ctx, "", 10, 0)
	if err != nil {
		t.Fatalf("ListIdeas failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(ideas) != 2 {
		t.Errorf("expected 2 ideas, got %d", len(ideas))
	}

	// List by status
	ideas, total, err = s.ListIdeas(ctx, "open", 10, 0)
	if err != nil {
		t.Fatalf("ListIdeas with status filter failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 open ideas, got %d", total)
	}

	// List with pagination
	ideas, _, err = s.ListIdeas(ctx, "", 1, 0)
	if err != nil {
		t.Fatalf("ListIdeas with limit failed: %v", err)
	}
	if len(ideas) != 1 {
		t.Errorf("expected 1 idea with limit=1, got %d", len(ideas))
	}
}

func TestListExpiredOpenIdeas(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	// Create one expired idea and one future idea
	testutil.CreateTestIdeaWithDeadline(t, pool, userID, "Expired", time.Now().Add(-1*time.Hour))
	testutil.CreateTestIdeaWithDeadline(t, pool, userID, "Future", time.Now().Add(7*24*time.Hour))

	ideas, err := s.ListExpiredOpenIdeas(ctx)
	if err != nil {
		t.Fatalf("ListExpiredOpenIdeas failed: %v", err)
	}
	if len(ideas) != 1 {
		t.Fatalf("expected 1 expired idea, got %d", len(ideas))
	}
	if ideas[0].Title != "Expired" {
		t.Errorf("expected expired idea title 'Expired', got %q", ideas[0].Title)
	}
}

func TestCloseIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	idea := testutil.CreateTestIdea(t, pool, userID, "To Close")

	err := s.CloseIdea(ctx, idea.ID)
	if err != nil {
		t.Fatalf("CloseIdea failed: %v", err)
	}

	got, err := s.GetIdeaByID(ctx, idea.ID)
	if err != nil {
		t.Fatalf("GetIdeaByID after close failed: %v", err)
	}
	if got.Status != model.IdeaStatusClosed {
		t.Errorf("expected status 'closed', got %q", got.Status)
	}
	if !got.RevealedAt.Valid {
		t.Error("expected revealed_at to be set")
	}
}

func TestCloseIdea_AlreadyClosed(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	idea := testutil.CreateTestIdea(t, pool, userID, "Already Closed")
	_ = s.CloseIdea(ctx, idea.ID)

	// Second close should return ErrNotFound (no rows affected)
	err := s.CloseIdea(ctx, idea.ID)
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound for already-closed idea, got %v", err)
	}
}

func TestUpdateIdeaStatus(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")

	idea := testutil.CreateTestIdea(t, pool, userID, "To Cancel")

	err := s.UpdateIdeaStatus(ctx, idea.ID, "cancelled")
	if err != nil {
		t.Fatalf("UpdateIdeaStatus failed: %v", err)
	}

	got, _ := s.GetIdeaByID(ctx, idea.ID)
	if got.Status != model.IdeaStatusCancelled {
		t.Errorf("expected status 'cancelled', got %q", got.Status)
	}
}

func TestUpdateIdeaStatus_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	err := s.UpdateIdeaStatus(ctx, 99999, "cancelled")
	if err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
