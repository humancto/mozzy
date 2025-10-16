package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/proxy"
)

var (
	proxyVerbose bool
)

var proxyCmd = &cobra.Command{
	Use:   "proxy [port]",
	Short: "Start an HTTP proxy server",
	Long: `Start an HTTP proxy server to intercept and inspect traffic.

Perfect for:
- Debugging mobile apps
- Inspecting API calls
- Testing network conditions
- Analyzing third-party integrations

Examples:
  mozzy proxy 8888                    # Start on port 8888
  mozzy proxy 8888 --verbose          # With detailed logging

Configure your browser or app:
  HTTP Proxy: localhost:8888

For mobile devices, use your computer's IP address.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProxy,
}

func runProxy(cmd *cobra.Command, args []string) error {
	port := 8888 // default port

	if len(args) > 0 {
		p, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid port: %s", args[0])
		}
		port = p
	}

	server := proxy.NewServer(port, proxyVerbose)
	return server.Start()
}

func init() {
	proxyCmd.Flags().BoolVarP(&proxyVerbose, "verbose", "v", false, "Show detailed request/response information")

	rootCmd.AddCommand(proxyCmd)
}
