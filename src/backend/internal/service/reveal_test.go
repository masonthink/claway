package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/service"
	"github.com/claway/server/internal/store"
	"github.com/claway/server/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

// --- Helpers ---

// createUserLocal inserts a user and returns the ID.
func createUserLocal(t *testing.T, pool *pgxpool.Pool, username string) int64 {
	t.Helper()
	return testutil.CreateTestUser(t, pool, "", username)
}

// createIdeaDirectly inserts an idea with a given deadline, bypassing service rate limits.
func createIdeaDirectly(t *testing.T, pool *pgxpool.Pool, initiatorID int64, deadline time.Time) int64 {
	t.Helper()
	idea := testutil.CreateTestIdeaWithDeadline(t, pool, initiatorID, "Test Idea", deadline)
	return idea.ID
}

// createAndSubmitContribution creates a draft and submits it via the service layer.
func createAndSubmitContribution(t *testing.T, svc *service.Service, authorID, ideaID int64, content string) int64 {
	t.Helper()
	ctx := context.Background()

	resp, err := svc.CreateContribution(ctx, authorID, ideaID, service.CreateContributionRequest{
		Content: content,
	})
	if err != nil {
		t.Fatalf("create contribution: %v", err)
	}

	_, err = svc.SubmitContribution(ctx, authorID, resp.ID)
	if err != nil {
		t.Fatalf("submit contribution: %v", err)
	}

	// Tiny sleep to ensure distinct submitted_at timestamps for ordering tests.
	time.Sleep(5 * time.Millisecond)
	return resp.ID
}

// castVote casts a vote directly via the store, bypassing rate-limit for test convenience.
func castVote(t *testing.T, st *store.Store, voterID, ideaID, contribID int64) {
	t.Helper()
	_, err := st.CreateVote(context.Background(), &model.Vote{
		IdeaID:         ideaID,
		VoterID:        voterID,
		ContributionID: contribID,
	})
	if err != nil {
		t.Fatalf("cast vote (voter=%d, contrib=%d): %v", voterID, contribID, err)
	}
}

// --- Tests ---

func TestReveal_FullLifecycle(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	// Scenario: 1 initiator, 5 contributors, 8 voters.
	initiator := createUserLocal(t, pool, "initiator")
	contributors := make([]int64, 5)
	for i := range contributors {
		contributors[i] = createUserLocal(t, pool, fmt.Sprintf("contributor%d", i))
	}
	voters := make([]int64, 8)
	for i := range voters {
		voters[i] = createUserLocal(t, pool, fmt.Sprintf("voter%d", i))
	}

	deadline := time.Now().Add(-1 * time.Hour)
	ideaID := createIdeaDirectly(t, pool, initiator, deadline)

	contribIDs := make([]int64, 5)
	for i, uid := range contributors {
		contribIDs[i] = createAndSubmitContribution(t, svc, uid, ideaID,
			fmt.Sprintf("My solution #%d with detailed analysis.", i))
	}

	// Vote distribution: contrib[0]=3, contrib[1]=2, contrib[2]=2, contrib[3]=1, contrib[4]=0
	castVote(t, st, voters[0], ideaID, contribIDs[0])
	castVote(t, st, voters[1], ideaID, contribIDs[0])
	castVote(t, st, voters[2], ideaID, contribIDs[0])
	castVote(t, st, voters[3], ideaID, contribIDs[1])
	castVote(t, st, voters[4], ideaID, contribIDs[1])
	castVote(t, st, voters[5], ideaID, contribIDs[2])
	castVote(t, st, voters[6], ideaID, contribIDs[2])
	castVote(t, st, voters[7], ideaID, contribIDs[3])

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	// Verify idea is now closed.
	ideaResp, err := svc.GetIdea(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetIdea after reveal: %v", err)
	}
	if ideaResp.Status != model.IdeaStatusClosed {
		t.Errorf("expected idea status 'closed', got %q", ideaResp.Status)
	}

	// Verify reveal result.
	result, err := svc.GetRevealResult(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealResult: %v", err)
	}
	if result.TotalVotes != 8 {
		t.Errorf("expected total_votes=8, got %d", result.TotalVotes)
	}
	if len(result.Results) != 5 {
		t.Fatalf("expected 5 ranked results, got %d", len(result.Results))
	}

	// Rank 1: contrib[0] (3 votes), featured
	assertRankedEntry(t, result.Results[0], contribIDs[0], 3, 1, true)
	// Rank 2 tie: contrib[1] and contrib[2] both have 2 votes.
	assertRankedEntry(t, result.Results[1], contribIDs[1], 2, 2, true)
	assertRankedEntry(t, result.Results[2], contribIDs[2], 2, 2, true)
	// Rank 4: contrib[3] (1 vote), not featured
	assertRankedEntry(t, result.Results[3], contribIDs[3], 1, 4, false)
	// Rank 5: contrib[4] (0 votes), not featured
	assertRankedEntry(t, result.Results[4], contribIDs[4], 0, 5, false)

	// Author info should be populated.
	if result.Results[0].AuthorUsername == "" {
		t.Error("expected author_username to be populated after reveal")
	}
}

func TestReveal_InsufficientVotes(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	contributors := make([]int64, 3)
	for i := range contributors {
		contributors[i] = createUserLocal(t, pool, fmt.Sprintf("contributor%d", i))
	}
	voters := make([]int64, 4)
	for i := range voters {
		voters[i] = createUserLocal(t, pool, fmt.Sprintf("voter%d", i))
	}

	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(-1*time.Hour))

	contribIDs := make([]int64, 3)
	for i, uid := range contributors {
		contribIDs[i] = createAndSubmitContribution(t, svc, uid, ideaID,
			fmt.Sprintf("Solution %d", i))
	}

	// 4 votes spread across contributions (< 5 threshold).
	castVote(t, st, voters[0], ideaID, contribIDs[0])
	castVote(t, st, voters[1], ideaID, contribIDs[0])
	castVote(t, st, voters[2], ideaID, contribIDs[1])
	castVote(t, st, voters[3], ideaID, contribIDs[2])

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	result, err := svc.GetRevealResult(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealResult: %v", err)
	}

	if result.TotalVotes != 4 {
		t.Errorf("expected total_votes=4, got %d", result.TotalVotes)
	}

	// None should be featured because total votes < 5.
	for i, entry := range result.Results {
		if entry.IsFeatured {
			t.Errorf("result[%d] (contrib=%d) should NOT be featured with only %d total votes",
				i, entry.ContributionID, result.TotalVotes)
		}
	}
}

func TestReveal_NoContributions(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(-1*time.Hour))

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	ideaResp, err := svc.GetIdea(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetIdea: %v", err)
	}
	if ideaResp.Status != model.IdeaStatusClosed {
		t.Errorf("expected status 'closed', got %q", ideaResp.Status)
	}

	result, err := svc.GetRevealResult(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealResult: %v", err)
	}
	if result.TotalVotes != 0 {
		t.Errorf("expected total_votes=0, got %d", result.TotalVotes)
	}
	if len(result.Results) != 0 {
		t.Errorf("expected 0 ranked results, got %d", len(result.Results))
	}
}

func TestReveal_TiedVotes_OrderBySubmissionTime(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	author1 := createUserLocal(t, pool, "author_early")
	author2 := createUserLocal(t, pool, "author_late")
	voters := make([]int64, 6)
	for i := range voters {
		voters[i] = createUserLocal(t, pool, fmt.Sprintf("voter%d", i))
	}

	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(-1*time.Hour))

	contribEarly := createAndSubmitContribution(t, svc, author1, ideaID, "Early solution")
	time.Sleep(10 * time.Millisecond)
	contribLate := createAndSubmitContribution(t, svc, author2, ideaID, "Late solution")

	// Both get 3 votes each (total=6, >= 5 so featured logic applies).
	castVote(t, st, voters[0], ideaID, contribEarly)
	castVote(t, st, voters[1], ideaID, contribEarly)
	castVote(t, st, voters[2], ideaID, contribEarly)
	castVote(t, st, voters[3], ideaID, contribLate)
	castVote(t, st, voters[4], ideaID, contribLate)
	castVote(t, st, voters[5], ideaID, contribLate)

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	result, err := svc.GetRevealResult(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealResult: %v", err)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	// Both should have rank 1 and be featured (tied at 3 votes).
	assertRankedEntry(t, result.Results[0], contribEarly, 3, 1, true)
	assertRankedEntry(t, result.Results[1], contribLate, 3, 1, true)
}

func TestReveal_FewerThanThreeContributions(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	author1 := createUserLocal(t, pool, "author1")
	author2 := createUserLocal(t, pool, "author2")
	voters := make([]int64, 5)
	for i := range voters {
		voters[i] = createUserLocal(t, pool, fmt.Sprintf("voter%d", i))
	}

	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(-1*time.Hour))

	contrib1 := createAndSubmitContribution(t, svc, author1, ideaID, "Solution A")
	contrib2 := createAndSubmitContribution(t, svc, author2, ideaID, "Solution B")

	// 5 total votes: contrib1=3, contrib2=2.
	castVote(t, st, voters[0], ideaID, contrib1)
	castVote(t, st, voters[1], ideaID, contrib1)
	castVote(t, st, voters[2], ideaID, contrib1)
	castVote(t, st, voters[3], ideaID, contrib2)
	castVote(t, st, voters[4], ideaID, contrib2)

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	result, err := svc.GetRevealResult(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealResult: %v", err)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	assertRankedEntry(t, result.Results[0], contrib1, 3, 1, true)
	assertRankedEntry(t, result.Results[1], contrib2, 2, 2, true)
}

func TestReveal_BlindVoting_NoVoteCountPreReveal(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	author1 := createUserLocal(t, pool, "author1")
	author2 := createUserLocal(t, pool, "author2")
	voter := createUserLocal(t, pool, "voter")

	// Deadline in the future (idea stays open).
	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(24*time.Hour))

	createAndSubmitContribution(t, svc, author1, ideaID, "Solution A")
	contrib2 := createAndSubmitContribution(t, svc, author2, ideaID, "Solution B")

	castVote(t, st, voter, ideaID, contrib2)

	// List contributions while idea is open.
	contribs, err := svc.ListContributions(ctx, ideaID)
	if err != nil {
		t.Fatalf("ListContributions: %v", err)
	}

	// Verify that author info is hidden (pre-reveal).
	for _, c := range contribs {
		if c.AuthorID != nil {
			t.Errorf("author_id should be nil pre-reveal, got %d for contrib %d", *c.AuthorID, c.ID)
		}
	}

	// Verify GetRevealResult fails for open idea.
	_, err = svc.GetRevealResult(ctx, ideaID)
	if err == nil {
		t.Error("GetRevealResult should fail for an open idea")
	}

	// Now expire the idea and process reveal.
	_, err = pool.Exec(ctx,
		`UPDATE ideas SET deadline = $1 WHERE id = $2`,
		time.Now().Add(-1*time.Hour), ideaID,
	)
	if err != nil {
		t.Fatalf("expire idea: %v", err)
	}

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	// After reveal, contributions should show author info.
	contribs, err = svc.ListContributions(ctx, ideaID)
	if err != nil {
		t.Fatalf("ListContributions post-reveal: %v", err)
	}
	for _, c := range contribs {
		if c.AuthorID == nil {
			t.Errorf("author_id should be visible post-reveal for contrib %d", c.ID)
		}
	}

	result, err := svc.GetRevealResult(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealResult post-reveal: %v", err)
	}
	if result.TotalVotes != 1 {
		t.Errorf("expected total_votes=1, got %d", result.TotalVotes)
	}
}

func TestReveal_ProcessExpiredIdeas(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	author := createUserLocal(t, pool, "author")

	expired1 := createIdeaDirectly(t, pool, initiator, time.Now().Add(-2*time.Hour))
	expired2 := createIdeaDirectly(t, pool, initiator, time.Now().Add(-1*time.Hour))
	notExpired := createIdeaDirectly(t, pool, initiator, time.Now().Add(24*time.Hour))

	createAndSubmitContribution(t, svc, author, expired1, "Solution")

	expiredIdeas, err := st.ListExpiredOpenIdeas(ctx)
	if err != nil {
		t.Fatalf("ListExpiredOpenIdeas: %v", err)
	}
	if len(expiredIdeas) != 2 {
		t.Fatalf("expected 2 expired ideas, got %d", len(expiredIdeas))
	}

	for _, idea := range expiredIdeas {
		if err := svc.ProcessReveal(ctx, idea.ID); err != nil {
			t.Fatalf("ProcessReveal for idea %d: %v", idea.ID, err)
		}
	}

	for _, eid := range []int64{expired1, expired2} {
		resp, err := svc.GetIdea(ctx, eid)
		if err != nil {
			t.Fatalf("GetIdea %d: %v", eid, err)
		}
		if resp.Status != model.IdeaStatusClosed {
			t.Errorf("idea %d: expected 'closed', got %q", eid, resp.Status)
		}
	}

	resp, err := svc.GetIdea(ctx, notExpired)
	if err != nil {
		t.Fatalf("GetIdea non-expired: %v", err)
	}
	if resp.Status != model.IdeaStatusOpen {
		t.Errorf("non-expired idea: expected 'open', got %q", resp.Status)
	}
}

func TestReveal_GetRevealResult_FailsForOpenIdea(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(24*time.Hour))

	_, err := svc.GetRevealResult(ctx, ideaID)
	if err == nil {
		t.Error("GetRevealResult should fail for an open idea")
	}
}

func TestReveal_SnapshotData_Integrity(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	t.Cleanup(func() { testutil.CleanupDB(t, pool); pool.Close() })

	svc := testutil.CreateTestService(t, pool)
	st := testutil.CreateTestStore(t, pool)
	ctx := context.Background()

	initiator := createUserLocal(t, pool, "initiator")
	author := createUserLocal(t, pool, "author")
	voters := make([]int64, 5)
	for i := range voters {
		voters[i] = createUserLocal(t, pool, fmt.Sprintf("voter%d", i))
	}

	ideaID := createIdeaDirectly(t, pool, initiator, time.Now().Add(-1*time.Hour))
	contribID := createAndSubmitContribution(t, svc, author, ideaID, "The one and only solution")

	for _, v := range voters {
		castVote(t, st, v, ideaID, contribID)
	}

	if err := svc.ProcessReveal(ctx, ideaID); err != nil {
		t.Fatalf("ProcessReveal: %v", err)
	}

	snap, err := st.GetRevealSnapshotByIdeaID(ctx, ideaID)
	if err != nil {
		t.Fatalf("GetRevealSnapshotByIdeaID: %v", err)
	}

	var ranked []model.RankedResult
	if err := json.Unmarshal(snap.RankedResults, &ranked); err != nil {
		t.Fatalf("unmarshal ranked_results: %v", err)
	}

	if len(ranked) != 1 {
		t.Fatalf("expected 1 ranked result, got %d", len(ranked))
	}
	if ranked[0].ContributionID != contribID {
		t.Errorf("expected contribution_id=%d, got %d", contribID, ranked[0].ContributionID)
	}
	if ranked[0].VoteCount != 5 {
		t.Errorf("expected vote_count=5, got %d", ranked[0].VoteCount)
	}
	if ranked[0].Rank != 1 {
		t.Errorf("expected rank=1, got %d", ranked[0].Rank)
	}
	if !ranked[0].IsFeatured {
		t.Error("expected is_featured=true for single contribution with 5 votes")
	}
	if snap.TotalVotes != 5 {
		t.Errorf("expected total_votes=5, got %d", snap.TotalVotes)
	}
	if snap.RevealedAt.IsZero() {
		t.Error("expected revealed_at to be set")
	}
}

// --- Assertion helpers ---

func assertRankedEntry(t *testing.T, entry service.RevealResultEntry, wantContribID int64, wantVotes, wantRank int, wantFeatured bool) {
	t.Helper()
	if entry.ContributionID != wantContribID {
		t.Errorf("contribution_id: want %d, got %d", wantContribID, entry.ContributionID)
	}
	if entry.VoteCount != wantVotes {
		t.Errorf("vote_count for contrib %d: want %d, got %d", wantContribID, wantVotes, entry.VoteCount)
	}
	if entry.Rank != wantRank {
		t.Errorf("rank for contrib %d: want %d, got %d", wantContribID, wantRank, entry.Rank)
	}
	if entry.IsFeatured != wantFeatured {
		t.Errorf("is_featured for contrib %d: want %v, got %v", wantContribID, wantFeatured, entry.IsFeatured)
	}
}
