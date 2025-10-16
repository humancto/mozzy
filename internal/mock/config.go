package mock

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the mock server configuration
type Config struct {
	Port      int               `yaml:"port" json:"port"`
	Host      string            `yaml:"host" json:"host"`
	Routes    []Route           `yaml:"routes" json:"routes"`
	Defaults  RouteDefaults     `yaml:"defaults,omitempty" json:"defaults,omitempty"`
	CORS      CORSConfig        `yaml:"cors,omitempty" json:"cors,omitempty"`
}

// Route represents a single mock endpoint
type Route struct {
	Path        string            `yaml:"path" json:"path"`
	Method      string            `yaml:"method" json:"method"`
	StatusCode  int               `yaml:"status" json:"status"`
	Response    interface{}       `yaml:"response" json:"response"`
	Headers     map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Delay       int               `yaml:"delay,omitempty" json:"delay,omitempty"` // milliseconds
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
}

// RouteDefaults contains default values for routes
type RouteDefaults struct {
	StatusCode int               `yaml:"status" json:"status"`
	Headers    map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
}

// CORSConfig contains CORS settings
type CORSConfig struct {
	Enabled bool     `yaml:"enabled" json:"enabled"`
	Origins []string `yaml:"origins,omitempty" json:"origins,omitempty"`
	Methods []string `yaml:"methods,omitempty" json:"methods,omitempty"`
	Headers []string `yaml:"headers,omitempty" json:"headers,omitempty"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Apply defaults
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	if config.Defaults.StatusCode == 0 {
		config.Defaults.StatusCode = 200
	}

	// Apply route defaults
	for i := range config.Routes {
		if config.Routes[i].StatusCode == 0 {
			config.Routes[i].StatusCode = config.Defaults.StatusCode
		}
		if config.Routes[i].Method == "" {
			config.Routes[i].Method = "GET"
		}
		// Merge default headers with route headers
		if config.Defaults.Headers != nil {
			if config.Routes[i].Headers == nil {
				config.Routes[i].Headers = make(map[string]string)
			}
			for k, v := range config.Defaults.Headers {
				if _, exists := config.Routes[i].Headers[k]; !exists {
					config.Routes[i].Headers[k] = v
				}
			}
		}
	}

	return &config, nil
}

// DefaultConfig returns a default configuration
func DefaultConfig(port int) *Config {
	return &Config{
		Port: port,
		Host: "localhost",
		Defaults: RouteDefaults{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		Routes: []Route{},
		CORS: CORSConfig{
			Enabled: true,
			Origins: []string{"*"},
			Methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			Headers: []string{"*"},
		},
	}
}

// SaveConfig saves configuration to a YAML file
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ToJSON converts response to JSON bytes
func (r *Route) ToJSON() ([]byte, error) {
	// If response is already a string, return as-is
	if str, ok := r.Response.(string); ok {
		return []byte(str), nil
	}
	// Otherwise marshal to JSON
	return json.Marshal(r.Response)
}
