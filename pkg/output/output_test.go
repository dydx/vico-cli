package output

import (
	"reflect"
	"testing"

	"github.com/dydx/vico-cli/pkg/output/stdout"
)

// TestFactory tests the Factory function to ensure it returns the correct handler type
func TestFactory(t *testing.T) {
	tests := []struct {
		name          string
		format        string
		expectedType  string
		expectError   bool
	}{
		{
			name:         "Factory returns JSONHandler for 'json' format",
			format:       "json",
			expectedType: "*stdout.JSONHandler",
			expectError:  false,
		},
		{
			name:         "Factory returns JSONHandler for 'JSON' format (case insensitive)",
			format:       "JSON",
			expectedType: "*stdout.JSONHandler",
			expectError:  false,
		},
		{
			name:         "Factory returns JSONHandler for 'JsOn' format (case insensitive)",
			format:       "JsOn",
			expectedType: "*stdout.JSONHandler",
			expectError:  false,
		},
		{
			name:         "Factory returns TableHandler for empty format",
			format:       "",
			expectedType: "*stdout.TableHandler",
			expectError:  false,
		},
		{
			name:         "Factory returns TableHandler for unknown format",
			format:       "unknown",
			expectedType: "*stdout.TableHandler",
			expectError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := Factory(tc.format)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Error("Expected an error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check handler type
			if handler == nil {
				t.Fatal("Handler is nil")
			}
			actualType := reflect.TypeOf(handler).String()
			if actualType != tc.expectedType {
				t.Errorf("Expected handler type %s, got %s", tc.expectedType, actualType)
			}
		})
	}
}

// TestNewStdoutHandler tests the NewStdoutHandler function to ensure it returns the correct handler type
func TestNewStdoutHandler(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		expectedType string
	}{
		{
			name:         "NewStdoutHandler returns JSONHandler for 'json' format",
			format:       "json",
			expectedType: "*stdout.JSONHandler",
		},
		{
			name:         "NewStdoutHandler returns JSONHandler for 'JSON' format (case insensitive)",
			format:       "JSON",
			expectedType: "*stdout.JSONHandler",
		},
		{
			name:         "NewStdoutHandler returns JSONHandler for 'JsOn' format (case insensitive)",
			format:       "JsOn",
			expectedType: "*stdout.JSONHandler",
		},
		{
			name:         "NewStdoutHandler returns TableHandler for empty format",
			format:       "",
			expectedType: "*stdout.TableHandler",
		},
		{
			name:         "NewStdoutHandler returns TableHandler for unknown format",
			format:       "unknown",
			expectedType: "*stdout.TableHandler",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewStdoutHandler(tc.format)

			// Check handler type
			if handler == nil {
				t.Fatal("Handler is nil")
			}
			actualType := reflect.TypeOf(handler).String()
			if actualType != tc.expectedType {
				t.Errorf("Expected handler type %s, got %s", tc.expectedType, actualType)
			}
		})
	}
}

// TestHandlerSelectionBasedOnFormat tests that the correct handler is selected based on the format
func TestHandlerSelectionBasedOnFormat(t *testing.T) {
	// Test case 1: Factory + "json" format should return JSONHandler
	jsonHandler, err := Factory("json")
	if err != nil {
		t.Fatalf("Unexpected error from Factory with 'json' format: %v", err)
	}
	if _, ok := jsonHandler.(*stdout.JSONHandler); !ok {
		t.Errorf("Factory with 'json' format should return *stdout.JSONHandler, got %T", jsonHandler)
	}

	// Test case 2: Factory + "" format should return TableHandler
	tableHandler, err := Factory("")
	if err != nil {
		t.Fatalf("Unexpected error from Factory with empty format: %v", err)
	}
	if _, ok := tableHandler.(*stdout.TableHandler); !ok {
		t.Errorf("Factory with empty format should return *stdout.TableHandler, got %T", tableHandler)
	}

	// Test case 3: Factory + "unknown" format should default to TableHandler
	defaultHandler, err := Factory("unknown")
	if err != nil {
		t.Fatalf("Unexpected error from Factory with 'unknown' format: %v", err)
	}
	if _, ok := defaultHandler.(*stdout.TableHandler); !ok {
		t.Errorf("Factory with 'unknown' format should return *stdout.TableHandler, got %T", defaultHandler)
	}

	// Test case 4: Case insensitivity - Factory + "JSON" should return JSONHandler
	upperCaseHandler, err := Factory("JSON")
	if err != nil {
		t.Fatalf("Unexpected error from Factory with 'JSON' format: %v", err)
	}
	if _, ok := upperCaseHandler.(*stdout.JSONHandler); !ok {
		t.Errorf("Factory with 'JSON' format should return *stdout.JSONHandler (case insensitive), got %T", upperCaseHandler)
	}

	// Test case 5: Case insensitivity - Factory + mixed case "JsOn" should return JSONHandler
	mixedCaseHandler, err := Factory("JsOn")
	if err != nil {
		t.Fatalf("Unexpected error from Factory with 'JsOn' format: %v", err)
	}
	if _, ok := mixedCaseHandler.(*stdout.JSONHandler); !ok {
		t.Errorf("Factory with 'JsOn' format should return *stdout.JSONHandler (case insensitive), got %T", mixedCaseHandler)
	}
}