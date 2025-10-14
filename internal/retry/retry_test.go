package retry

import (
	"errors"
	"testing"
)

func TestShouldRetry_DefaultBehavior(t *testing.T) {
	policy := &Policy{
		MaxRetries: 3,
		Conditions: nil, // Default behavior
	}

	tests := []struct {
		name          string
		statusCode    int
		err           error
		wantRetry     bool
	}{
		{"network error", 0, errors.New("connection failed"), true},
		{"5xx server error", 500, nil, true},
		{"503 service unavailable", 503, nil, true},
		{"4xx client error", 404, nil, false},
		{"400 bad request", 400, nil, false},
		{"200 success", 200, nil, false},
		{"301 redirect", 301, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.statusCode, tt.err)
			if got != tt.wantRetry {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}

func TestShouldRetry_StatusExact(t *testing.T) {
	policy := &Policy{
		MaxRetries: 3,
		Conditions: []Condition{
			{Type: "status", Value: "503"},
		},
	}

	tests := []struct {
		name          string
		statusCode    int
		wantRetry     bool
	}{
		{"exact match 503", 503, true},
		{"no match 500", 500, false},
		{"no match 504", 504, false},
		{"no match 404", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.statusCode, nil)
			if got != tt.wantRetry {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}

func TestShouldRetry_StatusRange(t *testing.T) {
	tests := []struct {
		name       string
		conditions []Condition
		statusCode int
		wantRetry  bool
	}{
		{"5xx matches 500", []Condition{{Type: "status_range", Value: "5xx"}}, 500, true},
		{"5xx matches 503", []Condition{{Type: "status_range", Value: "5xx"}}, 503, true},
		{"5xx matches 599", []Condition{{Type: "status_range", Value: "5xx"}}, 599, true},
		{"5xx no match 404", []Condition{{Type: "status_range", Value: "5xx"}}, 404, false},
		{"4xx matches 404", []Condition{{Type: "status_range", Value: "4xx"}}, 404, true},
		{"4xx matches 429", []Condition{{Type: "status_range", Value: "4xx"}}, 429, true},
		{"4xx no match 500", []Condition{{Type: "status_range", Value: "4xx"}}, 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := &Policy{MaxRetries: 3, Conditions: tt.conditions}
			got := policy.ShouldRetry(tt.statusCode, nil)
			if got != tt.wantRetry {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}

func TestShouldRetry_Comparisons(t *testing.T) {
	tests := []struct {
		name       string
		rangeExpr  string
		statusCode int
		wantRetry  bool
	}{
		{">=500 matches 500", ">=500", 500, true},
		{">=500 matches 503", ">=500", 503, true},
		{">=500 no match 404", ">=500", 404, false},
		{"<400 matches 200", "<400", 200, true},
		{"<400 matches 301", "<400", 301, true},
		{"<400 no match 404", "<400", 404, false},
		{">500 matches 503", ">500", 503, true},
		{">500 no match 500", ">500", 500, false},
		{"<=400 matches 400", "<=400", 400, true},
		{"<=400 matches 200", "<=400", 200, true},
		{"<=400 no match 404", "<=400", 404, false},
		{"==503 matches 503", "==503", 503, true},
		{"==503 no match 500", "==503", 500, false},
		{"!=200 matches 404", "!=200", 404, true},
		{"!=200 no match 200", "!=200", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := &Policy{
				MaxRetries: 3,
				Conditions: []Condition{{Type: "status_range", Value: tt.rangeExpr}},
			}
			got := policy.ShouldRetry(tt.statusCode, nil)
			if got != tt.wantRetry {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}

func TestShouldRetry_MultipleConditions(t *testing.T) {
	policy := &Policy{
		MaxRetries: 3,
		Conditions: []Condition{
			{Type: "status", Value: "429"},      // Rate limit
			{Type: "status_range", Value: "5xx"}, // Server errors
		},
	}

	tests := []struct {
		name       string
		statusCode int
		wantRetry  bool
	}{
		{"matches first condition 429", 429, true},
		{"matches second condition 503", 503, true},
		{"matches second condition 500", 500, true},
		{"no match 404", 404, false},
		{"no match 200", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.statusCode, nil)
			if got != tt.wantRetry {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}

func TestShouldRetry_NetworkError(t *testing.T) {
	policy := &Policy{
		MaxRetries: 3,
		Conditions: []Condition{
			{Type: "network_error"},
		},
	}

	tests := []struct {
		name      string
		err       error
		wantRetry bool
	}{
		{"network error", errors.New("connection failed"), true},
		{"no error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(0, tt.err)
			if got != tt.wantRetry {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}

func TestShouldRetry_Always(t *testing.T) {
	policy := &Policy{
		MaxRetries: 3,
		Conditions: []Condition{
			{Type: "always"},
		},
	}

	tests := []struct {
		name       string
		statusCode int
		err        error
	}{
		{"success", 200, nil},
		{"client error", 404, nil},
		{"server error", 500, nil},
		{"network error", 0, errors.New("failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.statusCode, tt.err)
			if !got {
				t.Errorf("ShouldRetry() with 'always' = false, want true")
			}
		})
	}
}

func TestShouldRetry_Never(t *testing.T) {
	policy := &Policy{
		MaxRetries: 3,
		Conditions: []Condition{
			{Type: "never"},
		},
	}

	tests := []struct {
		name       string
		statusCode int
		err        error
	}{
		{"success", 200, nil},
		{"client error", 404, nil},
		{"server error", 500, nil},
		{"network error", 0, errors.New("failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.statusCode, tt.err)
			if got {
				t.Errorf("ShouldRetry() with 'never' = true, want false")
			}
		})
	}
}

func TestParseConditions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Condition
		wantErr bool
	}{
		{
			name:  "single status code",
			input: "503",
			want:  []Condition{{Type: "status", Value: "503"}},
		},
		{
			name:  "status range 5xx",
			input: "5xx",
			want:  []Condition{{Type: "status_range", Value: "5xx"}},
		},
		{
			name:  "multiple conditions",
			input: "429,5xx",
			want: []Condition{
				{Type: "status", Value: "429"},
				{Type: "status_range", Value: "5xx"},
			},
		},
		{
			name:  "comparison operator",
			input: ">=500",
			want:  []Condition{{Type: "status_range", Value: ">=500"}},
		},
		{
			name:  "network error",
			input: "network_error",
			want:  []Condition{{Type: "network_error"}},
		},
		{
			name:  "always",
			input: "always",
			want:  []Condition{{Type: "always"}},
		},
		{
			name:  "never",
			input: "never",
			want:  []Condition{{Type: "never"}},
		},
		{
			name:  "complex mix",
			input: "429,5xx,network_error",
			want: []Condition{
				{Type: "status", Value: "429"},
				{Type: "status_range", Value: "5xx"},
				{Type: "network_error"},
			},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "spaces",
			input: " 429 , 5xx ",
			want: []Condition{
				{Type: "status", Value: "429"},
				{Type: "status_range", Value: "5xx"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConditions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ParseConditions() got %d conditions, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].Type != tt.want[i].Type || got[i].Value != tt.want[i].Value {
					t.Errorf("ParseConditions()[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMatchesStatusRange(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		rangeExpr  string
		want       bool
	}{
		// Nxx notation
		{"5xx matches 500", 500, "5xx", true},
		{"5xx matches 503", 503, "5xx", true},
		{"5xx no match 404", 404, "5xx", false},
		{"4xx matches 404", 404, "4xx", true},
		{"2xx matches 200", 200, "2xx", true},

		// Greater than or equal
		{">=500 matches 500", 500, ">=500", true},
		{">=500 matches 600", 600, ">=500", true},
		{">=500 no match 400", 400, ">=500", false},

		// Less than or equal
		{"<=400 matches 400", 400, "<=400", true},
		{"<=400 matches 200", 200, "<=400", true},
		{"<=400 no match 500", 500, "<=400", false},

		// Greater than
		{">500 matches 501", 501, ">500", true},
		{">500 no match 500", 500, ">500", false},

		// Less than
		{"<400 matches 399", 399, "<400", true},
		{"<400 no match 400", 400, "<400", false},

		// Equal
		{"==200 matches 200", 200, "==200", true},
		{"==200 no match 201", 201, "==200", false},

		// Not equal
		{"!=200 matches 404", 404, "!=200", true},
		{"!=200 no match 200", 200, "!=200", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesStatusRange(tt.statusCode, tt.rangeExpr)
			if got != tt.want {
				t.Errorf("matchesStatusRange(%d, %q) = %v, want %v", tt.statusCode, tt.rangeExpr, got, tt.want)
			}
		})
	}
}
