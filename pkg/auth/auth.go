package auth

import (
	"context"
	"fmt"
	"os"
)

// GetAuthToken returns a valid API token, trying all sources in order:
// 1. ANTHROPIC_API_KEY env var
// 2. Keychain/secure storage (API key or OAuth token)
// 3. Interactive login flow
func GetAuthToken(ctx context.Context, interactive bool) (string, error) {
	// 1. Environment variable
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return key, nil
	}
	if key := os.Getenv("CLAUDE_API_KEY"); key != "" {
		return key, nil
	}

	storage := NewSecureStorage()

	// 2a. Try stored API key
	if key, err := LoadAPIKey(storage); err == nil && key != "" {
		return key, nil
	}

	// 2b. Try stored OAuth tokens
	if tokens, err := LoadTokens(storage); err == nil && tokens != nil {
		if !tokens.IsExpired() {
			return tokens.AccessToken, nil
		}
		// Try refresh
		if tokens.RefreshToken != "" {
			newTokens, err := RefreshTokens(tokens.RefreshToken)
			if err == nil {
				SaveTokens(storage, newTokens)
				return newTokens.AccessToken, nil
			}
		}
	}

	// 3. Interactive login
	if !interactive {
		return "", fmt.Errorf("no API key found. Set ANTHROPIC_API_KEY or run 'claude-go' interactively to login")
	}

	fmt.Println("No API key found. Starting login flow...")
	fmt.Println()

	tokens, err := StartOAuthFlow(ctx)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	// Save tokens
	if err := SaveTokens(storage, tokens); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not save tokens: %v\n", err)
	}

	fmt.Println("Login successful!")
	return tokens.AccessToken, nil
}

// Logout clears all stored credentials
func Logout() error {
	storage := NewSecureStorage()
	storage.Delete(KeychainAccount)
	storage.Delete(APIKeyAccount)
	return nil
}

// IsLoggedIn checks if there are valid credentials
func IsLoggedIn() bool {
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return true
	}
	storage := NewSecureStorage()
	if key, err := LoadAPIKey(storage); err == nil && key != "" {
		return true
	}
	if tokens, err := LoadTokens(storage); err == nil && tokens != nil {
		return true
	}
	return false
}
