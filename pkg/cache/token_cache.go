// Package cache provides caching functionality for authentication tokens.
//
// This package implements file-based caching of authentication tokens in the user's
// home directory, with support for token expiration and management operations.
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TokenCache represents the structure of the token cache file.
// It stores both the authentication token and its expiration time.
type TokenCache struct {
	Token     string    `json:"token"`      // The authentication token
	ExpiresAt time.Time `json:"expires_at"` // When the token expires
}

// TokenCacheManager handles reading and writing authentication tokens to a cache file.
// It provides methods for saving, retrieving, and clearing tokens, with support for
// automatic expiration checking.
type TokenCacheManager struct {
	CacheDir  string // Directory where cache files are stored
	CacheFile string // Full path to the token cache file
}

// NewTokenCacheManager creates a new token cache manager.
// It sets up the cache directory in the user's home directory if it doesn't already exist.
// The cache location is ~/.vicohome/auth.json
//
// Returns:
//   - *TokenCacheManager: The configured cache manager if successful
//   - error: Any error encountered during setup
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

// SaveToken saves an authentication token to the cache file with an expiration time.
// The token is stored along with its expiration time calculated from the current time
// plus the specified duration in hours.
//
// Parameters:
//   - token: The authentication token to save
//   - durationHours: How long the token should be considered valid, in hours
//
// Returns:
//   - error: Any error encountered during the save operation
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

// GetToken retrieves a token from the cache file if it exists and is not expired.
// It checks both for the existence of the cache file and whether the token inside
// has expired based on its expiration timestamp.
//
// Returns:
//   - string: The cached token if valid
//   - bool: True if a valid token was found, false otherwise
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

// ClearToken removes the token cache file if it exists.
// This is typically called when a token is known to be invalid or expired.
//
// Returns:
//   - error: Any error encountered during the removal operation
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
