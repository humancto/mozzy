package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Progress represents download progress
type Progress struct {
	TotalBytes      int64
	DownloadedBytes int64
	StartTime       time.Time
	LastUpdate      time.Time
}

// DownloadOptions configures file download
type DownloadOptions struct {
	URL             string
	OutputPath      string // If empty, derive from URL
	ShowProgress    bool
	OverwriteExist  bool
	FollowRedirects bool
}

// Download downloads a file with optional progress display
func Download(opts DownloadOptions) (string, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Minute, // Longer timeout for large files
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !opts.FollowRedirects && len(via) > 0 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Make request
	resp, err := client.Get(opts.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Determine output path
	outputPath := opts.OutputPath
	if outputPath == "" {
		outputPath = getFilenameFromURL(opts.URL)
		// Try Content-Disposition header for better filename
		if cd := resp.Header.Get("Content-Disposition"); cd != "" {
			if fn := parseFilenameFromContentDisposition(cd); fn != "" {
				outputPath = fn
			}
		}
	}

	// Check if file exists
	if !opts.OverwriteExist {
		if _, err := os.Stat(outputPath); err == nil {
			return "", fmt.Errorf("file already exists: %s (use --overwrite to replace)", outputPath)
		}
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress
	progress := &Progress{
		TotalBytes: resp.ContentLength,
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}

	if opts.ShowProgress && resp.ContentLength > 0 {
		// Progress bar for known size
		return outputPath, downloadWithProgress(resp.Body, out, progress)
	} else if opts.ShowProgress {
		// Simple progress for unknown size
		return outputPath, downloadWithSimpleProgress(resp.Body, out, progress)
	} else {
		// No progress display
		_, err = io.Copy(out, resp.Body)
		return outputPath, err
	}
}

func downloadWithProgress(src io.Reader, dst io.Writer, progress *Progress) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Fprintf(os.Stderr, "\n📥 Downloading: %s\n", formatBytes(progress.TotalBytes))

	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, errW := dst.Write(buf[:nr])
			if errW != nil {
				return errW
			}
			if nr != nw {
				return io.ErrShortWrite
			}

			progress.DownloadedBytes += int64(nw)

			// Update progress every 100ms
			if time.Since(progress.LastUpdate) > 100*time.Millisecond {
				percentage := float64(progress.DownloadedBytes) / float64(progress.TotalBytes) * 100
				speed := calculateSpeed(progress)
				eta := calculateETA(progress)

				// Progress bar
				barWidth := 40
				filled := int(float64(barWidth) * percentage / 100)
				bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

				fmt.Fprintf(os.Stderr, "\r%s [%s] %.1f%% | %s/%s | %s/s | ETA: %s",
					cyan("Progress:"),
					green(bar),
					percentage,
					formatBytes(progress.DownloadedBytes),
					formatBytes(progress.TotalBytes),
					formatBytes(speed),
					eta,
				)

				progress.LastUpdate = time.Now()
			}
		}

		if err != nil {
			if err == io.EOF {
				// Final progress update
				fmt.Fprintf(os.Stderr, "\r%s [%s] %.1f%% | %s/%s | Done!%s\n",
					cyan("Progress:"),
					green(strings.Repeat("█", 40)),
					100.0,
					formatBytes(progress.DownloadedBytes),
					formatBytes(progress.TotalBytes),
					strings.Repeat(" ", 30), // Clear ETA
				)
				return nil
			}
			return err
		}
	}
}

func downloadWithSimpleProgress(src io.Reader, dst io.Writer, progress *Progress) error {
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Fprintf(os.Stderr, "\n📥 Downloading (size unknown)...\n")

	buf := make([]byte, 32*1024)

	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, errW := dst.Write(buf[:nr])
			if errW != nil {
				return errW
			}
			progress.DownloadedBytes += int64(nw)

			// Update progress every 500ms
			if time.Since(progress.LastUpdate) > 500*time.Millisecond {
				speed := calculateSpeed(progress)
				elapsed := time.Since(progress.StartTime)

				fmt.Fprintf(os.Stderr, "\r%s %s | %s/s | %s elapsed",
					cyan("Downloaded:"),
					formatBytes(progress.DownloadedBytes),
					formatBytes(speed),
					formatDuration(elapsed),
				)

				progress.LastUpdate = time.Now()
			}
		}

		if err != nil {
			if err == io.EOF {
				elapsed := time.Since(progress.StartTime)
				avgSpeed := float64(progress.DownloadedBytes) / elapsed.Seconds()

				fmt.Fprintf(os.Stderr, "\r%s %s | Avg: %s/s | %s%s\n",
					cyan("Downloaded:"),
					formatBytes(progress.DownloadedBytes),
					formatBytes(int64(avgSpeed)),
					formatDuration(elapsed),
					strings.Repeat(" ", 20),
				)
				return nil
			}
			return err
		}
	}
}

func calculateSpeed(p *Progress) int64 {
	elapsed := time.Since(p.StartTime).Seconds()
	if elapsed == 0 {
		return 0
	}
	return int64(float64(p.DownloadedBytes) / elapsed)
}

func calculateETA(p *Progress) string {
	if p.DownloadedBytes == 0 {
		return "calculating..."
	}

	elapsed := time.Since(p.StartTime).Seconds()
	speed := float64(p.DownloadedBytes) / elapsed
	remaining := float64(p.TotalBytes - p.DownloadedBytes)

	if speed == 0 {
		return "unknown"
	}

	etaSeconds := remaining / speed
	return formatDuration(time.Duration(etaSeconds * float64(time.Second)))
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

func getFilenameFromURL(urlStr string) string {
	parts := strings.Split(urlStr, "/")
	filename := parts[len(parts)-1]

	// Remove query parameters
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	// Default if empty
	if filename == "" {
		filename = "download"
	}

	return filename
}

func parseFilenameFromContentDisposition(cd string) string {
	// Simple parser for Content-Disposition: attachment; filename="..."
	parts := strings.Split(cd, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "filename=") {
			filename := strings.TrimPrefix(part, "filename=")
			filename = strings.Trim(filename, "\"")
			return filepath.Base(filename) // Security: prevent path traversal
		}
	}
	return ""
}
