package upload

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestUpload_SingleFile(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("test file content")
	err := os.WriteFile(testFile, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test server
	var receivedContent []byte
	var receivedFilename string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Failed to get form file: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		receivedFilename = header.Filename
		receivedContent, _ = io.ReadAll(file)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	// Upload file
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile},
		},
		ShowProgress: false,
	}

	resp, body, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Upload() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if string(body) != `{"status":"success"}` {
		t.Errorf("Upload() response = %q, want success message", body)
	}

	if string(receivedContent) != string(testContent) {
		t.Errorf("Received content = %q, want %q", receivedContent, testContent)
	}

	if receivedFilename != "test.txt" {
		t.Errorf("Received filename = %q, want %q", receivedFilename, "test.txt")
	}
}

func TestUpload_MultipleFiles(t *testing.T) {
	// Create test files
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)

	// Create test server
	receivedFiles := make(map[string]string)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)

		// Read all files from form
		for fieldName := range r.MultipartForm.File {
			file, header, _ := r.FormFile(fieldName)
			content, _ := io.ReadAll(file)
			receivedFiles[fieldName] = string(content)
			file.Close()
			_ = header // Use header to avoid unused variable
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Upload files
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file1", FilePath: file1},
			{FieldName: "file2", FilePath: file2},
		},
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if receivedFiles["file1"] != "content1" {
		t.Errorf("file1 content = %q, want %q", receivedFiles["file1"], "content1")
	}

	if receivedFiles["file2"] != "content2" {
		t.Errorf("file2 content = %q, want %q", receivedFiles["file2"], "content2")
	}
}

func TestUpload_WithFormFields(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Create test server
	var receivedName, receivedEmail string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		receivedName = r.FormValue("name")
		receivedEmail = r.FormValue("email")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Upload with form fields
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile},
		},
		Fields: []FormField{
			{Name: "name", Value: "Alice"},
			{Name: "email", Value: "alice@example.com"},
		},
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if receivedName != "Alice" {
		t.Errorf("Received name = %q, want %q", receivedName, "Alice")
	}

	if receivedEmail != "alice@example.com" {
		t.Errorf("Received email = %q, want %q", receivedEmail, "alice@example.com")
	}
}

func TestUpload_CustomFilename(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "original.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Create test server
	var receivedFilename string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		_, header, _ := r.FormFile("file")
		receivedFilename = header.Filename
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Upload with custom filename
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile, FileName: "custom-name.txt"},
		},
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if receivedFilename != "custom-name.txt" {
		t.Errorf("Received filename = %q, want %q", receivedFilename, "custom-name.txt")
	}
}

func TestUpload_WithHeaders(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Create test server
	var receivedAPIKey, receivedCustom string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAPIKey = r.Header.Get("X-API-Key")
		receivedCustom = r.Header.Get("X-Custom-Header")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Upload with headers
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile},
		},
		Headers: map[string]string{
			"X-API-Key":       "secret123",
			"X-Custom-Header": "value",
		},
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if receivedAPIKey != "secret123" {
		t.Errorf("Received API key = %q, want %q", receivedAPIKey, "secret123")
	}

	if receivedCustom != "value" {
		t.Errorf("Received custom header = %q, want %q", receivedCustom, "value")
	}
}

func TestUpload_WithAuth(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Create test server
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Upload with auth token
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile},
		},
		AuthToken:    "token123",
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if receivedAuth != "Bearer token123" {
		t.Errorf("Received auth = %q, want %q", receivedAuth, "Bearer token123")
	}
}

func TestUpload_FileNotFound(t *testing.T) {
	opts := UploadOptions{
		URL: "http://example.com",
		Files: []FileUpload{
			{FieldName: "file", FilePath: "/nonexistent/file.txt"},
		},
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err == nil {
		t.Error("Upload() expected error for nonexistent file, got nil")
	}

	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("Upload() error = %q, want 'file not found' message", err)
	}
}

func TestUpload_ServerError(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	// Upload
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile},
		},
		ShowProgress: false,
	}

	resp, body, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Upload() status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	if string(body) != "server error" {
		t.Errorf("Upload() response = %q, want error message", body)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.bytes), func(t *testing.T) {
			got := formatBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{5 * time.Second, "5s"},
		{65 * time.Second, "1m5s"},
		{3665 * time.Second, "1h1m5s"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestUpload_ContentType(t *testing.T) {
	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Create test server
	var receivedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Upload
	opts := UploadOptions{
		URL: server.URL,
		Files: []FileUpload{
			{FieldName: "file", FilePath: testFile},
		},
		ShowProgress: false,
	}

	_, _, err := Upload(opts)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	if !strings.HasPrefix(receivedContentType, "multipart/form-data; boundary=") {
		t.Errorf("Content-Type = %q, want multipart/form-data with boundary", receivedContentType)
	}
}
