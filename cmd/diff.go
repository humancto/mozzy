package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <file1.json> <file2.json>",
	Short: "Compare two JSON responses and show differences",
	Long: `Compare two JSON files and highlight differences.

Useful for:
- Comparing API responses between environments
- Detecting changes in API contracts
- Validating migrations

Example:
  mozzy diff prod-response.json staging-response.json
  mozzy GET /api/users/1 --env prod > prod.json
  mozzy GET /api/users/1 --env staging > staging.json
  mozzy diff prod.json staging.json`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	file1, file2 := args[0], args[1]

	// Read both files
	b1, err := os.ReadFile(file1)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1, err)
	}

	b2, err := os.ReadFile(file2)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2, err)
	}

	// Parse as JSON
	var data1, data2 interface{}
	if err := json.Unmarshal(b1, &data1); err != nil {
		return fmt.Errorf("failed to parse %s as JSON: %w", file1, err)
	}
	if err := json.Unmarshal(b2, &data2); err != nil {
		return fmt.Errorf("failed to parse %s as JSON: %w", file2, err)
	}

	// Compare
	printDiffHeader(file1, file2)

	diffs := compareJSON("", data1, data2)

	if len(diffs) == 0 {
		printNoDifferences()
		return nil
	}

	printDiffSummary(len(diffs))
	for _, diff := range diffs {
		printDiffLine(diff)
	}
	fmt.Println()

	return nil
}

type jsonDiff struct {
	Path      string
	LeftVal   interface{}
	RightVal  interface{}
	DiffType  string // "added", "removed", "changed", "type-mismatch"
}

func compareJSON(path string, left, right interface{}) []jsonDiff {
	var diffs []jsonDiff

	// Type mismatch
	if fmt.Sprintf("%T", left) != fmt.Sprintf("%T", right) {
		diffs = append(diffs, jsonDiff{
			Path:      path,
			LeftVal:   left,
			RightVal:  right,
			DiffType:  "type-mismatch",
		})
		return diffs
	}

	switch l := left.(type) {
	case map[string]interface{}:
		r := right.(map[string]interface{})

		// Check all keys in left
		for key, lval := range l {
			newPath := path + "." + key
			if path == "" {
				newPath = key
			}

			rval, exists := r[key]
			if !exists {
				diffs = append(diffs, jsonDiff{
					Path:     newPath,
					LeftVal:  lval,
					RightVal: nil,
					DiffType: "removed",
				})
				continue
			}

			diffs = append(diffs, compareJSON(newPath, lval, rval)...)
		}

		// Check for added keys in right
		for key, rval := range r {
			newPath := path + "." + key
			if path == "" {
				newPath = key
			}

			if _, exists := l[key]; !exists {
				diffs = append(diffs, jsonDiff{
					Path:     newPath,
					LeftVal:  nil,
					RightVal: rval,
					DiffType: "added",
				})
			}
		}

	case []interface{}:
		r := right.([]interface{})

		if len(l) != len(r) {
			diffs = append(diffs, jsonDiff{
				Path:     path + ".length",
				LeftVal:  len(l),
				RightVal: len(r),
				DiffType: "changed",
			})
		}

		// Compare elements
		minLen := len(l)
		if len(r) < minLen {
			minLen = len(r)
		}

		for i := 0; i < minLen; i++ {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			diffs = append(diffs, compareJSON(newPath, l[i], r[i])...)
		}

	default:
		// Primitive values - compare directly
		if fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right) {
			diffs = append(diffs, jsonDiff{
				Path:     path,
				LeftVal:  left,
				RightVal: right,
				DiffType: "changed",
			})
		}
	}

	return diffs
}

func printDiffHeader(file1, file2 string) {
	cyan := color.New(color.FgCyan, color.Bold)
	gray := color.New(color.FgHiBlack)

	fmt.Println()
	cyan.Println("╔════════════════════════════════════════════════════════════════════")
	cyan.Print("║ ")
	fmt.Print("📊 JSON Diff Comparison\n")
	cyan.Println("╚════════════════════════════════════════════════════════════════════")
	fmt.Println()

	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)

	red.Print("  - ")
	gray.Printf("Left:  %s\n", file1)
	green.Print("  + ")
	gray.Printf("Right: %s\n", file2)
	fmt.Println()

	gray.Println(strings.Repeat("─", 72))
	fmt.Println()
}

func printNoDifferences() {
	green := color.New(color.FgGreen, color.Bold)
	fmt.Println()
	green.Println("  ✅ No differences found")
	fmt.Println()
	gray := color.New(color.FgHiBlack)
	gray.Println("  Both JSON files are identical.")
	fmt.Println()
}

func printDiffSummary(count int) {
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("  Found %d difference(s):\n\n", count)
}

func printDiffLine(diff jsonDiff) {
	pathColor := color.New(color.FgCyan, color.Bold)
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta)
	gray := color.New(color.FgHiBlack)

	switch diff.DiffType {
	case "added":
		fmt.Print("  ")
		green.Print("+ ")
		pathColor.Print(diff.Path)
		gray.Print(" │ ")
		fmt.Printf("%v", formatValue(diff.RightVal))
		fmt.Println()

	case "removed":
		fmt.Print("  ")
		red.Print("- ")
		pathColor.Print(diff.Path)
		gray.Print(" │ ")
		fmt.Printf("%v", formatValue(diff.LeftVal))
		fmt.Println()

	case "changed":
		fmt.Print("  ")
		yellow.Print("~ ")
		pathColor.Println(diff.Path)
		fmt.Print("    ")
		red.Print("- ")
		fmt.Println(formatValue(diff.LeftVal))
		fmt.Print("    ")
		green.Print("+ ")
		fmt.Println(formatValue(diff.RightVal))

	case "type-mismatch":
		fmt.Print("  ")
		magenta.Print("! ")
		pathColor.Print(diff.Path)
		gray.Println(" (type mismatch)")
		fmt.Print("    ")
		red.Printf("- %T: %v\n", diff.LeftVal, formatValue(diff.LeftVal))
		fmt.Print("    ")
		green.Printf("+ %T: %v\n", diff.RightVal, formatValue(diff.RightVal))
	}
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val)
	case map[string]interface{}, []interface{}:
		b, _ := json.Marshal(val)
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}
