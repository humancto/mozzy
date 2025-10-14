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
// Example: "token=.access_token" or "firstId=.[0].id"
func Capture(body []byte, spec string) error {
	parts := strings.SplitN(spec, "=", 2)
	if len(parts) != 2 { return fmt.Errorf("invalid capture %q (want name=.json.path)", spec) }
	name, path := parts[0], parts[1]
	path = strings.TrimPrefix(path, ".")

	var data any
	if err := json.Unmarshal(body, &data); err != nil { return err }

	cur := data
	if path != "" {
		// Parse path segments supporting both object keys and array indices
		// e.g., "[0].id", "data.users[1].name"
		segments := parsePath(path)
		for _, seg := range segments {
			if seg.isArray {
				// Array index access
				switch node := cur.(type) {
				case []any:
					if seg.index >= 0 && seg.index < len(node) {
						cur = node[seg.index]
					} else {
						return fmt.Errorf("array index %d out of bounds (length %d)", seg.index, len(node))
					}
				default:
					return fmt.Errorf("expected array at index %d, got %T", seg.index, cur)
				}
			} else {
				// Object key access
				switch node := cur.(type) {
				case map[string]any:
					cur = node[seg.key]
				default:
					return fmt.Errorf("capture path not found at %q", seg.key)
				}
			}
		}
	}
	switch v := cur.(type) {
	case string:
		store[name] = v
	case float64:
		// Convert numbers to string without JSON encoding
		store[name] = fmt.Sprintf("%.0f", v)
	case bool:
		store[name] = fmt.Sprintf("%t", v)
	case nil:
		store[name] = "null"
	default:
		// store as JSON string for complex types
		b, _ := json.Marshal(v)
		store[name] = string(b)
	}
	return nil
}

type pathSegment struct {
	isArray bool
	index   int
	key     string
}

func parsePath(path string) []pathSegment {
	if path == "" {
		return nil
	}

	var segments []pathSegment
	current := ""
	inBracket := false

	for i := 0; i < len(path); i++ {
		ch := path[i]

		switch ch {
		case '[':
			// Save current segment as object key if not empty
			if current != "" {
				segments = append(segments, pathSegment{key: current})
				current = ""
			}
			inBracket = true
		case ']':
			if inBracket {
				// Parse array index
				idx := 0
				fmt.Sscanf(current, "%d", &idx)
				segments = append(segments, pathSegment{isArray: true, index: idx})
				current = ""
				inBracket = false
			}
		case '.':
			if !inBracket {
				// Segment boundary
				if current != "" {
					segments = append(segments, pathSegment{key: current})
					current = ""
				}
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}

	// Add final segment if exists
	if current != "" {
		segments = append(segments, pathSegment{key: current})
	}

	return segments
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
