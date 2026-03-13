package store_test

import (
	"context"
	"testing"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
	"github.com/claway/server/internal/testutil"
)

func TestCreateVote(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Vote Idea")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Contribution")

	vote, err := s.CreateVote(ctx, &model.Vote{
		IdeaID:         idea.ID,
		VoterID:        voter,
		ContributionID: contrib.ID,
	})
	if err != nil {
		t.Fatalf("CreateVote failed: %v", err)
	}
	if vote.ID == 0 {
		t.Error("expected non-zero vote ID")
	}
	if vote.VoterID != voter {
		t.Errorf("expected voter_id %d, got %d", voter, vote.VoterID)
	}
	if vote.VotedAt.IsZero() {
		t.Error("expected non-zero voted_at")
	}
}

func TestCreateVote_UniqueConstraint(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Vote Idea")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Contribution")

	_, err := s.CreateVote(ctx, &model.Vote{
		IdeaID:         idea.ID,
		VoterID:        voter,
		ContributionID: contrib.ID,
	})
	if err != nil {
		t.Fatalf("first vote failed: %v", err)
	}

	// Same voter + same idea should return ErrConflict
	_, err = s.CreateVote(ctx, &model.Vote{
		IdeaID:         idea.ID,
		VoterID:        voter,
		ContributionID: contrib.ID,
	})
	if err != store.ErrConflict {
		t.Errorf("expected ErrConflict for duplicate vote, got %v", err)
	}
}

func TestCountVotesByIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	v1 := testutil.CreateTestUser(t, pool, "oc2", "voter1")
	v2 := testutil.CreateTestUser(t, pool, "oc3", "voter2")
	idea := testutil.CreateTestIdea(t, pool, author, "Count Votes")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Contribution")

	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: v1, ContributionID: contrib.ID})
	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: v2, ContributionID: contrib.ID})

	count, err := s.CountVotesByIdea(ctx, idea.ID)
	if err != nil {
		t.Fatalf("CountVotesByIdea failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 votes, got %d", count)
	}
}

func TestCountVotersByIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "Voters Count")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Contribution")

	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: voter, ContributionID: contrib.ID})

	count, err := s.CountVotersByIdea(ctx, idea.ID)
	if err != nil {
		t.Fatalf("CountVotersByIdea failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 distinct voter, got %d", count)
	}
}

func TestHasUserVotedOnIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	bystander := testutil.CreateTestUser(t, pool, "oc3", "bystander")
	idea := testutil.CreateTestIdea(t, pool, author, "Check Vote")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Contribution")

	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: voter, ContributionID: contrib.ID})

	voted, err := s.HasUserVotedOnIdea(ctx, idea.ID, voter)
	if err != nil {
		t.Fatalf("HasUserVotedOnIdea failed: %v", err)
	}
	if !voted {
		t.Error("expected voter to have voted")
	}

	voted, err = s.HasUserVotedOnIdea(ctx, idea.ID, bystander)
	if err != nil {
		t.Fatalf("HasUserVotedOnIdea for bystander failed: %v", err)
	}
	if voted {
		t.Error("expected bystander to not have voted")
	}
}

func TestGetVoteCountsByIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author1 := testutil.CreateTestUser(t, pool, "oc1", "author1")
	author2 := testutil.CreateTestUser(t, pool, "oc2", "author2")
	v1 := testutil.CreateTestUser(t, pool, "oc3", "voter1")
	v2 := testutil.CreateTestUser(t, pool, "oc4", "voter2")
	v3 := testutil.CreateTestUser(t, pool, "oc5", "voter3")

	idea := testutil.CreateTestIdea(t, pool, author1, "Ranking")
	c1 := testutil.SubmitTestContribution(t, pool, idea.ID, author1, "Contrib 1")
	c2 := testutil.SubmitTestContribution(t, pool, idea.ID, author2, "Contrib 2")

	// c1 gets 2 votes, c2 gets 1 vote
	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: v1, ContributionID: c1.ID})
	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: v2, ContributionID: c1.ID})
	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: v3, ContributionID: c2.ID})

	counts, totalVotes, err := s.GetVoteCountsByIdea(ctx, idea.ID)
	if err != nil {
		t.Fatalf("GetVoteCountsByIdea failed: %v", err)
	}
	if totalVotes != 3 {
		t.Errorf("expected 3 total votes, got %d", totalVotes)
	}
	if len(counts) != 2 {
		t.Fatalf("expected 2 contributions, got %d", len(counts))
	}
	// First should be c1 with 2 votes (ordered DESC)
	if counts[0].ContributionID != c1.ID || counts[0].VoteCount != 2 {
		t.Errorf("expected first entry: contrib=%d votes=2, got contrib=%d votes=%d",
			c1.ID, counts[0].ContributionID, counts[0].VoteCount)
	}
	if counts[1].ContributionID != c2.ID || counts[1].VoteCount != 1 {
		t.Errorf("expected second entry: contrib=%d votes=1, got contrib=%d votes=%d",
			c2.ID, counts[1].ContributionID, counts[1].VoteCount)
	}
}

func TestGetVoteCountsByIdea_NoVotes(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	idea := testutil.CreateTestIdea(t, pool, author, "No Votes")
	testutil.SubmitTestContribution(t, pool, idea.ID, author, "Lonely contrib")

	counts, totalVotes, err := s.GetVoteCountsByIdea(ctx, idea.ID)
	if err != nil {
		t.Fatalf("GetVoteCountsByIdea failed: %v", err)
	}
	if totalVotes != 0 {
		t.Errorf("expected 0 total votes, got %d", totalVotes)
	}
	if len(counts) != 1 {
		t.Fatalf("expected 1 contribution entry, got %d", len(counts))
	}
	if counts[0].VoteCount != 0 {
		t.Errorf("expected 0 votes for contribution, got %d", counts[0].VoteCount)
	}
}

func TestGetUserVoteForIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")
	idea := testutil.CreateTestIdea(t, pool, author, "User Vote")
	contrib := testutil.SubmitTestContribution(t, pool, idea.ID, author, "Contribution")

	s.CreateVote(ctx, &model.Vote{IdeaID: idea.ID, VoterID: voter, ContributionID: contrib.ID})

	vote, err := s.GetUserVoteForIdea(ctx, idea.ID, voter)
	if err != nil {
		t.Fatalf("GetUserVoteForIdea failed: %v", err)
	}
	if vote.ContributionID != contrib.ID {
		t.Errorf("expected contribution_id %d, got %d", contrib.ID, vote.ContributionID)
	}
}

func TestGetUserVoteForIdea_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	idea := testutil.CreateTestIdea(t, pool, author, "No Vote")

	_, err := s.GetUserVoteForIdea(ctx, idea.ID, author)
	if err == nil {
		t.Fatal("expected error for non-existent vote")
	}
}

func TestListVotesByUser(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	s := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	author := testutil.CreateTestUser(t, pool, "oc1", "author")
	voter := testutil.CreateTestUser(t, pool, "oc2", "voter")

	idea1 := testutil.CreateTestIdea(t, pool, author, "Idea 1")
	idea2 := testutil.CreateTestIdea(t, pool, author, "Idea 2")
	c1 := testutil.SubmitTestContribution(t, pool, idea1.ID, author, "C1")
	c2 := testutil.SubmitTestContribution(t, pool, idea2.ID, author, "C2")

	s.CreateVote(ctx, &model.Vote{IdeaID: idea1.ID, VoterID: voter, ContributionID: c1.ID})
	s.CreateVote(ctx, &model.Vote{IdeaID: idea2.ID, VoterID: voter, ContributionID: c2.ID})

	votes, total, err := s.ListVotesByUser(ctx, voter, 10, 0)
	if err != nil {
		t.Fatalf("ListVotesByUser failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(votes) != 2 {
		t.Errorf("expected 2 votes, got %d", len(votes))
	}
}
