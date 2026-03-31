package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	KeychainService = "claude-code"
	KeychainAccount = "oauth-tokens"
	APIKeyAccount   = "api-key"
)

// SecureStorage interface for storing credentials
type SecureStorage interface {
	Get(account string) (string, error)
	Set(account string, value string) error
	Delete(account string) error
}

// NewSecureStorage returns the platform-appropriate storage
func NewSecureStorage() SecureStorage {
	if runtime.GOOS == "darwin" {
		return &MacOSKeychain{}
	}
	return &FileStorage{}
}

// MacOSKeychain uses the macOS security CLI
type MacOSKeychain struct{}

func (k *MacOSKeychain) Get(account string) (string, error) {
	cmd := exec.Command("security", "find-generic-password",
		"-s", KeychainService,
		"-a", account,
		"-w")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("keychain read: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (k *MacOSKeychain) Set(account string, value string) error {
	// Delete first (ignore error if not exists)
	exec.Command("security", "delete-generic-password",
		"-s", KeychainService,
		"-a", account).Run()

	cmd := exec.Command("security", "add-generic-password",
		"-s", KeychainService,
		"-a", account,
		"-w", value,
		"-U")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("keychain write: %w (%s)", err, string(out))
	}
	return nil
}

func (k *MacOSKeychain) Delete(account string) error {
	cmd := exec.Command("security", "delete-generic-password",
		"-s", KeychainService,
		"-a", account)
	cmd.Run() // Ignore error if not exists
	return nil
}

// FileStorage is a fallback for non-macOS platforms
type FileStorage struct{}

func (f *FileStorage) storagePath(account string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "credentials", account+".json")
}

func (f *FileStorage) Get(account string) (string, error) {
	data, err := os.ReadFile(f.storagePath(account))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (f *FileStorage) Set(account string, value string) error {
	path := f.storagePath(account)
	os.MkdirAll(filepath.Dir(path), 0o700)
	return os.WriteFile(path, []byte(value), 0o600)
}

func (f *FileStorage) Delete(account string) error {
	return os.Remove(f.storagePath(account))
}

// SaveTokens stores OAuth tokens in secure storage
func SaveTokens(storage SecureStorage, tokens *OAuthTokens) error {
	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	return storage.Set(KeychainAccount, string(data))
}

// LoadTokens retrieves OAuth tokens from secure storage
func LoadTokens(storage SecureStorage) (*OAuthTokens, error) {
	data, err := storage.Get(KeychainAccount)
	if err != nil {
		return nil, err
	}
	var tokens OAuthTokens
	if err := json.Unmarshal([]byte(data), &tokens); err != nil {
		return nil, err
	}
	return &tokens, nil
}

// SaveAPIKey stores an API key in secure storage
func SaveAPIKey(storage SecureStorage, key string) error {
	return storage.Set(APIKeyAccount, key)
}

// LoadAPIKey retrieves an API key from secure storage
func LoadAPIKey(storage SecureStorage) (string, error) {
	return storage.Get(APIKeyAccount)
}
