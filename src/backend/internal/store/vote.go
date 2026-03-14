package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateVote records a vote. Returns ErrConflict if the user already voted on this idea.
func (s *Store) CreateVote(ctx context.Context, v *model.Vote) (*model.Vote, error) {
	var result model.Vote
	err := s.db.QueryRow(ctx,
		`INSERT INTO votes (idea_id, voter_id, contribution_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, idea_id, voter_id, contribution_id, voted_at`,
		v.IdeaID, v.VoterID, v.ContributionID,
	).Scan(&result.ID, &result.IdeaID, &result.VoterID, &result.ContributionID, &result.VotedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, fmt.Errorf("create vote: %w", err)
	}
	return &result, nil
}

// CountVotersByIdeaIDs returns distinct voter counts for multiple ideas in one query.
func (s *Store) CountVotersByIdeaIDs(ctx context.Context, ideaIDs []int64) (map[int64]int, error) {
	if len(ideaIDs) == 0 {
		return make(map[int64]int), nil
	}
	rows, err := s.db.Query(ctx,
		`SELECT idea_id, COUNT(DISTINCT voter_id) FROM votes
		 WHERE idea_id = ANY($1)
		 GROUP BY idea_id`, ideaIDs)
	if err != nil {
		return nil, fmt.Errorf("count voters by idea ids: %w", err)
	}
	defer rows.Close()

	result := make(map[int64]int, len(ideaIDs))
	for rows.Next() {
		var id int64
		var count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, fmt.Errorf("count voters by idea ids scan: %w", err)
		}
		result[id] = count
	}
	return result, nil
}

// CountVotesByIdea returns the total number of votes for an idea.
func (s *Store) CountVotesByIdea(ctx context.Context, ideaID int64) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM votes WHERE idea_id = $1`, ideaID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count votes by idea: %w", err)
	}
	return count, nil
}

// CountVotersByIdea returns the number of distinct voters for an idea.
func (s *Store) CountVotersByIdea(ctx context.Context, ideaID int64) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(DISTINCT voter_id) FROM votes WHERE idea_id = $1`, ideaID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count voters by idea: %w", err)
	}
	return count, nil
}

// HasUserVotedOnIdea checks if a user already voted on an idea.
func (s *Store) HasUserVotedOnIdea(ctx context.Context, ideaID, userID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM votes WHERE idea_id = $1 AND voter_id = $2)`,
		ideaID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user voted: %w", err)
	}
	return exists, nil
}

// GetVoteCountsByIdea returns vote counts per contribution for a given idea,
// ordered by vote_count DESC, then by earliest submission time.
type ContributionVoteCount struct {
	ContributionID int64 `json:"contribution_id"`
	VoteCount      int   `json:"vote_count"`
}

func (s *Store) GetVoteCountsByIdea(ctx context.Context, ideaID int64) ([]ContributionVoteCount, int, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, COUNT(v.id) AS vote_count
		 FROM contributions c
		 LEFT JOIN votes v ON v.contribution_id = c.id
		 WHERE c.idea_id = $1 AND c.status = 'submitted'
		 GROUP BY c.id, c.submitted_at
		 ORDER BY vote_count DESC, c.submitted_at ASC`,
		ideaID,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("get vote counts by idea: %w", err)
	}
	defer rows.Close()

	var results []ContributionVoteCount
	totalVotes := 0
	for rows.Next() {
		var r ContributionVoteCount
		if err := rows.Scan(&r.ContributionID, &r.VoteCount); err != nil {
			return nil, 0, fmt.Errorf("vote counts scan: %w", err)
		}
		totalVotes += r.VoteCount
		results = append(results, r)
	}
	return results, totalVotes, nil
}

// ListVotesByUser returns all votes cast by a user.
func (s *Store) ListVotesByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.Vote, int, error) {
	var total int
	if err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM votes WHERE voter_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("list votes by user count: %w", err)
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, idea_id, voter_id, contribution_id, voted_at
		 FROM votes WHERE voter_id = $1
		 ORDER BY voted_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list votes by user query: %w", err)
	}
	defer rows.Close()

	var votes []*model.Vote
	for rows.Next() {
		var v model.Vote
		if err := rows.Scan(&v.ID, &v.IdeaID, &v.VoterID, &v.ContributionID, &v.VotedAt); err != nil {
			return nil, 0, fmt.Errorf("list votes by user scan: %w", err)
		}
		votes = append(votes, &v)
	}
	return votes, total, nil
}

// GetUserVoteForIdea returns the user's vote for a specific idea, or ErrNotFound.
func (s *Store) GetUserVoteForIdea(ctx context.Context, ideaID, userID int64) (*model.Vote, error) {
	var v model.Vote
	err := s.db.QueryRow(ctx,
		`SELECT id, idea_id, voter_id, contribution_id, voted_at
		 FROM votes WHERE idea_id = $1 AND voter_id = $2`,
		ideaID, userID,
	).Scan(&v.ID, &v.IdeaID, &v.VoterID, &v.ContributionID, &v.VotedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user vote for idea: %w", err)
	}
	return &v, nil
}
