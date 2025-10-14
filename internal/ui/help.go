package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderBanner renders the mozzy ASCII banner
func RenderBanner() string {
	banner := `
   __  __    ___    ____  ____  __  __
  |  \/  |  / _ \  |_  / |_  / \ \/ /
  | |\/| | | | | |  / /   / /   \  /
  | |  | | | |_| | / /_  / /_   /  \
  |_|  |_|  \___/ /____||____| /_/\_\

  🦣  Postman for your Terminal
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Beautiful HTTP client with superpowers
  `

	// Create gradient effect with different colors
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF")). // Bright cyan
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB86C")). // Orange
		Italic(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")) // Light cyan

	lines := strings.Split(banner, "\n")
	var result strings.Builder

	for i, line := range lines {
		if i >= 1 && i <= 5 {
			// ASCII art lines
			result.WriteString(titleStyle.Render(line))
		} else if i == 7 {
			// Emoji and title line
			result.WriteString(subtitleStyle.Render(line))
		} else if i == 8 {
			// Separator line
			result.WriteString(lipgloss.NewStyle().Foreground(ColorBorder).Render(line))
		} else if i == 9 {
			// Description
			result.WriteString(descStyle.Render(line))
		} else {
			result.WriteString(line)
		}
		result.WriteString("\n")
	}

	return result.String()
}

// RenderQuickStart renders quick start examples
func RenderQuickStart() string {
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B")). // Green
		Bold(true)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BD93F9")) // Purple

	commentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4")). // Gray
		Italic(true)

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF79C6")). // Pink
		Bold(true).
		Render("⚡ Quick Start")

	examples := []string{
		fmt.Sprintf("%s %s  %s",
			promptStyle.Render("$"),
			commandStyle.Render("mozzy GET https://api.example.com --color"),
			commentStyle.Render("# Colorized JSON")),
		fmt.Sprintf("%s %s  %s",
			promptStyle.Render("$"),
			commandStyle.Render("mozzy POST /users --json '{\"name\":\"Alice\"}'"),
			commentStyle.Render("# Send JSON data")),
		fmt.Sprintf("%s %s  %s",
			promptStyle.Render("$"),
			commandStyle.Render("mozzy save my-api GET /api"),
			commentStyle.Render("# Save for later")),
	}

	return fmt.Sprintf("\n%s\n  %s\n", title, strings.Join(examples, "\n  "))
}

// RenderCommandGroup renders a group of commands with a title
func RenderCommandGroup(title string, commands map[string]string) string {
	var sb strings.Builder

	// Box with title
	titleStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		MarginBottom(1)

	sb.WriteString(titleStyle.Render("╭─ " + title + " "))
	sb.WriteString(strings.Repeat("─", 50-len(title)))
	sb.WriteString("╮\n")

	// Commands
	cmdStyle := lipgloss.NewStyle().
		Foreground(ColorHighlight).
		Bold(true).
		Width(15)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D4D4D4"))

	for cmd, desc := range commands {
		sb.WriteString("│  ")
		sb.WriteString(cmdStyle.Render(cmd))
		sb.WriteString(descStyle.Render(desc))
		sb.WriteString("\n")
	}

	sb.WriteString("╰")
	sb.WriteString(strings.Repeat("─", 58))
	sb.WriteString("╯\n")

	return sb.String()
}

// RenderHelpSections renders categorized help sections
func RenderHelpSections() string {
	var sections []string

	// HTTP Commands
	httpCmds := map[string]string{
		"GET":    "Send HTTP GET request",
		"POST":   "Send HTTP POST request",
		"PUT":    "Send HTTP PUT request",
		"PATCH":  "Send HTTP PATCH request",
		"DELETE": "Send HTTP DELETE request",
	}
	sections = append(sections, RenderCommandGroup("HTTP Commands", httpCmds))

	// Collection Management
	collectionCmds := map[string]string{
		"save": "Save request to collection",
		"list": "List saved requests",
		"exec": "Execute saved request",
	}
	sections = append(sections, RenderCommandGroup("Collection Management", collectionCmds))

	// Advanced Features
	advancedCmds := map[string]string{
		"run":      "Execute YAML workflows",
		"test":     "Run workflow as test suite",
		"jwt":      "JWT decode/verify/sign",
		"diff":     "Compare JSON responses",
		"download": "Download files with progress",
		"upload":   "Upload files with multipart",
	}
	sections = append(sections, RenderCommandGroup("Advanced Features", advancedCmds))

	return strings.Join(sections, "\n")
}

// RenderTip renders a helpful tip
func RenderTip(message string) string {
	tipStyle := lipgloss.NewStyle().
		Foreground(ColorWarning).
		Italic(true)

	return fmt.Sprintf("\n%s %s\n", tipStyle.Render("💡 Tip:"), message)
}
