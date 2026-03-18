package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/store"
)

// --- OAuth state management (in-memory, good enough for single-instance MVP) ---

type oauthState struct {
	CodeVerifier string
	CLIPort      string // empty for web flow
	SessionID    string // non-empty for agent session flow
	CreatedAt    time.Time
}

var (
	oauthStates   = make(map[string]*oauthState)
	oauthStatesMu sync.Mutex
)

func saveOAuthState(state, codeVerifier, cliPort, sessionID string) {
	oauthStatesMu.Lock()
	defer oauthStatesMu.Unlock()
	oauthStates[state] = &oauthState{
		CodeVerifier: codeVerifier,
		CLIPort:      cliPort,
		SessionID:    sessionID,
		CreatedAt:    time.Now(),
	}
}

func popOAuthState(state string) (*oauthState, bool) {
	oauthStatesMu.Lock()
	defer oauthStatesMu.Unlock()
	s, ok := oauthStates[state]
	if ok {
		delete(oauthStates, state)
	}
	return s, ok
}

// --- X (Twitter) OAuth 2.0 ---

type xTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type xUserResponse struct {
	Data struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Username        string `json:"username"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

// GetXAuthURL generates the X OAuth authorization URL with PKCE.
func (s *Service) GetXAuthURL(cliPort string) (string, error) {
	if s.cfg.XClientID == "" {
		return "", fmt.Errorf("X_CLIENT_ID not configured")
	}

	// Generate PKCE code verifier and challenge
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return "", fmt.Errorf("generate code verifier: %w", err)
	}
	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
	challengeHash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(challengeHash[:])

	// Generate state
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	saveOAuthState(state, codeVerifier, cliPort, "")

	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {s.cfg.XClientID},
		"redirect_uri":          {s.cfg.XRedirectURI},
		"scope":                 {"tweet.read users.read offline.access"},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	return "https://twitter.com/i/oauth2/authorize?" + params.Encode(), nil
}

// HandleXCallback processes the X OAuth callback.
func (s *Service) HandleXCallback(ctx context.Context, code, state string) (redirectURL string, err error) {
	saved, ok := popOAuthState(state)
	if !ok {
		return "", fmt.Errorf("invalid or expired OAuth state")
	}

	// Exchange code for token
	tokenResp, err := s.exchangeXCode(code, saved.CodeVerifier)
	if err != nil {
		return "", err
	}

	// Fetch user profile
	profile, err := s.fetchXProfile(tokenResp.AccessToken)
	if err != nil {
		return "", err
	}

	// Find or create user
	tokenData := &oauthTokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	}
	user, err := s.findOrCreateOAuthUser(ctx, "x", profile.Data.ID, profile.Data.Username, profile.Data.Name, profile.Data.ProfileImageURL, "", tokenData)
	if err != nil {
		return "", err
	}

	// Issue JWT
	jwtToken, err := s.issueJWT(user.ID)
	if err != nil {
		return "", err
	}

	return s.buildAuthRedirect(saved, jwtToken)
}

func (s *Service) exchangeXCode(code, codeVerifier string) (*xTokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {s.cfg.XRedirectURI},
		"code_verifier": {codeVerifier},
	}

	req, err := http.NewRequest("POST", "https://api.x.com/2/oauth2/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// X requires Basic Auth (client_id:client_secret) for confidential clients
	req.SetBasicAuth(s.cfg.XClientID, s.cfg.XClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp xTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &tokenResp, nil
}

func (s *Service) fetchXProfile(accessToken string) (*xUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.x.com/2/users/me?user.fields=profile_image_url", nil)
	if err != nil {
		return nil, fmt.Errorf("create profile request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("profile fetch failed (status %d): %s", resp.StatusCode, string(body))
	}

	var profile xUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	return &profile, nil
}

// sanitizeAvatarURL validates that the URL uses https and is well-formed.
// Returns empty string for invalid URLs to fall back to initials avatar on frontend.
func sanitizeAvatarURL(raw string) string {
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return ""
	}
	return u.String()
}

// --- GitHub OAuth 2.0 ---

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type githubUserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// GetGitHubAuthURL generates the GitHub OAuth authorization URL.
func (s *Service) GetGitHubAuthURL(cliPort string) (string, error) {
	if s.cfg.GitHubClientID == "" {
		return "", fmt.Errorf("GITHUB_CLIENT_ID not configured")
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	saveOAuthState(state, "", cliPort, "")

	params := url.Values{
		"client_id":    {s.cfg.GitHubClientID},
		"redirect_uri": {s.cfg.GitHubRedirectURI},
		"scope":        {"read:user user:email"},
		"state":        {state},
	}

	return "https://github.com/login/oauth/authorize?" + params.Encode(), nil
}

// HandleGitHubCallback processes the GitHub OAuth callback.
func (s *Service) HandleGitHubCallback(ctx context.Context, code, state string) (redirectURL string, err error) {
	saved, ok := popOAuthState(state)
	if !ok {
		return "", fmt.Errorf("invalid or expired OAuth state")
	}

	tokenResp, err := s.exchangeGitHubCode(code)
	if err != nil {
		return "", err
	}

	profile, err := s.fetchGitHubProfile(tokenResp.AccessToken)
	if err != nil {
		return "", err
	}

	displayName := profile.Name
	if displayName == "" {
		displayName = profile.Login
	}

	tokenData := &oauthTokenData{
		AccessToken: tokenResp.AccessToken,
	}
	user, err := s.findOrCreateOAuthUser(ctx, "github", strconv.FormatInt(profile.ID, 10), profile.Login, displayName, profile.AvatarURL, profile.Email, tokenData)
	if err != nil {
		return "", err
	}

	jwtToken, err := s.issueJWT(user.ID)
	if err != nil {
		return "", err
	}

	return s.buildAuthRedirect(saved, jwtToken)
}

func (s *Service) exchangeGitHubCode(code string) (*githubTokenResponse, error) {
	data := url.Values{
		"client_id":     {s.cfg.GitHubClientID},
		"client_secret": {s.cfg.GitHubClientSecret},
		"code":          {code},
		"redirect_uri":  {s.cfg.GitHubRedirectURI},
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp githubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("GitHub token exchange returned empty access token")
	}

	return &tokenResp, nil
}

func (s *Service) fetchGitHubProfile(accessToken string) (*githubUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("create profile request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("profile fetch failed (status %d): %s", resp.StatusCode, string(body))
	}

	var profile githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	return &profile, nil
}

// --- Google OAuth 2.0 ---

type googleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type googleUserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GetGoogleAuthURL generates the Google OAuth authorization URL.
func (s *Service) GetGoogleAuthURL(cliPort string) (string, error) {
	if s.cfg.GoogleClientID == "" {
		return "", fmt.Errorf("GOOGLE_CLIENT_ID not configured")
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	saveOAuthState(state, "", cliPort, "")

	params := url.Values{
		"client_id":     {s.cfg.GoogleClientID},
		"redirect_uri":  {s.cfg.GoogleRedirectURI},
		"response_type": {"code"},
		"scope":         {"openid profile email"},
		"state":         {state},
		"access_type":   {"offline"},
	}

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode(), nil
}

// HandleGoogleCallback processes the Google OAuth callback.
func (s *Service) HandleGoogleCallback(ctx context.Context, code, state string) (redirectURL string, err error) {
	saved, ok := popOAuthState(state)
	if !ok {
		return "", fmt.Errorf("invalid or expired OAuth state")
	}

	tokenResp, err := s.exchangeGoogleCode(code)
	if err != nil {
		return "", err
	}

	profile, err := s.fetchGoogleProfile(tokenResp.AccessToken)
	if err != nil {
		return "", err
	}

	// Use email prefix as username
	username := profile.Email
	if idx := strings.Index(profile.Email, "@"); idx > 0 {
		username = profile.Email[:idx]
	}

	displayName := profile.Name
	if displayName == "" {
		displayName = username
	}

	tokenData := &oauthTokenData{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	}
	user, err := s.findOrCreateOAuthUser(ctx, "google", profile.ID, username, displayName, profile.Picture, profile.Email, tokenData)
	if err != nil {
		return "", err
	}

	jwtToken, err := s.issueJWT(user.ID)
	if err != nil {
		return "", err
	}

	return s.buildAuthRedirect(saved, jwtToken)
}

func (s *Service) exchangeGoogleCode(code string) (*googleTokenResponse, error) {
	data := url.Values{
		"client_id":     {s.cfg.GoogleClientID},
		"client_secret": {s.cfg.GoogleClientSecret},
		"code":          {code},
		"redirect_uri":  {s.cfg.GoogleRedirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp googleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}

	return &tokenResp, nil
}

func (s *Service) fetchGoogleProfile(accessToken string) (*googleUserResponse, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("create profile request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("profile fetch failed (status %d): %s", resp.StatusCode, string(body))
	}

	var profile googleUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	return &profile, nil
}

// buildAuthRedirect constructs the redirect URL based on the OAuth flow type.
func (s *Service) buildAuthRedirect(saved *oauthState, jwtToken string) (string, error) {
	if saved.SessionID != "" {
		ctx := context.Background()
		if err := s.store.CompleteAuthSession(ctx, saved.SessionID, jwtToken); err != nil {
			return "", fmt.Errorf("complete auth session: %w", err)
		}
		return fmt.Sprintf("%s/auth/session-success", s.cfg.FrontendURL), nil
	}

	if saved.CLIPort != "" {
		return fmt.Sprintf("http://127.0.0.1:%s/callback?token=%s", saved.CLIPort, url.QueryEscape(jwtToken)), nil
	}

	return fmt.Sprintf("%s/auth/callback?token=%s", s.cfg.FrontendURL, url.QueryEscape(jwtToken)), nil
}

// oauthTokenData is a provider-agnostic token container used by findOrCreateOAuthUser.
type oauthTokenData struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

func (s *Service) findOrCreateOAuthUser(ctx context.Context, provider, providerUserID, username, displayName, avatarURL, email string, tokenData *oauthTokenData) (*model.User, error) {
	avatarURL = sanitizeAvatarURL(avatarURL)
	// Check if OAuth account already exists
	oauthAccount, err := s.store.GetOAuthAccount(ctx, provider, providerUserID)
	if err == nil {
		// Existing account — update tokens and return user
		var expiresAt interface{}
		if tokenData.ExpiresIn > 0 {
			t := time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)
			expiresAt = t
		}
		_ = s.store.UpdateOAuthTokens(ctx, oauthAccount.ID, tokenData.AccessToken, tokenData.RefreshToken, expiresAt)

		user, err := s.store.GetUserByID(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("get user for oauth account: %w", err)
		}
		return user, nil
	}

	if !errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("lookup oauth account: %w", err)
	}

	// New user — create user + oauth account
	user, err := s.store.CreateUserFromOAuth(ctx, username, displayName, avatarURL)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	account := &model.OAuthAccount{
		UserID:           user.ID,
		Provider:         provider,
		ProviderUserID:   providerUserID,
		ProviderUsername:  username,
		ProviderEmail:    email,
		AccessToken:      tokenData.AccessToken,
		RefreshToken:     tokenData.RefreshToken,
	}
	if tokenData.ExpiresIn > 0 {
		account.TokenExpiresAt = model.NullTime{NullTime: sql.NullTime{
			Time:  time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second),
			Valid: true,
		}}
	}
	_, err = s.store.CreateOAuthAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("create oauth account: %w", err)
	}

	return user, nil
}

// --- Legacy OpenClaw OAuth (kept for backward compatibility) ---

type openClawTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type openClawUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthCallbackResponse is the response returned after successful OAuth callback.
type AuthCallbackResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *Service) HandleOpenClawCallback(ctx context.Context, code string) (*AuthCallbackResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	tokenURL := fmt.Sprintf("%s/oauth/token", s.cfg.OpenClawBaseURL)
	resp, err := http.PostForm(tokenURL, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {s.cfg.OpenClawClientID},
		"client_secret": {s.cfg.OpenClawClientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp openClawTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	profileURL := fmt.Sprintf("%s/api/v1/me", s.cfg.OpenClawBaseURL)
	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, profileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile request: %w", err)
	}
	profileReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}
	defer profileResp.Body.Close()

	if profileResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(profileResp.Body)
		return nil, fmt.Errorf("profile fetch failed (status %d): %s", profileResp.StatusCode, string(body))
	}

	var profile openClawUserResponse
	if err := json.NewDecoder(profileResp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile response: %w", err)
	}

	user, err := s.store.GetUserByOpenClawID(ctx, profile.ID)
	if errors.Is(err, store.ErrNotFound) {
		user, err = s.store.CreateUser(ctx, profile.ID, profile.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to look up user: %w", err)
	}

	token, err := s.issueJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to issue JWT: %w", err)
	}

	return &AuthCallbackResponse{
		Token: token,
		User:  user,
	}, nil
}

// --- Agent Auth Sessions ---

const authSessionTTL = 5 * time.Minute

// CreateAuthSession generates a new pending auth session and returns it along
// with the X OAuth authorization URL that includes the session ID.
func (s *Service) CreateAuthSession(ctx context.Context) (*model.AuthSession, string, error) {
	// Generate session ID (UUID v4 via crypto/rand)
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return nil, "", fmt.Errorf("generate session id: %w", err)
	}
	// Format as UUID v4
	idBytes[6] = (idBytes[6] & 0x0f) | 0x40 // version 4
	idBytes[8] = (idBytes[8] & 0x3f) | 0x80 // variant 10
	sessionID := fmt.Sprintf("%x-%x-%x-%x-%x",
		idBytes[0:4], idBytes[4:6], idBytes[6:8], idBytes[8:10], idBytes[10:16])

	now := time.Now()
	session := &model.AuthSession{
		ID:        sessionID,
		Status:    "pending",
		ExpiresAt: now.Add(authSessionTTL),
		CreatedAt: now,
	}

	if err := s.store.CreateAuthSession(ctx, session); err != nil {
		return nil, "", fmt.Errorf("create auth session: %w", err)
	}

	// Generate an X OAuth URL that carries the session ID through the state flow.
	authURL, err := s.getXAuthURLForSession(sessionID)
	if err != nil {
		return nil, "", err
	}

	return session, authURL, nil
}

// getXAuthURLForSession is like GetXAuthURL but tags the OAuth state with a session ID.
func (s *Service) getXAuthURLForSession(sessionID string) (string, error) {
	if s.cfg.XClientID == "" {
		return "", fmt.Errorf("X_CLIENT_ID not configured")
	}

	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return "", fmt.Errorf("generate code verifier: %w", err)
	}
	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
	challengeHash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(challengeHash[:])

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("generate state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	saveOAuthState(state, codeVerifier, "", sessionID)

	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {s.cfg.XClientID},
		"redirect_uri":          {s.cfg.XRedirectURI},
		"scope":                 {"tweet.read users.read offline.access"},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	return "https://twitter.com/i/oauth2/authorize?" + params.Encode(), nil
}

// GetAuthSession returns the current state of an auth session.
func (s *Service) GetAuthSession(ctx context.Context, id string) (*model.AuthSession, error) {
	session, err := s.store.GetAuthSession(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// StartAuthSessionCleanup starts a background goroutine that periodically
// removes expired auth sessions. It stops when the context is cancelled.
func (s *Service) StartAuthSessionCleanup(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.store.CleanupExpiredAuthSessions()
			}
		}
	}()
}

// --- JWT ---

// GetMe returns the current user's profile.
func (s *Service) GetMe(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// issueJWT creates a signed JWT token for the given user ID.
func (s *Service) issueJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signed, nil
}
