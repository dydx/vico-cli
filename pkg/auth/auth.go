package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
func Authenticate() (string, error) {
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