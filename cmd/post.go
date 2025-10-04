package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	postJSON        string
	postBodyFile    string
	postFormData    []string
	postContentType string
)

var postCmd = &cobra.Command{
	Use:   "POST <url-or-path>",
	Short: "Send an HTTP POST request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().StringArray("capture", nil, "Capture variables: name=.json.path (repeatable)")
		var body []byte
		if postBodyFile != "" {
			b, err := os.ReadFile(postBodyFile); if err != nil { return err }
			body = b
		} else if postJSON != "" {
			if strings.HasPrefix(postJSON, "@") {
				p := strings.TrimPrefix(postJSON, "@")
				b, err := os.ReadFile(p); if err != nil { return err }
				body = b
			} else {
				body = []byte(postJSON)
			}
		}
		return runVerb(cmd, "POST", args[0], body, postJSON != "" && postContentType == "")
	},
}

func init() {
	postCmd.Flags().StringVar(&postJSON, "json", "", "JSON payload (string or @file.json)")
	postCmd.Flags().StringVar(&postBodyFile, "file", "", "Raw body from file")
	postCmd.Flags().StringVar(&postContentType, "content-type", "", "Override Content-Type header")
	rootCmd.AddCommand(postCmd)
}
