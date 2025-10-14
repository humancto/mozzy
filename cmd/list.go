package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/collection"
	"github.com/humancto/mozzy/internal/ui"
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
			fmt.Println(ui.WarningBanner("No saved requests yet"))
			fmt.Println("\n" + ui.InfoStyle.Render("💡 Use 'mozzy save <name> <method> <url>' to add one"))
			return nil
		}

		fmt.Printf("\n%s\n\n", ui.TitleStyle.Render("📚 Saved Requests"))

		// Create table
		table := ui.NewTable([]string{"Name", "Method", "URL", "Description"})

		for _, req := range reqs {
			desc := req.Description
			if desc == "" {
				desc = "-"
			}
			// Truncate URL if too long
			url := req.URL
			if len(url) > 50 {
				url = url[:47] + "..."
			}
			table.AddRow([]string{req.Name, req.Method, url, desc})
		}

		fmt.Println(table.Render())

		fmt.Println(ui.InfoStyle.Render("💡 Tip: Run 'mozzy exec <name>' to execute a saved request"))
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
