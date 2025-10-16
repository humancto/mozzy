package httpclient

import (
	"fmt"
	"time"
)

// Grade represents a performance grade from A to F
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// PerformanceGrade contains individual grades for each timing metric
type PerformanceGrade struct {
	DNS      Grade
	TCP      Grade
	TLS      Grade
	TTFB     Grade
	Overall  Grade
	Total    time.Duration
}

// Thresholds for grading (in milliseconds)
var (
	dnsThresholds = []int{20, 50, 100, 200}    // A: <20ms, B: <50ms, C: <100ms, D: <200ms, F: >=200ms
	tcpThresholds = []int{30, 80, 150, 300}    // A: <30ms, B: <80ms, C: <150ms, D: <300ms, F: >=300ms
	tlsThresholds = []int{50, 100, 200, 400}   // A: <50ms, B: <100ms, C: <200ms, D: <400ms, F: >=400ms
	ttfbThresholds = []int{100, 300, 800, 1500} // A: <100ms, B: <300ms, C: <800ms, D: <1500ms, F: >=1500ms
)

// gradeFromDuration assigns a grade based on duration and thresholds
func gradeFromDuration(d time.Duration, thresholds []int) Grade {
	ms := d.Milliseconds()

	if ms < int64(thresholds[0]) {
		return GradeA
	}
	if ms < int64(thresholds[1]) {
		return GradeB
	}
	if ms < int64(thresholds[2]) {
		return GradeC
	}
	if ms < int64(thresholds[3]) {
		return GradeD
	}
	return GradeF
}

// calculateOverallGrade computes overall grade from individual grades
func calculateOverallGrade(grades []Grade) Grade {
	if len(grades) == 0 {
		return GradeC
	}

	points := 0
	count := 0

	for _, g := range grades {
		if g == "" {
			continue
		}
		count++
		switch g {
		case GradeA:
			points += 4
		case GradeB:
			points += 3
		case GradeC:
			points += 2
		case GradeD:
			points += 1
		case GradeF:
			points += 0
		}
	}

	if count == 0 {
		return GradeC
	}

	avg := float64(points) / float64(count)

	if avg >= 3.5 {
		return GradeA
	}
	if avg >= 2.5 {
		return GradeB
	}
	if avg >= 1.5 {
		return GradeC
	}
	if avg >= 0.5 {
		return GradeD
	}
	return GradeF
}

// CalculateGrade computes performance grades for all timing metrics
func CalculateGrade(timings TimingInfo) PerformanceGrade {
	pg := PerformanceGrade{
		Total: timings.Total,
	}

	var grades []Grade

	// Grade DNS lookup
	if timings.DNSLookup > 0 {
		pg.DNS = gradeFromDuration(timings.DNSLookup, dnsThresholds)
		grades = append(grades, pg.DNS)
	}

	// Grade TCP connection
	if timings.TCPConnection > 0 {
		pg.TCP = gradeFromDuration(timings.TCPConnection, tcpThresholds)
		grades = append(grades, pg.TCP)
	}

	// Grade TLS handshake
	if timings.TLSHandshake > 0 {
		pg.TLS = gradeFromDuration(timings.TLSHandshake, tlsThresholds)
		grades = append(grades, pg.TLS)
	}

	// Grade Server Processing (Time to First Byte)
	if timings.ServerProcessing > 0 {
		pg.TTFB = gradeFromDuration(timings.ServerProcessing, ttfbThresholds)
		grades = append(grades, pg.TTFB)
	}

	// Calculate overall grade
	pg.Overall = calculateOverallGrade(grades)

	return pg
}

// GetRecommendations returns performance improvement suggestions
func (pg PerformanceGrade) GetRecommendations() []string {
	var recommendations []string

	if pg.DNS == GradeD || pg.DNS == GradeF {
		recommendations = append(recommendations, "• DNS resolution is slow. Consider using a faster DNS provider (e.g., 1.1.1.1, 8.8.8.8)")
	}

	if pg.TCP == GradeD || pg.TCP == GradeF {
		recommendations = append(recommendations, "• TCP connection time is high. Check network latency or use a CDN closer to your location")
	}

	if pg.TLS == GradeD || pg.TLS == GradeF {
		recommendations = append(recommendations, "• TLS handshake is slow. Server may need TLS session resumption or HTTP/2")
	}

	if pg.TTFB == GradeD || pg.TTFB == GradeF {
		recommendations = append(recommendations, "• Time to first byte is high. Server processing may be slow or needs caching")
	}

	if pg.Overall == GradeA {
		recommendations = append(recommendations, "✨ Excellent performance! This endpoint is well-optimized")
	} else if pg.Overall == GradeB {
		recommendations = append(recommendations, "👍 Good performance overall. Minor optimizations possible")
	}

	return recommendations
}

// FormatGradeWithColor returns a colored grade string for terminal display
func FormatGradeWithColor(g Grade) string {
	if g == "" {
		return ""
	}

	switch g {
	case GradeA:
		return fmt.Sprintf("\033[32m%s\033[0m", g) // green
	case GradeB:
		return fmt.Sprintf("\033[36m%s\033[0m", g) // cyan
	case GradeC:
		return fmt.Sprintf("\033[33m%s\033[0m", g) // yellow
	case GradeD:
		return fmt.Sprintf("\033[38;5;208m%s\033[0m", g) // orange
	case GradeF:
		return fmt.Sprintf("\033[31m%s\033[0m", g) // red
	default:
		return string(g)
	}
}
