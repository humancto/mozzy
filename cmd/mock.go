package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"github.com/humancto/mozzy/internal/collection"
	"github.com/humancto/mozzy/internal/mock"
	"github.com/humancto/mozzy/internal/ui"
)

var (
	mockConfigFile     string
	mockFromCollection bool
	mockGenerate       bool
)

var mockCmd = &cobra.Command{
	Use:   "mock [port]",
	Short: "Start a mock HTTP server",
	Long: `Start a mock HTTP server for testing and development.

The mock server can load routes from:
- A YAML configuration file (--config)
- Your saved collection (--from-collection)
- Generate a sample config (--generate)

Examples:
  # Start mock server on port 8080 with config file
  mozzy mock 8080 --config mock.yaml

  # Start mock server using saved collection
  mozzy mock 3000 --from-collection

  # Generate a sample configuration file
  mozzy mock --generate > mock.yaml

  # Start on default port (8080) with inline config
  mozzy mock`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle --generate flag
		if mockGenerate {
			return generateMockConfig()
		}

		// Determine port
		port := 8080
		if len(args) > 0 {
			p, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid port number: %s", args[0])
			}
			port = p
		}

		var config *mock.Config
		var err error

		// Load config based on flags
		if mockConfigFile != "" {
			// Load from YAML file
			config, err = mock.LoadConfig(mockConfigFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			config.Port = port
		} else if mockFromCollection {
			// Load from saved collection
			config, err = loadConfigFromCollection(port)
			if err != nil {
				return err
			}
		} else {
			// Use default config
			config = mock.DefaultConfig(port)
			fmt.Println(ui.WarningBanner("No config specified. Using default empty configuration."))
			fmt.Println(ui.DimStyle.Render("Use --config <file> or --from-collection to add routes."))
			fmt.Println(ui.DimStyle.Render("Use --generate to create a sample config file.\n"))
		}

		// Create and start server
		server := mock.NewServer(config, mockConfigFile)

		// Print server info
		fmt.Println(ui.TitleStyle.Render("🎭 Starting Mock Server"))
		fmt.Printf("\n%s\n", color.CyanString("Routes:"))

		if len(config.Routes) == 0 {
			fmt.Println(ui.DimStyle.Render("  No routes configured"))
		}

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		errChan := make(chan error, 1)
		go func() {
			errChan <- server.Start()
		}()

		// Wait for interrupt or error
		select {
		case <-sigChan:
			fmt.Println(color.YellowString("\n\n🛑 Shutting down mock server..."))
			return server.Stop()
		case err := <-errChan:
			return err
		}
	},
}

func loadConfigFromCollection(port int) (*mock.Config, error) {
	coll, err := collection.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load collection: %w", err)
	}

	requests := coll.List()
	if len(requests) == 0 {
		return nil, fmt.Errorf("no saved requests found in collection")
	}

	config := mock.DefaultConfig(port)

	// Convert saved requests to mock routes
	for _, req := range requests {
		route := mock.Route{
			Path:        "/" + req.Name,
			Method:      req.Method,
			StatusCode:  200,
			Description: req.Description,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		// Use body as response if available
		if req.Body != "" {
			var bodyData interface{}
			// Try to parse as JSON
			if err := json.Unmarshal([]byte(req.Body), &bodyData); err == nil {
				route.Response = bodyData
			} else {
				route.Response = req.Body
			}
		} else {
			// Default response
			route.Response = map[string]interface{}{
				"message": fmt.Sprintf("Mock response for %s", req.Name),
				"method":  req.Method,
				"url":     req.URL,
			}
		}

		config.Routes = append(config.Routes, route)
	}

	return config, nil
}

func generateMockConfig() error {
	sampleConfig := &mock.Config{
		Port: 8080,
		Host: "localhost",
		Defaults: mock.RouteDefaults{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		CORS: mock.CORSConfig{
			Enabled: true,
			Origins: []string{"*"},
			Methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			Headers: []string{"*"},
		},
		Routes: []mock.Route{
			{
				Path:        "/api/users",
				Method:      "GET",
				StatusCode:  200,
				Description: "Get all users",
				Response: []map[string]interface{}{
					{"id": 1, "name": "Alice", "email": "alice@example.com"},
					{"id": 2, "name": "Bob", "email": "bob@example.com"},
				},
			},
			{
				Path:        "/api/users/1",
				Method:      "GET",
				StatusCode:  200,
				Description: "Get user by ID",
				Response: map[string]interface{}{
					"id":    1,
					"name":  "Alice",
					"email": "alice@example.com",
				},
			},
			{
				Path:        "/api/users",
				Method:      "POST",
				StatusCode:  201,
				Description: "Create a new user",
				Response: map[string]interface{}{
					"id":      3,
					"name":    "New User",
					"email":   "newuser@example.com",
					"created": true,
				},
			},
			{
				Path:        "/api/slow",
				Method:      "GET",
				StatusCode:  200,
				Delay:       2000,
				Description: "Slow endpoint (2s delay)",
				Response: map[string]interface{}{
					"message": "This response was delayed by 2 seconds",
				},
			},
			{
				Path:        "/api/error",
				Method:      "GET",
				StatusCode:  500,
				Description: "Simulate server error",
				Response: map[string]interface{}{
					"error":   "Internal server error",
					"message": "Something went wrong",
				},
			},
		},
	}

	// Save to YAML and print
	data, err := yaml.Marshal(sampleConfig)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func init() {
	mockCmd.Flags().StringVarP(&mockConfigFile, "config", "c", "", "Path to mock configuration YAML file")
	mockCmd.Flags().BoolVar(&mockFromCollection, "from-collection", false, "Generate mock routes from saved collection")
	mockCmd.Flags().BoolVar(&mockGenerate, "generate", false, "Generate a sample mock configuration")

	rootCmd.AddCommand(mockCmd)
}
