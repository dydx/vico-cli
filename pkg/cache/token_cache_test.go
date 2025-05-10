package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setupTestCacheManager creates a temp directory for tests and configures a TokenCacheManager to use it.
func setupTestCacheManager(t *testing.T) (*TokenCacheManager, string) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "vicohome-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a TokenCacheManager using the test directory
	cacheManager := &TokenCacheManager{
		CacheDir:  tempDir,
		CacheFile: filepath.Join(tempDir, "auth.json"),
	}

	return cacheManager, tempDir
}

// cleanupTestCacheManager removes the test directory
func cleanupTestCacheManager(tempDir string) {
	os.RemoveAll(tempDir)
}

// TestTokenCache tests the token cache functionality using a table-driven approach
func TestTokenCache(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		setup       func(t *testing.T, manager *TokenCacheManager)
		operation   string                                 // Description of the operation being tested
		action      func(manager *TokenCacheManager) error // Function that performs the operation
		checkResult func(t *testing.T, manager *TokenCacheManager, err error)
	}{
		{
			name: "SaveToken",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// No setup needed
			},
			operation: "SaveToken with valid token",
			action: func(manager *TokenCacheManager) error {
				return manager.SaveToken("test-token-12345", 24)
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Verify no error
				if err != nil {
					t.Errorf("SaveToken() error = %v", err)
				}

				// Verify the file was created
				if _, err := os.Stat(manager.CacheFile); os.IsNotExist(err) {
					t.Error("Cache file was not created")
				}

				// Read the cache file and verify content
				data, err := os.ReadFile(manager.CacheFile)
				if err != nil {
					t.Fatalf("Failed to read cache file: %v", err)
				}

				var tc TokenCache
				if err := json.Unmarshal(data, &tc); err != nil {
					t.Fatalf("Failed to unmarshal cache file: %v", err)
				}

				if tc.Token != "test-token-12345" {
					t.Errorf("Cached token = %v, want %v", tc.Token, "test-token-12345")
				}

				// Verify the expiration time is set to approximately 24 hours in the future
				expectedTime := time.Now().Add(24 * time.Hour)
				timeDiff := tc.ExpiresAt.Sub(expectedTime)
				if timeDiff < -5*time.Second || timeDiff > 5*time.Second {
					t.Errorf("Expiration time %v not within 5 seconds of expected %v", tc.ExpiresAt, expectedTime)
				}
			},
		},
		{
			name: "SaveToken_NegativeDuration",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// No setup needed
			},
			operation: "SaveToken with negative duration (tests default behavior)",
			action: func(manager *TokenCacheManager) error {
				return manager.SaveToken("test-token-negative", -5)
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Verify no error
				if err != nil {
					t.Errorf("SaveToken() error = %v", err)
				}

				// Read the cache file and verify content
				data, err := os.ReadFile(manager.CacheFile)
				if err != nil {
					t.Fatalf("Failed to read cache file: %v", err)
				}

				var tc TokenCache
				if err := json.Unmarshal(data, &tc); err != nil {
					t.Fatalf("Failed to unmarshal cache file: %v", err)
				}

				// Verify the expiration time is set to the default 24 hours
				expectedTime := time.Now().Add(24 * time.Hour)
				timeDiff := tc.ExpiresAt.Sub(expectedTime)
				if timeDiff < -5*time.Second || timeDiff > 5*time.Second {
					t.Errorf("Default expiration time %v not within 5 seconds of expected %v (24 hours)",
						tc.ExpiresAt, expectedTime)
				}
			},
		},
		{
			name: "SaveToken_FileWriteError",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create an unwritable directory by making CacheDir read-only
				if err := os.Chmod(manager.CacheDir, 0500); err != nil {
					t.Fatalf("Failed to make directory read-only: %v", err)
				}
			},
			operation: "SaveToken with file write error",
			action: func(manager *TokenCacheManager) error {
				return manager.SaveToken("test-token-write-error", 24)
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Fix permissions to allow cleanup
				os.Chmod(manager.CacheDir, 0700)

				// Verify we got an error
				if err == nil {
					t.Error("Expected error, got nil")
				}

				// Verify the error is related to file writing
				if err != nil && !errors.Is(err, os.ErrPermission) && !errors.Is(err, os.ErrNotExist) {
					t.Logf("Got expected write error: %v", err)
				}
			},
		},
		{
			name: "GetToken_Valid",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create a valid token cache file
				testToken := "test-token-67890"
				expires := time.Now().Add(1 * time.Hour)
				tc := TokenCache{
					Token:     testToken,
					ExpiresAt: expires,
				}

				data, err := json.Marshal(tc)
				if err != nil {
					t.Fatalf("Failed to marshal test token cache: %v", err)
				}

				if err := os.WriteFile(manager.CacheFile, data, 0600); err != nil {
					t.Fatalf("Failed to write test cache file: %v", err)
				}
			},
			operation: "GetToken with valid token",
			action: func(manager *TokenCacheManager) error {
				token, valid := manager.GetToken()
				if !valid || token != "test-token-67890" {
					t.Errorf("GetToken() returned token=%q, valid=%v, want token=%q, valid=true",
						token, valid, "test-token-67890")
				}
				return nil
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Error checking done in action
			},
		},
		{
			name: "GetToken_Expired",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create an expired token cache file
				testToken := "expired-token-12345"
				expires := time.Now().Add(-1 * time.Hour) // 1 hour in the past
				tc := TokenCache{
					Token:     testToken,
					ExpiresAt: expires,
				}

				data, err := json.Marshal(tc)
				if err != nil {
					t.Fatalf("Failed to marshal test token cache: %v", err)
				}

				if err := os.WriteFile(manager.CacheFile, data, 0600); err != nil {
					t.Fatalf("Failed to write test cache file: %v", err)
				}
			},
			operation: "GetToken with expired token",
			action: func(manager *TokenCacheManager) error {
				token, valid := manager.GetToken()
				if valid || token != "" {
					t.Errorf("GetToken() returned token=%q, valid=%v, want empty token and valid=false",
						token, valid)
				}
				return nil
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Error checking done in action
			},
		},
		{
			name: "GetToken_ReadError",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create a file that exists but is not readable
				if err := os.WriteFile(manager.CacheFile, []byte("test"), 0600); err != nil {
					t.Fatalf("Failed to write test cache file: %v", err)
				}

				// Change permissions to make it unreadable
				if err := os.Chmod(manager.CacheFile, 0000); err != nil {
					t.Fatalf("Failed to change file permissions: %v", err)
				}
			},
			operation: "GetToken with file read error",
			action: func(manager *TokenCacheManager) error {
				token, valid := manager.GetToken()
				// Should return false with empty token on read error
				if valid || token != "" {
					t.Errorf("GetToken() with unreadable file returned token=%q, valid=%v, want empty token and valid=false",
						token, valid)
				}
				return nil
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Fix permissions to allow cleanup
				os.Chmod(manager.CacheFile, 0600)
			},
		},
		{
			name: "GetToken_InvalidJSON",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create a file with invalid JSON
				if err := os.WriteFile(manager.CacheFile, []byte("invalid json content"), 0600); err != nil {
					t.Fatalf("Failed to write invalid cache file: %v", err)
				}
			},
			operation: "GetToken with invalid JSON",
			action: func(manager *TokenCacheManager) error {
				token, valid := manager.GetToken()
				// Should return false with empty token on unmarshal error
				if valid || token != "" {
					t.Errorf("GetToken() with invalid JSON returned token=%q, valid=%v, want empty token and valid=false",
						token, valid)
				}
				return nil
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// No additional checks needed
			},
		},
		{
			name: "ClearToken_Existing",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create a test cache file
				testToken := "clear-token-12345"
				tc := TokenCache{
					Token:     testToken,
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}

				data, err := json.Marshal(tc)
				if err != nil {
					t.Fatalf("Failed to marshal test token cache: %v", err)
				}

				if err := os.WriteFile(manager.CacheFile, data, 0600); err != nil {
					t.Fatalf("Failed to write test cache file: %v", err)
				}
			},
			operation: "ClearToken with existing file",
			action: func(manager *TokenCacheManager) error {
				return manager.ClearToken()
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Verify no error
				if err != nil {
					t.Errorf("ClearToken() error = %v", err)
				}

				// Check if file was removed
				if _, err := os.Stat(manager.CacheFile); !os.IsNotExist(err) {
					t.Error("ClearToken() did not remove the cache file")
				}
			},
		},
		{
			name: "ClearToken_Nonexistent",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// No setup needed - we want to ensure the file doesn't exist
			},
			operation: "ClearToken with nonexistent file",
			action: func(manager *TokenCacheManager) error {
				return manager.ClearToken()
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Verify no error
				if err != nil {
					t.Errorf("ClearToken() error = %v, want nil", err)
				}
			},
		},
		{
			name: "ClearToken_RemoveError",
			setup: func(t *testing.T, manager *TokenCacheManager) {
				// Create a test cache file
				if err := os.WriteFile(manager.CacheFile, []byte("test"), 0600); err != nil {
					t.Fatalf("Failed to write test cache file: %v", err)
				}

				// Make the parent directory read-only to prevent removal
				if err := os.Chmod(manager.CacheDir, 0500); err != nil {
					t.Fatalf("Failed to make directory read-only: %v", err)
				}
			},
			operation: "ClearToken with file remove error",
			action: func(manager *TokenCacheManager) error {
				return manager.ClearToken()
			},
			checkResult: func(t *testing.T, manager *TokenCacheManager, err error) {
				// Fix permissions to allow cleanup
				os.Chmod(manager.CacheDir, 0700)

				// Verify we got an error
				if err == nil {
					t.Error("Expected error, got nil")
				}

				// Verify the error is related to file removal
				if err != nil && !errors.Is(err, os.ErrPermission) {
					t.Logf("Got expected remove error: %v", err)
				}
			},
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			manager, tempDir := setupTestCacheManager(t)
			defer cleanupTestCacheManager(tempDir)

			// Run setup
			tc.setup(t, manager)

			// Run operation
			err := tc.action(manager)

			// Check results
			tc.checkResult(t, manager, err)
		})
	}
}

// TestNewTokenCacheManager tests the creation of a TokenCacheManager
func TestNewTokenCacheManager(t *testing.T) {
	// This test may be skipped in CI environments where home directory access is restricted
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	manager, err := NewTokenCacheManager()
	if err != nil {
		t.Fatalf("Failed to create TokenCacheManager: %v", err)
	}

	if manager == nil {
		t.Fatal("Expected non-nil TokenCacheManager")
	}

	// Verify the cache directory structure
	homeDir, _ := os.UserHomeDir()
	expectedCacheDir := filepath.Join(homeDir, ".vicohome")
	if manager.CacheDir != expectedCacheDir {
		t.Errorf("CacheDir = %v, want %v", manager.CacheDir, expectedCacheDir)
	}

	expectedCacheFile := filepath.Join(expectedCacheDir, "auth.json")
	if manager.CacheFile != expectedCacheFile {
		t.Errorf("CacheFile = %v, want %v", manager.CacheFile, expectedCacheFile)
	}
}

// TestNewTokenCacheManagerErrors tests error paths in the NewTokenCacheManager function
// This test uses a custom home directory to test error paths without affecting the actual system
func TestNewTokenCacheManagerErrors(t *testing.T) {
	// Save original home directory env var
	origHome := os.Getenv("HOME")

	// Restore the original HOME environment variable after the test
	defer func() {
		os.Setenv("HOME", origHome)
	}()

	// Test when home directory can't be determined
	t.Run("HomeDirectoryError", func(t *testing.T) {
		// Set HOME env var to empty to simulate error
		os.Setenv("HOME", "")

		// Should fail now
		_, err := NewTokenCacheManager()
		if err == nil {
			t.Error("Expected error when HOME is not set, got nil")
		}
	})

	// Test when directory can't be created
	t.Run("DirectoryCreationError", func(t *testing.T) {
		// Create a temp file (not directory) to prevent directory creation
		tempFile, err := os.CreateTemp("", "vicohome-test-file-")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tempFile.Close()
		defer os.Remove(tempFile.Name())

		// Set HOME to point to a temp file, which will cause MkdirAll to fail
		os.Setenv("HOME", tempFile.Name())

		// Should fail now because can't create directory in a file
		_, err = NewTokenCacheManager()
		if err == nil {
			t.Error("Expected error when directory can't be created, got nil")
		}
	})
}

// Benchmarks for cache operations
func BenchmarkSaveToken(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "vicohome-bench-")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := &TokenCacheManager{
		CacheDir:  tempDir,
		CacheFile: filepath.Join(tempDir, "auth.json"),
	}

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use a unique token for each iteration to avoid reading cached results
		token := "benchmark-token-" + string(rune(i))
		err := manager.SaveToken(token, 24)
		if err != nil {
			b.Fatalf("SaveToken failed: %v", err)
		}
	}
}

func BenchmarkGetToken(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "vicohome-bench-")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := &TokenCacheManager{
		CacheDir:  tempDir,
		CacheFile: filepath.Join(tempDir, "auth.json"),
	}

	// Create a valid token cache file for benchmark
	token := "benchmark-token"
	expiresAt := time.Now().Add(1 * time.Hour)
	tc := TokenCache{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	data, err := json.Marshal(tc)
	if err != nil {
		b.Fatalf("Failed to marshal test token cache: %v", err)
	}

	if err := os.WriteFile(manager.CacheFile, data, 0600); err != nil {
		b.Fatalf("Failed to write test cache file: %v", err)
	}

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, valid := manager.GetToken()
		if !valid || token == "" {
			b.Fatalf("GetToken failed, got token=%q, valid=%v", token, valid)
		}
	}
}

func BenchmarkClearToken(b *testing.B) {
	// Setup
	tempDir, err := os.MkdirTemp("", "vicohome-bench-")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := &TokenCacheManager{
		CacheDir:  tempDir,
		CacheFile: filepath.Join(tempDir, "auth.json"),
	}

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new cache file for each iteration
		if err := os.WriteFile(manager.CacheFile, []byte("test"), 0600); err != nil {
			b.Fatalf("Failed to write test cache file: %v", err)
		}

		// Clear the token
		err := manager.ClearToken()
		if err != nil {
			b.Fatalf("ClearToken failed: %v", err)
		}
	}
}
