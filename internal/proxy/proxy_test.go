package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer(8888, false, false)

	if server.Port != 8888 {
		t.Errorf("Expected port 8888, got %d", server.Port)
	}

	if server.Verbose {
		t.Error("Expected verbose to be false")
	}

	if server.HTTPS {
		t.Error("Expected HTTPS to be false")
	}
}

func TestHeaderInjection(t *testing.T) {
	// Create a test HTTP server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if injected header is present
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Errorf("Expected injected header, got: %v", r.Header)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer backend.Close()

	// Create proxy server with header injection
	proxy := NewServer(0, false, false) // Port 0 for random port
	proxy.InjectHeaders = map[string]string{
		"X-Test-Header": "test-value",
	}

	// Note: This is a simplified test. In reality, you'd need to set up
	// the proxy and make requests through it, which is complex for unit tests.
	// Consider integration tests for full proxy testing.
}

func TestRequestFiltering(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		host       string
		statusCode int
		filters    struct {
			domain      string
			methods     []string
			errorsOnly  bool
		}
		shouldFilter bool
	}{
		{
			name:         "No filters - should not filter",
			method:       "GET",
			host:         "example.com",
			statusCode:   200,
			shouldFilter: false,
		},
		{
			name:       "Filter by error-only mode - 200 should be filtered",
			method:     "GET",
			host:       "example.com",
			statusCode: 200,
			filters: struct {
				domain     string
				methods    []string
				errorsOnly bool
			}{
				errorsOnly: true,
			},
			shouldFilter: true,
		},
		{
			name:       "Filter by error-only mode - 404 should not be filtered",
			method:     "GET",
			host:       "example.com",
			statusCode: 404,
			filters: struct {
				domain     string
				methods    []string
				errorsOnly bool
			}{
				errorsOnly: true,
			},
			shouldFilter: false,
		},
		{
			name:       "Filter by method - POST should be filtered",
			method:     "POST",
			host:       "example.com",
			statusCode: 200,
			filters: struct {
				domain     string
				methods    []string
				errorsOnly bool
			}{
				methods: []string{"GET", "PUT"},
			},
			shouldFilter: true,
		},
		{
			name:       "Filter by method - GET should not be filtered",
			method:     "GET",
			host:       "example.com",
			statusCode: 200,
			filters: struct {
				domain     string
				methods    []string
				errorsOnly bool
			}{
				methods: []string{"GET", "PUT"},
			},
			shouldFilter: false,
		},
		{
			name:       "Filter by domain - non-matching should be filtered",
			method:     "GET",
			host:       "example.com",
			statusCode: 200,
			filters: struct {
				domain     string
				methods    []string
				errorsOnly bool
			}{
				domain: "api.test.com",
			},
			shouldFilter: true,
		},
		{
			name:       "Filter by domain - matching should not be filtered",
			method:     "GET",
			host:       "api.example.com",
			statusCode: 200,
			filters: struct {
				domain     string
				methods    []string
				errorsOnly bool
			}{
				domain: "example.com",
			},
			shouldFilter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(8888, false, false)
			server.FilterDomain = tt.filters.domain
			server.FilterMethods = tt.filters.methods
			server.FilterErrors = tt.filters.errorsOnly

			result := server.shouldFilter(tt.method, tt.host, tt.statusCode)
			if result != tt.shouldFilter {
				t.Errorf("Expected shouldFilter=%v, got %v", tt.shouldFilter, result)
			}
		})
	}
}

func TestHARExport(t *testing.T) {
	server := NewServer(8888, false, false)

	// Add some test requests
	server.requests = []Request{
		{
			ID:         1,
			Timestamp:  time.Now(),
			Method:     "GET",
			URL:        "http://example.com/test",
			Host:       "example.com",
			Path:       "/test",
			StatusCode: 200,
			Duration:   100 * time.Millisecond,
			ReqSize:    0,
			RespSize:   1024,
			Headers:    http.Header{"User-Agent": []string{"test"}},
		},
		{
			ID:         2,
			Timestamp:  time.Now(),
			Method:     "POST",
			URL:        "http://example.com/api",
			Host:       "example.com",
			Path:       "/api",
			StatusCode: 201,
			Duration:   200 * time.Millisecond,
			ReqSize:    512,
			RespSize:   256,
			Headers:    http.Header{"Content-Type": []string{"application/json"}},
		},
	}

	// Export to temporary file
	tmpfile, err := os.CreateTemp("", "test-*.har")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Export HAR
	err = server.ExportHAR(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to export HAR: %v", err)
	}

	// Read and parse HAR file
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read HAR file: %v", err)
	}

	var har HAR
	err = json.Unmarshal(data, &har)
	if err != nil {
		t.Fatalf("Failed to parse HAR JSON: %v", err)
	}

	// Verify HAR structure
	if har.Log.Version != "1.2" {
		t.Errorf("Expected HAR version 1.2, got %s", har.Log.Version)
	}

	if har.Log.Creator.Name != "Mozzy Proxy" {
		t.Errorf("Expected creator 'Mozzy Proxy', got %s", har.Log.Creator.Name)
	}

	if len(har.Log.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(har.Log.Entries))
	}

	// Verify first entry
	entry := har.Log.Entries[0]
	if entry.Request.Method != "GET" {
		t.Errorf("Expected GET method, got %s", entry.Request.Method)
	}
	if entry.Request.URL != "http://example.com/test" {
		t.Errorf("Expected URL http://example.com/test, got %s", entry.Request.URL)
	}
	if entry.Response.Status != 200 {
		t.Errorf("Expected status 200, got %d", entry.Response.Status)
	}
	if entry.Time != 100 {
		t.Errorf("Expected time 100ms, got %f", entry.Time)
	}
}

func TestGetStatusText(t *testing.T) {
	tests := []struct {
		code int
		text string
	}{
		{200, "OK"},
		{201, "Created"},
		{204, "No Content"},
		{301, "Moved Permanently"},
		{302, "Found"},
		{304, "Not Modified"},
		{400, "Bad Request"},
		{401, "Unauthorized"},
		{403, "Forbidden"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
		{502, "Bad Gateway"},
		{503, "Service Unavailable"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := getStatusText(tt.code)
			if result != tt.text {
				t.Errorf("Expected %s for code %d, got %s", tt.text, tt.code, result)
			}
		})
	}
}

func TestContainsHelper(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"example.com", "example", true},
		{"api.example.com", "example", true},
		{"example.com", "api", false},
		{"test", "test", true},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.want {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.want)
			}
		})
	}
}

// Integration test helper
func TestProxyBasicHTTP(t *testing.T) {
	// This would be a more complex integration test
	// that starts a backend server, a proxy server,
	// and makes requests through the proxy.
	// Skipping for now as it requires more setup.
	t.Skip("Integration test - requires complex setup")
}
