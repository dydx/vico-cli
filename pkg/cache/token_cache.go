package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TokenCache represents the structure of our token cache file
type TokenCache struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TokenCacheManager handles reading and writing the authentication token to a cache file
type TokenCacheManager struct {
	CacheDir  string
	CacheFile string
}

// NewTokenCacheManager creates a new cache manager
func NewTokenCacheManager() (*TokenCacheManager, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	// Create cache directory if it doesn't exist
	cacheDir := filepath.Join(homeDir, ".vicohome")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("error creating cache directory: %w", err)
	}

	return &TokenCacheManager{
		CacheDir:  cacheDir,
		CacheFile: filepath.Join(cacheDir, "auth.json"),
	}, nil
}

// SaveToken saves a token to the cache file with an expiration time
func (m *TokenCacheManager) SaveToken(token string, durationHours int) error {
	// Default to 24 hours if not specified
	if durationHours <= 0 {
		durationHours = 24
	}

	// Create token cache object
	tokenCache := TokenCache{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(durationHours) * time.Hour),
	}

	// Marshal to JSON
	cacheData, err := json.Marshal(tokenCache)
	if err != nil {
		return fmt.Errorf("error marshaling token cache: %w", err)
	}

	// Write to file
	if err := os.WriteFile(m.CacheFile, cacheData, 0600); err != nil {
		return fmt.Errorf("error writing token cache: %w", err)
	}

	return nil
}

// GetToken retrieves a token from the cache file if it exists and is not expired
func (m *TokenCacheManager) GetToken() (string, bool) {
	// Check if cache file exists
	_, err := os.Stat(m.CacheFile)
	if os.IsNotExist(err) {
		return "", false
	}

	// Read cache file
	cacheData, err := os.ReadFile(m.CacheFile)
	if err != nil {
		// If there's an error reading, treat as if no cache exists
		return "", false
	}

	// Unmarshal data
	var tokenCache TokenCache
	if err := json.Unmarshal(cacheData, &tokenCache); err != nil {
		// If there's an error unmarshaling, treat as if no cache exists
		return "", false
	}

	// Check if token is expired
	if time.Now().After(tokenCache.ExpiresAt) {
		return "", false
	}

	// Return valid token
	return tokenCache.Token, true
}

// ClearToken removes the token cache file
func (m *TokenCacheManager) ClearToken() error {
	_, err := os.Stat(m.CacheFile)
	if os.IsNotExist(err) {
		// File doesn't exist, nothing to clear
		return nil
	}

	// Remove the cache file
	if err := os.Remove(m.CacheFile); err != nil {
		return fmt.Errorf("error removing token cache: %w", err)
	}

	return nil
}