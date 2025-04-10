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

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginType int    `json:"loginType"`
}

type LoginResponse struct {
	Result  int    `json:"result"`
	Msg     string `json:"msg"`
	Data    struct {
		Token struct {
			Token string `json:"token"`
		} `json:"token"`
	} `json:"data"`
}

// Authenticate gets the token for Vicohome API
// It tries to use a cached token first, then falls back to authentication if needed
func Authenticate() (string, error) {
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

// authenticateDirectly performs direct authentication to the API without using the cache
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

// ValidateResponse checks if a response contains an authentication error
// and returns true if the token needs to be refreshed
func ValidateResponse(respBody []byte) (bool, error) {
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

// ExecuteWithRetry executes an HTTP request with automatic token refresh on auth errors
func ExecuteWithRetry(req *http.Request) ([]byte, error) {
	// First attempt with current token
	client := &http.Client{}
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
	needsRefresh, _ := ValidateResponse(respBody)
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
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request after token refresh: %w", err)
		}
		defer resp.Body.Close()

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body after token refresh: %w", err)
		}
	}

	return respBody, nil
}