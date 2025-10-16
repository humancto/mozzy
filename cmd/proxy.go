package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/proxy"
)

var (
	proxyVerbose      bool
	proxyHTTPS        bool
	exportCert        bool
	certInfo          bool
	recordFile        string
	injectHeaders     []string
	filterDomain      string
	filterMethods     string
	filterErrorsOnly  bool
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

	// Set recording file
	server.RecordFile = recordFile

	// Parse and set inject headers
	if len(injectHeaders) > 0 {
		server.InjectHeaders = make(map[string]string)
		for _, h := range injectHeaders {
			parts := splitHeader(h)
			if len(parts) == 2 {
				server.InjectHeaders[parts[0]] = parts[1]
			}
		}
	}

	// Set filters
	server.FilterDomain = filterDomain
	server.FilterErrors = filterErrorsOnly
	if filterMethods != "" {
		server.FilterMethods = splitMethods(filterMethods)
	}

	return server.Start()
}

// splitHeader splits "Key: Value" into ["Key", "Value"]
func splitHeader(h string) []string {
	for i := 0; i < len(h); i++ {
		if h[i] == ':' {
			key := h[:i]
			value := h[i+1:]
			// Trim spaces
			if len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}
			return []string{key, value}
		}
	}
	return []string{h}
}

// splitMethods splits "GET,POST,PUT" into ["GET", "POST", "PUT"]
func splitMethods(m string) []string {
	var result []string
	current := ""
	for i := 0; i < len(m); i++ {
		if m[i] == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(m[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func init() {
	proxyCmd.Flags().BoolVarP(&proxyVerbose, "verbose", "v", false, "Show detailed request/response information")
	proxyCmd.Flags().BoolVar(&proxyHTTPS, "https", false, "Enable HTTPS interception (requires CA certificate installation)")
	proxyCmd.Flags().BoolVar(&exportCert, "export-cert", false, "Export CA certificate in PEM format")
	proxyCmd.Flags().BoolVar(&certInfo, "cert-info", false, "Show CA certificate information")

	// Recording and filtering
	proxyCmd.Flags().StringVarP(&recordFile, "record", "r", "", "Record all traffic to HAR file")
	proxyCmd.Flags().StringArrayVarP(&injectHeaders, "inject-header", "H", []string{}, "Inject header into requests (can be used multiple times)")
	proxyCmd.Flags().StringVar(&filterDomain, "filter-domain", "", "Only log requests matching domain (substring match)")
	proxyCmd.Flags().StringVar(&filterMethods, "filter-methods", "", "Only log specific methods (comma-separated: GET,POST)")
	proxyCmd.Flags().BoolVar(&filterErrorsOnly, "errors-only", false, "Only log requests with 4xx/5xx status codes")

	rootCmd.AddCommand(proxyCmd)
}
