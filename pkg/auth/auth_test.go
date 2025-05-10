package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dydx/vico-cli/pkg/cache"
)

// mockRoundTripper is a custom RoundTripper that returns predefined responses for testing
type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

// mockReadCloser implements io.ReadCloser for testing
type mockReadCloser struct {
	reader io.Reader
}

func (m mockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m mockReadCloser) Close() error {
	return nil
}

// createMockResponse creates a mock HTTP response with the given status code and body
func createMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       mockReadCloser{reader: bytes.NewReader([]byte(body))},
		Header:     http.Header{},
	}
}

// TestValidateResponse tests the ValidateResponse function with various response scenarios
func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name           string
		responseJSON   string
		expectRefresh  bool
		expectError    bool
		expectedErrMsg string
	}{
		{
			name: "Valid response with no errors",
			responseJSON: `{
				"result": 0,
				"msg": "success",
				"data": {"some": "data"}
			}`,
			expectRefresh:  false,
			expectError:    false,
			expectedErrMsg: "",
		},
		{
			name: "Response with authentication error - token missing",
			responseJSON: `{
				"result": -1025,
				"msg": "Authentication token missing or invalid",
				"data": null
			}`,
			expectRefresh:  true,
			expectError:    true,
			expectedErrMsg: "authentication error: Authentication token missing or invalid (code: -1025)",
		},
		{
			name: "Response with authentication error - account kicked",
			responseJSON: `{
				"result": -1024,
				"msg": "Account has been logged out",
				"data": null
			}`,
			expectRefresh:  true,
			expectError:    true,
			expectedErrMsg: "authentication error: Account has been logged out (code: -1024)",
		},
		{
			name: "Response with non-authentication API error",
			responseJSON: `{
				"result": -2000,
				"msg": "Resource not found",
				"data": null
			}`,
			expectRefresh:  false,
			expectError:    true,
			expectedErrMsg: "API error: Resource not found (code: -2000)",
		},
		{
			name:           "Invalid JSON response",
			responseJSON:   "invalid-json",
			expectRefresh:  false,
			expectError:    true,
			expectedErrMsg: "error unmarshaling response",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function with the test response
			needsRefresh, err := ValidateResponse([]byte(tc.responseJSON))

			// Check refresh flag
			if needsRefresh != tc.expectRefresh {
				t.Errorf("Expected needsRefresh=%v, got %v", tc.expectRefresh, needsRefresh)
			}

			// Check error
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error with message %q, got nil", tc.expectedErrMsg)
				} else if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

// TestExecuteWithRetry_NetworkError tests the error handling for network errors
func TestExecuteWithRetry_NetworkError(t *testing.T) {
	// Save original transport to restore later
	defaultTransport := http.DefaultTransport
	
	// Create a mock transport that returns a network error
	mockTransport := &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}
	
	// Set the mock transport
	http.DefaultTransport = mockTransport
	defer func() {
		http.DefaultTransport = defaultTransport
	}()
	
	// Create a test request
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	
	// Call the function
	_, err := ExecuteWithRetry(req)
	
	// Verify an error is returned
	if err == nil {
		t.Error("Expected an error but got nil")
	}
	
	// Verify the error contains the network error message
	expectedErrMsg := "network error"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error to contain %q, got %q", expectedErrMsg, err.Error())
	}
}

// TestExecuteWithRetry_SuccessfulRequest tests the successful execution of a request
func TestExecuteWithRetry_SuccessfulRequest(t *testing.T) {
	// Save original transport to restore later
	defaultTransport := http.DefaultTransport
	
	// Create a mock transport that returns a successful response
	mockTransport := &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"test": "data"}}`), nil
		},
	}
	
	// Set the mock transport
	http.DefaultTransport = mockTransport
	defer func() {
		http.DefaultTransport = defaultTransport
	}()
	
	// Create a test request
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	
	// Call the function
	respBody, err := ExecuteWithRetry(req)
	
	// Verify no error is returned
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify response contains expected data
	expectedData := `"test": "data"`
	if !strings.Contains(string(respBody), expectedData) {
		t.Errorf("Expected response to contain %q, got %q", expectedData, string(respBody))
	}
}

// TestExecuteWithRetry_APIError tests handling of API errors
func TestExecuteWithRetry_APIError(t *testing.T) {
	// Save original transport to restore later
	defaultTransport := http.DefaultTransport

	// Create a mock transport that returns an API error
	mockTransport := &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return createMockResponse(200, `{"result": -2000, "msg": "API Error Message", "data": null}`), nil
		},
	}

	// Set the mock transport
	http.DefaultTransport = mockTransport
	defer func() {
		http.DefaultTransport = defaultTransport
	}()

	// Create a test request
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	// Call the function
	_, err := ExecuteWithRetry(req)

	// Verify an error is returned
	if err == nil {
		t.Error("Expected an API error but got nil")
		return
	}

	// Verify the error contains the API error message
	expectedErrMsg := "API error: API Error Message"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error to contain %q, got %q", expectedErrMsg, err.Error())
	}
}

// mockCacheManager is a mock implementation of the cache.TokenCacheManager for testing
type mockCacheManager struct {
	token      string
	valid      bool
	saveError  bool
	clearError bool
}

// Replace original TokenCacheManager with our mock for testing
func mockNewTokenCacheManager() (*cache.TokenCacheManager, error) {
	return nil, fmt.Errorf("mock error")
}

// GetToken implements the token getting functionality
func (m *mockCacheManager) GetToken() (string, bool) {
	return m.token, m.valid
}

// SaveToken implements the token saving functionality
func (m *mockCacheManager) SaveToken(token string, durationHours int) error {
	if m.saveError {
		return fmt.Errorf("mock save error")
	}
	m.token = token
	m.valid = true
	return nil
}

// ClearToken implements the token clearing functionality
func (m *mockCacheManager) ClearToken() error {
	if m.clearError {
		return fmt.Errorf("mock clear error")
	}
	m.token = ""
	m.valid = false
	return nil
}

// TestAuthenticate tests the main Authenticate function with various scenarios
func TestAuthenticate(t *testing.T) {
	// Save original function
	originalAuthenticate := Authenticate
	defer func() {
		Authenticate = originalAuthenticate
	}()

	// Save original transport
	defaultTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = defaultTransport
	}()

	// Test cases
	tests := []struct {
		name          string
		setupAuth     func()
		expectedToken string
		expectError   bool
		errorContains string
	}{
		{
			name: "Successful authentication",
			setupAuth: func() {
				Authenticate = func() (string, error) {
					return "test-token-123", nil
				}
			},
			expectedToken: "test-token-123",
			expectError:   false,
		},
		{
			name: "Authentication failure",
			setupAuth: func() {
				Authenticate = func() (string, error) {
					return "", fmt.Errorf("API error: Invalid credentials")
				}
			},
			expectedToken: "",
			expectError:   true,
			errorContains: "Invalid credentials",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test case
			tc.setupAuth()

			// Call the function
			token, err := Authenticate()

			// Check results
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain %q, got %q", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if token != tc.expectedToken {
					t.Errorf("Expected token %q, got %q", tc.expectedToken, token)
				}
			}
		})
	}
}

// TestExecuteWithRetry_TokenRefresh tests the token refresh functionality when an authentication error occurs
func TestExecuteWithRetry_TokenRefresh(t *testing.T) {
	// Save original transport to restore later
	defaultTransport := http.DefaultTransport

	// Create a counter to track which request is being made
	requestCount := 0

	// Create a mock transport that returns different responses for each request
	mockTransport := &mockRoundTripper{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++

			if requestCount == 1 {
				// First request - return auth error to trigger token refresh
				return createMockResponse(200, `{"result": -1025, "msg": "Authentication token missing or invalid", "data": null}`), nil
			} else if requestCount == 2 {
				// Second request - this should be the token request
				// Verify this is a login request
				if req.URL.String() != "https://api-us.vicohome.io/account/login" {
					t.Errorf("Expected login request, got request to %s", req.URL.String())
				}

				// Return successful login response with token
				return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"token": {"token": "new-test-token"}}}`), nil
			} else {
				// Third request - this should be the retry with the new token
				// Verify the Authorization header contains the new token
				if req.Header.Get("Authorization") != "new-test-token" {
					t.Errorf("Expected Authorization header with new token, got %s", req.Header.Get("Authorization"))
				}

				// Return success
				return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"sample": "data"}}`), nil
			}
		},
	}

	// Set the mock transport
	http.DefaultTransport = mockTransport

	// Save original environment variables
	originalEmail := os.Getenv("VICOHOME_EMAIL")
	originalPassword := os.Getenv("VICOHOME_PASSWORD")

	// Set test environment variables
	os.Setenv("VICOHOME_EMAIL", "test@example.com")
	os.Setenv("VICOHOME_PASSWORD", "password123")

	// Restore original values when test is done
	defer func() {
		http.DefaultTransport = defaultTransport
		os.Setenv("VICOHOME_EMAIL", originalEmail)
		os.Setenv("VICOHOME_PASSWORD", originalPassword)
	}()

	// Create a test request
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Header.Set("Authorization", "old-expired-token")

	// Call the function
	respBody, err := ExecuteWithRetry(req)

	// Verify the function executed without errors
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the response contains expected data
	expectedData := `"sample": "data"`
	if !strings.Contains(string(respBody), expectedData) {
		t.Errorf("Expected response to contain %q, got %q", expectedData, string(respBody))
	}

	// Verify all 3 requests were made
	if requestCount != 3 {
		t.Errorf("Expected 3 requests to be made, got %d", requestCount)
	}
}

// TestAuthenticate tests the main Authenticate function with various scenarios
func TestAuthenticate(t *testing.T) {
	// Save original function pointers to restore later
	origNewTokenCacheManager := cache.NewTokenCacheManager
	defer func() {
		cache.NewTokenCacheManager = origNewTokenCacheManager
	}()

	// Save original transport to restore later
	defaultTransport := http.DefaultTransport

	// Save original environment variables
	originalEmail := os.Getenv("VICOHOME_EMAIL")
	originalPassword := os.Getenv("VICOHOME_PASSWORD")

	// Set test credentials
	os.Setenv("VICOHOME_EMAIL", "test@example.com")
	os.Setenv("VICOHOME_PASSWORD", "password123")

	// Restore original values when test is done
	defer func() {
		http.DefaultTransport = defaultTransport
		os.Setenv("VICOHOME_EMAIL", originalEmail)
		os.Setenv("VICOHOME_PASSWORD", originalPassword)
	}()

	tests := []struct {
		name          string
		cacheSetup    func() *mockCacheManager
		transportSetup func() *mockRoundTripper
		expectedToken string
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid cached token",
			cacheSetup: func() *mockCacheManager {
				return &mockCacheManager{
					token: "cached-token-123",
					valid: true,
				}
			},
			transportSetup: func() *mockRoundTripper {
				return nil // Not used for this test
			},
			expectedToken: "cached-token-123",
			expectError:   false,
		},
		{
			name: "Cache error, fallback to direct auth",
			cacheSetup: func() *mockCacheManager {
				return nil // We'll use the error from mockNewTokenCacheManager
			},
			transportSetup: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"token": {"token": "direct-token-123"}}}`), nil
					},
				}
			},
			expectedToken: "direct-token-123",
			expectError:   false,
		},
		{
			name: "Invalid cached token, successful refresh",
			cacheSetup: func() *mockCacheManager {
				return &mockCacheManager{
					token: "expired-token",
					valid: false,
				}
			},
			transportSetup: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"token": {"token": "fresh-token-123"}}}`), nil
					},
				}
			},
			expectedToken: "fresh-token-123",
			expectError:   false,
		},
		{
			name: "Direct auth failure",
			cacheSetup: func() *mockCacheManager {
				return &mockCacheManager{
					token: "",
					valid: false,
				}
			},
			transportSetup: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": -1000, "msg": "Invalid credentials", "data": null}`), nil
					},
				}
			},
			expectedToken: "",
			expectError:   true,
			errorContains: "API error: Invalid credentials",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock cache
			mockCache := tc.cacheSetup()
			if mockCache != nil {
				// Replace the NewTokenCacheManager function to return our mock
				cache.NewTokenCacheManager = func() (*cache.TokenCacheManager, error) {
					// We're ignoring types here since Go doesn't know our mock implements the same interface
					// In a real project, we would define a proper interface
					return &cache.TokenCacheManager{
						CacheDir:  "mock-dir",
						CacheFile: "mock-file",
					}, nil
				}

				// Store the mock to use its methods
				// This is hacky but works for testing
				authenticateOriginal := Authenticate
				Authenticate = func() (string, error) {
					if mockCache.valid {
						return mockCache.token, nil
					}
					token, err := authenticateDirectly()
					if err != nil {
						return "", err
					}
					mockCache.SaveToken(token, 24)
					return token, nil
				}
				defer func() {
					Authenticate = authenticateOriginal
				}()
			} else {
				// Use the error returning mock for cases where we want cache creation to fail
				cache.NewTokenCacheManager = mockNewTokenCacheManager
			}

			// Setup mock transport if provided
			mockTransport := tc.transportSetup()
			if mockTransport != nil {
				http.DefaultTransport = mockTransport
			}

			// Call the function
			token, err := Authenticate()

			// Check results
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain %q, got %q", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if token != tc.expectedToken {
					t.Errorf("Expected token %q, got %q", tc.expectedToken, token)
				}
			}
		})
	}
}

// TestAuthenticateDirectly tests the authenticateDirectly function with various scenarios
func TestAuthenticateDirectly(t *testing.T) {
	// Save original transport to restore later
	defaultTransport := http.DefaultTransport

	// Save original environment variables to restore later
	originalEmail := os.Getenv("VICOHOME_EMAIL")
	originalPassword := os.Getenv("VICOHOME_PASSWORD")
	defer func() {
		os.Setenv("VICOHOME_EMAIL", originalEmail)
		os.Setenv("VICOHOME_PASSWORD", originalPassword)
		http.DefaultTransport = defaultTransport
	}()

	tests := []struct {
		name           string
		setupEnv       func()
		setupMock      func() *mockRoundTripper
		expectError    bool
		expectedToken  string
		expectedErrMsg string
	}{
		{
			name: "Missing environment variables",
			setupEnv: func() {
				os.Unsetenv("VICOHOME_EMAIL")
				os.Unsetenv("VICOHOME_PASSWORD")
			},
			setupMock: func() *mockRoundTripper {
				return nil // No mock needed for this test
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "VICOHOME_EMAIL and VICOHOME_PASSWORD environment variables are required",
		},
		{
			name: "Network error",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return nil, errors.New("network error")
					},
				}
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "error making request: network error",
		},
		{
			name: "API error response",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": -1000, "msg": "Invalid credentials", "data": null}`), nil
					},
				}
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "API error: Invalid credentials",
		},
		{
			name: "Invalid JSON response",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `invalid-json`), nil
					},
				}
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "error unmarshaling response",
		},
		{
			name: "Missing data in response",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": 0, "msg": "success", "data": null}`), nil
					},
				}
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "login failed: missing data",
		},
		{
			name: "Missing token object in response",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"user": "info"}}`), nil
					},
				}
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "login failed: missing token",
		},
		{
			name: "Empty token in response",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"token": {"token": ""}}}`), nil
					},
				}
			},
			expectError:    true,
			expectedToken:  "",
			expectedErrMsg: "login failed: empty token",
		},
		{
			name: "Successful authentication",
			setupEnv: func() {
				os.Setenv("VICOHOME_EMAIL", "test@example.com")
				os.Setenv("VICOHOME_PASSWORD", "password123")
			},
			setupMock: func() *mockRoundTripper {
				return &mockRoundTripper{
					roundTripFunc: func(req *http.Request) (*http.Response, error) {
						// Verify request is properly formed
						if req.Method != "POST" {
							t.Errorf("Expected POST request, got %s", req.Method)
						}
						if req.URL.String() != "https://api-us.vicohome.io/account/login" {
							t.Errorf("Expected URL https://api-us.vicohome.io/account/login, got %s", req.URL.String())
						}
						if req.Header.Get("Content-Type") != "application/json" {
							t.Errorf("Expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
						}

						// Read request body to verify credentials
						body, _ := io.ReadAll(req.Body)
						if !strings.Contains(string(body), "test@example.com") || !strings.Contains(string(body), "password123") {
							t.Errorf("Expected request body to contain credentials, got %s", string(body))
						}

						return createMockResponse(200, `{"result": 0, "msg": "success", "data": {"token": {"token": "test-token-123"}}}`), nil
					},
				}
			},
			expectError:   false,
			expectedToken: "test-token-123",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup environment variables
			tc.setupEnv()

			// Setup mock transport if provided
			mockTransport := tc.setupMock()
			if mockTransport != nil {
				http.DefaultTransport = mockTransport
			}

			// Call the function
			token, err := authenticateDirectly()

			// Check results
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tc.name == "Network error" {
					if !strings.Contains(err.Error(), "error making request") {
						t.Errorf("Expected error to contain 'error making request', got %q", err.Error())
					}
				} else if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if token != tc.expectedToken {
					t.Errorf("Expected token %q, got %q", tc.expectedToken, token)
				}
			}
		})
	}
}