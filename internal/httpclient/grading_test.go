package httpclient

import (
	"testing"
	"time"
)

func TestGradeFromDuration(t *testing.T) {
	tests := []struct {
		name       string
		duration   time.Duration
		thresholds []int
		want       Grade
	}{
		// DNS thresholds: [20, 50, 100, 200]
		{"DNS Grade A - 10ms", 10 * time.Millisecond, dnsThresholds, GradeA},
		{"DNS Grade A - boundary 19ms", 19 * time.Millisecond, dnsThresholds, GradeA},
		{"DNS Grade B - 20ms", 20 * time.Millisecond, dnsThresholds, GradeB},
		{"DNS Grade B - 49ms", 49 * time.Millisecond, dnsThresholds, GradeB},
		{"DNS Grade C - 50ms", 50 * time.Millisecond, dnsThresholds, GradeC},
		{"DNS Grade C - 99ms", 99 * time.Millisecond, dnsThresholds, GradeC},
		{"DNS Grade D - 100ms", 100 * time.Millisecond, dnsThresholds, GradeD},
		{"DNS Grade D - 199ms", 199 * time.Millisecond, dnsThresholds, GradeD},
		{"DNS Grade F - 200ms", 200 * time.Millisecond, dnsThresholds, GradeF},
		{"DNS Grade F - 500ms", 500 * time.Millisecond, dnsThresholds, GradeF},

		// TCP thresholds: [30, 80, 150, 300]
		{"TCP Grade A - 15ms", 15 * time.Millisecond, tcpThresholds, GradeA},
		{"TCP Grade B - 50ms", 50 * time.Millisecond, tcpThresholds, GradeB},
		{"TCP Grade C - 100ms", 100 * time.Millisecond, tcpThresholds, GradeC},
		{"TCP Grade D - 200ms", 200 * time.Millisecond, tcpThresholds, GradeD},
		{"TCP Grade F - 400ms", 400 * time.Millisecond, tcpThresholds, GradeF},

		// TLS thresholds: [50, 100, 200, 400]
		{"TLS Grade A - 30ms", 30 * time.Millisecond, tlsThresholds, GradeA},
		{"TLS Grade B - 70ms", 70 * time.Millisecond, tlsThresholds, GradeB},
		{"TLS Grade C - 150ms", 150 * time.Millisecond, tlsThresholds, GradeC},
		{"TLS Grade D - 300ms", 300 * time.Millisecond, tlsThresholds, GradeD},
		{"TLS Grade F - 500ms", 500 * time.Millisecond, tlsThresholds, GradeF},

		// TTFB thresholds: [100, 300, 800, 1500]
		{"TTFB Grade A - 50ms", 50 * time.Millisecond, ttfbThresholds, GradeA},
		{"TTFB Grade B - 200ms", 200 * time.Millisecond, ttfbThresholds, GradeB},
		{"TTFB Grade C - 500ms", 500 * time.Millisecond, ttfbThresholds, GradeC},
		{"TTFB Grade D - 1000ms", 1000 * time.Millisecond, ttfbThresholds, GradeD},
		{"TTFB Grade F - 2000ms", 2000 * time.Millisecond, ttfbThresholds, GradeF},

		// Edge cases
		{"Zero duration", 0, dnsThresholds, GradeA},
		{"Very large duration", 10 * time.Second, dnsThresholds, GradeF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gradeFromDuration(tt.duration, tt.thresholds)
			if got != tt.want {
				t.Errorf("gradeFromDuration(%v, %v) = %v, want %v", tt.duration, tt.thresholds, got, tt.want)
			}
		})
	}
}

func TestCalculateOverallGrade(t *testing.T) {
	tests := []struct {
		name   string
		grades []Grade
		want   Grade
	}{
		{"All A grades", []Grade{GradeA, GradeA, GradeA, GradeA}, GradeA},
		{"All B grades", []Grade{GradeB, GradeB, GradeB, GradeB}, GradeB},
		{"All C grades", []Grade{GradeC, GradeC, GradeC, GradeC}, GradeC},
		{"All D grades", []Grade{GradeD, GradeD, GradeD, GradeD}, GradeD},
		{"All F grades", []Grade{GradeF, GradeF, GradeF, GradeF}, GradeF},
		{"Mixed A and B - avg 3.5+", []Grade{GradeA, GradeA, GradeB, GradeB}, GradeA},
		{"Mixed B and C - avg 2.5+", []Grade{GradeB, GradeB, GradeC, GradeC}, GradeB},
		{"Mixed C and D - avg 1.5+", []Grade{GradeC, GradeC, GradeD, GradeD}, GradeC},
		{"Mixed D and F - avg 0.5+", []Grade{GradeD, GradeD, GradeF, GradeF}, GradeD},
		{"Mostly F with one D", []Grade{GradeF, GradeF, GradeF, GradeD}, GradeF},
		{"One grade only - A", []Grade{GradeA}, GradeA},
		{"One grade only - F", []Grade{GradeF}, GradeF},
		{"Empty grades", []Grade{}, GradeC},
		{"Grades with empty strings", []Grade{GradeA, "", GradeB, ""}, GradeA},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateOverallGrade(tt.grades)
			if got != tt.want {
				t.Errorf("calculateOverallGrade(%v) = %v, want %v", tt.grades, got, tt.want)
			}
		})
	}
}

func TestCalculateGrade(t *testing.T) {
	tests := []struct {
		name    string
		timings TimingInfo
		want    PerformanceGrade
	}{
		{
			name: "All excellent timings",
			timings: TimingInfo{
				DNSLookup:        10 * time.Millisecond,
				TCPConnection:    15 * time.Millisecond,
				TLSHandshake:     30 * time.Millisecond,
				ServerProcessing: 50 * time.Millisecond,
				Total:            105 * time.Millisecond,
			},
			want: PerformanceGrade{
				DNS:     GradeA,
				TCP:     GradeA,
				TLS:     GradeA,
				TTFB:    GradeA,
				Overall: GradeA,
				Total:   105 * time.Millisecond,
			},
		},
		{
			name: "All poor timings",
			timings: TimingInfo{
				DNSLookup:        300 * time.Millisecond,
				TCPConnection:    400 * time.Millisecond,
				TLSHandshake:     500 * time.Millisecond,
				ServerProcessing: 2000 * time.Millisecond,
				Total:            3200 * time.Millisecond,
			},
			want: PerformanceGrade{
				DNS:     GradeF,
				TCP:     GradeF,
				TLS:     GradeF,
				TTFB:    GradeF,
				Overall: GradeF,
				Total:   3200 * time.Millisecond,
			},
		},
		{
			name: "Mixed performance",
			timings: TimingInfo{
				DNSLookup:        15 * time.Millisecond,  // A
				TCPConnection:    200 * time.Millisecond, // D
				TLSHandshake:     60 * time.Millisecond,  // B
				ServerProcessing: 150 * time.Millisecond, // B
				Total:            425 * time.Millisecond,
			},
			want: PerformanceGrade{
				DNS:     GradeA,
				TCP:     GradeD,
				TLS:     GradeB,
				TTFB:    GradeB,
				Overall: GradeB, // avg: (4+1+3+3)/4 = 2.75 -> B
				Total:   425 * time.Millisecond,
			},
		},
		{
			name: "HTTP only - no TLS",
			timings: TimingInfo{
				DNSLookup:        10 * time.Millisecond,
				TCPConnection:    20 * time.Millisecond,
				TLSHandshake:     0, // No TLS
				ServerProcessing: 60 * time.Millisecond,
				Total:            90 * time.Millisecond,
			},
			want: PerformanceGrade{
				DNS:     GradeA,
				TCP:     GradeA,
				TLS:     "",
				TTFB:    GradeA,
				Overall: GradeA,
				Total:   90 * time.Millisecond,
			},
		},
		{
			name: "Zero DNS - cached",
			timings: TimingInfo{
				DNSLookup:        0, // Cached DNS
				TCPConnection:    25 * time.Millisecond,
				TLSHandshake:     45 * time.Millisecond,
				ServerProcessing: 80 * time.Millisecond,
				Total:            150 * time.Millisecond,
			},
			want: PerformanceGrade{
				DNS:     "",
				TCP:     GradeA,
				TLS:     GradeA,
				TTFB:    GradeA,
				Overall: GradeA,
				Total:   150 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateGrade(tt.timings)

			if got.DNS != tt.want.DNS {
				t.Errorf("DNS grade = %v, want %v", got.DNS, tt.want.DNS)
			}
			if got.TCP != tt.want.TCP {
				t.Errorf("TCP grade = %v, want %v", got.TCP, tt.want.TCP)
			}
			if got.TLS != tt.want.TLS {
				t.Errorf("TLS grade = %v, want %v", got.TLS, tt.want.TLS)
			}
			if got.TTFB != tt.want.TTFB {
				t.Errorf("TTFB grade = %v, want %v", got.TTFB, tt.want.TTFB)
			}
			if got.Overall != tt.want.Overall {
				t.Errorf("Overall grade = %v, want %v", got.Overall, tt.want.Overall)
			}
			if got.Total != tt.want.Total {
				t.Errorf("Total duration = %v, want %v", got.Total, tt.want.Total)
			}
		})
	}
}

func TestGetRecommendations(t *testing.T) {
	tests := []struct {
		name          string
		grade         PerformanceGrade
		wantMinCount  int
		wantContains  []string
		wantNotContains []string
	}{
		{
			name: "All excellent - Grade A",
			grade: PerformanceGrade{
				DNS:     GradeA,
				TCP:     GradeA,
				TLS:     GradeA,
				TTFB:    GradeA,
				Overall: GradeA,
			},
			wantMinCount: 1,
			wantContains: []string{"Excellent performance"},
		},
		{
			name: "All poor - Grade F",
			grade: PerformanceGrade{
				DNS:     GradeF,
				TCP:     GradeF,
				TLS:     GradeF,
				TTFB:    GradeF,
				Overall: GradeF,
			},
			wantMinCount: 4,
			wantContains: []string{"DNS resolution", "TCP connection", "TLS handshake", "Time to first byte"},
		},
		{
			name: "Only DNS slow",
			grade: PerformanceGrade{
				DNS:     GradeD,
				TCP:     GradeA,
				TLS:     GradeA,
				TTFB:    GradeA,
				Overall: GradeB,
			},
			wantMinCount:    2,
			wantContains:    []string{"DNS resolution", "Good performance"},
			wantNotContains: []string{"TCP connection", "TLS handshake"},
		},
		{
			name: "TLS and TTFB slow",
			grade: PerformanceGrade{
				DNS:     GradeA,
				TCP:     GradeB,
				TLS:     GradeF,
				TTFB:    GradeD,
				Overall: GradeC,
			},
			wantMinCount:    2,
			wantContains:    []string{"TLS handshake", "Time to first byte"},
			wantNotContains: []string{"DNS resolution", "Excellent performance"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.grade.GetRecommendations()

			if len(got) < tt.wantMinCount {
				t.Errorf("GetRecommendations() returned %d recommendations, want at least %d", len(got), tt.wantMinCount)
			}

			recommendations := ""
			for _, r := range got {
				recommendations += r + " "
			}

			for _, want := range tt.wantContains {
				found := false
				for _, r := range got {
					if contains(r, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetRecommendations() should contain '%s', got: %v", want, got)
				}
			}

			for _, notWant := range tt.wantNotContains {
				for _, r := range got {
					if contains(r, notWant) {
						t.Errorf("GetRecommendations() should NOT contain '%s', got: %v", notWant, got)
					}
				}
			}
		})
	}
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
