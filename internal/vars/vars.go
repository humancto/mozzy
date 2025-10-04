package vars

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var store = map[string]string{}

// Interpolate replaces {{name}} occurrences using the in-memory store
func Interpolate(s string) string {
	re := regexp.MustCompile(`\{\{([a-zA-Z0-9_.-]+)\}\}`)
	return re.ReplaceAllStringFunc(s, func(m string) string {
		key := re.FindStringSubmatch(m)[1]
		if v, ok := store[key]; ok {
			return v
		}
		return m
	})
}

// Capture parses JSON body and stores `name=path` where path is dot-notation
// Example: "token=.access_token"
func Capture(body []byte, spec string) error {
	parts := strings.SplitN(spec, "=", 2)
	if len(parts) != 2 { return fmt.Errorf("invalid capture %q (want name=.json.path)", spec) }
	name, path := parts[0], parts[1]
	path = strings.TrimPrefix(path, ".")

	var data any
	if err := json.Unmarshal(body, &data); err != nil { return err }

	cur := data
	if path != "" {
		for _, seg := range strings.Split(path, ".") {
			switch node := cur.(type) {
			case map[string]any:
				cur = node[seg]
			default:
				return fmt.Errorf("capture path not found at %q", seg)
			}
		}
	}
	switch v := cur.(type) {
	case string:
		store[name] = v
	default:
		// store as JSON string
		b, _ := json.Marshal(v)
		store[name] = string(b)
	}
	return nil
}

// ResolveBase picks base from env file or CLI flag
func ResolveBase(cliBase, envName string) string {
	if envName == "" && cliBase != "" { return cliBase }
	// load .mozzy.json if present
	wd, _ := os.Getwd()
	cfgPath := filepath.Join(wd, ".mozzy.json")
	b, err := os.ReadFile(cfgPath)
	if err != nil { return cliBase }

	var cfg struct {
		DefaultEnv    string            `json:"default_env"`
		Environments  map[string]string `json:"environments"`
	}
	_ = json.Unmarshal(b, &cfg)

	target := envName
	if target == "" && cfg.DefaultEnv != "" { target = cfg.DefaultEnv }
	if target != "" {
		if base, ok := cfg.Environments[target]; ok {
			return base
		}
	}
	return cliBase
}
