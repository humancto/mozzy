package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/collection"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved requests in your collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		coll, err := collection.Load()
		if err != nil {
			return err
		}

		reqs := coll.List()
		if len(reqs) == 0 {
			fmt.Println(color.YellowString("📭 No saved requests yet. Use 'mozzy save <name> <method> <url>' to add one."))
			return nil
		}

		titleColor := color.New(color.FgCyan, color.Bold)
		methodColor := color.New(color.FgMagenta, color.Bold)
		nameColor := color.New(color.FgGreen, color.Bold)

		fmt.Printf("\n%s\n\n", titleColor.Sprint("📚 Saved Requests"))

		for _, req := range reqs {
			fmt.Printf("%s %s %s\n",
				nameColor.Sprint(req.Name),
				methodColor.Sprint(req.Method),
				req.URL,
			)
			if req.Description != "" {
				fmt.Printf("   %s\n", color.New(color.FgWhite, color.Faint).Sprint(req.Description))
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
