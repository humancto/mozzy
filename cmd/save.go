package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/collection"
)

var (
	saveDescription string
)

var saveCmd = &cobra.Command{
	Use:   "save <name> <method> <url>",
	Short: "Save a request to your collection",
	Long: `Save a named request for later reuse.

Examples:
  mozzy save login POST https://api.example.com/auth --json '{"user":"alice"}'
  mozzy save get-users GET https://api.example.com/users --header "X-API-Key: secret"`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		method := strings.ToUpper(args[1])
		url := args[2]

		// Collect headers
		headerMap := make(map[string]string)
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Get body from flags
		body := ""
		if method != "GET" && method != "DELETE" {
			// Check for JSON flag from parent commands
			if postJSON != "" {
				body = postJSON
			}
		}

		coll, err := collection.Load()
		if err != nil {
			return err
		}

		req := collection.Request{
			Name:        name,
			Method:      method,
			URL:         url,
			Headers:     headerMap,
			Body:        body,
			Description: saveDescription,
		}

		if err := coll.Add(req); err != nil {
			return err
		}

		successColor := color.New(color.FgGreen, color.Bold)
		fmt.Printf("%s %s\n",
			color.GreenString("✓"),
			successColor.Sprintf("Saved request '%s' to collection", name),
		)
		return nil
	},
}

func init() {
	saveCmd.Flags().StringVar(&saveDescription, "desc", "", "Description of this request")
	rootCmd.AddCommand(saveCmd)
}
