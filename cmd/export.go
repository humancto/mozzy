package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/humancto/mozzy/internal/chain"
)

var (
	exportFormat string
)

var exportCmd = &cobra.Command{
	Use:   "export <collection-name>",
	Short: "Export a saved request to various formats",
	Long: `Export saved requests to curl, Postman, or other formats.

Formats:
  - curl: Generate curl command
  - postman: Generate Postman collection JSON

Examples:
  mozzy export my-request --format curl
  mozzy export my-request --format postman > collection.json`,
	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVar(&exportFormat, "format", "curl", "Export format (curl, postman)")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Try to load from collection first
	request, err := loadRequestFromCollection(name)
	if err != nil {
		// If not in collection, try as workflow
		return exportWorkflow(name, exportFormat)
	}

	switch strings.ToLower(exportFormat) {
	case "curl":
		return exportToCurl(request)
	case "postman":
		return exportToPostman(request)
	default:
		return fmt.Errorf("unsupported format: %s (supported: curl, postman)", exportFormat)
	}
}

type SavedRequest struct {
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        string            `json:"body,omitempty"`
	Description string            `json:"description,omitempty"`
}

func loadRequestFromCollection(name string) (*SavedRequest, error) {
	homeDir, _ := os.UserHomeDir()
	collectionFile := homeDir + "/.mozzy/collections.json"

	b, err := os.ReadFile(collectionFile)
	if err != nil {
		return nil, err
	}

	var data struct {
		Requests map[string]SavedRequest `json:"requests"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	req, exists := data.Requests[name]
	if !exists {
		return nil, fmt.Errorf("request '%s' not found in collection", name)
	}

	return &req, nil
}

func exportToCurl(req *SavedRequest) error {
	curl := fmt.Sprintf("curl -X %s '%s'", req.Method, req.URL)

	// Add headers
	for key, value := range req.Headers {
		curl += fmt.Sprintf(" \\\n  -H '%s: %s'", key, value)
	}

	// Add body
	if req.Body != "" {
		curl += fmt.Sprintf(" \\\n  -d '%s'", strings.ReplaceAll(req.Body, "'", "'\\''"))
	}

	fmt.Println(curl)
	return nil
}

func exportToPostman(req *SavedRequest) error {
	collection := map[string]interface{}{
		"info": map[string]interface{}{
			"name":        req.Name,
			"description": req.Description,
			"schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"item": []map[string]interface{}{
			{
				"name": req.Name,
				"request": map[string]interface{}{
					"method": req.Method,
					"url":    req.URL,
					"header": convertHeadersToPostman(req.Headers),
					"body":   convertBodyToPostman(req.Body),
				},
			},
		},
	}

	b, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func convertHeadersToPostman(headers map[string]string) []map[string]string {
	var result []map[string]string
	for key, value := range headers {
		result = append(result, map[string]string{
			"key":   key,
			"value": value,
		})
	}
	return result
}

func convertBodyToPostman(body string) map[string]interface{} {
	if body == "" {
		return nil
	}
	return map[string]interface{}{
		"mode": "raw",
		"raw":  body,
	}
}

func exportWorkflow(name string, format string) error {
	// Try to load as workflow file
	b, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("not found in collection or as workflow file: %w", err)
	}

	var flow chain.Flow
	if err := yaml.Unmarshal(b, &flow); err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	fmt.Printf("# Workflow: %s\n", flow.Name)
	fmt.Printf("# Description: %s\n\n", flow.Description)

	for i, step := range flow.Steps {
		fmt.Printf("# Step %d: %s\n", i+1, step.Name)

		req := &SavedRequest{
			Name:   step.Name,
			Method: step.Method,
			URL:    step.URL,
		}

		if format == "curl" {
			exportToCurl(req)
		}
		fmt.Println()
	}

	return nil
}
