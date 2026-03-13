package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/service"
	"github.com/claway/server/internal/testutil"
)

func TestCastVote_Success(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Vote Test")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Good solution")

	vote, err := svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{
		ContributionID: contrib.ID,
	})
	if err != nil {
		t.Fatalf("CastVote failed: %v", err)
	}
	if vote.ID == 0 {
		t.Error("expected non-zero vote ID")
	}
	if vote.VoterID != voter {
		t.Errorf("expected voter_id %d, got %d", voter, vote.VoterID)
	}
}

func TestCastVote_SelfVoteRejected(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	idea := testutil.CreateTestIdea(t, pool, author, "Self Vote")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "My solution")

	// Author tries to vote for their own contribution
	_, err := svc.CastVote(ctx, author, idea.ID, service.CastVoteRequest{
		ContributionID: contrib.ID,
	})
	if err == nil {
		t.Fatal("expected error for self-vote")
	}
	if !strings.Contains(err.Error(), "own contribution") {
		t.Errorf("expected error about own contribution, got %q", err.Error())
	}
}

func TestCastVote_DuplicateVoteRejected(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Dup Vote")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Solution")

	_, err := svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{
		ContributionID: contrib.ID,
	})
	if err != nil {
		t.Fatalf("first vote failed: %v", err)
	}

	// Second vote on same idea should be rejected
	_, err = svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{
		ContributionID: contrib.ID,
	})
	if err == nil {
		t.Fatal("expected error for duplicate vote")
	}
	if !strings.Contains(err.Error(), "already voted") {
		t.Errorf("expected 'already voted' error, got %q", err.Error())
	}
}

func TestCastVote_IdeaClosedRejected(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Closed Idea")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Solution")

	// Close the idea
	if err := s.CloseIdea(ctx, idea.ID); err != nil {
		t.Fatalf("CloseIdea failed: %v", err)
	}

	_, err := svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{
		ContributionID: contrib.ID,
	})
	if err == nil {
		t.Fatal("expected error for voting on closed idea")
	}
	if !strings.Contains(err.Error(), "ended") {
		t.Errorf("expected error about voting ended, got %q", err.Error())
	}
}

func TestCastVote_DraftContributionRejected(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Draft Vote")
	// Create draft only (not submitted)
	draft := testutil.CreateTestContribution(t, pool, idea.ID, author, "Draft")

	_, err := svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{
		ContributionID: draft.ID,
	})
	if err == nil {
		t.Fatal("expected error for voting on draft contribution")
	}
	if !strings.Contains(err.Error(), "draft") {
		t.Errorf("expected error about draft, got %q", err.Error())
	}
}

func TestCastVote_ContributionWrongIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea1 := testutil.CreateTestIdea(t, pool, author, "Idea 1")
	idea2 := testutil.CreateTestIdea(t, pool, author, "Idea 2")
	contrib := testutil.SubmitTestContribution(t, pool, idea1.ID, author, "For idea 1")

	// Try to vote on idea2 with contribution from idea1
	_, err := svc.CastVote(ctx, voter, idea2.ID, service.CastVoteRequest{
		ContributionID: contrib.ID,
	})
	if err == nil {
		t.Fatal("expected error for mismatched idea/contribution")
	}
	if !strings.Contains(err.Error(), "does not belong") {
		t.Errorf("expected 'does not belong' error, got %q", err.Error())
	}
}

func TestCastVote_MissingContributionID(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	idea := testutil.CreateTestIdea(t, pool, author, "No ContribID")

	_, err := svc.CastVote(ctx, author, idea.ID, service.CastVoteRequest{})
	if err == nil {
		t.Fatal("expected error for missing contribution_id")
	}
	if !strings.Contains(err.Error(), "contribution_id is required") {
		t.Errorf("expected 'contribution_id is required' error, got %q", err.Error())
	}
}

func TestCastVote_DailyLimitReached(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	// Create many authors and ideas so the voter can cast 10+ votes
	voter := testutil.CreateTestUser(t, pool, "voter-oc", "voter")

	for i := 1; i <= 11; i++ {
		author := testutil.CreateTestUser(t, pool,
			testutil.MustFormat("oc-author-%d", i),
			testutil.MustFormat("author%d", i))
		idea := testutil.CreateTestIdea(t, pool, author, testutil.MustFormat("Idea %d", i))
		contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, testutil.MustFormat("Solution %d", i))

		_, err := svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{
			ContributionID: contrib.ID,
		})

		if i <= 10 {
			if err != nil {
				t.Fatalf("vote %d should succeed, got: %v", i, err)
			}
		} else {
			// 11th vote should be rejected by rate limit
			if err == nil {
				t.Fatal("expected error for exceeding daily vote limit")
			}
			if !strings.Contains(err.Error(), "daily vote limit") {
				t.Errorf("expected daily limit error, got %q", err.Error())
			}
		}
	}
}

func TestListMyVotes(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "My Votes")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Solution")

	svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{ContributionID: contrib.ID})

	votes, total, err := svc.ListMyVotes(ctx, voter, 10, 0)
	if err != nil {
		t.Fatalf("ListMyVotes failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(votes) != 1 {
		t.Errorf("expected 1 vote, got %d", len(votes))
	}

	// Verify default limit is applied
	_, _, err = svc.ListMyVotes(ctx, voter, 0, 0)
	if err != nil {
		t.Fatalf("ListMyVotes with default limit failed: %v", err)
	}
}

// Ensure CastVote_IdeaNotFound returns a clear error.
func TestCastVote_IdeaNotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	voter := testutil.CreateTestUser(t, pool, "oc1", "voter")

	_, err := svc.CastVote(ctx, voter, 99999, service.CastVoteRequest{ContributionID: 1})
	if err == nil {
		t.Fatal("expected error for non-existent idea")
	}
	if !strings.Contains(err.Error(), "idea not found") {
		t.Errorf("expected 'idea not found' error, got %q", err.Error())
	}
}

// Verify that a vote for a non-existent contribution is rejected.
func TestCastVote_ContributionNotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "No Contrib")

	// Contribution 99999 doesn't exist
	_, err := svc.CastVote(ctx, voter, idea.ID, service.CastVoteRequest{ContributionID: 99999})
	if err == nil {
		t.Fatal("expected error for non-existent contribution")
	}
	if !strings.Contains(err.Error(), "contribution not found") {
		t.Errorf("expected 'contribution not found' error, got %q", err.Error())
	}
}

// Suppress unused import warning.
var _ = model.IdeaStatusOpen
