package download

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetFilenameFromURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://example.com/file.txt", "file.txt"},
		{"https://example.com/path/to/file.zip", "file.zip"},
		{"https://example.com/file.txt?version=1", "file.txt"},
		{"https://example.com/", "download"},
		{"https://example.com", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := getFilenameFromURL(tt.url)
			if got != tt.want {
				t.Errorf("getFilenameFromURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestParseFilenameFromContentDisposition(t *testing.T) {
	tests := []struct {
		name string
		cd   string
		want string
	}{
		{
			name: "simple attachment",
			cd:   `attachment; filename="file.txt"`,
			want: "file.txt",
		},
		{
			name: "without quotes",
			cd:   `attachment; filename=file.txt`,
			want: "file.txt",
		},
		{
			name: "with extra params",
			cd:   `attachment; filename="document.pdf"; size=12345`,
			want: "document.pdf",
		},
		{
			name: "no filename",
			cd:   `attachment`,
			want: "",
		},
		{
			name: "path traversal attempt",
			cd:   `attachment; filename="../../../etc/passwd"`,
			want: "passwd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFilenameFromContentDisposition(tt.cd)
			if got != tt.want {
				t.Errorf("parseFilenameFromContentDisposition() = %q, want %q", got, tt.want)
			}
		})
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
		{5368709120, "5.0 GB"},
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
		{125 * time.Second, "2m5s"},
		{3665 * time.Second, "1h1m5s"},
		{0, "0s"},
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

func TestCalculateSpeed(t *testing.T) {
	tests := []struct {
		name            string
		downloadedBytes int64
		elapsed         time.Duration
		wantMin         int64
		wantMax         int64
	}{
		{
			name:            "1MB in 1 second",
			downloadedBytes: 1048576,
			elapsed:         1 * time.Second,
			wantMin:         1000000,
			wantMax:         1100000,
		},
		{
			name:            "minimal time elapsed",
			downloadedBytes: 1000,
			elapsed:         1 * time.Millisecond,
			wantMin:         0,
			wantMax:         10000000, // Very fast, but acceptable
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Progress{
				DownloadedBytes: tt.downloadedBytes,
				StartTime:       time.Now().Add(-tt.elapsed),
			}
			got := calculateSpeed(p)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("calculateSpeed() = %d, want between %d and %d", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestDownload_Success(t *testing.T) {
	// Create test server
	content := []byte("test file content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}))
	defer server.Close()

	// Create temp directory
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.txt")

	// Download
	opts := DownloadOptions{
		URL:          server.URL + "/test.txt",
		OutputPath:   outputPath,
		ShowProgress: false,
	}

	result, err := Download(opts)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}

	if result != outputPath {
		t.Errorf("Download() returned path = %q, want %q", result, outputPath)
	}

	// Verify file content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("Downloaded content = %q, want %q", data, content)
	}
}

func TestDownload_FileExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("new content"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "existing.txt")

	// Create existing file
	err := os.WriteFile(outputPath, []byte("existing content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try download without overwrite
	opts := DownloadOptions{
		URL:            server.URL + "/file.txt",
		OutputPath:     outputPath,
		ShowProgress:   false,
		OverwriteExist: false,
	}

	_, err = Download(opts)
	if err == nil {
		t.Error("Download() expected error for existing file, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Download() error = %q, want 'already exists' message", err)
	}
}

func TestDownload_WithOverwrite(t *testing.T) {
	newContent := []byte("new content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(newContent)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "existing.txt")

	// Create existing file
	err := os.WriteFile(outputPath, []byte("old content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Download with overwrite
	opts := DownloadOptions{
		URL:            server.URL + "/file.txt",
		OutputPath:     outputPath,
		ShowProgress:   false,
		OverwriteExist: true,
	}

	_, err = Download(opts)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}

	// Verify new content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(data) != string(newContent) {
		t.Errorf("File content = %q, want %q", data, newContent)
	}
}

func TestDownload_ContentDisposition(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", `attachment; filename="custom-name.pdf"`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pdf content"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	// Download without specifying output path
	opts := DownloadOptions{
		URL:          server.URL + "/some/path/file.txt",
		OutputPath:   "", // Let it derive from Content-Disposition
		ShowProgress: false,
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	result, err := Download(opts)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}

	if result != "custom-name.pdf" {
		t.Errorf("Download() used filename = %q, want %q", result, "custom-name.pdf")
	}
}

func TestDownload_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.txt")

	opts := DownloadOptions{
		URL:          server.URL + "/missing.txt",
		OutputPath:   outputPath,
		ShowProgress: false,
	}

	_, err := Download(opts)
	if err == nil {
		t.Error("Download() expected error for 404, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Download() error = %q, want '404' in message", err)
	}
}
