package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/humancto/mozzy/internal/retry"
)

type Request struct {
	Method         string
	URL            string
	Headers        []string // "Key: Value"
	Token          string
	Body           []byte
	JSON           bool
	Verbose        bool
	RetryCount     int
	RetryCondition string // e.g., "5xx", ">=500", "429,5xx"
	CookieJar      string
}

type TimingInfo struct {
	DNSLookup        time.Duration
	TCPConnection    time.Duration
	TLSHandshake     time.Duration
	ServerProcessing time.Duration
	ContentTransfer  time.Duration
	Total            time.Duration
}

type VerboseInfo struct {
	ResolvedIP      string
	DNSLatency      time.Duration
	Protocol        string
	TLSVersion      string
	TLSCipher       string
	CertSubject     string
	CertIssuer      string
	CertExpiry      time.Time
	RequestSize     int64
	ResponseSize    int64
	Compressed      bool
	CompressionRatio float64
}

var globalCookieJar *cookiejar.Jar

func Do(ctx context.Context, r Request) (*http.Response, []byte, time.Duration, error) {
	var timings TimingInfo
	var verboseInfo VerboseInfo
	var res *http.Response
	var body []byte
	var err error

	// Parse retry conditions
	conditions, parseErr := retry.ParseConditions(r.RetryCondition)
	if parseErr != nil {
		return nil, nil, 0, fmt.Errorf("invalid retry condition: %w", parseErr)
	}

	policy := &retry.Policy{
		MaxRetries: r.RetryCount,
		Conditions: conditions,
	}

	// Retry logic with exponential backoff
	maxRetries := r.RetryCount
	if maxRetries < 0 {
		maxRetries = 0
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			if r.Verbose {
				fmt.Fprintf(os.Stderr, "⏱  Retry %d/%d after %v backoff...\n", attempt, maxRetries, backoff)
			}
			time.Sleep(backoff)
		}

		res, body, timings, verboseInfo, err = doRequest(ctx, r)

		// Check if we should retry based on policy
		shouldRetry := false
		if attempt < maxRetries {
			if err != nil {
				shouldRetry = policy.ShouldRetry(0, err)
			} else {
				shouldRetry = policy.ShouldRetry(res.StatusCode, nil)
			}
		}

		if !shouldRetry {
			break
		}

		if r.Verbose && shouldRetry {
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Request failed: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "❌ Status %d matches retry condition\n", res.StatusCode)
			}
		}
	}

	if r.Verbose && err == nil {
		printVerbose(r, res, timings, verboseInfo)
	}

	return res, body, timings.Total, err
}

func doRequest(ctx context.Context, r Request) (*http.Response, []byte, TimingInfo, VerboseInfo, error) {
	var timings TimingInfo
	var verboseInfo VerboseInfo
	var dnsStart, connectStart, tlsStart, reqStart time.Time

	// Request tracing for timing and connection info
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			timings.DNSLookup = time.Since(dnsStart)
			verboseInfo.DNSLatency = timings.DNSLookup
			if len(info.Addrs) > 0 {
				verboseInfo.ResolvedIP = info.Addrs[0].String()
			}
		},
		ConnectStart: func(_, _ string) {
			connectStart = time.Now()
		},
		ConnectDone: func(_, addr string, _ error) {
			timings.TCPConnection = time.Since(connectStart)
			// If DNS didn't resolve an IP, capture it from connect
			if verboseInfo.ResolvedIP == "" && addr != "" {
				verboseInfo.ResolvedIP = addr
			}
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			timings.TLSHandshake = time.Since(tlsStart)
		},
		GotFirstResponseByte: func() {
			if !reqStart.IsZero() {
				timings.ServerProcessing = time.Since(reqStart)
			}
		},
	}

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(ctx, trace), r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil {
		return nil, nil, timings, verboseInfo, err
	}

	// Calculate request size
	verboseInfo.RequestSize = int64(len(r.Body))
	if req.Header != nil {
		for k, vv := range req.Header {
			for _, v := range vv {
				verboseInfo.RequestSize += int64(len(k) + len(v) + 4) // key: value\r\n
			}
		}
	}

	if r.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Token)
	}
	if r.JSON {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}
	for _, h := range r.Headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// Cookie jar setup
	client := &http.Client{Timeout: 30 * time.Second}
	if r.CookieJar != "" {
		if globalCookieJar == nil {
			globalCookieJar, _ = cookiejar.New(nil)
			loadCookies(r.CookieJar, globalCookieJar)
		}
		client.Jar = globalCookieJar
	}

	start := time.Now()
	reqStart = start
	res, err := client.Do(req)
	timings.Total = time.Since(start)

	if err != nil {
		return nil, nil, timings, verboseInfo, err
	}

	// Capture protocol info
	verboseInfo.Protocol = res.Proto

	// Capture TLS info if HTTPS
	if res.TLS != nil {
		verboseInfo.TLSVersion = tlsVersionString(res.TLS.Version)
		verboseInfo.TLSCipher = tls.CipherSuiteName(res.TLS.CipherSuite)

		if len(res.TLS.PeerCertificates) > 0 {
			cert := res.TLS.PeerCertificates[0]
			verboseInfo.CertSubject = cert.Subject.CommonName
			verboseInfo.CertIssuer = cert.Issuer.CommonName
			verboseInfo.CertExpiry = cert.NotAfter
		}
	}

	bodyStart := time.Now()
	body, err := io.ReadAll(res.Body)
	timings.ContentTransfer = time.Since(bodyStart)

	if err != nil {
		return res, nil, timings, verboseInfo, err
	}

	// Calculate response size
	verboseInfo.ResponseSize = int64(len(body))
	if res.Header != nil {
		for k, vv := range res.Header {
			for _, v := range vv {
				verboseInfo.ResponseSize += int64(len(k) + len(v) + 4)
			}
		}
	}

	// Check if response was compressed
	contentEncoding := res.Header.Get("Content-Encoding")
	if contentEncoding == "gzip" || contentEncoding == "deflate" || contentEncoding == "br" {
		verboseInfo.Compressed = true
		// ContentLength is uncompressed size if available
		if res.ContentLength > 0 {
			verboseInfo.CompressionRatio = (1.0 - float64(len(body))/float64(res.ContentLength)) * 100
		}
	}

	// Save cookies if jar is enabled
	if r.CookieJar != "" && globalCookieJar != nil {
		saveCookies(r.CookieJar, globalCookieJar)
	}

	return res, body, timings, verboseInfo, nil
}

func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

func printVerbose(r Request, res *http.Response, t TimingInfo, v VerboseInfo) {
	cyan := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	// DNS & Connection Info
	if v.ResolvedIP != "" || v.Protocol != "" {
		fmt.Fprintf(os.Stderr, "\n%s\n", cyan("🌐 Connection Info:"))
		if v.ResolvedIP != "" {
			fmt.Fprintf(os.Stderr, "%s Resolved IP:  %s\n", gray("•"), blue(v.ResolvedIP))
		}
		if v.Protocol != "" {
			fmt.Fprintf(os.Stderr, "%s Protocol:     %s\n", gray("•"), blue(v.Protocol))
		}
		if v.TLSVersion != "" {
			fmt.Fprintf(os.Stderr, "%s TLS Version:  %s\n", gray("•"), green(v.TLSVersion))
			if v.TLSCipher != "" {
				fmt.Fprintf(os.Stderr, "%s TLS Cipher:   %s\n", gray("•"), gray(v.TLSCipher))
			}
		}
	}

	// TLS Certificate Info
	if v.CertSubject != "" {
		fmt.Fprintf(os.Stderr, "\n%s\n", cyan("🔐 TLS Certificate:"))
		fmt.Fprintf(os.Stderr, "%s Subject:      %s\n", gray("•"), yellow(v.CertSubject))
		if v.CertIssuer != "" {
			fmt.Fprintf(os.Stderr, "%s Issuer:       %s\n", gray("•"), gray(v.CertIssuer))
		}
		if !v.CertExpiry.IsZero() {
			daysRemaining := int(time.Until(v.CertExpiry).Hours() / 24)
			expiryColor := green
			if daysRemaining < 30 {
				expiryColor = yellow
			}
			fmt.Fprintf(os.Stderr, "%s Expires:      %s (%s)\n",
				gray("•"),
				v.CertExpiry.Format("2006-01-02"),
				expiryColor(fmt.Sprintf("%d days remaining", daysRemaining)))
		}
	}

	// Transfer Sizes
	if v.RequestSize > 0 || v.ResponseSize > 0 {
		fmt.Fprintf(os.Stderr, "\n%s\n", cyan("📦 Transfer Details:"))
		if v.RequestSize > 0 {
			fmt.Fprintf(os.Stderr, "%s Request:      %s\n", gray("•"), formatBytes(v.RequestSize))
		}
		if v.ResponseSize > 0 {
			sizeStr := formatBytes(v.ResponseSize)
			if v.Compressed {
				sizeStr += fmt.Sprintf(" (%s compressed, %.1f%% reduction)",
					yellow("gzip"),
					v.CompressionRatio)
			}
			fmt.Fprintf(os.Stderr, "%s Response:     %s\n", gray("•"), sizeStr)
		}
	}

	// Request headers
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("→ Request Headers:"))
	fmt.Fprintf(os.Stderr, "%s %s %s\n", gray(">"), r.Method, r.URL)
	fmt.Fprintf(os.Stderr, "%s Host: %s\n", gray(">"), res.Request.Host)
	fmt.Fprintf(os.Stderr, "%s User-Agent: mozzy/1.6.0\n", gray(">"))
	for k, v := range res.Request.Header {
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", gray(">"), k, strings.Join(v, ", "))
	}

	// Response headers
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("← Response Headers:"))
	fmt.Fprintf(os.Stderr, "%s HTTP/%d.%d %s\n", gray("<"), res.ProtoMajor, res.ProtoMinor, res.Status)
	for k, v := range res.Header {
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", gray("<"), k, strings.Join(v, ", "))
	}

	// Enhanced Timing breakdown with visual bars
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("⏱  Request Timeline:"))
	totalMs := float64(t.Total.Milliseconds())

	if t.DNSLookup > 0 {
		pct := float64(t.DNSLookup.Milliseconds()) / totalMs * 100
		bar := renderTimingBar(pct)
		fmt.Fprintf(os.Stderr, "%s DNS Lookup       %6s  %s  %s\n",
			gray("├─"), green(formatDuration(t.DNSLookup)), bar, gray(fmt.Sprintf("%.0f%%", pct)))
	}
	if t.TCPConnection > 0 {
		pct := float64(t.TCPConnection.Milliseconds()) / totalMs * 100
		bar := renderTimingBar(pct)
		fmt.Fprintf(os.Stderr, "%s TCP Connect      %6s  %s  %s\n",
			gray("├─"), green(formatDuration(t.TCPConnection)), bar, gray(fmt.Sprintf("%.0f%%", pct)))
	}
	if t.TLSHandshake > 0 {
		pct := float64(t.TLSHandshake.Milliseconds()) / totalMs * 100
		bar := renderTimingBar(pct)
		fmt.Fprintf(os.Stderr, "%s TLS Handshake    %6s  %s  %s\n",
			gray("├─"), green(formatDuration(t.TLSHandshake)), bar, gray(fmt.Sprintf("%.0f%%", pct)))
	}
	if t.ServerProcessing > 0 {
		pct := float64(t.ServerProcessing.Milliseconds()) / totalMs * 100
		bar := renderTimingBar(pct)
		fmt.Fprintf(os.Stderr, "%s Server Response  %6s  %s  %s\n",
			gray("├─"), green(formatDuration(t.ServerProcessing)), bar, gray(fmt.Sprintf("%.0f%%", pct)))
	}
	if t.ContentTransfer > 0 {
		pct := float64(t.ContentTransfer.Milliseconds()) / totalMs * 100
		bar := renderTimingBar(pct)
		fmt.Fprintf(os.Stderr, "%s Content Transfer %6s  %s  %s\n",
			gray("├─"), green(formatDuration(t.ContentTransfer)), bar, gray(fmt.Sprintf("%.0f%%", pct)))
	}
	fmt.Fprintf(os.Stderr, "%s Total            %6s  %s\n",
		gray("└─"), green(formatDuration(t.Total)), yellow("████████████████████"))

	// Performance Grading
	grade := CalculateGrade(t)
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("📊 Performance Grade:"))

	if grade.DNS != "" {
		fmt.Fprintf(os.Stderr, "%s DNS Lookup:      %s\n",
			gray("├─"), FormatGradeWithColor(grade.DNS))
	}
	if grade.TCP != "" {
		fmt.Fprintf(os.Stderr, "%s TCP Connect:     %s\n",
			gray("├─"), FormatGradeWithColor(grade.TCP))
	}
	if grade.TLS != "" {
		fmt.Fprintf(os.Stderr, "%s TLS Handshake:   %s\n",
			gray("├─"), FormatGradeWithColor(grade.TLS))
	}
	if grade.TTFB != "" {
		fmt.Fprintf(os.Stderr, "%s Server Response: %s\n",
			gray("├─"), FormatGradeWithColor(grade.TTFB))
	}
	fmt.Fprintf(os.Stderr, "%s Overall:         %s\n",
		gray("└─"), FormatGradeWithColor(grade.Overall))

	// Performance Recommendations
	recommendations := grade.GetRecommendations()
	if len(recommendations) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s\n", cyan("💡 Performance Insights:"))
		for _, rec := range recommendations {
			fmt.Fprintf(os.Stderr, "  %s\n", gray(rec))
		}
	}

	fmt.Fprintf(os.Stderr, "\n")
}

func renderTimingBar(percentage float64) string {
	barWidth := 20
	filled := int(percentage / 100 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return bar
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
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func loadCookies(filename string, jar *cookiejar.Jar) {
	// Simple cookie persistence - would need full implementation
	// For now, this is a placeholder
}

func saveCookies(filename string, jar *cookiejar.Jar) {
	// Simple cookie persistence - would need full implementation
	// For now, this is a placeholder
}
