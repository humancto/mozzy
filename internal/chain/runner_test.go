package chain

import (
	"testing"
)

func TestHandleStepResult_Success(t *testing.T) {
	stepIndex := map[string]int{
		"step1": 0,
		"step2": 1,
		"step3": 2,
		"cleanup": 3,
	}

	tests := []struct {
		name         string
		currentIndex int
		step         Step
		success      bool
		wantNext     int
	}{
		{
			name:         "success - default continue",
			currentIndex: 0,
			step:         Step{Name: "step1"},
			success:      true,
			wantNext:     1,
		},
		{
			name:         "success - explicit continue",
			currentIndex: 0,
			step:         Step{Name: "step1", OnSuccess: "continue"},
			success:      true,
			wantNext:     1,
		},
		{
			name:         "success - stop",
			currentIndex: 0,
			step:         Step{Name: "step1", OnSuccess: "stop"},
			success:      true,
			wantNext:     -1,
		},
		{
			name:         "success - jump to step",
			currentIndex: 0,
			step:         Step{Name: "step1", OnSuccess: "step3"},
			success:      true,
			wantNext:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleStepResult(tt.currentIndex, tt.step, tt.success, stepIndex)
			if got != tt.wantNext {
				t.Errorf("handleStepResult() = %v, want %v", got, tt.wantNext)
			}
		})
	}
}

func TestHandleStepResult_Failure(t *testing.T) {
	stepIndex := map[string]int{
		"step1": 0,
		"step2": 1,
		"step3": 2,
		"cleanup": 3,
	}

	tests := []struct {
		name         string
		currentIndex int
		step         Step
		success      bool
		wantNext     int
	}{
		{
			name:         "failure - default stop",
			currentIndex: 0,
			step:         Step{Name: "step1"},
			success:      false,
			wantNext:     -1,
		},
		{
			name:         "failure - explicit stop",
			currentIndex: 0,
			step:         Step{Name: "step1", OnFailure: "stop"},
			success:      false,
			wantNext:     -1,
		},
		{
			name:         "failure - continue",
			currentIndex: 0,
			step:         Step{Name: "step1", OnFailure: "continue"},
			success:      false,
			wantNext:     1,
		},
		{
			name:         "failure - jump to cleanup",
			currentIndex: 0,
			step:         Step{Name: "step1", OnFailure: "cleanup"},
			success:      false,
			wantNext:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleStepResult(tt.currentIndex, tt.step, tt.success, stepIndex)
			if got != tt.wantNext {
				t.Errorf("handleStepResult() = %v, want %v", got, tt.wantNext)
			}
		})
	}
}

func TestHandleStepResult_UnknownStep(t *testing.T) {
	stepIndex := map[string]int{
		"step1": 0,
		"step2": 1,
	}

	step := Step{Name: "step1", OnSuccess: "nonexistent"}
	got := handleStepResult(0, step, true, stepIndex)

	// Should continue to next step when unknown step referenced
	if got != 1 {
		t.Errorf("handleStepResult() with unknown step = %v, want 1 (continue)", got)
	}
}

func TestHandleStepResult_ComplexFlow(t *testing.T) {
	// Simulates a flow like:
	// 1. Try API call
	// 2. On success -> continue to next
	// 3. On failure -> jump to retry
	// 4. Retry step
	// 5. On retry success -> jump to success handler
	// 6. On retry failure -> jump to error handler

	stepIndex := map[string]int{
		"try_api":        0,
		"process_result": 1,
		"retry":          2,
		"success":        3,
		"error":          4,
	}

	tests := []struct {
		name         string
		currentIndex int
		step         Step
		success      bool
		wantNext     int
		description  string
	}{
		{
			name:         "initial try succeeds",
			currentIndex: 0,
			step:         Step{Name: "try_api", OnSuccess: "process_result", OnFailure: "retry"},
			success:      true,
			wantNext:     1,
			description:  "Success on first try, go to process_result",
		},
		{
			name:         "initial try fails",
			currentIndex: 0,
			step:         Step{Name: "try_api", OnSuccess: "process_result", OnFailure: "retry"},
			success:      false,
			wantNext:     2,
			description:  "Failure on first try, go to retry",
		},
		{
			name:         "retry succeeds",
			currentIndex: 2,
			step:         Step{Name: "retry", OnSuccess: "success", OnFailure: "error"},
			success:      true,
			wantNext:     3,
			description:  "Retry successful, go to success handler",
		},
		{
			name:         "retry fails",
			currentIndex: 2,
			step:         Step{Name: "retry", OnSuccess: "success", OnFailure: "error"},
			success:      false,
			wantNext:     4,
			description:  "Retry failed, go to error handler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleStepResult(tt.currentIndex, tt.step, tt.success, stepIndex)
			if got != tt.wantNext {
				t.Errorf("%s: handleStepResult() = %v, want %v", tt.description, got, tt.wantNext)
			}
		})
	}
}

func TestHandleStepResult_StopExecution(t *testing.T) {
	stepIndex := map[string]int{
		"step1": 0,
		"step2": 1,
		"step3": 2,
	}

	tests := []struct {
		name    string
		step    Step
		success bool
	}{
		{
			name:    "explicit stop on success",
			step:    Step{Name: "step1", OnSuccess: "stop"},
			success: true,
		},
		{
			name:    "explicit stop on failure",
			step:    Step{Name: "step1", OnFailure: "stop"},
			success: false,
		},
		{
			name:    "default stop on failure",
			step:    Step{Name: "step1"}, // No OnFailure means default stop
			success: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleStepResult(0, tt.step, tt.success, stepIndex)
			if got != -1 {
				t.Errorf("handleStepResult() should stop (-1), got %v", got)
			}
		})
	}
}

func TestHandleStepResult_ContinueExecution(t *testing.T) {
	stepIndex := map[string]int{
		"step1": 0,
		"step2": 1,
		"step3": 2,
	}

	tests := []struct {
		name         string
		currentIndex int
		step         Step
		success      bool
		wantNext     int
	}{
		{
			name:         "explicit continue on success",
			currentIndex: 0,
			step:         Step{Name: "step1", OnSuccess: "continue"},
			success:      true,
			wantNext:     1,
		},
		{
			name:         "explicit continue on failure",
			currentIndex: 0,
			step:         Step{Name: "step1", OnFailure: "continue"},
			success:      false,
			wantNext:     1,
		},
		{
			name:         "default continue on success",
			currentIndex: 1,
			step:         Step{Name: "step2"}, // No OnSuccess means default continue
			success:      true,
			wantNext:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleStepResult(tt.currentIndex, tt.step, tt.success, stepIndex)
			if got != tt.wantNext {
				t.Errorf("handleStepResult() = %v, want %v", got, tt.wantNext)
			}
		})
	}
}

func TestHandleStepResult_BackwardJump(t *testing.T) {
	// Test jumping backward (for retry loops)
	stepIndex := map[string]int{
		"try":     0,
		"process": 1,
		"verify":  2,
	}

	step := Step{Name: "verify", OnFailure: "try"} // Jump back to try
	got := handleStepResult(2, step, false, stepIndex)

	if got != 0 {
		t.Errorf("handleStepResult() with backward jump = %v, want 0 (jump to 'try')", got)
	}
}

func TestHandleStepResult_ForwardJump(t *testing.T) {
	// Test jumping forward (skipping steps)
	stepIndex := map[string]int{
		"check":   0,
		"step1":   1,
		"step2":   2,
		"cleanup": 3,
	}

	step := Step{Name: "check", OnSuccess: "cleanup"} // Skip step1 and step2
	got := handleStepResult(0, step, true, stepIndex)

	if got != 3 {
		t.Errorf("handleStepResult() with forward jump = %v, want 3 (jump to 'cleanup')", got)
	}
}
