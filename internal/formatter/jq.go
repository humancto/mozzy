package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ApplyJQ applies a simple JSONPath query to the data
// Supports: .key, .key.subkey, .[0], .[0].key, etc.
func ApplyJQ(data []byte, query string) ([]byte, error) {
	if query == "" || query == "." {
		return data, nil
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	result, err := navigate(parsed, query)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(result, "", "  ")
}

func navigate(data any, path string) (any, error) {
	path = strings.TrimPrefix(path, ".")
	if path == "" {
		return data, nil
	}

	segments := parseJQPath(path)
	current := data

	for _, seg := range segments {
		if seg.isArray {
			arr, ok := current.([]any)
			if !ok {
				return nil, fmt.Errorf("expected array, got %T", current)
			}
			if seg.index < 0 || seg.index >= len(arr) {
				return nil, fmt.Errorf("index %d out of bounds", seg.index)
			}
			current = arr[seg.index]
		} else {
			obj, ok := current.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected object, got %T at key %q", current, seg.key)
			}
			val, exists := obj[seg.key]
			if !exists {
				return nil, fmt.Errorf("key %q not found", seg.key)
			}
			current = val
		}
	}

	return current, nil
}

type jqSegment struct {
	isArray bool
	index   int
	key     string
}

func parseJQPath(path string) []jqSegment {
	if path == "" {
		return nil
	}

	var segments []jqSegment
	current := ""
	inBracket := false

	for i := 0; i < len(path); i++ {
		ch := path[i]

		switch ch {
		case '[':
			if current != "" {
				segments = append(segments, jqSegment{key: current})
				current = ""
			}
			inBracket = true
		case ']':
			if inBracket {
				idx := 0
				fmt.Sscanf(current, "%d", &idx)
				segments = append(segments, jqSegment{isArray: true, index: idx})
				current = ""
				inBracket = false
			}
		case '.':
			if !inBracket {
				if current != "" {
					segments = append(segments, jqSegment{key: current})
					current = ""
				}
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}

	if current != "" {
		segments = append(segments, jqSegment{key: current})
	}

	return segments
}
