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
	"strings"

	"github.com/dydx/vico-cli/pkg/cache"
)

// Error codes from the API
const (
	ErrorAccountKicked = -1024 // Account has been kicked offline
	ErrorTokenMissing  = -1025 // Token is missing
	ErrorTokenInvalid  = -1026 // Token is invalid
	ErrorTokenExpired  = -1027 // Token has expired
)

// isDebugMode returns true if debug logging is enabled
func isDebugMode() bool {
	debug := os.Getenv("VICOHOME_DEBUG")
	return debug == "true" || debug == "1"
}

// logDebug prints a message only if debug mode is enabled
func logDebug(format string, args ...interface{}) {
	if isDebugMode() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

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

// Authenticate obtains an authentication token for the Vicohome API.
// It first tries to retrieve a valid cached token. If no valid token is found,
// it falls back to direct authentication using credentials from environment variables.
// Successfully acquired tokens are cached for future use to minimize authentication requests.
//
// Returns:
//   - string: The authentication token if successful
//   - error: Any error encountered during the authentication process
func Authenticate() (string, error) {
	// Try to get a cached token first
	cacheManager, err := cache.NewTokenCacheManager()
	if err != nil {
		// If we can't create a cache manager, fall back to direct authentication
		logDebug("Warning: Could not create token cache manager: %v\n", err)
		return authenticateDirectly()
	}

	token, valid := cacheManager.GetToken()
	if valid {
		logDebug("Using cached token\n")
		// We have a valid cached token, return it
		return token, nil
	}

	// No valid cached token, authenticate and cache the new token
	logDebug("No valid cached token found, authenticating directly\n")
	token, err = authenticateDirectly()
	if err != nil {
		return "", err
	}

	// Cache the token for future use (24 hours validity)
	if err := cacheManager.SaveToken(token, 24); err != nil {
		// Non-fatal error, we can still return the token
		logDebug("Warning: failed to cache token: %v\n", err)
	} else {
		logDebug("Successfully cached new token\n")
	}

	return token, nil
}
func GetBaseURL() string {
if v := os.Getenv("VICOHOME_BASE_URL"); v != "" {
return v
}
return "https://api-us.vicohome.io"
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

	req, err := http.NewRequest("POST", GetBaseURL()+"/account/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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

// ValidateResponse checks if an API response contains an authentication error
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
func ValidateResponse(respBody []byte) (bool, error) {
	// Check if we have a non-JSON response (probably HTML error page)
	if len(respBody) > 0 && (respBody[0] == '<' || respBody[0] == '\r' || respBody[0] == '\n') {
		// Print a preview for debugging
		if isDebugMode() {
			preview := string(respBody)
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			logDebug("Warning: Received non-JSON response (likely auth issue): %s\n", preview)
		}
		return true, fmt.Errorf("received non-JSON response (likely authentication issue)")
	}

	// Try to parse the response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return false, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Print the response for debugging
	if isDebugMode() {
		prettyJSON, _ := json.MarshalIndent(responseMap, "", "  ")
		logDebug("API Response: %s\n", string(prettyJSON))
	}

	// Check for authentication errors - try both result and code fields (API inconsistency)
	var errorCode float64
	var errorMsg string
	var hasError bool

	// Check the "result" field first (main API)
	if result, ok := responseMap["result"].(float64); ok {
		errorCode = result
		msg, _ := responseMap["msg"].(string)
		errorMsg = msg
		hasError = result != 0
	}

	// Also check the "code" field (some endpoints use this instead)
	if code, ok := responseMap["code"].(float64); ok {
		errorCode = code
		msg, _ := responseMap["msg"].(string)
		errorMsg = msg
		hasError = code != 0
	}

	// If we found an error code
	if hasError {
		// Check if it's an auth error requiring token refresh
		isAuthError := errorCode == ErrorAccountKicked ||
			errorCode == ErrorTokenMissing ||
			errorCode == ErrorTokenInvalid ||
			errorCode == ErrorTokenExpired

		// Also check message strings for auth-related errors
		if !isAuthError && errorMsg != "" {
			errorMsgLower := strings.ToLower(errorMsg)
			isAuthError = strings.Contains(errorMsgLower, "token") ||
				strings.Contains(errorMsgLower, "auth") ||
				strings.Contains(errorMsgLower, "login") ||
				strings.Contains(errorMsgLower, "账号") || // account in Chinese
				strings.Contains(errorMsgLower, "踢下线") // kicked offline in Chinese
		}

		if isAuthError {
			if isDebugMode() {
				logDebug("Auth error detected: %s (code: %.0f)\n", errorMsg, errorCode)
			}
			// Don't clear cache here, let the caller handle it
			return true, fmt.Errorf("authentication error: %s (code: %.0f)", errorMsg, errorCode)
		}

		// Otherwise it's a regular API error
		return false, fmt.Errorf("API error: %s (code: %.0f)", errorMsg, errorCode)
	}

	return false, nil
}

// ExecuteWithRetry executes an HTTP request with automatic token refresh on authentication errors.
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
func ExecuteWithRetry(req *http.Request) ([]byte, error) {
	// First attempt with current token
	client := &http.Client{}

	// Make sure we can reuse the request body if needed
	var requestBodyBytes []byte
	if req.Body != nil {
		var err error
		requestBodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}

	// First attempt
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check if we need to refresh the token
	needsRefresh, apiErr := ValidateResponse(respBody)
	if needsRefresh {
		// Only show detailed logs in debug mode
		logDebug("Token refresh needed: %v\n", apiErr)

		// Clear the cache and get a new token
		cacheManager, err := cache.NewTokenCacheManager()
		if err == nil {
			cacheManager.ClearToken()
		}

		// Get a new token directly (bypass cache)
		token, err := authenticateDirectly()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Cache the new token with verbose error handling
		if cacheManager != nil {
			if err := cacheManager.SaveToken(token, 24); err != nil {
				logDebug("Warning: failed to cache refreshed token: %v\n", err)
			} else {
				logDebug("Successfully refreshed and cached new token\n")
			}
		}

		// Create a new request with the same parameters but new token
		newReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request for retry: %w", err)
		}

		// Copy all headers from the original request
		for key, values := range req.Header {
			for _, value := range values {
				newReq.Header.Add(key, value)
			}
		}

		// Set the new Authorization header
		newReq.Header.Set("Authorization", token)

		// Add back the body if there was one
		if len(requestBodyBytes) > 0 {
			newReq.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
			newReq.ContentLength = int64(len(requestBodyBytes))
		}

		// Retry the request with the new token
		logDebug("Retrying request with refreshed token\n")
		resp, err = client.Do(newReq)
		if err != nil {
			return nil, fmt.Errorf("error making request after token refresh: %w", err)
		}
		defer resp.Body.Close()

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body after token refresh: %w", err)
		}

		// Check if we still got an error after refreshing the token
		needsRefresh2, apiError := ValidateResponse(respBody)
		if needsRefresh2 || apiError != nil {
			// This is a critical error worth showing - authentication failed even after token refresh
			return nil, fmt.Errorf("authentication failed even after token refresh: %v", apiError)
		}
	}

	return respBody, nil
}
