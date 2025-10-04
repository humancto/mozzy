package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method  string
	URL     string
	Headers []string // "Key: Value"
	Token   string
	Body    []byte
	JSON    bool
}

func Do(ctx context.Context, r Request) (*http.Response, []byte, time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, bytes.NewReader(r.Body))
	if err != nil { return nil, nil, 0, err }

	if r.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Token)
	}
	if r.JSON {
		// don't override if provided through --header
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

	client := &http.Client{ Timeout: 30 * time.Second }
	start := time.Now()
	res, err := client.Do(req)
	dur := time.Since(start)
	if err != nil { return nil, nil, dur, err }

	body, err := io.ReadAll(res.Body)
	if err != nil { return res, nil, dur, err }
	return res, body, dur, nil
}
