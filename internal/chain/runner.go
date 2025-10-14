package chain

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/humancto/mozzy/internal/assertions"
	"github.com/humancto/mozzy/internal/vars"
	"github.com/humancto/mozzy/internal/httpclient"
	"github.com/humancto/mozzy/internal/formatter"
)

type Step struct {
	Name      string            `yaml:"name"`
	Method    string            `yaml:"method"`
	URL       string            `yaml:"url"`
	Headers   map[string]string `yaml:"headers,omitempty"`
	JSON      any               `yaml:"json,omitempty"`
	File      string            `yaml:"file,omitempty"`
	Capture   map[string]string `yaml:"capture,omitempty"`
	Assert    []string          `yaml:"assert,omitempty"`
	OnSuccess string            `yaml:"on_success,omitempty"` // Step name or "continue" (default) or "stop"
	OnFailure string            `yaml:"on_failure,omitempty"` // Step name or "stop" (default) or "continue"
}

type Flow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Env         string `yaml:"env"`
	Steps       []Step `yaml:"steps"`

	// populated by cmd/run
	EnvName    string `yaml:"-"`
	BaseURL    string `yaml:"-"`
	GlobalAuth string `yaml:"-"`
}

func Run(ctx context.Context, f Flow) error {
	base := vars.ResolveBase(f.BaseURL, firstNonEmpty(f.EnvName, f.Env))

	// Build step name index for jumps
	stepIndex := make(map[string]int)
	for i, s := range f.Steps {
		stepIndex[s.Name] = i
	}

	i := 0
	for i < len(f.Steps) {
		s := f.Steps[i]
		stepSuccess := true

		method := strings.ToUpper(s.Method)
		url := s.URL
		if base != "" && strings.HasPrefix(url, "/") {
			url = strings.TrimRight(base, "/") + url
		}
		url = vars.Interpolate(url)

		// headers
		hdrs := []string{}
		for k, v := range s.Headers {
			hdrs = append(hdrs, fmt.Sprintf("%s: %s", k, vars.Interpolate(v)))
		}
		token := "" // prefer explicit header if provided
		if f.GlobalAuth != "" && !hasAuthHeader(hdrs) {
			token = vars.Interpolate(f.GlobalAuth)
		}

		// body
		var body []byte
		var isJSON bool
		switch {
		case s.File != "":
			b, err := os.ReadFile(s.File)
			if err != nil {
				stepSuccess = false
				if nextStep := handleStepResult(i, s, false, stepIndex); nextStep >= 0 {
					i = nextStep
					continue
				}
				return err
			}
			body = b
		case s.JSON != nil:
			b, err := formatter.MarshalJSON(s.JSON)
			if err != nil {
				stepSuccess = false
				if nextStep := handleStepResult(i, s, false, stepIndex); nextStep >= 0 {
					i = nextStep
					continue
				}
				return err
			}
			body = b
			isJSON = true
		}

		res, resBody, ms, err := httpclient.Do(ctx, httpclient.Request{
			Method:  method,
			URL:     url,
			Headers: hdrs,
			Token:   token,
			Body:    body,
			JSON:    isJSON,
		})
		if err != nil {
			stepSuccess = false
			fmt.Fprintf(os.Stderr, "\n📋 Step %d/%d: %s\n", i+1, len(f.Steps), s.Name)
			fmt.Fprintf(os.Stderr, "❌ Request failed: %v\n", err)
			if nextStep := handleStepResult(i, s, false, stepIndex); nextStep >= 0 {
				i = nextStep
				continue
			}
			return err
		}
		defer res.Body.Close()

		fmt.Fprintf(os.Stderr, "\n📋 Step %d/%d: %s\n", i+1, len(f.Steps), s.Name)
		formatter.PrintStatusLine(method, url, res.StatusCode, ms)
		if err := formatter.PrintJSONOrText(resBody, ""); err != nil { return err }

		// Check HTTP status
		if res.StatusCode >= 400 {
			stepSuccess = false
		}

		// captures
		for name, path := range s.Capture {
			spec := fmt.Sprintf("%s=%s", name, path)
			if err := vars.Capture(resBody, spec); err != nil {
				fmt.Fprintf(os.Stderr, "warn: capture failed %q: %v\n", spec, err)
			}
		}

		// assertions
		if len(s.Assert) > 0 {
			fmt.Fprintf(os.Stderr, "\n🧪 Running assertions...\n")
			allPassed := true
			for _, assertExpr := range s.Assert {
				result, err := assertions.Evaluate(assertExpr, res.StatusCode, resBody, ms)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  ⚠️  Error: %v\n", err)
					allPassed = false
					continue
				}
				fmt.Fprintf(os.Stderr, "  %s\n", result.Message)
				if !result.Passed {
					allPassed = false
				}
			}
			if !allPassed {
				stepSuccess = false
				fmt.Fprintf(os.Stderr, "❌ Assertions failed\n")
				if nextStep := handleStepResult(i, s, false, stepIndex); nextStep >= 0 {
					i = nextStep
					continue
				}
				return fmt.Errorf("❌ assertions failed for step: %s", s.Name)
			}
			fmt.Fprintf(os.Stderr, "✅ All assertions passed\n")
		}

		// Handle conditional execution
		nextStep := handleStepResult(i, s, stepSuccess, stepIndex)
		if nextStep < 0 {
			break // Stop execution
		}
		i = nextStep
	}
	return nil
}

// handleStepResult determines the next step based on success/failure and on_success/on_failure
// Returns -1 to stop execution, or the index of the next step
func handleStepResult(currentIndex int, step Step, success bool, stepIndex map[string]int) int {
	var action string
	if success {
		action = step.OnSuccess
		if action == "" {
			action = "continue" // Default: continue to next step
		}
	} else {
		action = step.OnFailure
		if action == "" {
			action = "stop" // Default: stop on failure
		}
	}

	switch action {
	case "continue":
		return currentIndex + 1
	case "stop":
		return -1
	default:
		// Jump to named step
		if idx, exists := stepIndex[action]; exists {
			return idx
		}
		fmt.Fprintf(os.Stderr, "⚠️  Unknown step reference: %q, continuing...\n", action)
		return currentIndex + 1
	}
}

func hasAuthHeader(h []string) bool {
	for _, v := range h {
		if strings.HasPrefix(strings.ToLower(v), "authorization:") { return true }
	}
	return false
}

func firstNonEmpty(a, b string) string {
	if a != "" { return a }
	return b
}
