package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Table represents a simple table for CLI output
type Table struct {
	Headers []string
	Rows    [][]string
	Width   int
}

// NewTable creates a new table
func NewTable(headers []string) *Table {
	return &Table{
		Headers: headers,
		Rows:    [][]string{},
		Width:   0,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(cells []string) {
	t.Rows = append(t.Rows, cells)
}

// Render renders the table with box drawing characters
func (t *Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		colWidths[i] = len(header)
	}

	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Add padding
	for i := range colWidths {
		colWidths[i] += 2 // Add 1 space padding on each side
	}

	var sb strings.Builder

	// Top border
	sb.WriteString("╭")
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("─", width))
		if i < len(colWidths)-1 {
			sb.WriteString("┬")
		}
	}
	sb.WriteString("╮\n")

	// Header row
	sb.WriteString("│")
	for i, header := range t.Headers {
		paddedHeader := fmt.Sprintf(" %-*s", colWidths[i]-1, header)
		sb.WriteString(TableHeaderStyle.Render(paddedHeader))
		sb.WriteString("│")
	}
	sb.WriteString("\n")

	// Header separator
	sb.WriteString("├")
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("─", width))
		if i < len(colWidths)-1 {
			sb.WriteString("┼")
		}
	}
	sb.WriteString("┤\n")

	// Data rows
	for _, row := range t.Rows {
		sb.WriteString("│")
		for i, cell := range row {
			if i >= len(colWidths) {
				break
			}
			// Truncate if too long
			displayCell := cell
			if len(cell) > colWidths[i]-2 {
				displayCell = cell[:colWidths[i]-5] + "..."
			}
			paddedCell := fmt.Sprintf(" %-*s", colWidths[i]-1, displayCell)
			sb.WriteString(TableRowStyle.Render(paddedCell))
			sb.WriteString("│")
		}
		sb.WriteString("\n")
	}

	// Bottom border
	sb.WriteString("╰")
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("─", width))
		if i < len(colWidths)-1 {
			sb.WriteString("┴")
		}
	}
	sb.WriteString("╯\n")

	return sb.String()
}

// RenderSimple renders a simpler version without heavy borders
func (t *Table) RenderSimple() string {
	if len(t.Headers) == 0 {
		return ""
	}

	// Calculate column widths
	colWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		colWidths[i] = len(header)
	}

	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder

	// Header
	for i, header := range t.Headers {
		sb.WriteString(TableHeaderStyle.Render(fmt.Sprintf("%-*s", colWidths[i]+2, header)))
		if i < len(t.Headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Header underline
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("─", width+2))
		if i < len(colWidths)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Rows
	for _, row := range t.Rows {
		for i, cell := range row {
			if i >= len(colWidths) {
				break
			}
			displayCell := cell
			if len(cell) > colWidths[i] {
				displayCell = cell[:colWidths[i]-3] + "..."
			}
			sb.WriteString(fmt.Sprintf("%-*s", colWidths[i]+2, displayCell))
			if i < len(row)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// KeyValue renders a simple key-value list
func KeyValue(pairs map[string]string) string {
	var sb strings.Builder

	// Find max key length
	maxLen := 0
	for key := range pairs {
		if len(key) > maxLen {
			maxLen = len(key)
		}
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(ColorInfo).
		Bold(true).
		Width(maxLen + 2).
		Align(lipgloss.Right)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	for key, value := range pairs {
		sb.WriteString(keyStyle.Render(key + ":"))
		sb.WriteString(" ")
		sb.WriteString(valueStyle.Render(value))
		sb.WriteString("\n")
	}

	return sb.String()
}
