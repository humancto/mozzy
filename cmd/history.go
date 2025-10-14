package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/history"
)

var (
	historyLimit int
	historyJSON  bool
)

var histCmd = &cobra.Command{
	Use:   "history",
	Short: "Show recent requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := history.Load()
		if err != nil {
			return err
		}

		// Limit results
		if historyLimit > 0 && historyLimit < len(entries) {
			entries = entries[len(entries)-historyLimit:]
		}

		// Raw JSON output
		if historyJSON {
			b, _ := json.MarshalIndent(entries, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		// Pretty formatted output
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println(cyan("📜 Request History"))
		fmt.Println()

		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			timeStr := entry.Timestamp.Format("Jan 02 15:04:05")

			// Status color
			statusStr := fmt.Sprintf("%d", entry.Status)
			if entry.Status >= 200 && entry.Status < 300 {
				statusStr = green(statusStr)
			} else if entry.Status >= 400 {
				statusStr = red(statusStr)
			} else {
				statusStr = yellow(statusStr)
			}

			// Duration formatting
			duration := time.Duration(entry.Duration)
			durationStr := gray(formatDuration(duration))

			fmt.Printf("%s %s %-6s %s %s\n",
				gray(timeStr),
				statusStr,
				cyan(entry.Method),
				entry.URL,
				durationStr,
			)
		}

		fmt.Println()
		fmt.Printf(gray("Showing %d requests. Use --limit to change or --json for raw output.\n"), len(entries))
		return nil
	},
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func init() {
	histCmd.Flags().IntVar(&historyLimit, "limit", 20, "Number of recent requests to show")
	histCmd.Flags().BoolVar(&historyJSON, "json", false, "Output raw JSON")
	rootCmd.AddCommand(histCmd)
}
