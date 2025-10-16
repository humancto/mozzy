package proxy

import (
	"bufio"
	"crypto/tls"
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
	Port      int
	Verbose   bool
	HTTPS     bool
	CA        *CA
	certCache map[string]*tls.Certificate
	requests  []Request
	mu        sync.RWMutex
	reqID     int
}

// NewServer creates a new proxy server
func NewServer(port int, verbose bool, https bool) *Server {
	return &Server{
		Port:      port,
		Verbose:   verbose,
		HTTPS:     https,
		certCache: make(map[string]*tls.Certificate),
		requests:  make([]Request, 0),
	}
}

// Start starts the proxy server
func (s *Server) Start() error {
	// Load or generate CA if HTTPS is enabled
	if s.HTTPS {
		ca, err := GetCA()
		if err != nil {
			return fmt.Errorf("failed to get CA: %w", err)
		}
		s.CA = ca
		fmt.Println()
	}

	handler := http.HandlerFunc(s.handleRequest)

	addr := fmt.Sprintf(":%d", s.Port)

	// Get local IP for display
	localIP := getLocalIP()

	fmt.Println()
	color.Cyan("╔════════════════════════════════════════════════════════════════════")
	if s.HTTPS {
		color.Cyan("║ 🔐 Mozzy HTTPS Proxy Server")
	} else {
		color.Cyan("║ 🔄 Mozzy HTTP Proxy Server")
	}
	color.Cyan("╚════════════════════════════════════════════════════════════════════")
	fmt.Println()
	color.Green("📡 Listening on:  0.0.0.0:%d", s.Port)
	color.Green("🌐 Local IP:      %s:%d", localIP, s.Port)
	fmt.Println()
	color.HiBlack("Configure your browser or app to use this proxy:")
	color.HiBlack("  HTTP%s Proxy:  %s:%d", map[bool]string{true: "S", false: ""}[s.HTTPS], localIP, s.Port)

	if s.HTTPS {
		fmt.Println()
		color.Yellow("⚠️  HTTPS Mode: You must install the CA certificate")
		color.HiBlack("  Run: mozzy proxy --export-cert > mozzy-ca.pem")
		color.HiBlack("  Then install mozzy-ca.pem in your system")
	}

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
	// Handle CONNECT for HTTPS
	if r.Method == http.MethodConnect {
		s.handleHTTPS(w, r)
		return
	}

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

// handleHTTPS handles HTTPS CONNECT requests
func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	// Extract host from request
	host := r.Host
	if host == "" {
		color.Red("✗ CONNECT - missing host")
		http.Error(w, "Missing host", http.StatusBadRequest)
		return
	}

	if s.Verbose {
		color.Cyan("→ CONNECT %s", host)
	}

	// Generate or retrieve cached certificate for this host
	s.mu.Lock()
	cert, ok := s.certCache[host]
	if !ok {
		// Extract hostname without port for certificate generation
		hostname := host
		if h, _, err := net.SplitHostPort(host); err == nil {
			hostname = h
		}

		if s.Verbose {
			color.Yellow("  Generating certificate for %s", hostname)
		}

		// Generate certificate for this host
		serverCert, serverKey, err := s.CA.GenerateServerCert(hostname)
		if err != nil {
			s.mu.Unlock()
			color.Red("✗ Failed to generate certificate for %s: %v", hostname, err)
			http.Error(w, "Certificate generation failed", http.StatusInternalServerError)
			return
		}

		// Create TLS certificate
		tlsCert := &tls.Certificate{
			Certificate: [][]byte{serverCert.Raw, s.CA.Cert.Raw},
			PrivateKey:  serverKey,
			Leaf:        serverCert,
		}
		s.certCache[host] = tlsCert
		cert = tlsCert
	}
	s.mu.Unlock()

	// Hijack the connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		color.Red("✗ CONNECT - hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		color.Red("✗ CONNECT - hijack failed: %v", err)
		return
	}
	defer clientConn.Close()

	if s.Verbose {
		color.Yellow("  Connection hijacked")
	}

	// Send 200 Connection Established
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		if s.Verbose {
			color.Red("✗ CONNECT - failed to send 200: %v", err)
		}
		return
	}

	if s.Verbose {
		color.Yellow("  Sent 200 Connection Established")
	}

	// Wrap connection with TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
	}
	tlsConn := tls.Server(clientConn, tlsConfig)
	defer tlsConn.Close()

	if s.Verbose {
		color.Yellow("  Performing TLS handshake...")
	}

	// Perform TLS handshake
	if err := tlsConn.Handshake(); err != nil {
		if s.Verbose {
			color.Red("✗ TLS handshake failed for %s: %v", host, err)
		}
		return
	}

	if s.Verbose {
		color.Green("  TLS handshake successful")
	}

	// Read the actual HTTPS request
	reader := bufio.NewReader(tlsConn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		if s.Verbose {
			color.Red("✗ Failed to read HTTPS request: %v", err)
		}
		return
	}

	// Reconstruct the full URL
	req.URL.Scheme = "https"
	req.URL.Host = r.Host

	s.mu.Lock()
	s.reqID++
	reqID := s.reqID
	s.mu.Unlock()

	start := time.Now()

	// Forward the request to the target server
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Verify the actual target server's cert
			},
		},
	}

	// Create proxy request
	proxyReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		s.logError(reqID, req, err)
		return
	}

	// Copy headers
	for key, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}
	removeHopHeaders(proxyReq.Header)

	// Send request
	resp, err := client.Do(proxyReq)
	if err != nil {
		s.logError(reqID, req, err)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	// Write response back to client
	err = resp.Write(tlsConn)
	if err != nil {
		if s.Verbose {
			color.Red("✗ Failed to write response: %v", err)
		}
		return
	}

	// Log the request
	reqLog := Request{
		ID:         reqID,
		Timestamp:  start,
		Method:     req.Method,
		URL:        req.URL.String(),
		Host:       req.Host,
		Path:       req.URL.Path,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		ReqSize:    req.ContentLength,
		Headers:    req.Header,
	}

	s.mu.Lock()
	s.requests = append(s.requests, reqLog)
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
		color.CyanString(req.Method),
		truncate(req.URL.String(), 50),
		statusColor("%d", resp.StatusCode),
		duration.Milliseconds(),
	)
}
