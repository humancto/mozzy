package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/proxy"
)

var (
	proxyVerbose bool
	proxyHTTPS   bool
	exportCert   bool
	certInfo     bool
)

var proxyCmd = &cobra.Command{
	Use:   "proxy [port]",
	Short: "Start an HTTP/HTTPS proxy server",
	Long: `Start an HTTP/HTTPS proxy server to intercept and inspect traffic.

Perfect for:
- Debugging mobile apps
- Inspecting API calls
- Testing network conditions
- Analyzing third-party integrations

Examples:
  mozzy proxy 8888                    # HTTP proxy on port 8888
  mozzy proxy 8888 --https            # HTTPS proxy with SSL interception
  mozzy proxy 8888 --verbose          # With detailed logging
  mozzy proxy --export-cert           # Export CA certificate for installation
  mozzy proxy --cert-info             # Show CA certificate information

Configure your browser or app:
  HTTP Proxy: localhost:8888
  HTTPS Proxy: localhost:8888

For HTTPS mode, you must install the CA certificate:
  mozzy proxy --export-cert > mozzy-ca.pem
  Then install mozzy-ca.pem in your system's trusted certificates.

For mobile devices, use your computer's IP address.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProxy,
}

func runProxy(cmd *cobra.Command, args []string) error {
	// Handle certificate export
	if exportCert {
		ca, err := proxy.GetCA()
		if err != nil {
			return fmt.Errorf("failed to get CA: %w", err)
		}
		certPEM, err := ca.ExportCert()
		if err != nil {
			return fmt.Errorf("failed to export certificate: %w", err)
		}
		fmt.Print(string(certPEM))
		return nil
	}

	// Handle certificate info
	if certInfo {
		ca, err := proxy.GetCA()
		if err != nil {
			return fmt.Errorf("failed to get CA: %w", err)
		}
		fmt.Println(ca.GetInfo())
		return nil
	}

	port := 8888 // default port

	if len(args) > 0 {
		p, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid port: %s", args[0])
		}
		port = p
	}

	server := proxy.NewServer(port, proxyVerbose, proxyHTTPS)
	return server.Start()
}

func init() {
	proxyCmd.Flags().BoolVarP(&proxyVerbose, "verbose", "v", false, "Show detailed request/response information")
	proxyCmd.Flags().BoolVar(&proxyHTTPS, "https", false, "Enable HTTPS interception (requires CA certificate installation)")
	proxyCmd.Flags().BoolVar(&exportCert, "export-cert", false, "Export CA certificate in PEM format")
	proxyCmd.Flags().BoolVar(&certInfo, "cert-info", false, "Show CA certificate information")

	rootCmd.AddCommand(proxyCmd)
}
