package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/claway/server/internal/model"
	"github.com/jackc/pgx/v5"
)

const userColumns = `id, openclaw_id, username, display_name, avatar_url, created_at, updated_at`

func scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.OpenClawID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetUserByOpenClawID retrieves a user by their OpenClaw ID (legacy).
func (s *Store) GetUserByOpenClawID(ctx context.Context, openclawID string) (*model.User, error) {
	u, err := scanUser(s.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE openclaw_id = $1`, openclawID))
	if err != nil {
		return nil, fmt.Errorf("get user by openclaw_id: %w", err)
	}
	return u, nil
}

// GetUserByID retrieves a user by their internal ID.
func (s *Store) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	u, err := scanUser(s.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`, id))
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

// GetUserByUsername retrieves a user by their username.
func (s *Store) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	u, err := scanUser(s.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE username = $1`, username))
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return u, nil
}

// GetUsersByIDs retrieves multiple users by their IDs in a single query.
func (s *Store) GetUsersByIDs(ctx context.Context, ids []int64) (map[int64]*model.User, error) {
	if len(ids) == 0 {
		return make(map[int64]*model.User), nil
	}
	rows, err := s.db.Query(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, fmt.Errorf("get users by ids: %w", err)
	}
	defer rows.Close()

	result := make(map[int64]*model.User, len(ids))
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.OpenClawID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("get users by ids scan: %w", err)
		}
		result[u.ID] = &u
	}
	return result, nil
}

// CreateUser inserts a new user and returns the created record.
func (s *Store) CreateUser(ctx context.Context, openclawID, username string) (*model.User, error) {
	u, err := scanUser(s.db.QueryRow(ctx,
		`INSERT INTO users (openclaw_id, username)
		 VALUES ($1, $2)
		 RETURNING `+userColumns,
		openclawID, username))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

// CreateUserFromOAuth creates a user from an OAuth profile.
func (s *Store) CreateUserFromOAuth(ctx context.Context, username, displayName, avatarURL string) (*model.User, error) {
	u, err := scanUser(s.db.QueryRow(ctx,
		`INSERT INTO users (openclaw_id, username, display_name, avatar_url)
		 VALUES ('', $1, $2, $3)
		 RETURNING `+userColumns,
		username, displayName, avatarURL))
	if err != nil {
		return nil, fmt.Errorf("create user from oauth: %w", err)
	}
	return u, nil
}

// --- OAuth Accounts ---

// GetOAuthAccount finds a linked OAuth account by provider and provider user ID.
func (s *Store) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*model.OAuthAccount, error) {
	var a model.OAuthAccount
	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, provider, provider_user_id, provider_username, provider_email,
		        access_token, refresh_token, token_expires_at, created_at, updated_at
		 FROM user_oauth_accounts
		 WHERE provider = $1 AND provider_user_id = $2`,
		provider, providerUserID,
	).Scan(&a.ID, &a.UserID, &a.Provider, &a.ProviderUserID, &a.ProviderUsername, &a.ProviderEmail,
		&a.AccessToken, &a.RefreshToken, &a.TokenExpiresAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get oauth account: %w", err)
	}
	return &a, nil
}

// CreateOAuthAccount links an OAuth account to a user.
func (s *Store) CreateOAuthAccount(ctx context.Context, account *model.OAuthAccount) (*model.OAuthAccount, error) {
	var a model.OAuthAccount
	err := s.db.QueryRow(ctx,
		`INSERT INTO user_oauth_accounts (user_id, provider, provider_user_id, provider_username, provider_email, access_token, refresh_token, token_expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, user_id, provider, provider_user_id, provider_username, provider_email,
		           access_token, refresh_token, token_expires_at, created_at, updated_at`,
		account.UserID, account.Provider, account.ProviderUserID, account.ProviderUsername, account.ProviderEmail,
		account.AccessToken, account.RefreshToken, account.TokenExpiresAt,
	).Scan(&a.ID, &a.UserID, &a.Provider, &a.ProviderUserID, &a.ProviderUsername, &a.ProviderEmail,
		&a.AccessToken, &a.RefreshToken, &a.TokenExpiresAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create oauth account: %w", err)
	}
	return &a, nil
}

// UpdateOAuthTokens refreshes stored tokens for an OAuth account.
func (s *Store) UpdateOAuthTokens(ctx context.Context, id int64, accessToken, refreshToken string, expiresAt interface{}) error {
	_, err := s.db.Exec(ctx,
		`UPDATE user_oauth_accounts
		 SET access_token = $1, refresh_token = $2, token_expires_at = $3, updated_at = NOW()
		 WHERE id = $4`,
		accessToken, refreshToken, expiresAt, id,
	)
	if err != nil {
		return fmt.Errorf("update oauth tokens: %w", err)
	}
	return nil
}
