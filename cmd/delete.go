package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "DELETE <url-or-path>",
	Short: "Send an HTTP DELETE request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().StringArray("capture", nil, "Capture variables: name=.json.path (repeatable)")
		return runVerb(cmd, "DELETE", args[0], nil, false)
	},
}

func init() { rootCmd.AddCommand(deleteCmd) }
