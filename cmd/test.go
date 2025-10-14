package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/humancto/mozzy/internal/chain"
)

var (
	testJUnitOutput string
)

var testCmd = &cobra.Command{
	Use:   "test <workflow.yaml>",
	Short: "Run a workflow as a test suite with pass/fail summary",
	Long: `Run a YAML workflow as a test suite.

All assertions must pass for the test to succeed.
Exit code 0 if all tests pass, 1 if any fail.

Perfect for CI/CD pipelines.

Example:
  mozzy test api-tests.yaml
  mozzy test api-tests.yaml --junit-output results.xml`,
	Args: cobra.ExactArgs(1),
	RunE: runTest,
}

func init() {
	testCmd.Flags().StringVar(&testJUnitOutput, "junit-output", "", "Write JUnit XML report to file")
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
	workflowFile := args[0]

	// Read workflow file
	b, err := os.ReadFile(workflowFile)
	if err != nil {
		return fmt.Errorf("failed to read workflow: %w", err)
	}

	var flow chain.Flow
	if err := yaml.Unmarshal(b, &flow); err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Populate from flags
	flow.BaseURL = baseURL
	flow.EnvName = envName
	flow.GlobalAuth = authToken

	fmt.Printf("🧪 Running test suite: %s\n", flow.Name)
	if flow.Description != "" {
		fmt.Printf("   %s\n", flow.Description)
	}
	fmt.Printf("   Steps: %d\n\n", len(flow.Steps))

	startTime := time.Now()

	// Run the workflow
	ctx := context.Background()
	err = chain.Run(ctx, flow)

	duration := time.Since(startTime)

	// Print summary
	separator := "============================================================"
	fmt.Println("\n" + separator)
	if err != nil {
		fmt.Printf("❌ TEST SUITE FAILED\n")
		fmt.Printf("   Duration: %s\n", duration.Round(time.Millisecond))
		fmt.Printf("   Error: %v\n", err)
		fmt.Println(separator)

		// TODO: Write JUnit XML if requested
		if testJUnitOutput != "" {
			writeJUnitXML(testJUnitOutput, flow.Name, len(flow.Steps), duration, err)
		}

		os.Exit(1)
	}

	fmt.Printf("✅ TEST SUITE PASSED\n")
	fmt.Printf("   Duration: %s\n", duration.Round(time.Millisecond))
	fmt.Printf("   Steps: %d\n", len(flow.Steps))
	fmt.Println(separator)

	// Write JUnit XML if requested
	if testJUnitOutput != "" {
		writeJUnitXML(testJUnitOutput, flow.Name, len(flow.Steps), duration, nil)
	}

	return nil
}

func writeJUnitXML(filename, suiteName string, tests int, duration time.Duration, err error) {
	failures := 0
	if err != nil {
		failures = 1
	}

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="%s" tests="%d" failures="%d" time="%.3f">
    <testcase name="%s" time="%.3f">
`, suiteName, tests, failures, duration.Seconds(), suiteName, duration.Seconds())

	if err != nil {
		xml += fmt.Sprintf(`      <failure message="%s">%s</failure>
`, err.Error(), err.Error())
	}

	xml += `    </testcase>
  </testsuite>
</testsuites>
`

	if writeErr := os.WriteFile(filename, []byte(xml), 0644); writeErr != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Failed to write JUnit XML: %v\n", writeErr)
	} else {
		fmt.Printf("📝 JUnit XML written to: %s\n", filename)
	}
}
