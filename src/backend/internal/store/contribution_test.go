package store_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/testutil"
)

func TestCreateContribution(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	c, err := s.CreateContribution(ctx, &model.Contribution{
		IdeaID:      idea.ID,
		AuthorID:    userID,
		Content:     "My solution",
		DecisionLog: json.RawMessage("[]"),
	})
	if err != nil {
		t.Fatalf("CreateContribution failed: %v", err)
	}

	if c.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if c.Status != model.ContributionStatusDraft {
		t.Errorf("expected status 'draft', got %q", c.Status)
	}
	if c.Content != "My solution" {
		t.Errorf("expected content 'My solution', got %q", c.Content)
	}
	if c.AuthorID != userID {
		t.Errorf("expected author_id %d, got %d", userID, c.AuthorID)
	}
}

func TestCreateContribution_UniqueConstraint(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	_, err := s.CreateContribution(ctx, &model.Contribution{
		IdeaID:      idea.ID,
		AuthorID:    userID,
		Content:     "First",
		DecisionLog: json.RawMessage("[]"),
	})
	if err != nil {
		t.Fatalf("first CreateContribution failed: %v", err)
	}

	// Same user + same idea should fail due to UNIQUE(idea_id, author_id)
	_, err = s.CreateContribution(ctx, &model.Contribution{
		IdeaID:      idea.ID,
		AuthorID:    userID,
		Content:     "Duplicate",
		DecisionLog: json.RawMessage("[]"),
	})
	if err == nil {
		t.Fatal("expected error for duplicate contribution, got nil")
	}
}

func TestGetContributionByID(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	created := testutil.CreateTestContribution(t, pool, idea.ID, userID, "My content")

	got, err := s.GetContributionByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetContributionByID failed: %v", err)
	}
	if got.Content != "My content" {
		t.Errorf("expected content 'My content', got %q", got.Content)
	}
}

func TestGetContributionByID_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	_, err := s.GetContributionByID(ctx, 99999)
	if err == nil {
		t.Fatal("expected error for non-existent contribution")
	}
}

func TestGetContributionByIdeaAndAuthor(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	testutil.CreateTestContribution(t, pool, idea.ID, userID, "Found me")

	got, err := s.GetContributionByIdeaAndAuthor(ctx, idea.ID, userID)
	if err != nil {
		t.Fatalf("GetContributionByIdeaAndAuthor failed: %v", err)
	}
	if got.Content != "Found me" {
		t.Errorf("expected content 'Found me', got %q", got.Content)
	}
}

func TestGetContributionByIdeaAndAuthor_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	_, err := s.GetContributionByIdeaAndAuthor(ctx, idea.ID, 99999)
	if err == nil {
		t.Fatal("expected error for non-existent author")
	}
}

func TestUpdateContributionDraft(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	c := testutil.CreateTestContribution(t, pool, idea.ID, userID, "Original")

	newLog := json.RawMessage(`[{"decision": "use Go"}]`)
	updated, err := s.UpdateContributionDraft(ctx, c.ID, "Updated content", newLog)
	if err != nil {
		t.Fatalf("UpdateContributionDraft failed: %v", err)
	}
	if updated.Content != "Updated content" {
		t.Errorf("expected content 'Updated content', got %q", updated.Content)
	}
}

func TestUpdateContributionDraft_SubmittedFails(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	c := testutil.SubmitTestContribution(t, pool, idea.ID, userID, "Submitted content")

	// Updating a submitted contribution should fail (no rows match WHERE status='draft')
	_, err := s.UpdateContributionDraft(ctx, c.ID, "Should fail", nil)
	if err == nil {
		t.Fatal("expected error when updating submitted contribution")
	}
}

func TestSubmitContribution(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	c := testutil.CreateTestContribution(t, pool, idea.ID, userID, "Draft content")

	submitted, err := s.SubmitContribution(ctx, c.ID)
	if err != nil {
		t.Fatalf("SubmitContribution failed: %v", err)
	}
	if submitted.Status != model.ContributionStatusSubmitted {
		t.Errorf("expected status 'submitted', got %q", submitted.Status)
	}
	if !submitted.SubmittedAt.Valid {
		t.Error("expected submitted_at to be set")
	}
}

func TestSubmitContribution_AlreadySubmitted(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	c := testutil.SubmitTestContribution(t, pool, idea.ID, userID, "Already submitted")

	// Submitting again should fail
	_, err := s.SubmitContribution(ctx, c.ID)
	if err == nil {
		t.Fatal("expected error when submitting already-submitted contribution")
	}
}

func TestListContributionsByIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	alice := testutil.CreateTestUser(t, pool, "oc1", "alice")
	bob := testutil.CreateTestUser(t, pool, "oc2", "bob")
	idea := testutil.CreateTestIdea(t, pool, alice, "Test Idea")

	// Create one draft and one submitted
	testutil.CreateTestContribution(t, pool, idea.ID, alice, "Draft by alice")
	testutil.SubmitTestContribution(t, pool, idea.ID, bob, "Submitted by bob")

	// ListContributionsByIdea only returns submitted contributions
	contribs, err := s.ListContributionsByIdea(ctx, idea.ID, "open")
	if err != nil {
		t.Fatalf("ListContributionsByIdea failed: %v", err)
	}
	if len(contribs) != 1 {
		t.Fatalf("expected 1 submitted contribution, got %d", len(contribs))
	}
	if contribs[0].AuthorID != bob {
		t.Errorf("expected author_id %d, got %d", bob, contribs[0].AuthorID)
	}
}

func TestListContributionsByAuthor(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	alice := testutil.CreateTestUser(t, pool, "oc1", "alice")
	bob := testutil.CreateTestUser(t, pool, "oc2", "bob")
	idea1 := testutil.CreateTestIdea(t, pool, alice, "Idea 1")
	idea2 := testutil.CreateTestIdea(t, pool, alice, "Idea 2")

	testutil.CreateTestContribution(t, pool, idea1.ID, bob, "Bob's first")
	testutil.CreateTestContribution(t, pool, idea2.ID, bob, "Bob's second")

	contribs, total, err := s.ListContributionsByAuthor(ctx, bob, 10, 0)
	if err != nil {
		t.Fatalf("ListContributionsByAuthor failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(contribs) != 2 {
		t.Errorf("expected 2 contributions, got %d", len(contribs))
	}
}

func TestCountContributionsByIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	alice := testutil.CreateTestUser(t, pool, "oc1", "alice")
	bob := testutil.CreateTestUser(t, pool, "oc2", "bob")
	idea := testutil.CreateTestIdea(t, pool, alice, "Test Idea")

	// Draft should not be counted
	testutil.CreateTestContribution(t, pool, idea.ID, alice, "Draft")
	testutil.SubmitTestContribution(t, pool, idea.ID, bob, "Submitted")

	count, err := s.CountContributionsByIdea(ctx, idea.ID)
	if err != nil {
		t.Fatalf("CountContributionsByIdea failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 submitted contribution, got %d", count)
	}
}

func TestIncrementContributionViewCount(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()
	userID := testutil.CreateTestUser(t, pool, "oc1", "alice")
	idea := testutil.CreateTestIdea(t, pool, userID, "Test Idea")

	c := testutil.CreateTestContribution(t, pool, idea.ID, userID, "Content")

	for i := 0; i < 3; i++ {
		if err := s.IncrementContributionViewCount(ctx, c.ID); err != nil {
			t.Fatalf("IncrementContributionViewCount failed: %v", err)
		}
	}

	got, _ := s.GetContributionByID(ctx, c.ID)
	if got.ViewCount != 3 {
		t.Errorf("expected view_count 3, got %d", got.ViewCount)
	}
}
