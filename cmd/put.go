package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	putJSON        string
	putBodyFile    string
	putFormData    []string
	putContentType string
)

var putCmd = &cobra.Command{
	Use:   "PUT <url-or-path>",
	Short: "Send an HTTP PUT request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var body []byte
		if putBodyFile != "" {
			b, err := os.ReadFile(putBodyFile); if err != nil { return err }
			body = b
		} else if putJSON != "" {
			if strings.HasPrefix(putJSON, "@") {
				p := strings.TrimPrefix(putJSON, "@")
				b, err := os.ReadFile(p); if err != nil { return err }
				body = b
			} else {
				body = []byte(putJSON)
			}
		}
		return runVerb(cmd, "PUT", args[0], body, putJSON != "" && putContentType == "")
	},
}

func init() {
	putCmd.Flags().StringVar(&putJSON, "json", "", "JSON payload (string or @file.json)")
	putCmd.Flags().StringVar(&putBodyFile, "file", "", "Raw body from file")
	putCmd.Flags().StringVar(&putContentType, "content-type", "", "Override Content-Type header")
	putCmd.Flags().StringArray("capture", nil, "Capture variables: name=.json.path (repeatable)")
	rootCmd.AddCommand(putCmd)
}
