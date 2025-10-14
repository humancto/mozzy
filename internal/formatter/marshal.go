package formatter

import (
	"encoding/json"

	"github.com/humancto/mozzy/internal/vars"
)

func MarshalJSON(v any) ([]byte, error) {
	// First marshal to JSON
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Then interpolate {{variables}} in the JSON string
	interpolated := vars.Interpolate(string(b))

	// If interpolation changed the string, we need to re-parse to ensure valid JSON
	// This handles cases where variables expand to non-string values
	if interpolated != string(b) {
		// Try to unmarshal and re-marshal to validate/normalize the JSON
		var temp any
		if err := json.Unmarshal([]byte(interpolated), &temp); err == nil {
			return json.Marshal(temp)
		}
		// If that fails, return the interpolated string as-is
		// (might have string values that look like {{var}})
		return []byte(interpolated), nil
	}

	return b, nil
}
