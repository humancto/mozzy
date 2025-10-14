package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderBanner renders the mozzy ASCII banner
func RenderBanner() string {
	banner := `
╭─────────────────────────────────────────────────────────╮
│                                                         │
│  🦣  MOZZY - Postman for your Terminal                 │
│      Beautiful HTTP client with superpowers            │
│                                                         │
╰─────────────────────────────────────────────────────────╯`

	bannerStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	return bannerStyle.Render(banner)
}

// RenderQuickStart renders quick start examples
func RenderQuickStart() string {
	title := TitleStyle.Render("Quick Start")
	examples := []string{
		CodeStyle.Render("mozzy GET https://api.example.com --color"),
		CodeStyle.Render("mozzy POST /users --json '{\"name\":\"Alice\"}'"),
		CodeStyle.Render("mozzy save my-request GET /api"),
	}

	return fmt.Sprintf("%s\n  %s\n", title, strings.Join(examples, "\n  "))
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
