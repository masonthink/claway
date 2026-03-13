package service_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/service"
	"github.com/claway/server/internal/testutil"
)

func TestCreateContribution_Success(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Contrib Test")

	resp, err := svc.CreateContribution(ctx, author, idea.ID, service.CreateContributionRequest{
		Content: "My solution document",
	})
	if err != nil {
		t.Fatalf("CreateContribution failed: %v", err)
	}
	if resp.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if resp.Status != model.ContributionStatusDraft {
		t.Errorf("expected draft status, got %q", resp.Status)
	}
	if resp.Content != "My solution document" {
		t.Errorf("expected content 'My solution document', got %q", resp.Content)
	}
}

func TestCreateContribution_DuplicateReturnsDraft(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Dup Contrib")

	first, err := svc.CreateContribution(ctx, author, idea.ID, service.CreateContributionRequest{
		Content: "First draft",
	})
	if err != nil {
		t.Fatalf("first CreateContribution failed: %v", err)
	}

	// Second create should return existing draft (not error)
	second, err := svc.CreateContribution(ctx, author, idea.ID, service.CreateContributionRequest{
		Content: "Should not overwrite",
	})
	if err != nil {
		t.Fatalf("second CreateContribution failed: %v", err)
	}
	if second.ID != first.ID {
		t.Errorf("expected same ID %d, got %d", first.ID, second.ID)
	}
}

func TestCreateContribution_SubmittedDuplicate(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Submitted Dup")

	// Create and submit
	testutil.SubmitTestContribution(t, pool, idea.ID, author, "Already submitted")

	// Try to create again after submission
	_, err := svc.CreateContribution(ctx, author, idea.ID, service.CreateContributionRequest{
		Content: "Try again",
	})
	if err == nil {
		t.Fatal("expected error for duplicate after submission")
	}
	if !strings.Contains(err.Error(), "already submitted") {
		t.Errorf("expected 'already submitted' error, got %q", err.Error())
	}
}

func TestCreateContribution_IdeaClosed(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Closed Idea")

	s.CloseIdea(ctx, idea.ID)

	_, err := svc.CreateContribution(ctx, author, idea.ID, service.CreateContributionRequest{
		Content: "Too late",
	})
	if err == nil {
		t.Fatal("expected error for closed idea")
	}
	if !strings.Contains(err.Error(), "not open") {
		t.Errorf("expected 'not open' error, got %q", err.Error())
	}
}

func TestCreateContribution_DefaultDecisionLog(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Default Log")

	resp, err := svc.CreateContribution(ctx, author, idea.ID, service.CreateContributionRequest{
		Content: "No decision log",
		// DecisionLog is nil
	})
	if err != nil {
		t.Fatalf("CreateContribution failed: %v", err)
	}
	// Verify decision log defaults to empty array
	if resp.DecisionLog == nil {
		t.Error("expected non-nil decision_log default")
	}
}

func TestUpdateContribution_Success(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Update Test")

	c := testutil.CreateTestContribution(t, pool, idea.ID, author, "Original")

	resp, err := svc.UpdateContribution(ctx, author, c.ID, service.UpdateContributionRequest{
		Content:     "Updated content",
		DecisionLog: json.RawMessage(`[{"key":"val"}]`),
	})
	if err != nil {
		t.Fatalf("UpdateContribution failed: %v", err)
	}
	if resp.Content != "Updated content" {
		t.Errorf("expected updated content, got %q", resp.Content)
	}
}

func TestUpdateContribution_NotAuthor(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	other := testutil.CreateTestUser(t, pool, "oc2", "other")
	initiator := testutil.CreateTestUser(t, pool, "oc3", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Auth Test")

	c := testutil.CreateTestContribution(t, pool, idea.ID, author, "Mine")

	_, err := svc.UpdateContribution(ctx, other, c.ID, service.UpdateContributionRequest{
		Content: "Hijack",
	})
	if err == nil {
		t.Fatal("expected error for unauthorized update")
	}
	if !strings.Contains(err.Error(), "not authorized") {
		t.Errorf("expected 'not authorized' error, got %q", err.Error())
	}
}

func TestUpdateContribution_AlreadySubmitted(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Locked Test")

	c := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Submitted")

	_, err := svc.UpdateContribution(ctx, author, c.ID, service.UpdateContributionRequest{
		Content: "Try update",
	})
	if err == nil {
		t.Fatal("expected error for updating submitted contribution")
	}
	if !strings.Contains(err.Error(), "draft") {
		t.Errorf("expected 'draft' error, got %q", err.Error())
	}
}

func TestSubmitContribution_Success(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Submit Test")

	c := testutil.CreateTestContribution(t, pool, idea.ID, author, "Ready to submit")

	resp, err := svc.SubmitContribution(ctx, author, c.ID)
	if err != nil {
		t.Fatalf("SubmitContribution failed: %v", err)
	}
	if resp.Status != model.ContributionStatusSubmitted {
		t.Errorf("expected submitted status, got %q", resp.Status)
	}
	if resp.SubmittedAt == nil {
		t.Error("expected submitted_at to be set")
	}
}

func TestSubmitContribution_EmptyContent(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Empty Submit")

	c := testutil.CreateTestContribution(t, pool, idea.ID, author, "") // empty content

	_, err := svc.SubmitContribution(ctx, author, c.ID)
	if err == nil {
		t.Fatal("expected error for submitting empty contribution")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("expected 'empty' error, got %q", err.Error())
	}
}

func TestSubmitContribution_NotAuthor(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	other := testutil.CreateTestUser(t, pool, "oc2", "other")
	initiator := testutil.CreateTestUser(t, pool, "oc3", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Auth Submit")

	c := testutil.CreateTestContribution(t, pool, idea.ID, author, "Content")

	_, err := svc.SubmitContribution(ctx, other, c.ID)
	if err == nil {
		t.Fatal("expected error for unauthorized submit")
	}
	if !strings.Contains(err.Error(), "not authorized") {
		t.Errorf("expected 'not authorized' error, got %q", err.Error())
	}
}

func TestSubmitContribution_AlreadySubmitted(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Double Submit")

	c := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Submitted")

	_, err := svc.SubmitContribution(ctx, author, c.ID)
	if err == nil {
		t.Fatal("expected error for double submit")
	}
	if !strings.Contains(err.Error(), "already submitted") {
		t.Errorf("expected 'already submitted' error, got %q", err.Error())
	}
}

func TestSubmitContribution_IdeaClosed(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	initiator := testutil.CreateTestUser(t, pool, "oc2", "initiator")
	idea := testutil.CreateTestIdea(t, pool, initiator, "Closed Submit")

	c := testutil.CreateTestContribution(t, pool, idea.ID, author, "Content")

	// Close the idea
	s.CloseIdea(ctx, idea.ID)

	_, err := svc.SubmitContribution(ctx, author, c.ID)
	if err == nil {
		t.Fatal("expected error for submitting after idea closed")
	}
	if !strings.Contains(err.Error(), "no longer open") {
		t.Errorf("expected 'no longer open' error, got %q", err.Error())
	}
}
