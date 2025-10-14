package upload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// FileUpload represents a file to upload
type FileUpload struct {
	FieldName string // Form field name
	FilePath  string // Local file path
	FileName  string // Optional custom filename (defaults to basename)
}

// FormField represents a form field value
type FormField struct {
	Name  string
	Value string
}

// UploadOptions configures file upload
type UploadOptions struct {
	URL           string
	Files         []FileUpload
	Fields        []FormField
	Headers       map[string]string
	AuthToken     string
	ShowProgress  bool
	Timeout       time.Duration
}

// Progress tracks upload progress
type Progress struct {
	TotalBytes    int64
	UploadedBytes int64
	StartTime     time.Time
	LastUpdate    time.Time
}

// Upload performs multipart file upload
func Upload(opts UploadOptions) (*http.Response, []byte, error) {
	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Track total size for progress
	var totalSize int64

	// Add files
	for _, file := range opts.Files {
		fileInfo, err := os.Stat(file.FilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("file not found: %s: %w", file.FilePath, err)
		}
		totalSize += fileInfo.Size()

		// Determine filename
		filename := file.FileName
		if filename == "" {
			filename = filepath.Base(file.FilePath)
		}

		part, err := writer.CreateFormFile(file.FieldName, filename)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create form file: %w", err)
		}

		f, err := os.Open(file.FilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open file: %w", err)
		}

		if _, err := io.Copy(part, f); err != nil {
			f.Close()
			return nil, nil, fmt.Errorf("failed to copy file: %w", err)
		}
		f.Close()
	}

	// Add form fields
	for _, field := range opts.Fields {
		if err := writer.WriteField(field.Name, field.Value); err != nil {
			return nil, nil, fmt.Errorf("failed to write field: %w", err)
		}
	}

	// Close multipart writer
	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Show upload info
	if opts.ShowProgress {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Fprintf(os.Stderr, "\n📤 Uploading: %s (%d file(s), %s)\n",
			cyan(opts.URL),
			len(opts.Files),
			formatBytes(totalSize),
		)
	}

	// Create request
	req, err := http.NewRequest("POST", opts.URL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add custom headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Add auth token
	if opts.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.AuthToken)
	}

	// Create client with timeout
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Minute // Default 10 min for uploads
	}

	client := &http.Client{
		Timeout: timeout,
	}

	// Execute request
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Show completion
	if opts.ShowProgress {
		elapsed := time.Since(startTime)
		avgSpeed := float64(totalSize) / elapsed.Seconds()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Fprintf(os.Stderr, "%s Upload complete! (%s, avg: %s/s)\n",
			green("✅"),
			formatDuration(elapsed),
			formatBytes(int64(avgSpeed)),
		)
	}

	return resp, respBody, nil
}

// UploadWithProgress performs upload with real-time progress tracking
// Note: This requires server support for tracking upload progress
// For now, we show a simple spinner since standard HTTP doesn't track upload progress easily
func UploadWithProgress(opts UploadOptions) (*http.Response, []byte, error) {
	if !opts.ShowProgress {
		return Upload(opts)
	}

	// For real progress tracking, we'd need to wrap the body reader
	// or use chunked transfer encoding with progress callbacks
	// For now, we'll just show a spinner
	done := make(chan bool)
	go func() {
		cyan := color.New(color.FgCyan).SprintFunc()
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Fprintf(os.Stderr, "\r%s\n", "                    ")
				return
			default:
				fmt.Fprintf(os.Stderr, "\r%s Uploading... %s",
					cyan(spinner[i%len(spinner)]),
					"",
				)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()

	resp, body, err := Upload(opts)
	done <- true
	return resp, body, err
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
