package schema

import (
	"encoding/json"
	"testing"
)

func TestValidate_TypeValidation(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		schema    Schema
		wantError bool
	}{
		{
			name:      "valid string",
			data:      `"hello"`,
			schema:    Schema{Type: "string"},
			wantError: false,
		},
		{
			name:      "invalid string type",
			data:      `123`,
			schema:    Schema{Type: "string"},
			wantError: true,
		},
		{
			name:      "valid number",
			data:      `42.5`,
			schema:    Schema{Type: "number"},
			wantError: false,
		},
		{
			name:      "invalid number type",
			data:      `"42"`,
			schema:    Schema{Type: "number"},
			wantError: true,
		},
		{
			name:      "valid integer",
			data:      `42`,
			schema:    Schema{Type: "integer"},
			wantError: false,
		},
		{
			name:      "invalid integer (float)",
			data:      `42.5`,
			schema:    Schema{Type: "integer"},
			wantError: true,
		},
		{
			name:      "valid boolean",
			data:      `true`,
			schema:    Schema{Type: "boolean"},
			wantError: false,
		},
		{
			name:      "valid array",
			data:      `[1,2,3]`,
			schema:    Schema{Type: "array"},
			wantError: false,
		},
		{
			name:      "valid object",
			data:      `{"key":"value"}`,
			schema:    Schema{Type: "object"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), tt.schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}

func TestValidate_ObjectProperties(t *testing.T) {
	schema := Schema{
		Type: "object",
		Properties: map[string]Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
		Required: []string{"name"},
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
		errorPath string
	}{
		{
			name:      "valid object",
			data:      `{"name":"Alice","age":30}`,
			wantError: false,
		},
		{
			name:      "missing required field",
			data:      `{"age":30}`,
			wantError: true,
			errorPath: "name",
		},
		{
			name:      "wrong property type",
			data:      `{"name":"Alice","age":"thirty"}`,
			wantError: true,
			errorPath: "age",
		},
		{
			name:      "valid with missing optional field",
			data:      `{"name":"Alice"}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
			if tt.wantError && hasError && tt.errorPath != "" {
				found := false
				for _, err := range errors {
					if err.Path == tt.errorPath {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error at path %q, got errors: %v", tt.errorPath, errors)
				}
			}
		})
	}
}

func TestValidate_NestedObjects(t *testing.T) {
	schema := Schema{
		Type: "object",
		Properties: map[string]Schema{
			"user": {
				Type: "object",
				Properties: map[string]Schema{
					"name": {Type: "string"},
					"email": {Type: "string"},
				},
				Required: []string{"email"},
			},
		},
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
		errorPath string
	}{
		{
			name:      "valid nested object",
			data:      `{"user":{"name":"Alice","email":"alice@example.com"}}`,
			wantError: false,
		},
		{
			name:      "missing required nested field",
			data:      `{"user":{"name":"Alice"}}`,
			wantError: true,
			errorPath: "user.email",
		},
		{
			name:      "wrong nested field type",
			data:      `{"user":{"name":"Alice","email":123}}`,
			wantError: true,
			errorPath: "user.email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
			if tt.wantError && hasError && tt.errorPath != "" {
				found := false
				for _, err := range errors {
					if err.Path == tt.errorPath {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error at path %q, got errors: %v", tt.errorPath, errors)
				}
			}
		})
	}
}

func TestValidate_Arrays(t *testing.T) {
	schema := Schema{
		Type: "array",
		Items: &Schema{
			Type: "object",
			Properties: map[string]Schema{
				"id":   {Type: "integer"},
				"name": {Type: "string"},
			},
			Required: []string{"id"},
		},
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
	}{
		{
			name:      "valid array of objects",
			data:      `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`,
			wantError: false,
		},
		{
			name:      "array with missing required field",
			data:      `[{"id":1,"name":"Alice"},{"name":"Bob"}]`,
			wantError: true,
		},
		{
			name:      "array with wrong type",
			data:      `[{"id":"one","name":"Alice"}]`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}

func TestValidate_StringConstraints(t *testing.T) {
	minLen := 3
	maxLen := 10
	schema := Schema{
		Type:      "string",
		MinLength: &minLen,
		MaxLength: &maxLen,
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
	}{
		{
			name:      "valid length",
			data:      `"hello"`,
			wantError: false,
		},
		{
			name:      "too short",
			data:      `"hi"`,
			wantError: true,
		},
		{
			name:      "too long",
			data:      `"hello world!"`,
			wantError: true,
		},
		{
			name:      "minimum length",
			data:      `"abc"`,
			wantError: false,
		},
		{
			name:      "maximum length",
			data:      `"abcdefghij"`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}

func TestValidate_NumberConstraints(t *testing.T) {
	min := 0.0
	max := 100.0
	schema := Schema{
		Type:    "number",
		Minimum: &min,
		Maximum: &max,
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
	}{
		{
			name:      "valid number",
			data:      `50`,
			wantError: false,
		},
		{
			name:      "below minimum",
			data:      `-1`,
			wantError: true,
		},
		{
			name:      "above maximum",
			data:      `101`,
			wantError: true,
		},
		{
			name:      "at minimum",
			data:      `0`,
			wantError: false,
		},
		{
			name:      "at maximum",
			data:      `100`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}

func TestValidate_Enum(t *testing.T) {
	schema := Schema{
		Type: "string",
		Enum: []interface{}{"red", "green", "blue"},
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
	}{
		{
			name:      "valid enum value",
			data:      `"red"`,
			wantError: false,
		},
		{
			name:      "another valid enum value",
			data:      `"blue"`,
			wantError: false,
		},
		{
			name:      "invalid enum value",
			data:      `"yellow"`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}

func TestValidate_AdditionalProperties(t *testing.T) {
	schemaStrict := Schema{
		Type: "object",
		Properties: map[string]Schema{
			"name": {Type: "string"},
		},
		AdditionalProperties: false,
	}

	schemaPermissive := Schema{
		Type: "object",
		Properties: map[string]Schema{
			"name": {Type: "string"},
		},
		AdditionalProperties: true,
	}

	tests := []struct {
		name      string
		schema    Schema
		data      string
		wantError bool
	}{
		{
			name:      "strict - no additional properties",
			schema:    schemaStrict,
			data:      `{"name":"Alice"}`,
			wantError: false,
		},
		{
			name:      "strict - with additional properties",
			schema:    schemaStrict,
			data:      `{"name":"Alice","age":30}`,
			wantError: true,
		},
		{
			name:      "permissive - with additional properties",
			schema:    schemaPermissive,
			data:      `{"name":"Alice","age":30}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := Validate([]byte(tt.data), tt.schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}

func TestLoadSchema(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "integer"}
		},
		"required": ["name"]
	}`

	schema, err := LoadSchema([]byte(schemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error = %v", err)
	}

	if schema.Type != "object" {
		t.Errorf("LoadSchema() Type = %v, want object", schema.Type)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("LoadSchema() Properties length = %v, want 2", len(schema.Properties))
	}

	if len(schema.Required) != 1 || schema.Required[0] != "name" {
		t.Errorf("LoadSchema() Required = %v, want [name]", schema.Required)
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ValidationError
		want string
	}{
		{
			name: "with path",
			err:  ValidationError{Path: "user.email", Message: "required field missing"},
			want: "user.email: required field missing",
		},
		{
			name: "without path",
			err:  ValidationError{Path: "", Message: "invalid JSON"},
			want: "invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("ValidationError.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidate_ComplexSchema(t *testing.T) {
	// Real-world-like API response schema
	schemaJSON := `{
		"type": "object",
		"properties": {
			"users": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"id": {"type": "integer"},
						"name": {"type": "string", "minLength": 1},
						"email": {"type": "string"},
						"role": {"type": "string", "enum": ["admin", "user", "guest"]},
						"active": {"type": "boolean"}
					},
					"required": ["id", "name", "email"]
				}
			},
			"total": {"type": "integer", "minimum": 0}
		},
		"required": ["users", "total"]
	}`

	schema, err := LoadSchema([]byte(schemaJSON))
	if err != nil {
		t.Fatalf("LoadSchema() error = %v", err)
	}

	tests := []struct {
		name      string
		data      string
		wantError bool
	}{
		{
			name: "valid complex data",
			data: `{
				"users": [
					{"id": 1, "name": "Alice", "email": "alice@example.com", "role": "admin", "active": true},
					{"id": 2, "name": "Bob", "email": "bob@example.com", "role": "user", "active": false}
				],
				"total": 2
			}`,
			wantError: false,
		},
		{
			name: "missing required field in array item",
			data: `{
				"users": [
					{"id": 1, "name": "Alice", "role": "admin"}
				],
				"total": 1
			}`,
			wantError: true,
		},
		{
			name: "invalid enum value",
			data: `{
				"users": [
					{"id": 1, "name": "Alice", "email": "alice@example.com", "role": "superuser"}
				],
				"total": 1
			}`,
			wantError: true,
		},
		{
			name: "negative total",
			data: `{
				"users": [],
				"total": -1
			}`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compact the JSON to remove formatting
			var compacted interface{}
			json.Unmarshal([]byte(tt.data), &compacted)
			compactedBytes, _ := json.Marshal(compacted)

			errors := Validate(compactedBytes, *schema)
			hasError := len(errors) > 0
			if hasError != tt.wantError {
				t.Errorf("Validate() hasError = %v, wantError %v, errors: %v", hasError, tt.wantError, errors)
			}
		})
	}
}
