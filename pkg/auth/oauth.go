package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	// Production OAuth endpoints - from old-src/constants/oauth.ts PROD_OAUTH_CONFIG
	AuthorizationEndpoint = "https://platform.claude.com/oauth/authorize"
	ClaudeAIAuthorizeURL  = "https://claude.com/cai/oauth/authorize"
	TokenEndpoint         = "https://platform.claude.com/v1/oauth/token"
	APIKeyEndpoint        = "https://api.anthropic.com/api/oauth/claude_cli/create_api_key"
	SuccessURL            = "https://platform.claude.com/oauth/code/success?app=claude-code"
	ManualRedirectURL     = "https://platform.claude.com/oauth/code/callback"
	ClientID              = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	RedirectPath          = "/oauth/callback"
	DefaultScopes         = "org:create_api_key user:profile user:inference user:sessions:claude_code user:mcp_servers user:file_upload"
	OAuthBetaHeader       = "oauth-2025-04-20"
)

// OAuthTokens holds the token pair
type OAuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope,omitempty"`
}

// IsExpired checks if the access token is expired
func (t *OAuthTokens) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-30 * time.Second))
}

// generateCodeVerifier creates a random PKCE code verifier
func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge creates the S256 code challenge from a verifier
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// generateState creates a random state parameter
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// StartOAuthFlow initiates the OAuth PKCE authorization flow
// Returns tokens on success
func StartOAuthFlow(ctx context.Context) (*OAuthTokens, error) {
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("generate verifier: %w", err)
	}

	challenge := generateCodeChallenge(verifier)

	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("generate state: %w", err)
	}

	// Start local callback server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("start callback server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	redirectURI := fmt.Sprintf("http://127.0.0.1:%d%s", port, RedirectPath)

	// Build authorization URL
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {ClientID},
		"redirect_uri":          {redirectURI},
		"scope":                 {DefaultScopes},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}

	authURL := fmt.Sprintf("%s?%s", AuthorizationEndpoint, params.Encode())

	// Channel to receive the auth code
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Handle callback
	mux := http.NewServeMux()
	mux.HandleFunc(RedirectPath, func(w http.ResponseWriter, r *http.Request) {
		receivedState := r.URL.Query().Get("state")
		if receivedState != state {
			errCh <- fmt.Errorf("state mismatch")
			http.Error(w, "State mismatch", http.StatusBadRequest)
			return
		}

		if errParam := r.URL.Query().Get("error"); errParam != "" {
			errCh <- fmt.Errorf("OAuth error: %s - %s", errParam, r.URL.Query().Get("error_description"))
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>%s</p><p>You can close this window.</p></body></html>", errParam)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no authorization code received")
			http.Error(w, "No code", http.StatusBadRequest)
			return
		}

		codeCh <- code
		fmt.Fprint(w, "<html><body><h1>Authentication Successful!</h1><p>You can close this window and return to Claude Code.</p></body></html>")
	})

	server := &http.Server{Handler: mux}
	go server.Serve(listener)
	defer server.Shutdown(context.Background())

	// Open browser
	fmt.Printf("\nOpening browser for authentication...\n")
	fmt.Printf("If the browser doesn't open, visit:\n  %s\n\n", authURL)
	openBrowser(authURL)

	// Wait for callback or timeout
	select {
	case code := <-codeCh:
		return exchangeCodeForTokens(code, verifier, redirectURI)
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authentication timed out after 5 minutes")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// exchangeCodeForTokens exchanges an authorization code for tokens
func exchangeCodeForTokens(code, verifier, redirectURI string) (*OAuthTokens, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {ClientID},
		"code_verifier": {verifier},
	}

	resp, err := http.PostForm(TokenEndpoint, data)
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokens OAuthTokens
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("parse tokens: %w", err)
	}

	tokens.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	return &tokens, nil
}

// RefreshTokens uses a refresh token to get new tokens
func RefreshTokens(refreshToken string) (*OAuthTokens, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {ClientID},
	}

	resp, err := http.PostForm(TokenEndpoint, data)
	if err != nil {
		return nil, fmt.Errorf("token refresh: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokens OAuthTokens
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("parse tokens: %w", err)
	}

	tokens.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	return &tokens, nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", strings.ReplaceAll(url, "&", "^&"))
	}
	if cmd != nil {
		cmd.Start()
	}
}
