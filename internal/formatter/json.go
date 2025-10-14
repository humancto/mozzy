package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
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
	// Check if colors should be disabled
	// Allow CLICOLOR_FORCE to override TTY check
	forceColor := os.Getenv("CLICOLOR_FORCE") != "" && os.Getenv("CLICOLOR_FORCE") != "0"
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	// Disable colors if:
	// - NO_COLOR is set
	// - Not a TTY (unless forced)
	// - color.NoColor flag is set
	shouldDisableColors := color.NoColor || (!isTTY && !forceColor)

	if shouldDisableColors {
		fmt.Println(jsonStr)
		return
	}

	lines := strings.Split(jsonStr, "\n")

	// Direct ANSI color codes for maximum compatibility
	const (
		cyan    = "\033[36;1m"
		green   = "\033[32m"
		yellow  = "\033[33m"
		magenta = "\033[35m"
		red     = "\033[31m"
		white   = "\033[37;1m"
		reset   = "\033[0m"
	)

	// Regex patterns
	keyPattern := regexp.MustCompile(`("[\w\-\.]+"):`)
	stringPattern := regexp.MustCompile(`: ("[^"]*")`)
	numberPattern := regexp.MustCompile(`: (-?\d+\.?\d*([eE][+-]?\d+)?)`)
	boolPattern := regexp.MustCompile(`: (true|false)`)
	nullPattern := regexp.MustCompile(`: (null)`)

	for _, line := range lines {
		// Colorize keys
		line = keyPattern.ReplaceAllString(line, cyan+"$1"+reset+":")

		// Colorize string values
		line = stringPattern.ReplaceAllString(line, ": "+green+"$1"+reset)

		// Colorize numbers
		line = numberPattern.ReplaceAllString(line, ": "+yellow+"$1"+reset)

		// Colorize booleans
		line = boolPattern.ReplaceAllString(line, ": "+magenta+"$1"+reset)

		// Colorize null
		line = nullPattern.ReplaceAllString(line, ": "+red+"$1"+reset)

		// Colorize brackets and braces
		line = strings.ReplaceAll(line, "{", white+"{"+reset)
		line = strings.ReplaceAll(line, "}", white+"}"+reset)
		line = strings.ReplaceAll(line, "[", white+"["+reset)
		line = strings.ReplaceAll(line, "]", white+"]"+reset)

		fmt.Println(line)
	}
}
