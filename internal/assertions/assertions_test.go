package assertions

import (
	"testing"
	"time"
)

func TestEvaluateStatus(t *testing.T) {
	tests := []struct {
		name       string
		expr       string
		statusCode int
		want       bool
		wantErr    bool
	}{
		{"status equals 200", "status == 200", 200, true, false},
		{"status equals 404", "status == 200", 404, false, false},
		{"status not equals", "status != 404", 200, true, false},
		{"status greater than", "status > 199", 200, true, false},
		{"status less than", "status < 300", 200, true, false},
		{"status greater or equal", "status >= 200", 200, true, false},
		{"status less or equal", "status <= 299", 200, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expr, tt.statusCode, nil, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Passed != tt.want {
				t.Errorf("Evaluate() = %v, want %v (message: %s)", result.Passed, tt.want, result.Message)
			}
		})
	}
}

func TestEvaluateResponseTime(t *testing.T) {
	tests := []struct {
		name         string
		expr         string
		responseTime time.Duration
		want         bool
		wantErr      bool
	}{
		{"response time less than", "response_time < 500ms", 300 * time.Millisecond, true, false},
		{"response time greater than", "response_time > 100ms", 300 * time.Millisecond, true, false},
		{"response time exceeds limit", "response_time < 100ms", 300 * time.Millisecond, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expr, 200, nil, tt.responseTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Passed != tt.want {
				t.Errorf("Evaluate() = %v, want %v (message: %s)", result.Passed, tt.want, result.Message)
			}
		})
	}
}

func TestEvaluateJSONPath(t *testing.T) {
	jsonBody := []byte(`{
		"name": "Alice",
		"email": "alice@example.com",
		"age": 30,
		"items": [
			{"id": 1, "title": "First"},
			{"id": 2, "title": "Second"}
		]
	}`)

	tests := []struct {
		name    string
		expr    string
		want    bool
		wantErr bool
	}{
		{"name equals", ".name == Alice", true, false},
		{"name not equals", ".name == Bob", false, false},
		{"email contains", ".email contains @example.com", true, false},
		{"email not contains", ".email contains @test.com", false, false},
		{"field exists", ".name exists", true, false},
		{"field not exists", ".missing exists", false, false},
		{"array index access", ".items[0].id == 1", true, false},
		{"array index title", ".items[1].title == Second", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expr, 200, jsonBody, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Passed != tt.want {
				t.Errorf("Evaluate() = %v, want %v (message: %s)", result.Passed, tt.want, result.Message)
			}
		})
	}
}

func TestEvaluateLength(t *testing.T) {
	jsonBody := []byte(`{
		"items": [1, 2, 3, 4, 5],
		"users": [
			{"name": "Alice"},
			{"name": "Bob"}
		]
	}`)

	tests := []struct {
		name    string
		expr    string
		want    bool
		wantErr bool
	}{
		{"length greater than", "length(.items) > 3", true, false},
		{"length less than", "length(.items) < 10", true, false},
		{"length equals", "length(.users) == 2", true, false},
		{"length not equals", "length(.items) != 3", true, false},
		{"length fails", "length(.items) > 10", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expr, 200, jsonBody, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Passed != tt.want {
				t.Errorf("Evaluate() = %v, want %v (message: %s)", result.Passed, tt.want, result.Message)
			}
		})
	}
}

func TestNavigateJSONPath(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
			"address": map[string]interface{}{
				"city": "NYC",
			},
		},
		"items": []interface{}{
			map[string]interface{}{"id": 1},
			map[string]interface{}{"id": 2},
		},
	}

	tests := []struct {
		name   string
		path   string
		want   interface{}
		exists bool
	}{
		{"simple path", ".user.name", "Alice", true},
		{"nested path", ".user.address.city", "NYC", true},
		{"array index", ".items[0].id", 1, true},  // Go map uses int, not float64
		{"missing path", ".user.missing", nil, false},
		{"invalid array index", ".items[10]", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, exists := navigateJSONPath(tt.path, data)
			if exists != tt.exists {
				t.Errorf("navigateJSONPath() exists = %v, want %v", exists, tt.exists)
				return
			}
			if exists && got != tt.want {
				t.Errorf("navigateJSONPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
