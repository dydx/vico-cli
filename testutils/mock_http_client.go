package testutils

import (
	"bytes"
	"io"
	"net/http"
)

// MockTransport implements the http.RoundTripper interface
// for mocking HTTP requests during testing.
type MockTransport struct {
	Responses map[string]MockResponse
	Requests  map[string]*http.Request // Stores requests for later inspection
}

// MockResponse represents a mock HTTP response for testing.
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// RoundTrip implements the http.RoundTripper interface.
// It returns a mocked response based on the request method and URL.
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()
	
	// Save request for later inspection
	if m.Requests == nil {
		m.Requests = make(map[string]*http.Request)
	}
	
	// Create a copy of the request body for saving
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
		
		// Create a new reader for downstream handling
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	
	// Clone the request for inspection
	reqCopy, _ := http.NewRequest(req.Method, req.URL.String(), io.NopCloser(bytes.NewBuffer(bodyBytes)))
	for k, v := range req.Header {
		for _, vv := range v {
			reqCopy.Header.Add(k, vv)
		}
	}
	m.Requests[key] = reqCopy
	
	// Get mock response for this request
	mockResp, exists := m.Responses[key]
	if !exists {
		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString("Not Found: " + key)),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	}
	
	// Create response headers
	header := http.Header{}
	for k, v := range mockResp.Headers {
		header.Add(k, v)
	}
	
	// Return mocked response
	return &http.Response{
		StatusCode: mockResp.StatusCode,
		Body:       io.NopCloser(bytes.NewBufferString(mockResp.Body)),
		Header:     header,
		Request:    req,
	}, nil
}

// NewMockClient creates an http.Client with the mock transport.
func NewMockClient(responses map[string]MockResponse) (*http.Client, *MockTransport) {
	transport := &MockTransport{
		Responses: responses,
		Requests:  make(map[string]*http.Request),
	}
	
	return &http.Client{
		Transport: transport,
	}, transport
}

// GetRequestBody retrieves the body of a request that was sent through the mock transport.
func (m *MockTransport) GetRequestBody(method, url string) []byte {
	key := method + " " + url
	req, exists := m.Requests[key]
	if !exists {
		return nil
	}
	
	if req.Body == nil {
		return nil
	}
	
	bodyBytes, _ := io.ReadAll(req.Body)
	req.Body.Close()
	
	// Restore the body for potential future reads
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	
	return bodyBytes
}

// GetRequestHeader retrieves a specific header from a request that was sent through the mock transport.
func (m *MockTransport) GetRequestHeader(method, url, header string) string {
	key := method + " " + url
	req, exists := m.Requests[key]
	if !exists {
		return ""
	}
	
	return req.Header.Get(header)
}