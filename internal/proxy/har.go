package proxy

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// HAR (HTTP Archive) format structures
// Spec: http://www.softwareishard.com/blog/har-12-spec/

type HAR struct {
	Log HARLog `json:"log"`
}

type HARLog struct {
	Version string      `json:"version"`
	Creator HARCreator  `json:"creator"`
	Entries []HAREntry  `json:"entries"`
}

type HARCreator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type HAREntry struct {
	StartedDateTime string      `json:"startedDateTime"`
	Time            float64     `json:"time"`
	Request         HARRequest  `json:"request"`
	Response        HARResponse `json:"response"`
	Cache           HARCache    `json:"cache"`
	Timings         HARTimings  `json:"timings"`
}

type HARRequest struct {
	Method      string       `json:"method"`
	URL         string       `json:"url"`
	HTTPVersion string       `json:"httpVersion"`
	Headers     []HARHeader  `json:"headers"`
	QueryString []HARQuery   `json:"queryString"`
	HeadersSize int64        `json:"headersSize"`
	BodySize    int64        `json:"bodySize"`
}

type HARResponse struct {
	Status      int          `json:"status"`
	StatusText  string       `json:"statusText"`
	HTTPVersion string       `json:"httpVersion"`
	Headers     []HARHeader  `json:"headers"`
	Content     HARContent   `json:"content"`
	RedirectURL string       `json:"redirectURL"`
	HeadersSize int64        `json:"headersSize"`
	BodySize    int64        `json:"bodySize"`
}

type HARHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HARQuery struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HARContent struct {
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
}

type HARCache struct{}

type HARTimings struct {
	Send    float64 `json:"send"`
	Wait    float64 `json:"wait"`
	Receive float64 `json:"receive"`
}

// ExportHAR exports captured requests to HAR format
func (s *Server) ExportHAR(filename string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]HAREntry, 0, len(s.requests))

	for _, req := range s.requests {
		// Skip failed requests
		if req.Error != "" {
			continue
		}

		// Convert headers to HAR format
		headers := make([]HARHeader, 0, len(req.Headers))
		for name, values := range req.Headers {
			for _, value := range values {
				headers = append(headers, HARHeader{
					Name:  name,
					Value: value,
				})
			}
		}

		// Create HAR entry
		entry := HAREntry{
			StartedDateTime: req.Timestamp.Format(time.RFC3339),
			Time:            float64(req.Duration.Milliseconds()),
			Request: HARRequest{
				Method:      req.Method,
				URL:         req.URL,
				HTTPVersion: "HTTP/1.1",
				Headers:     headers,
				QueryString: []HARQuery{}, // TODO: Parse query string
				HeadersSize: -1,
				BodySize:    req.ReqSize,
			},
			Response: HARResponse{
				Status:      req.StatusCode,
				StatusText:  getStatusText(req.StatusCode),
				HTTPVersion: "HTTP/1.1",
				Headers:     []HARHeader{}, // TODO: Capture response headers
				Content: HARContent{
					Size:     req.RespSize,
					MimeType: "application/octet-stream", // TODO: Get from headers
				},
				RedirectURL: "",
				HeadersSize: -1,
				BodySize:    req.RespSize,
			},
			Cache: HARCache{},
			Timings: HARTimings{
				Send:    -1,
				Wait:    float64(req.Duration.Milliseconds()),
				Receive: -1,
			},
		}

		entries = append(entries, entry)
	}

	har := HAR{
		Log: HARLog{
			Version: "1.2",
			Creator: HARCreator{
				Name:    "Mozzy Proxy",
				Version: "1.14.0",
			},
			Entries: entries,
		},
	}

	// Write to file
	data, err := json.MarshalIndent(har, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal HAR: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write HAR file: %w", err)
	}

	return nil
}

// getStatusText returns HTTP status text for a status code
func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 301:
		return "Moved Permanently"
	case 302:
		return "Found"
	case 304:
		return "Not Modified"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	default:
		return "Unknown"
	}
}
