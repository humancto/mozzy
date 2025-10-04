package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	baseURL    string
	authToken  string
	headers    []string
	envName    string
	jqQuery    string
	timeoutStr string
	failOnErr  bool
)

var rootCmd = &cobra.Command{
	Use:   "mozzy",
	Short: "mozzy: Postman-level JSON HTTP client for your terminal",
	Long:  "mozzy is a JSON-native HTTP client with pretty output, inline querying, JWT tools, history, and request chaining.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "base", "", "Base URL, e.g. https://api.example.com (overridden by --env)")
	rootCmd.PersistentFlags().StringVar(&authToken, "auth", "", "Bearer token (Authorization: Bearer ...)")
	rootCmd.PersistentFlags().StringSliceVar(&headers, "header", nil, "Extra headers (repeat), e.g. --header 'X-Env: staging'")
	rootCmd.PersistentFlags().StringVar(&envName, "env", "", "Named environment from .mozzy.json")
	rootCmd.PersistentFlags().StringVar(&jqQuery, "jq", "", "Inline JSON query (JSONPath/JQ-lite; stubbed)")
	rootCmd.PersistentFlags().StringVar(&timeoutStr, "timeout", "30s", "Request timeout, e.g. 2s, 500ms")
	rootCmd.PersistentFlags().BoolVar(&failOnErr, "fail", false, "Exit non-zero on HTTP status >= 400 (CI-friendly)")
}
