package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func PrintJSONOrText(b []byte, jqQuery string) error {
	// Apply jq query if provided
	if jqQuery != "" && jqQuery != "." {
		filtered, err := ApplyJQ(b, jqQuery)
		if err != nil {
			return fmt.Errorf("jq query failed: %w", err)
		}
		b = filtered
	}

	// Pretty-print JSON
	var out bytes.Buffer
	if json.Indent(&out, b, "", "  ") == nil {
		colorizeJSON(out.String())
		return nil
	}
	// Not JSON — print raw
	fmt.Println(string(b))
	return nil
}

func colorizeJSON(jsonStr string) {
	if color.NoColor {
		fmt.Println(jsonStr)
		return
	}

	// Define colors
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta)
	red := color.New(color.FgRed)
	white := color.New(color.FgWhite, color.Bold)

	lines := strings.Split(jsonStr, "\n")

	// Patterns to identify JSON components
	keyPattern := regexp.MustCompile(`^(\s*)"([^"]+)":`)
	stringPattern := regexp.MustCompile(`"([^"]*)"`)
	numberPattern := regexp.MustCompile(`-?\d+\.?\d*([eE][+-]?\d+)?`)
	boolPattern := regexp.MustCompile(`\b(true|false)\b`)
	nullPattern := regexp.MustCompile(`\bnull\b`)

	for _, line := range lines {
		// Extract leading whitespace
		leadingSpace := ""
		trimmed := strings.TrimLeft(line, " ")
		if len(line) > len(trimmed) {
			leadingSpace = line[:len(line)-len(trimmed)]
		}

		// Check if line has a key
		if keyPattern.MatchString(line) {
			matches := keyPattern.FindStringSubmatch(line)
			if len(matches) >= 3 {
				// Print: whitespace + colored key + colon
				fmt.Print(matches[1])
				cyan.Printf("\"%s\"", matches[2])
				fmt.Print(":")

				// Get the rest of the line after the key
				rest := line[len(matches[0]):]
				printColoredValue(rest, green, yellow, magenta, red, white, stringPattern, numberPattern, boolPattern, nullPattern)
				fmt.Println()
				continue
			}
		}

		// No key, just print the line with appropriate colors
		fmt.Print(leadingSpace)
		printColoredValue(strings.TrimLeft(line, " "), green, yellow, magenta, red, white, stringPattern, numberPattern, boolPattern, nullPattern)
		fmt.Println()
	}
}

func printColoredValue(s string, green, yellow, magenta, red, white *color.Color, stringPattern, numberPattern, boolPattern, nullPattern *regexp.Regexp) {
	s = strings.TrimSpace(s)

	// Check for brackets/braces first
	if s == "{" || s == "}" || s == "[" || s == "]" {
		white.Print(s)
		return
	}

	// Remove trailing comma if exists
	hasComma := strings.HasSuffix(s, ",")
	if hasComma {
		s = strings.TrimSuffix(s, ",")
		s = strings.TrimSpace(s)
	}

	// Check what type of value this is
	if stringPattern.MatchString(s) && strings.HasPrefix(s, "\"") {
		green.Print(s)
	} else if boolPattern.MatchString(s) {
		magenta.Print(s)
	} else if nullPattern.MatchString(s) {
		red.Print(s)
	} else if numberPattern.MatchString(s) {
		yellow.Print(s)
	} else if s == "{" || s == "}" {
		white.Print(s)
	} else {
		// Default: print as-is
		fmt.Print(s)
	}

	if hasComma {
		fmt.Print(",")
	}
}
