package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	patchJSON        string
	patchBodyFile    string
	patchFormData    []string
	patchContentType string
)

var patchCmd = &cobra.Command{
	Use:   "PATCH <url-or-path>",
	Short: "Send an HTTP PATCH request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().StringArray("capture", nil, "Capture variables: name=.json.path (repeatable)")
		var body []byte
		if patchBodyFile != "" {
			b, err := os.ReadFile(patchBodyFile); if err != nil { return err }
			body = b
		} else if patchJSON != "" {
			if strings.HasPrefix(patchJSON, "@") {
				p := strings.TrimPrefix(patchJSON, "@")
				b, err := os.ReadFile(p); if err != nil { return err }
				body = b
			} else {
				body = []byte(patchJSON)
			}
		}
		return runVerb(cmd, "PATCH", args[0], body, patchJSON != "" && patchContentType == "")
	},
}

func init() {
	patchCmd.Flags().StringVar(&patchJSON, "json", "", "JSON payload (string or @file.json)")
	patchCmd.Flags().StringVar(&patchBodyFile, "file", "", "Raw body from file")
	patchCmd.Flags().StringVar(&patchContentType, "content-type", "", "Override Content-Type header")
	rootCmd.AddCommand(patchCmd)
}
