package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Server represents a mock HTTP server
type Server struct {
	config     *Config
	configPath string
	mu         sync.RWMutex
	server     *http.Server
	logger     *RequestLogger
}

// RequestLogger logs incoming requests
type RequestLogger struct {
	mu      sync.Mutex
	entries []LogEntry
}

// LogEntry represents a logged request
type LogEntry struct {
	Timestamp  time.Time
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	RemoteAddr string
}

// NewServer creates a new mock server
func NewServer(config *Config, configPath string) *Server {
	return &Server{
		config:     config,
		configPath: configPath,
		logger:     &RequestLogger{entries: make([]LogEntry, 0)},
	}
}

// Start starts the mock server
func (s *Server) Start() error {
	// Create a simple router that supports method-specific routes
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find matching route
		var matchedRoute *Route
		for i := range s.config.Routes {
			route := &s.config.Routes[i]
			if r.URL.Path == route.Path && r.Method == route.Method {
				matchedRoute = route
				break
			}
		}

		// Handle CORS preflight
		if r.Method == "OPTIONS" && s.config.CORS.Enabled {
			s.handleCORS(w, r)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// No matching route found
		if matchedRoute == nil {
			if s.config.CORS.Enabled {
				s.handleCORS(w, r)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Route not found",
				"path":  r.URL.Path,
			})
			s.printRequest(r.Method, r.URL.Path, http.StatusNotFound, 0)
			return
		}

		start := time.Now()

		// Apply CORS headers
		if s.config.CORS.Enabled {
			s.handleCORS(w, r)
		}

		// Apply custom headers
		for k, v := range matchedRoute.Headers {
			w.Header().Set(k, v)
		}

		// Apply delay if configured
		if matchedRoute.Delay > 0 {
			time.Sleep(time.Duration(matchedRoute.Delay) * time.Millisecond)
		}

		// Write response
		w.WriteHeader(matchedRoute.StatusCode)

		responseData, err := matchedRoute.ToJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(responseData)

		// Log request
		duration := time.Since(start)
		s.logger.Log(LogEntry{
			Timestamp:  time.Now(),
			Method:     r.Method,
			Path:       r.URL.Path,
			StatusCode: matchedRoute.StatusCode,
			Duration:   duration,
			RemoteAddr: r.RemoteAddr,
		})

		// Print to console
		s.printRequest(r.Method, r.URL.Path, matchedRoute.StatusCode, duration)
	})

	// Print registered routes
	for _, route := range s.config.Routes {
		fmt.Printf("  %s %s → %d\n",
			color.MagentaString(route.Method),
			color.CyanString(route.Path),
			route.StatusCode)
	}

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	fmt.Printf("\n%s\n", color.GreenString("🚀 Mock server started at http://%s", addr))
	fmt.Printf("%s\n\n", color.CyanString("Press Ctrl+C to stop"))

	return s.server.ListenAndServe()
}

// Stop stops the mock server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// handleCORS applies CORS headers
func (s *Server) handleCORS(w http.ResponseWriter, r *http.Request) {
	if len(s.config.CORS.Origins) > 0 {
		origin := "*"
		if len(s.config.CORS.Origins) == 1 {
			origin = s.config.CORS.Origins[0]
		} else if r.Header.Get("Origin") != "" {
			// Check if origin is allowed
			for _, allowed := range s.config.CORS.Origins {
				if allowed == "*" || allowed == r.Header.Get("Origin") {
					origin = r.Header.Get("Origin")
					break
				}
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if len(s.config.CORS.Methods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(s.config.CORS.Methods, ", "))
	}

	if len(s.config.CORS.Headers) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(s.config.CORS.Headers, ", "))
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// printRequest prints a formatted request log
func (s *Server) printRequest(method, path string, status int, duration time.Duration) {
	methodColor := color.New(color.FgMagenta, color.Bold)
	statusColor := color.GreenString
	if status >= 400 {
		statusColor = color.YellowString
	}
	if status >= 500 {
		statusColor = color.RedString
	}

	fmt.Printf("%s %s %s %s %s\n",
		color.HiBlackString(time.Now().Format("15:04:05")),
		methodColor.Sprint(method),
		color.CyanString(path),
		statusColor(fmt.Sprintf("(%d)", status)),
		color.HiBlackString(duration.String()))
}

// GetLogs returns all logged requests
func (rl *RequestLogger) GetLogs() []LogEntry {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Return a copy
	logs := make([]LogEntry, len(rl.entries))
	copy(logs, rl.entries)
	return logs
}

// Log adds a new log entry
func (rl *RequestLogger) Log(entry LogEntry) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.entries = append(rl.entries, entry)

	// Keep only last 100 entries
	if len(rl.entries) > 100 {
		rl.entries = rl.entries[len(rl.entries)-100:]
	}
}

// Clear clears all log entries
func (rl *RequestLogger) Clear() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.entries = make([]LogEntry, 0)
}
