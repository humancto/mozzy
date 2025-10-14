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

var globalCookieJar *cookiejar.Jar

func Do(ctx context.Context, r Request) (*http.Response, []byte, time.Duration, error) {
	var timings TimingInfo
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

		res, body, timings, err = doRequest(ctx, r)

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
		printVerbose(r, res, timings)
	}

	return res, body, timings.Total, err
}

func doRequest(ctx context.Context, r Request) (*http.Response, []byte, TimingInfo, error) {
	var timings TimingInfo
	var dnsStart, connectStart, tlsStart, reqStart time.Time

	// Request tracing for timing
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			timings.DNSLookup = time.Since(dnsStart)
		},
		ConnectStart: func(_, _ string) {
			connectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			timings.TCPConnection = time.Since(connectStart)
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
		return nil, nil, timings, err
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
		return nil, nil, timings, err
	}

	bodyStart := time.Now()
	body, err := io.ReadAll(res.Body)
	timings.ContentTransfer = time.Since(bodyStart)

	if err != nil {
		return res, nil, timings, err
	}

	// Save cookies if jar is enabled
	if r.CookieJar != "" && globalCookieJar != nil {
		saveCookies(r.CookieJar, globalCookieJar)
	}

	return res, body, timings, nil
}

func printVerbose(r Request, res *http.Response, t TimingInfo) {
	cyan := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	// Request headers
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("→ Request Headers:"))
	fmt.Fprintf(os.Stderr, "%s %s %s\n", gray(">"), r.Method, r.URL)
	fmt.Fprintf(os.Stderr, "%s Host: %s\n", gray(">"), res.Request.Host)
	fmt.Fprintf(os.Stderr, "%s User-Agent: mozzy/1.0\n", gray(">"))
	for k, v := range res.Request.Header {
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", gray(">"), k, strings.Join(v, ", "))
	}

	// Response headers
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("← Response Headers:"))
	fmt.Fprintf(os.Stderr, "%s HTTP/%d.%d %s\n", gray("<"), res.ProtoMajor, res.ProtoMinor, res.Status)
	for k, v := range res.Header {
		fmt.Fprintf(os.Stderr, "%s %s: %s\n", gray("<"), k, strings.Join(v, ", "))
	}

	// Timing breakdown
	fmt.Fprintf(os.Stderr, "\n%s\n", cyan("⏱  Timing Breakdown:"))
	if t.DNSLookup > 0 {
		fmt.Fprintf(os.Stderr, "%s DNS Lookup:        %s\n", gray("•"), green(formatDuration(t.DNSLookup)))
	}
	if t.TCPConnection > 0 {
		fmt.Fprintf(os.Stderr, "%s TCP Connection:    %s\n", gray("•"), green(formatDuration(t.TCPConnection)))
	}
	if t.TLSHandshake > 0 {
		fmt.Fprintf(os.Stderr, "%s TLS Handshake:     %s\n", gray("•"), green(formatDuration(t.TLSHandshake)))
	}
	if t.ServerProcessing > 0 {
		fmt.Fprintf(os.Stderr, "%s Server Processing: %s\n", gray("•"), green(formatDuration(t.ServerProcessing)))
	}
	if t.ContentTransfer > 0 {
		fmt.Fprintf(os.Stderr, "%s Content Transfer:  %s\n", gray("•"), green(formatDuration(t.ContentTransfer)))
	}
	fmt.Fprintf(os.Stderr, "%s Total:             %s\n", gray("•"), green(formatDuration(t.Total)))
	fmt.Fprintf(os.Stderr, "\n")
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
