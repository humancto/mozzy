package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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
	fmt.Printf("📊 Comparing:\n")
	fmt.Printf("  Left:  %s\n", file1)
	fmt.Printf("  Right: %s\n\n", file2)

	diffs := compareJSON("", data1, data2)

	if len(diffs) == 0 {
		color.Green("✅ No differences found - files are identical\n")
		return nil
	}

	fmt.Printf("Found %d difference(s):\n\n", len(diffs))
	for _, diff := range diffs {
		printDiff(diff)
	}

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

func printDiff(diff jsonDiff) {
	switch diff.DiffType {
	case "added":
		color.Green("  + %s: %v\n", diff.Path, diff.RightVal)
	case "removed":
		color.Red("  - %s: %v\n", diff.Path, diff.LeftVal)
	case "changed":
		color.Yellow("  ~ %s:\n", diff.Path)
		color.Red("    - %v\n", diff.LeftVal)
		color.Green("    + %v\n", diff.RightVal)
	case "type-mismatch":
		color.Magenta("  ! %s: type mismatch\n", diff.Path)
		color.Red("    - %T: %v\n", diff.LeftVal, diff.LeftVal)
		color.Green("    + %T: %v\n", diff.RightVal, diff.RightVal)
	}
}
