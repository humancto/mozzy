package retry

import (
	"fmt"
	"strconv"
	"strings"
)

// Condition represents a retry condition
type Condition struct {
	Type  string // "status", "status_range", "always", "never"
	Value string // e.g., "500", "5xx", "4xx"
}

// Policy defines retry behavior
type Policy struct {
	MaxRetries int
	Conditions []Condition
}

// ShouldRetry determines if a request should be retried based on conditions
func (p *Policy) ShouldRetry(statusCode int, err error) bool {
	// Always retry on network errors (if conditions allow)
	if err != nil {
		return p.matchesConditions(0, true)
	}

	return p.matchesConditions(statusCode, false)
}

func (p *Policy) matchesConditions(statusCode int, isNetworkError bool) bool {
	// If no conditions specified, use default behavior (retry on 5xx and network errors)
	if len(p.Conditions) == 0 {
		if isNetworkError {
			return true
		}
		return statusCode >= 500
	}

	// Check each condition
	for _, cond := range p.Conditions {
		if p.matchesCondition(cond, statusCode, isNetworkError) {
			return true
		}
	}

	return false
}

func (p *Policy) matchesCondition(cond Condition, statusCode int, isNetworkError bool) bool {
	switch cond.Type {
	case "always":
		return true

	case "never":
		return false

	case "network_error":
		return isNetworkError

	case "status":
		// Exact status code match
		if targetStatus, err := strconv.Atoi(cond.Value); err == nil {
			return statusCode == targetStatus
		}
		return false

	case "status_range":
		// Range like "5xx", "4xx", ">=500", "<400"
		return matchesStatusRange(statusCode, cond.Value)

	default:
		return false
	}
}

func matchesStatusRange(statusCode int, rangeExpr string) bool {
	rangeExpr = strings.TrimSpace(rangeExpr)

	// Handle Nxx notation (e.g., "5xx", "4xx")
	if strings.HasSuffix(rangeExpr, "xx") {
		prefix := strings.TrimSuffix(rangeExpr, "xx")
		if len(prefix) == 1 {
			statusPrefix := fmt.Sprintf("%d", statusCode)
			return strings.HasPrefix(statusPrefix, prefix)
		}
	}

	// Handle comparison operators (>=, <=, >, <, ==, !=)
	if strings.HasPrefix(rangeExpr, ">=") {
		val, err := strconv.Atoi(strings.TrimSpace(rangeExpr[2:]))
		return err == nil && statusCode >= val
	}
	if strings.HasPrefix(rangeExpr, "<=") {
		val, err := strconv.Atoi(strings.TrimSpace(rangeExpr[2:]))
		return err == nil && statusCode <= val
	}
	if strings.HasPrefix(rangeExpr, ">") {
		val, err := strconv.Atoi(strings.TrimSpace(rangeExpr[1:]))
		return err == nil && statusCode > val
	}
	if strings.HasPrefix(rangeExpr, "<") {
		val, err := strconv.Atoi(strings.TrimSpace(rangeExpr[1:]))
		return err == nil && statusCode < val
	}
	if strings.HasPrefix(rangeExpr, "==") {
		val, err := strconv.Atoi(strings.TrimSpace(rangeExpr[2:]))
		return err == nil && statusCode == val
	}
	if strings.HasPrefix(rangeExpr, "!=") {
		val, err := strconv.Atoi(strings.TrimSpace(rangeExpr[2:]))
		return err == nil && statusCode != val
	}

	return false
}

// ParseConditions parses retry condition strings
// Examples:
//   - "5xx"           -> retry on 5xx errors
//   - ">=500"         -> retry on status >= 500
//   - "503"           -> retry only on 503
//   - "4xx,5xx"       -> retry on 4xx or 5xx
//   - "network_error" -> retry on network errors only
func ParseConditions(condStr string) ([]Condition, error) {
	if condStr == "" {
		return nil, nil // Use default behavior
	}

	parts := strings.Split(condStr, ",")
	conditions := make([]Condition, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		var cond Condition

		switch part {
		case "always":
			cond = Condition{Type: "always"}
		case "never":
			cond = Condition{Type: "never"}
		case "network_error":
			cond = Condition{Type: "network_error"}
		default:
			// Check if it's a status code or range
			if _, err := strconv.Atoi(part); err == nil {
				// Exact status code
				cond = Condition{Type: "status", Value: part}
			} else {
				// Assume it's a range expression
				cond = Condition{Type: "status_range", Value: part}
			}
		}

		conditions = append(conditions, cond)
	}

	return conditions, nil
}
