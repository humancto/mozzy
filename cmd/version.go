package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mozzy version %s\n", version)
		if commit != "dev" {
			fmt.Printf("commit: %s\n", commit)
		}
		if date != "unknown" {
			fmt.Printf("built: %s\n", date)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
