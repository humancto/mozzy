package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Request represents a captured HTTP request
type Request struct {
	ID         int
	Timestamp  time.Time
	Method     string
	URL        string
	Host       string
	Path       string
	StatusCode int
	Duration   time.Duration
	ReqSize    int64
	RespSize   int64
	Headers    http.Header
	Error      string
}

// Server represents a proxy server
type Server struct {
	Port     int
	Verbose  bool
	requests []Request
	mu       sync.RWMutex
	reqID    int
}

// NewServer creates a new proxy server
func NewServer(port int, verbose bool) *Server {
	return &Server{
		Port:     port,
		Verbose:  verbose,
		requests: make([]Request, 0),
	}
}

// Start starts the proxy server
func (s *Server) Start() error {
	handler := http.HandlerFunc(s.handleRequest)

	addr := fmt.Sprintf(":%d", s.Port)

	// Get local IP for display
	localIP := getLocalIP()

	fmt.Println()
	color.Cyan("╔════════════════════════════════════════════════════════════════════")
	color.Cyan("║ 🔄 Mozzy Proxy Server")
	color.Cyan("╚════════════════════════════════════════════════════════════════════")
	fmt.Println()
	color.Green("📡 Listening on:  0.0.0.0:%d", s.Port)
	color.Green("🌐 Local IP:      %s:%d", localIP, s.Port)
	fmt.Println()
	color.HiBlack("Configure your browser or app to use this proxy:")
	color.HiBlack("  HTTP Proxy:  %s:%d", localIP, s.Port)
	fmt.Println()
	color.Yellow("📊 Waiting for connections...")
	fmt.Println()
	color.HiBlack("────────────────────────────────────────────────────────────────────────")
	fmt.Println()

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return server.ListenAndServe()
}

// handleRequest handles incoming proxy requests
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.reqID++
	reqID := s.reqID
	s.mu.Unlock()

	start := time.Now()

	// Log incoming request
	if s.Verbose {
		color.Cyan("→ %s %s", r.Method, r.URL.String())
	}

	// Create the request to forward
	targetURL := r.URL.String()
	if r.URL.Scheme == "" {
		// If no scheme, assume http and use the Host header
		targetURL = "http://" + r.Host + r.URL.Path
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}
	}

	// Create new request
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		s.logError(reqID, r, err)
		http.Error(w, "Proxy error", http.StatusBadGateway)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Remove hop-by-hop headers
	removeHopHeaders(proxyReq.Header)

	// Send the request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		s.logError(reqID, r, err)
		http.Error(w, "Failed to reach target", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	removeHopHeaders(w.Header())

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body and track size
	respSize, _ := io.Copy(w, resp.Body)

	// Log the request
	req := Request{
		ID:         reqID,
		Timestamp:  start,
		Method:     r.Method,
		URL:        targetURL,
		Host:       r.Host,
		Path:       r.URL.Path,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		ReqSize:    r.ContentLength,
		RespSize:   respSize,
		Headers:    r.Header,
	}

	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()

	// Print summary
	statusColor := color.GreenString
	if resp.StatusCode >= 400 {
		statusColor = color.RedString
	} else if resp.StatusCode >= 300 {
		statusColor = color.YellowString
	}

	fmt.Printf("%s  %-6s %-50s %s (%dms)\n",
		color.HiBlackString(start.Format("15:04:05")),
		color.CyanString(r.Method),
		truncate(targetURL, 50),
		statusColor("%d", resp.StatusCode),
		duration.Milliseconds(),
	)
}

// logError logs a proxy error
func (s *Server) logError(reqID int, r *http.Request, err error) {
	req := Request{
		ID:        reqID,
		Timestamp: time.Now(),
		Method:    r.Method,
		URL:       r.URL.String(),
		Host:      r.Host,
		Path:      r.URL.Path,
		Error:     err.Error(),
	}

	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()

	color.Red("✗ %s %s - %v", r.Method, r.URL.String(), err)
}

// GetRequests returns captured requests
func (s *Server) GetRequests() []Request {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy
	reqs := make([]Request, len(s.requests))
	copy(reqs, s.requests)
	return reqs
}

// removeHopHeaders removes hop-by-hop headers
func removeHopHeaders(h http.Header) {
	hopHeaders := []string{
		"Connection",
		"Proxy-Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}

	for _, header := range hopHeaders {
		h.Del(header)
	}
}

// getLocalIP gets the local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "localhost"
}

// truncate truncates a string to the specified length
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
