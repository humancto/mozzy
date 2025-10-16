package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/formatter"
	"github.com/humancto/mozzy/internal/history"
	"github.com/humancto/mozzy/internal/httpclient"
	"github.com/humancto/mozzy/internal/vars"
)

func runVerb(cmd *cobra.Command, method string, target string, body []byte, isJSON bool) error {
	// Resolve base/env
	resolvedBase := vars.ResolveBase(baseURL, envName)
	if resolvedBase != "" {
		u, err := url.Parse(resolvedBase)
		if err != nil { return err }
		p, err := url.Parse(target)
		if err != nil { return err }
		target = u.ResolveReference(p).String()
	}

	// Interpolate {{vars}} into URL and headers
	target = vars.Interpolate(target)
	hdrs := make([]string, len(headers))
	for i, h := range headers { hdrs[i] = vars.Interpolate(h) }
	token := vars.Interpolate(authToken)

	// Timeout
	dur, err := time.ParseDuration(timeoutStr)
	if err != nil { dur = 30 * time.Second }
	ctx, cancel := context.WithTimeout(cmd.Context(), dur)
	defer cancel()

	req := httpclient.Request{
		Method:         method,
		URL:            target,
		Headers:        hdrs,
		Token:          token,
		Body:           body,
		JSON:           isJSON,
		Verbose:        verbose,
		RetryCount:     retryCount,
		RetryCondition: retryCondition,
		CookieJar:      cookieJar,
		Throttle:       throttle,
	}

	res, resBody, ms, err := httpclient.Do(ctx, req)
	if err != nil { return err }
	defer res.Body.Close()

	_ = history.Append(history.Entry{
		Timestamp: time.Now(),
		Method:    method,
		URL:       target,
		Status:    res.StatusCode,
		Duration:  ms,
		BodySize:  len(body),
	})

	formatter.PrintStatusLine(method, target, res.StatusCode, ms)

	if err := formatter.PrintJSONOrText(resBody, jqQuery); err != nil { return err }

	if failOnErr && res.StatusCode >= 400 {
		os.Exit(1)
	}

	// Capture support: --capture name=.json.path
	caps, _ := cmd.Flags().GetStringArray("capture")
	for _, c := range caps {
		if err := vars.Capture(resBody, c); err != nil {
			fmt.Fprintf(os.Stderr, "warn: capture failed: %v\n", err)
		}
	}
	return nil
}

var getCmd = &cobra.Command{
	Use:   "GET <url-or-path>",
	Short: "Send an HTTP GET request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVerb(cmd, "GET", args[0], nil, false)
	},
}

func init() {
	getCmd.Flags().StringArray("capture", nil, "Capture variables: name=.json.path (repeatable)")
	rootCmd.AddCommand(getCmd)
}
