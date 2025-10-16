package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/collection"
	"github.com/humancto/mozzy/internal/formatter"
	"github.com/humancto/mozzy/internal/history"
	"github.com/humancto/mozzy/internal/httpclient"
	"github.com/humancto/mozzy/internal/ui"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive mode to browse and execute history or saved requests",
	Long: `Launch interactive mode to browse request history with arrow keys.

Use ↑/↓ to navigate, Enter to execute, Esc/Ctrl+C to quit.

Examples:
  mozzy interactive          # Browse history
  mozzy interactive --saved  # Browse saved requests
  mozzy i                    # short alias`,
	Aliases: []string{"i"},
	RunE: func(cmd *cobra.Command, args []string) error {
		useSaved, _ := cmd.Flags().GetBool("saved")

		if useSaved {
			return interactiveSavedRequests(cmd)
		}
		return interactiveHistory(cmd)
	},
}

func interactiveHistory(cmd *cobra.Command) error {
	entries, err := history.Load()
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("no history found - make some requests first!")
	}

	// Reverse to show most recent first
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Method | cyan }} {{ .URL | faint }} {{ .StatusDisplay | yellow }} {{ .DurationDisplay | faint }}",
		Inactive: "  {{ .Method | cyan }} {{ .URL | faint }} {{ .StatusDisplay | yellow }} {{ .DurationDisplay | faint }}",
		Selected: "✓ {{ .Method | cyan }} {{ .URL | faint }}",
	}

	// Create displayable entries
	type DisplayEntry struct {
		history.Entry
		StatusDisplay   string
		DurationDisplay string
	}

	displayEntries := make([]DisplayEntry, len(entries))
	for i, entry := range entries {
		displayEntries[i] = DisplayEntry{
			Entry:           entry,
			StatusDisplay:   fmt.Sprintf("(%d)", entry.Status),
			DurationDisplay: fmt.Sprintf("%dms", entry.Duration.Milliseconds()),
		}
	}

	prompt := promptui.Select{
		Label:     ui.TitleStyle.Render("🔍 Select a request from history"),
		Items:     displayEntries,
		Templates: templates,
		Size:      15,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return err
	}

	selected := entries[idx]
	return executeHistoryEntry(cmd, selected)
}

func interactiveSavedRequests(cmd *cobra.Command) error {
	coll, err := collection.Load()
	if err != nil {
		return err
	}

	requests := coll.List()
	if len(requests) == 0 {
		return fmt.Errorf("no saved requests - save one with 'mozzy save <name> <method> <url>'")
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Name | cyan | bold }} {{ .Method | green }} {{ .URL | faint }} {{ if .Description }}({{ .Description | faint }}){{ end }}",
		Inactive: "  {{ .Name | cyan }} {{ .Method | green }} {{ .URL | faint }} {{ if .Description }}({{ .Description | faint }}){{ end }}",
		Selected: "✓ {{ .Name | cyan | bold }}",
	}

	prompt := promptui.Select{
		Label:     ui.TitleStyle.Render("📚 Select a saved request"),
		Items:     requests,
		Templates: templates,
		Size:      15,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return err
	}

	selected := requests[idx]
	return executeSavedRequest(cmd, selected)
}

func executeSavedRequest(cmd *cobra.Command, req collection.Request) error {
	methodColor := color.New(color.FgMagenta, color.Bold)
	fmt.Printf("\n%s %s %s %s\n\n",
		ui.TitleStyle.Render("🚀 Executing:"),
		methodColor.Sprint(req.Method),
		req.URL,
		ui.DimStyle.Render(fmt.Sprintf("(%s)", req.Name)))

	dur, _ := time.ParseDuration(timeoutStr)
	if dur == 0 {
		dur = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), dur)
	defer cancel()

	// Build headers
	hdrs := []string{}
	for k, v := range req.Headers {
		hdrs = append(hdrs, fmt.Sprintf("%s: %s", k, v))
	}

	httpReq := httpclient.Request{
		Method:  req.Method,
		URL:     req.URL,
		Headers: hdrs,
		Body:    []byte(req.Body),
		JSON:    req.Body != "",
		Verbose: verbose,
	}

	res, resBody, ms, err := httpclient.Do(ctx, httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Record in history
	_ = history.Append(history.Entry{
		Timestamp: time.Now(),
		Method:    req.Method,
		URL:       req.URL,
		Status:    res.StatusCode,
		Duration:  ms,
	})

	formatter.PrintStatusLine(req.Method, req.URL, res.StatusCode, ms)

	if err := formatter.PrintJSONOrText(resBody, jqQuery); err != nil {
		return err
	}

	if failOnErr && res.StatusCode >= 400 {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}

	return nil
}

func executeHistoryEntry(cmd *cobra.Command, entry history.Entry) error {
	methodColor := color.New(color.FgMagenta, color.Bold)
	fmt.Printf("\n%s %s %s\n\n",
		ui.TitleStyle.Render("🚀 Re-executing:"),
		methodColor.Sprint(entry.Method),
		entry.URL)

	dur, _ := time.ParseDuration(timeoutStr)
	if dur == 0 {
		dur = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), dur)
	defer cancel()

	httpReq := httpclient.Request{
		Method:  entry.Method,
		URL:     entry.URL,
		Headers: []string{},
		Verbose: verbose,
	}

	res, resBody, ms, err := httpclient.Do(ctx, httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Record in history
	_ = history.Append(history.Entry{
		Timestamp: time.Now(),
		Method:    entry.Method,
		URL:       entry.URL,
		Status:    res.StatusCode,
		Duration:  ms,
	})

	formatter.PrintStatusLine(entry.Method, entry.URL, res.StatusCode, ms)

	if err := formatter.PrintJSONOrText(resBody, jqQuery); err != nil {
		return err
	}

	if failOnErr && res.StatusCode >= 400 {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}

	return nil
}

func init() {
	interactiveCmd.Flags().Bool("saved", false, "Browse saved requests instead of history")
	rootCmd.AddCommand(interactiveCmd)
}
