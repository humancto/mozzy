package cmd

import (
	"reflect"
	"testing"
)

func TestSplitHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple header",
			input:    "X-Test: value",
			expected: []string{"X-Test", "value"},
		},
		{
			name:     "Header without space",
			input:    "X-Test:value",
			expected: []string{"X-Test", "value"},
		},
		{
			name:     "Header with multiple colons",
			input:    "Authorization: Bearer: token",
			expected: []string{"Authorization", "Bearer: token"},
		},
		{
			name:     "No colon",
			input:    "InvalidHeader",
			expected: []string{"InvalidHeader"},
		},
		{
			name:     "Empty value",
			input:    "X-Empty:",
			expected: []string{"X-Empty", ""},
		},
		{
			name:     "Value with spaces",
			input:    "Content-Type: application/json; charset=utf-8",
			expected: []string{"Content-Type", "application/json; charset=utf-8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitHeader(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitHeader(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitMethods(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Single method",
			input:    "GET",
			expected: []string{"GET"},
		},
		{
			name:     "Multiple methods",
			input:    "GET,POST,PUT",
			expected: []string{"GET", "POST", "PUT"},
		},
		{
			name:     "Methods with spaces (no trimming)",
			input:    "GET, POST, PUT",
			expected: []string{"GET", " POST", " PUT"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "Single comma",
			input:    ",",
			expected: nil,
		},
		{
			name:     "Trailing comma",
			input:    "GET,POST,",
			expected: []string{"GET", "POST"},
		},
		{
			name:     "Leading comma",
			input:    ",GET,POST",
			expected: []string{"GET", "POST"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitMethods(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitMethods(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
