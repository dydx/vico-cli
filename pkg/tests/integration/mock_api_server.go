package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// MockAPIServer represents the mock Vicohome API server for integration testing.
type MockAPIServer struct {
	// The underlying test server
	Server *httptest.Server

	// Configurations
	AuthTokens         map[string]bool         // Valid auth tokens
	AuthFailures       int                     // Counter for simulating auth failures
	TimeoutFailures    int                     // Counter for simulating timeout failures
	EventResponses     map[string]interface{}  // Predefined event responses
	DeviceResponses    map[string]interface{}  // Predefined device responses
	RequestLog         []map[string]interface{} // Log of received requests
	SimulateExpiration bool                     // Whether to simulate token expiration
}

// NewMockAPIServer creates a new mock API server for testing.
func NewMockAPIServer() *MockAPIServer {
	mock := &MockAPIServer{
		AuthTokens:         make(map[string]bool),
		EventResponses:     make(map[string]interface{}),
		DeviceResponses:    make(map[string]interface{}),
		RequestLog:         make([]map[string]interface{}, 0),
		SimulateExpiration: false,
	}

	// Create a test server that will handle our mock API responses
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request for later analysis
		requestInfo := map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"token":  r.Header.Get("Authorization"),
			"time":   time.Now(),
		}
		mock.RequestLog = append(mock.RequestLog, requestInfo)

		// Check for timeout simulation
		if mock.TimeoutFailures > 0 {
			mock.TimeoutFailures--
			time.Sleep(3 * time.Second) // Simulate a slow response
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		// Authenticate endpoint handler
		if r.URL.Path == "/account/login" {
			mock.handleLogin(w, r)
			return
		}

		// Check auth token (for all authenticated endpoints)
		token := r.Header.Get("Authorization")
		isValid := mock.AuthTokens[token]
		
		// Handle token expiration simulation
		if mock.SimulateExpiration && isValid {
			// Token has "expired"
			mock.AuthTokens[token] = false
			isValid = false
		}

		if !isValid {
			// Token is invalid or expired
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"result": -1025,
				"msg":    "Authentication token missing or invalid",
				"data":   nil,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Handle different API endpoints
		switch r.URL.Path {
		case "/library/newselectlibrary":
			mock.handleEvents(w, r)
		case "/device/listuserdevices":
			mock.handleDevices(w, r)
		default:
			// Unknown endpoint
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return mock
}

// handleLogin handles login requests and issues auth tokens.
func (m *MockAPIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var loginReq map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": -1,
			"msg":    "Invalid request format",
			"data":   nil,
		})
		return
	}

	// Check for auth failures simulation
	if m.AuthFailures > 0 {
		m.AuthFailures--
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": -1000,
			"msg":    "Invalid credentials",
			"data":   nil,
		})
		return
	}

	// Generate token
	token := fmt.Sprintf("mock-token-%d", time.Now().UnixNano())
	m.AuthTokens[token] = true

	// Return successful login response
	resp := map[string]interface{}{
		"result": 0,
		"msg":    "success",
		"data": map[string]interface{}{
			"token": map[string]interface{}{
				"token": token,
			},
		},
	}
	json.NewEncoder(w).Encode(resp)
}

// handleEvents handles event listing requests.
func (m *MockAPIServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body to extract time range
	var eventReq map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&eventReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": -1,
			"msg":    "Invalid request format",
			"data":   nil,
		})
		return
	}

	// Get predefined response or use default
	responseKey := "default"
	if startTimestamp, ok := eventReq["startTimestamp"].(string); ok {
		if endTimestamp, ok := eventReq["endTimestamp"].(string); ok {
			rangeKey := fmt.Sprintf("%s-%s", startTimestamp, endTimestamp)
			if resp, exists := m.EventResponses[rangeKey]; exists {
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}

	if resp, exists := m.EventResponses[responseKey]; exists {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Default response with empty events
	defaultResp := map[string]interface{}{
		"result": 0,
		"msg":    "success",
		"data": map[string]interface{}{
			"list": []interface{}{},
		},
	}
	json.NewEncoder(w).Encode(defaultResp)
}

// handleDevices handles device listing requests.
func (m *MockAPIServer) handleDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var deviceReq map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&deviceReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": -1,
			"msg":    "Invalid request format",
			"data":   nil,
		})
		return
	}

	// Get predefined response or use default
	responseKey := "default"
	if language, ok := deviceReq["language"].(string); ok {
		if countryNo, ok := deviceReq["countryNo"].(string); ok {
			requestKey := fmt.Sprintf("%s-%s", language, countryNo)
			if resp, exists := m.DeviceResponses[requestKey]; exists {
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}

	if resp, exists := m.DeviceResponses[responseKey]; exists {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Default response with empty devices
	defaultResp := map[string]interface{}{
		"result": 0,
		"msg":    "success",
		"data": map[string]interface{}{
			"list": []interface{}{},
		},
	}
	json.NewEncoder(w).Encode(defaultResp)
}

// AddEventResponse adds a predefined event response for a specific time range.
func (m *MockAPIServer) AddEventResponse(startTimestamp, endTimestamp string, response map[string]interface{}) {
	responseKey := fmt.Sprintf("%s-%s", startTimestamp, endTimestamp)
	m.EventResponses[responseKey] = response
}

// SetDefaultEventResponse sets the default event response.
func (m *MockAPIServer) SetDefaultEventResponse(response map[string]interface{}) {
	m.EventResponses["default"] = response
}

// AddDeviceResponse adds a predefined device response for a specific language and country.
func (m *MockAPIServer) AddDeviceResponse(language, countryNo string, response map[string]interface{}) {
	responseKey := fmt.Sprintf("%s-%s", language, countryNo)
	m.DeviceResponses[responseKey] = response
}

// SetDefaultDeviceResponse sets the default device response.
func (m *MockAPIServer) SetDefaultDeviceResponse(response map[string]interface{}) {
	m.DeviceResponses["default"] = response
}

// GetRequestCount returns the number of requests received.
func (m *MockAPIServer) GetRequestCount() int {
	return len(m.RequestLog)
}

// GetLastRequest returns the last request received.
func (m *MockAPIServer) GetLastRequest() map[string]interface{} {
	if len(m.RequestLog) == 0 {
		return nil
	}
	return m.RequestLog[len(m.RequestLog)-1]
}

// Close shuts down the mock server.
func (m *MockAPIServer) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}