package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Schema represents a JSON Schema (simplified)
type Schema struct {
	Type                 string            `json:"type,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"`
	Required             []string          `json:"required,omitempty"`
	Items                *Schema           `json:"items,omitempty"`
	Enum                 []interface{}     `json:"enum,omitempty"`
	MinLength            *int              `json:"minLength,omitempty"`
	MaxLength            *int              `json:"maxLength,omitempty"`
	Minimum              *float64          `json:"minimum,omitempty"`
	Maximum              *float64          `json:"maximum,omitempty"`
	Pattern              string            `json:"pattern,omitempty"`
	AdditionalProperties bool              `json:"additionalProperties,omitempty"`
}

// ValidationError represents a schema validation error
type ValidationError struct {
	Path    string
	Message string
}

func (e ValidationError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// Validate validates data against a schema
func Validate(data []byte, schema Schema) []ValidationError {
	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return []ValidationError{{Path: "", Message: fmt.Sprintf("invalid JSON: %v", err)}}
	}

	return validate(parsed, schema, "")
}

func validate(data interface{}, schema Schema, path string) []ValidationError {
	var errors []ValidationError

	// Type validation
	if schema.Type != "" {
		if err := validateType(data, schema.Type, path); err != nil {
			errors = append(errors, *err)
			return errors // Stop if type is wrong
		}
	}

	// Type-specific validations
	switch schema.Type {
	case "object":
		if obj, ok := data.(map[string]interface{}); ok {
			errors = append(errors, validateObject(obj, schema, path)...)
		}
	case "array":
		if arr, ok := data.([]interface{}); ok {
			errors = append(errors, validateArray(arr, schema, path)...)
		}
	case "string":
		if str, ok := data.(string); ok {
			errors = append(errors, validateString(str, schema, path)...)
		}
	case "number", "integer":
		if num, ok := data.(float64); ok {
			errors = append(errors, validateNumber(num, schema, path)...)
		}
	}

	// Enum validation
	if len(schema.Enum) > 0 {
		if err := validateEnum(data, schema.Enum, path); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

func validateType(data interface{}, expectedType string, path string) *ValidationError {
	actualType := getJSONType(data)
	if actualType != expectedType {
		// Special case: integer is a subset of number
		if expectedType == "integer" && actualType == "number" {
			if num, ok := data.(float64); ok && num == float64(int64(num)) {
				return nil // It's an integer
			}
		}
		return &ValidationError{
			Path:    path,
			Message: fmt.Sprintf("expected type %s, got %s", expectedType, actualType),
		}
	}
	return nil
}

func getJSONType(data interface{}) string {
	if data == nil {
		return "null"
	}
	switch data.(type) {
	case bool:
		return "boolean"
	case float64:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

func validateObject(obj map[string]interface{}, schema Schema, path string) []ValidationError {
	var errors []ValidationError

	// Required fields
	for _, req := range schema.Required {
		if _, exists := obj[req]; !exists {
			errors = append(errors, ValidationError{
				Path:    joinPath(path, req),
				Message: "required field missing",
			})
		}
	}

	// Validate properties
	for key, value := range obj {
		propSchema, hasSchema := schema.Properties[key]
		if hasSchema {
			errors = append(errors, validate(value, propSchema, joinPath(path, key))...)
		} else if len(schema.Properties) > 0 && !schema.AdditionalProperties {
			// Only check additionalProperties if schema has defined properties
			errors = append(errors, ValidationError{
				Path:    joinPath(path, key),
				Message: "additional property not allowed",
			})
		}
	}

	return errors
}

func validateArray(arr []interface{}, schema Schema, path string) []ValidationError {
	var errors []ValidationError

	if schema.Items != nil {
		for i, item := range arr {
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			errors = append(errors, validate(item, *schema.Items, itemPath)...)
		}
	}

	return errors
}

func validateString(str string, schema Schema, path string) []ValidationError {
	var errors []ValidationError

	if schema.MinLength != nil && len(str) < *schema.MinLength {
		errors = append(errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("string length %d is less than minimum %d", len(str), *schema.MinLength),
		})
	}

	if schema.MaxLength != nil && len(str) > *schema.MaxLength {
		errors = append(errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("string length %d exceeds maximum %d", len(str), *schema.MaxLength),
		})
	}

	// Pattern matching would require regexp package
	// Skipping for now to keep dependencies minimal

	return errors
}

func validateNumber(num float64, schema Schema, path string) []ValidationError {
	var errors []ValidationError

	if schema.Minimum != nil && num < *schema.Minimum {
		errors = append(errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("number %v is less than minimum %v", num, *schema.Minimum),
		})
	}

	if schema.Maximum != nil && num > *schema.Maximum {
		errors = append(errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("number %v exceeds maximum %v", num, *schema.Maximum),
		})
	}

	return errors
}

func validateEnum(data interface{}, enum []interface{}, path string) *ValidationError {
	for _, allowed := range enum {
		if reflect.DeepEqual(data, allowed) {
			return nil
		}
	}
	return &ValidationError{
		Path:    path,
		Message: fmt.Sprintf("value must be one of %v", enum),
	}
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	if strings.HasPrefix(child, "[") {
		return parent + child
	}
	return parent + "." + child
}

// LoadSchema loads a schema from JSON bytes
func LoadSchema(data []byte) (*Schema, error) {
	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
