package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/collection"
	"github.com/humancto/mozzy/internal/formatter"
	"github.com/humancto/mozzy/internal/history"
	"github.com/humancto/mozzy/internal/httpclient"
	"github.com/humancto/mozzy/internal/vars"
)

var execCmd = &cobra.Command{
	Use:   "exec <name>",
	Short: "Execute a saved request from your collection",
	Long: `Run a previously saved request by name.

Examples:
  mozzy exec login
  mozzy exec get-users --auth $TOKEN`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		coll, err := collection.Load()
		if err != nil {
			return err
		}

		req, err := coll.Get(name)
		if err != nil {
			return err
		}

		// Print what we're running
		infoColor := color.New(color.FgCyan)
		fmt.Printf("%s %s\n\n", color.CyanString("🚀"), infoColor.Sprintf("Executing saved request: %s", name))

		// Build headers from saved request + CLI overrides
		hdrs := []string{}
		for k, v := range req.Headers {
			hdrs = append(hdrs, fmt.Sprintf("%s: %s", k, vars.Interpolate(v)))
		}
		// Add any additional headers from command line
		for _, h := range headers {
			hdrs = append(hdrs, vars.Interpolate(h))
		}

		// Auth token override
		token := vars.Interpolate(authToken)

		// Interpolate URL
		url := vars.Interpolate(req.URL)

		// Interpolate body
		body := []byte(vars.Interpolate(req.Body))

		// Timeout
		dur, err := time.ParseDuration(timeoutStr)
		if err != nil {
			dur = 30 * time.Second
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), dur)
		defer cancel()

		// Execute request
		httpReq := httpclient.Request{
			Method:  req.Method,
			URL:     url,
			Headers: hdrs,
			Token:   token,
			Body:    body,
			JSON:    req.Body != "" && strings.HasPrefix(strings.TrimSpace(req.Body), "{"),
		}

		res, resBody, ms, err := httpclient.Do(ctx, httpReq)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		_ = history.Append(history.Entry{
			Timestamp: time.Now(),
			Method:    req.Method,
			URL:       url,
			Status:    res.StatusCode,
			Duration:  ms,
			BodySize:  len(body),
		})

		formatter.PrintStatusLine(req.Method, url, res.StatusCode, ms)

		if err := formatter.PrintJSONOrText(resBody, jqQuery); err != nil {
			return err
		}

		if failOnErr && res.StatusCode >= 400 {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
