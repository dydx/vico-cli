// Package auth provides authentication functionality for the Vicohome API.
//
// This package handles user authentication, token management, and automated token
// refreshing when needed. It uses environment variables for credentials and supports
// token caching to minimize authentication requests.
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dydx/vico-cli/pkg/cache"
)

// Error codes from the API
const (
	ErrorAccountKicked = -1024
	ErrorTokenMissing  = -1025
)

// LoginRequest represents the JSON request body sent to the Vicohome API
// during authentication.
type LoginRequest struct {
	Email     string `json:"email"`     // User's email address
	Password  string `json:"password"`  // User's password
	LoginType int    `json:"loginType"` // Type of login (0 for standard login)
}

// LoginResponse represents the JSON response body received from the Vicohome API
// after an authentication request. It contains the authentication token and status information.
type LoginResponse struct {
	Result int    `json:"result"` // Result code (0 for success)
	Msg    string `json:"msg"`    // Response message
	Data   struct {
		Token struct {
			Token string `json:"token"` // Authentication token
		} `json:"token"`
	} `json:"data"`
}

// AuthenticateFunc is the function type for authentication to allow mocking in tests
type AuthenticateFunc func() (string, error)

// Authenticate is the main authentication function that can be replaced in tests
var Authenticate AuthenticateFunc = authenticate

// HTTPClient is the client used for API requests, which can be replaced in tests
var HTTPClient *http.Client = &http.Client{}

// MockAuthenticate sets up the global Authenticate function to return a static token and error.
// Returns a cleanup function to restore the original function.
func MockAuthenticate(token string, err error) func() {
	original := Authenticate

	mockFunc := func() (string, error) {
		return token, err
	}

	Authenticate = mockFunc

	return func() {
		Authenticate = original
	}
}

// authenticate obtains an authentication token for the Vicohome API.
// It first tries to retrieve a valid cached token. If no valid token is found,
// it falls back to direct authentication using credentials from environment variables.
// Successfully acquired tokens are cached for future use to minimize authentication requests.
//
// Returns:
//   - string: The authentication token if successful
//   - error: Any error encountered during the authentication process
func authenticate() (string, error) {
	// Try to get a cached token first
	cacheManager, err := cache.NewTokenCacheManager()
	if err != nil {
		// If we can't create a cache manager, fall back to direct authentication
		return authenticateDirectly()
	}

	token, valid := cacheManager.GetToken()
	if valid {
		// We have a valid cached token, return it
		return token, nil
	}

	// No valid cached token, authenticate and cache the new token
	token, err = authenticateDirectly()
	if err != nil {
		return "", err
	}

	// Cache the token for future use (24 hours validity)
	if err := cacheManager.SaveToken(token, 24); err != nil {
		// Non-fatal error, we can still return the token
		fmt.Fprintf(os.Stderr, "Warning: failed to cache token: %v\n", err)
	}

	return token, nil
}

// authenticateDirectly performs authentication to the Vicohome API without using the token cache.
// It retrieves credentials from environment variables (VICOHOME_EMAIL and VICOHOME_PASSWORD),
// makes an authentication request to the API, and parses the response to extract the token.
//
// Returns:
//   - string: The authentication token if successful
//   - error: Any error encountered during the authentication process
func authenticateDirectly() (string, error) {
	// Get credentials from environment variables
	email := os.Getenv("VICOHOME_EMAIL")
	password := os.Getenv("VICOHOME_PASSWORD")

	// Check if credentials are available
	if email == "" || password == "" {
		return "", fmt.Errorf("Error: VICOHOME_EMAIL and VICOHOME_PASSWORD environment variables are required")
	}
	// Use the proper JSON marshaling to avoid escaping issues
	loginReq := map[string]interface{}{
		"email":     email,
		"password":  password,
		"loginType": 0,
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling login request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-us.vicohome.io/account/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Try to parse as generic map first to handle all possible response formats
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check if there's a result code and error message in the API response
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		return "", fmt.Errorf("API error: %s (code: %.0f)", msg, result)
	}

	// Check if we have data.token.token in the response
	data, ok := responseMap["data"].(map[string]interface{})
	if !ok || len(data) == 0 {
		return "", fmt.Errorf("login failed: missing data in response")
	}

	tokenObj, ok := data["token"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("login failed: missing token in response")
	}

	tokenStr, ok := tokenObj["token"].(string)
	if !ok || tokenStr == "" {
		return "", fmt.Errorf("login failed: empty token in response")
	}

	return tokenStr, nil
}

// APIError represents an error from the Vicohome API.
type APIError struct {
	Code    string
	Message string
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	return e.Message
}

// ValidateResponseFunc is a function type for the ValidateResponse function
type ValidateResponseFunc func(respBody []byte) (bool, error)

// ValidateResponse is the implementation of ValidateResponseFunc that can be replaced in tests
var ValidateResponse ValidateResponseFunc = validateResponse

// validateResponse checks if an API response contains an authentication error
// and determines if the token needs to be refreshed. It analyzes the response body
// for specific error codes that indicate authentication issues. If such an error
// is found, it clears the token cache to force a new authentication.
//
// Parameters:
//   - respBody: The raw response body from the API
//
// Returns:
//   - bool: True if the token needs to be refreshed, false otherwise
//   - error: Any error found in the response, or nil if no error was found
func validateResponse(respBody []byte) (bool, error) {
	// Try to parse the response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return false, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Check for authentication errors
	if result, ok := responseMap["result"].(float64); ok {
		msg, _ := responseMap["msg"].(string)

		// Check if we need to refresh the token
		if result == ErrorAccountKicked || result == ErrorTokenMissing {
			// Clear the cache
			cacheManager, err := cache.NewTokenCacheManager()
			if err == nil {
				cacheManager.ClearToken()
			}
			return true, fmt.Errorf("authentication error: %s (code: %.0f)", msg, result)
		}

		// Check for other API errors
		if result != 0 {
			return false, fmt.Errorf("API error: %s (code: %.0f)", msg, result)
		}
	}

	return false, nil
}

// ExecuteWithRetryFunc is a function type for the ExecuteWithRetry function
type ExecuteWithRetryFunc func(req *http.Request) ([]byte, error)

// ExecuteWithRetry is the implementation of ExecuteWithRetryFunc that can be replaced in tests
var ExecuteWithRetry ExecuteWithRetryFunc = executeWithRetry

// executeWithRetry executes an HTTP request with automatic token refresh on authentication errors.
// If the initial request fails due to an authentication error, it refreshes the token and
// retries the request once with the new token. This handles cases where a token has expired
// or been invalidated since it was cached.
//
// Parameters:
//   - req: The HTTP request to execute
//
// Returns:
//   - []byte: The response body if successful
//   - error: Any error encountered during the request process
func executeWithRetry(req *http.Request) ([]byte, error) {
	// First attempt with current token
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check if we need to refresh the token and if there are any errors
	needsRefresh, validateErr := validateResponse(respBody)

	// If there's an API error (not an auth error), return it immediately
	if validateErr != nil && !needsRefresh {
		return nil, validateErr
	}

	if needsRefresh {
		// Clear the cache and get a new token
		cacheManager, err := cache.NewTokenCacheManager()
		if err == nil {
			cacheManager.ClearToken()
		}

		// Get a new token
		token, err := authenticateDirectly()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Cache the new token
		if cacheManager != nil {
			cacheManager.SaveToken(token, 24)
		}

		// Update the request with the new token
		req.Header.Set("Authorization", token)

		// Retry the request
		resp, err = HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request after token refresh: %w", err)
		}
		defer resp.Body.Close()

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body after token refresh: %w", err)
		}

		// Validate the response again after refresh
		_, validateErr = validateResponse(respBody)
		if validateErr != nil {
			return nil, validateErr
		}
	}

	return respBody, nil
}
