package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/humancto/mozzy/internal/ui"
)

var (
	baseURL        string
	authToken      string
	headers        []string
	envName        string
	jqQuery        string
	timeoutStr     string
	failOnErr      bool
	noColor        bool
	forceColor     bool
	verbose        bool
	retryCount     int
	retryCondition string
	cookieJar      string
)

var rootCmd = &cobra.Command{
	Use:   "mozzy",
	Short: "mozzy: Postman-level JSON HTTP client for your terminal",
	Long:  ui.RenderBanner() + "\n\nmozzy is a JSON-native HTTP client with pretty output, inline querying, JWT tools, history, and request chaining.\n\n" + ui.RenderQuickStart(),
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when no args provided
		cmd.Help()
		fmt.Println("\n" + ui.InfoStyle.Render("💡 Tip: Try 'mozzy GET https://jsonplaceholder.typicode.com/users/1 --color' for a quick demo"))
		fmt.Println(ui.DimStyle.Render("   Run 'mozzy update' to check for new versions"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Custom help template with colors
	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	rootCmd.PersistentFlags().StringVar(&baseURL, "base", "", "Base URL, e.g. https://api.example.com (overridden by --env)")
	rootCmd.PersistentFlags().StringVar(&authToken, "auth", "", "Bearer token (Authorization: Bearer ...)")
	rootCmd.PersistentFlags().StringSliceVar(&headers, "header", nil, "Extra headers (repeat), e.g. --header 'X-Env: staging'")
	rootCmd.PersistentFlags().StringVar(&envName, "env", "", "Named environment from .mozzy.json")
	rootCmd.PersistentFlags().StringVar(&jqQuery, "jq", "", "Inline JSON query (JSONPath/JQ-lite)")
	rootCmd.PersistentFlags().StringVar(&timeoutStr, "timeout", "30s", "Request timeout, e.g. 2s, 500ms")
	rootCmd.PersistentFlags().BoolVar(&failOnErr, "fail", false, "Exit non-zero on HTTP status >= 400 (CI-friendly)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&forceColor, "color", false, "Force colored output even when not in a TTY")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show request/response headers and timing details")
	rootCmd.PersistentFlags().IntVar(&retryCount, "retry", 0, "Number of retry attempts on failure (with exponential backoff)")
	rootCmd.PersistentFlags().StringVar(&retryCondition, "retry-on", "", "Retry condition: 5xx, 429, >=500, network_error, etc. (comma-separated)")
	rootCmd.PersistentFlags().StringVar(&cookieJar, "cookie-jar", "", "File to store/load cookies for session management")

	// Custom usage template with colors
	rootCmd.SetUsageFunc(customUsage)

	// Auto-detect color support - disable if:
	// 1. --no-color flag is set
	// 2. NO_COLOR env var exists
	// Enable if --color flag is set
	cobra.OnInitialize(func() {
		if noColor || os.Getenv("NO_COLOR") != "" {
			color.NoColor = true
		}
		if forceColor {
			color.NoColor = false
			os.Setenv("CLICOLOR_FORCE", "1")
		}
	})
}

func customUsage(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()

	// Style definitions using lipgloss for beautiful output
	sectionStyle := ui.TitleStyle.Copy().
		MarginTop(1).
		MarginBottom(1)

	cmdNameStyle := ui.CodeStyle.Copy().
		Bold(true).
		Width(15)

	flagNameStyle := ui.SuccessStyle.Copy().
		Width(28)

	descStyle := ui.DimStyle

	usageStyle := ui.InfoStyle.Copy().
		Italic(true)

	// Usage section
	fmt.Fprintf(out, "\n%s\n", sectionStyle.Render("📋 Usage"))
	if cmd.Runnable() {
		fmt.Fprintf(out, "  %s %s\n", cmd.CommandPath(), usageStyle.Render("[flags]"))
	}
	if cmd.HasAvailableSubCommands() {
		fmt.Fprintf(out, "  %s %s\n", cmd.CommandPath(), usageStyle.Render("[command]"))
	}

	// Available Commands section
	if cmd.HasAvailableSubCommands() {
		fmt.Fprintf(out, "\n%s\n", sectionStyle.Render("⚡ Available Commands"))
		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() {
				continue
			}
			fmt.Fprintf(out, "  %s %s\n",
				cmdNameStyle.Render(c.Name()),
				descStyle.Render(c.Short))
		}
	}

	// Flags section
	if cmd.HasAvailableLocalFlags() {
		fmt.Fprintf(out, "\n%s\n", sectionStyle.Render("🚩 Flags"))
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			if f.Hidden {
				return
			}
			flagName := "--" + f.Name
			if f.Shorthand != "" && f.ShorthandDeprecated == "" {
				flagName = fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
			}
			fmt.Fprintf(out, "  %s %s\n",
				flagNameStyle.Render(flagName),
				descStyle.Render(f.Usage))
		})
	}

	// Global Flags section
	if cmd.HasAvailableInheritedFlags() {
		fmt.Fprintf(out, "\n%s\n", sectionStyle.Render("🌐 Global Flags"))
		cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
			if f.Hidden {
				return
			}
			flagName := "--" + f.Name
			if f.Shorthand != "" && f.ShorthandDeprecated == "" {
				flagName = fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
			}
			fmt.Fprintf(out, "  %s %s\n",
				flagNameStyle.Render(flagName),
				descStyle.Render(f.Usage))
		})
	}

	// Footer
	if cmd.HasSubCommands() {
		fmt.Fprintf(out, "\n%s\n",
			descStyle.Render(fmt.Sprintf("Use \"%s [command] --help\" for more information about a command.", cmd.CommandPath())))
	}

	return nil
}
