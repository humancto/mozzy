package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#3B8EEA") // Blue
	ColorSuccess   = lipgloss.Color("#4EC9B0") // Green
	ColorError     = lipgloss.Color("#F44747") // Red
	ColorWarning   = lipgloss.Color("#FFCC66") // Yellow
	ColorInfo      = lipgloss.Color("#4FC1FF") // Cyan
	ColorMuted     = lipgloss.Color("#6A737D") // Gray
	ColorBorder    = lipgloss.Color("#444444") // Dark gray
	ColorHighlight = lipgloss.Color("#C586C0") // Purple
)

// Base styles
var (
	// BoxStyle for bordered content
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	// BoldStyle for emphasis
	BoldStyle = lipgloss.NewStyle().Bold(true)

	// DimStyle for de-emphasized text
	DimStyle = lipgloss.NewStyle().Foreground(ColorMuted)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	// WarningStyle for warnings
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	// InfoStyle for info messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// TitleStyle for section titles
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// CodeStyle for code/commands
	CodeStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight)
)

// Banner styles
var (
	SuccessBannerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		Foreground(ColorSuccess).
		Padding(0, 2).
		Bold(true)

	ErrorBannerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorError).
		Foreground(ColorError).
		Padding(0, 2).
		Bold(true)

	InfoBannerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorInfo).
		Foreground(ColorInfo).
		Padding(0, 2)
)

// Table styles
var (
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Padding(0, 1)

	TableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	TableBorderStyle = lipgloss.NewStyle().
				Foreground(ColorBorder)
)

// Helper functions

// SuccessBanner creates a success banner with a message
func SuccessBanner(message string) string {
	return SuccessBannerStyle.Render("✅ " + message)
}

// ErrorBanner creates an error banner with a message
func ErrorBanner(message string) string {
	return ErrorBannerStyle.Render("❌ " + message)
}

// WarningBanner creates a warning banner
func WarningBanner(message string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorWarning).
		Foreground(ColorWarning).
		Padding(0, 2).
		Render("⚠️  " + message)
}

// InfoBanner creates an info banner with a message
func InfoBanner(message string) string {
	return InfoBannerStyle.Render("ℹ️  " + message)
}

// Box wraps content in a rounded box
func Box(title, content string) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1).
		Width(60)

	return box.Render(
		titleStyle.Render(title) + "\n" + content,
	)
}

// Separator creates a visual separator
func Separator(width int) string {
	return lipgloss.NewStyle().
		Foreground(ColorBorder).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, lipgloss.NewStyle().Width(width).Render("─")))
}
