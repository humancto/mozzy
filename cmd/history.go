package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/history"
)

var histCmd = &cobra.Command{
	Use:   "history",
	Short: "Show recent requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := history.Load()
		if err != nil { return err }
		b, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(b))
		return nil
	},
}

func init() { rootCmd.AddCommand(histCmd) }
