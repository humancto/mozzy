package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show available environments from .mozzy.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(".mozzy.json")
		if err != nil {
			return fmt.Errorf("no .mozzy.json file found in current directory")
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("invalid .mozzy.json: %w", err)
		}

		envs, ok := config["environments"].(map[string]interface{})
		if !ok || len(envs) == 0 {
			fmt.Println("No environments defined in .mozzy.json")
			fmt.Println()
			fmt.Println("Example .mozzy.json:")
			fmt.Println(`{
  "environments": {
    "dev": {
      "base_url": "http://localhost:3000",
      "headers": {
        "X-Env": "development"
      }
    },
    "prod": {
      "base_url": "https://api.example.com",
      "auth_token": "your-token-here"
    }
  }
}`)
			return nil
		}

		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println(cyan("🌍 Available Environments"))
		fmt.Println()

		for name, env := range envs {
			envMap, ok := env.(map[string]interface{})
			if !ok {
				continue
			}

			fmt.Printf("%s %s\n", green("•"), cyan(name))

			if baseURL, ok := envMap["base_url"].(string); ok {
				fmt.Printf("  %s %s\n", gray("Base URL:"), baseURL)
			}
			if authToken, ok := envMap["auth_token"].(string); ok && authToken != "" {
				fmt.Printf("  %s %s\n", gray("Auth:"), "configured ✓")
			}
			if headers, ok := envMap["headers"].(map[string]interface{}); ok && len(headers) > 0 {
				fmt.Printf("  %s %d custom headers\n", gray("Headers:"), len(headers))
			}
			fmt.Println()
		}

		fmt.Println(gray("Use with: mozzy --env <name> GET /endpoint"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
