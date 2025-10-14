package assertions

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Assertion struct {
	Expression string
	Passed     bool
	Message    string
}

// Evaluate checks if an assertion passes
// Supported formats:
//   - status == 200
//   - status >= 200 && status < 300
//   - .name == "Alice"
//   - .email contains "@example.com"
//   - response_time < 500ms
//   - .items[0].id exists
//   - length(.items) > 0
func Evaluate(expr string, statusCode int, body []byte, responseTime time.Duration) (*Assertion, error) {
	expr = strings.TrimSpace(expr)

	result := &Assertion{
		Expression: expr,
		Passed:     false,
	}

	// Parse JSON body once
	var jsonData interface{}
	if len(body) > 0 {
		json.Unmarshal(body, &jsonData)
	}

	// Handle status assertions
	if strings.HasPrefix(expr, "status") {
		passed, err := evaluateStatus(expr, statusCode)
		if err != nil {
			return nil, err
		}
		result.Passed = passed
		if passed {
			result.Message = fmt.Sprintf("✓ Status code %d matches condition", statusCode)
		} else {
			result.Message = fmt.Sprintf("✗ Status code %d does not match condition", statusCode)
		}
		return result, nil
	}

	// Handle response_time assertions
	if strings.HasPrefix(expr, "response_time") {
		passed, err := evaluateResponseTime(expr, responseTime)
		if err != nil {
			return nil, err
		}
		result.Passed = passed
		if passed {
			result.Message = fmt.Sprintf("✓ Response time %s matches condition", responseTime)
		} else {
			result.Message = fmt.Sprintf("✗ Response time %s does not match condition", responseTime)
		}
		return result, nil
	}

	// Handle JSON path assertions (starts with .)
	if strings.HasPrefix(expr, ".") {
		passed, msg, err := evaluateJSONPath(expr, jsonData)
		if err != nil {
			return nil, err
		}
		result.Passed = passed
		result.Message = msg
		return result, nil
	}

	// Handle length() assertions
	if strings.HasPrefix(expr, "length(") {
		passed, msg, err := evaluateLength(expr, jsonData)
		if err != nil {
			return nil, err
		}
		result.Passed = passed
		result.Message = msg
		return result, nil
	}

	return nil, fmt.Errorf("unsupported assertion format: %s", expr)
}

func evaluateStatus(expr string, statusCode int) (bool, error) {
	// Remove "status" prefix and trim
	condition := strings.TrimSpace(strings.TrimPrefix(expr, "status"))

	// Parse operators: ==, !=, >, <, >=, <=
	if strings.HasPrefix(condition, "==") {
		expected, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(condition, "==")))
		if err != nil {
			return false, err
		}
		return statusCode == expected, nil
	}
	if strings.HasPrefix(condition, "!=") {
		expected, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(condition, "!=")))
		if err != nil {
			return false, err
		}
		return statusCode != expected, nil
	}
	if strings.HasPrefix(condition, ">=") {
		expected, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(condition, ">=")))
		if err != nil {
			return false, err
		}
		return statusCode >= expected, nil
	}
	if strings.HasPrefix(condition, "<=") {
		expected, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(condition, "<=")))
		if err != nil {
			return false, err
		}
		return statusCode <= expected, nil
	}
	if strings.HasPrefix(condition, ">") {
		expected, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(condition, ">")))
		if err != nil {
			return false, err
		}
		return statusCode > expected, nil
	}
	if strings.HasPrefix(condition, "<") {
		expected, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(condition, "<")))
		if err != nil {
			return false, err
		}
		return statusCode < expected, nil
	}

	return false, fmt.Errorf("invalid status condition: %s", condition)
}

func evaluateResponseTime(expr string, responseTime time.Duration) (bool, error) {
	// response_time < 500ms
	condition := strings.TrimSpace(strings.TrimPrefix(expr, "response_time"))

	var operator string
	var valueStr string

	if strings.Contains(condition, "<") {
		operator = "<"
		valueStr = strings.TrimSpace(strings.TrimPrefix(condition, "<"))
	} else if strings.Contains(condition, ">") {
		operator = ">"
		valueStr = strings.TrimSpace(strings.TrimPrefix(condition, ">"))
	} else {
		return false, fmt.Errorf("invalid response_time condition: %s", condition)
	}

	expectedDuration, err := time.ParseDuration(valueStr)
	if err != nil {
		return false, fmt.Errorf("invalid duration: %s", valueStr)
	}

	switch operator {
	case "<":
		return responseTime < expectedDuration, nil
	case ">":
		return responseTime > expectedDuration, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func evaluateJSONPath(expr string, data interface{}) (bool, string, error) {
	// Split on operators: ==, !=, contains, exists
	var path, operator, expected string

	if strings.Contains(expr, " contains ") {
		parts := strings.SplitN(expr, " contains ", 2)
		path = strings.TrimSpace(parts[0])
		operator = "contains"
		expected = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	} else if strings.Contains(expr, " exists") {
		path = strings.TrimSpace(strings.TrimSuffix(expr, " exists"))
		operator = "exists"
	} else if strings.Contains(expr, " == ") {
		parts := strings.SplitN(expr, " == ", 2)
		path = strings.TrimSpace(parts[0])
		operator = "=="
		expected = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	} else if strings.Contains(expr, " != ") {
		parts := strings.SplitN(expr, " != ", 2)
		path = strings.TrimSpace(parts[0])
		operator = "!="
		expected = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	} else {
		return false, "", fmt.Errorf("invalid JSON path expression: %s", expr)
	}

	// Navigate JSON path
	value, exists := navigateJSONPath(path, data)

	switch operator {
	case "exists":
		if exists {
			return true, fmt.Sprintf("✓ Path %s exists", path), nil
		}
		return false, fmt.Sprintf("✗ Path %s does not exist", path), nil

	case "==":
		if !exists {
			return false, fmt.Sprintf("✗ Path %s does not exist", path), nil
		}
		valueStr := fmt.Sprintf("%v", value)
		if valueStr == expected {
			return true, fmt.Sprintf("✓ %s == %s", path, expected), nil
		}
		return false, fmt.Sprintf("✗ %s is %s, expected %s", path, valueStr, expected), nil

	case "!=":
		if !exists {
			return true, fmt.Sprintf("✓ Path %s does not exist (!=)", path), nil
		}
		valueStr := fmt.Sprintf("%v", value)
		if valueStr != expected {
			return true, fmt.Sprintf("✓ %s != %s", path, expected), nil
		}
		return false, fmt.Sprintf("✗ %s is %s, expected not equal", path, valueStr), nil

	case "contains":
		if !exists {
			return false, fmt.Sprintf("✗ Path %s does not exist", path), nil
		}
		valueStr := fmt.Sprintf("%v", value)
		if strings.Contains(valueStr, expected) {
			return true, fmt.Sprintf("✓ %s contains '%s'", path, expected), nil
		}
		return false, fmt.Sprintf("✗ %s (%s) does not contain '%s'", path, valueStr, expected), nil
	}

	return false, "", fmt.Errorf("unsupported operator: %s", operator)
}

func evaluateLength(expr string, data interface{}) (bool, string, error) {
	// length(.items) > 0
	re := regexp.MustCompile(`length\(([^)]+)\)\s*([><=!]+)\s*(\d+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) != 4 {
		return false, "", fmt.Errorf("invalid length expression: %s", expr)
	}

	path := matches[1]
	operator := matches[2]
	expectedStr := matches[3]

	expected, err := strconv.Atoi(expectedStr)
	if err != nil {
		return false, "", err
	}

	value, exists := navigateJSONPath(path, data)
	if !exists {
		return false, fmt.Sprintf("✗ Path %s does not exist", path), nil
	}

	var length int
	switch v := value.(type) {
	case []interface{}:
		length = len(v)
	case map[string]interface{}:
		length = len(v)
	case string:
		length = len(v)
	default:
		return false, "", fmt.Errorf("cannot get length of %T", value)
	}

	var passed bool
	switch operator {
	case ">":
		passed = length > expected
	case "<":
		passed = length < expected
	case ">=":
		passed = length >= expected
	case "<=":
		passed = length <= expected
	case "==":
		passed = length == expected
	case "!=":
		passed = length != expected
	default:
		return false, "", fmt.Errorf("unsupported operator: %s", operator)
	}

	if passed {
		return true, fmt.Sprintf("✓ length(%s) is %d %s %d", path, length, operator, expected), nil
	}
	return false, fmt.Sprintf("✗ length(%s) is %d, expected %s %d", path, length, operator, expected), nil
}

func navigateJSONPath(path string, data interface{}) (interface{}, bool) {
	if data == nil {
		return nil, false
	}

	// Remove leading dot
	path = strings.TrimPrefix(path, ".")
	if path == "" {
		return data, true
	}

	// Split path by dots and brackets
	segments := parsePathSegments(path)

	current := data
	for _, seg := range segments {
		if seg.isIndex {
			// Array index
			arr, ok := current.([]interface{})
			if !ok || seg.index < 0 || seg.index >= len(arr) {
				return nil, false
			}
			current = arr[seg.index]
		} else {
			// Object key
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, false
			}
			val, exists := obj[seg.key]
			if !exists {
				return nil, false
			}
			current = val
		}
	}

	return current, true
}

type pathSegment struct {
	key     string
	isIndex bool
	index   int
}

func parsePathSegments(path string) []pathSegment {
	var segments []pathSegment
	current := ""
	inBracket := false

	for i := 0; i < len(path); i++ {
		ch := path[i]

		switch ch {
		case '[':
			if current != "" {
				segments = append(segments, pathSegment{key: current})
				current = ""
			}
			inBracket = true
		case ']':
			if inBracket {
				idx, _ := strconv.Atoi(current)
				segments = append(segments, pathSegment{isIndex: true, index: idx})
				current = ""
				inBracket = false
			}
		case '.':
			if !inBracket && current != "" {
				segments = append(segments, pathSegment{key: current})
				current = ""
			} else if !inBracket {
				// Skip leading dots
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}

	if current != "" {
		segments = append(segments, pathSegment{key: current})
	}

	return segments
}
