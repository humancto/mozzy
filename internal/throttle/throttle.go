package throttle

import (
	"fmt"
	"io"
	"time"
)

// Bandwidth presets in bytes per second
const (
	Dial56k  = 7 * 1024      // 56 Kbps (dial-up)
	Slow     = 12.5 * 1024   // 100 Kbps
	GPRS     = 12.5 * 1024   // 100 Kbps
	EDGE     = 30 * 1024     // 240 Kbps
	ThreeG   = 93.75 * 1024  // 750 Kbps
	FourG    = 1.25 * 1024 * 1024  // 10 Mbps
	LTE      = 6.25 * 1024 * 1024  // 50 Mbps
	FiveG    = 12.5 * 1024 * 1024  // 100 Mbps
)

// Profile represents a network throttling profile
type Profile struct {
	Name      string
	Bandwidth int64 // bytes per second
	Latency   time.Duration
}

// Predefined profiles
var Profiles = map[string]Profile{
	"56k": {
		Name:      "Dial-up (56k)",
		Bandwidth: Dial56k,
		Latency:   200 * time.Millisecond,
	},
	"slow": {
		Name:      "Slow Connection",
		Bandwidth: Slow,
		Latency:   150 * time.Millisecond,
	},
	"gprs": {
		Name:      "GPRS",
		Bandwidth: GPRS,
		Latency:   500 * time.Millisecond,
	},
	"edge": {
		Name:      "EDGE",
		Bandwidth: EDGE,
		Latency:   300 * time.Millisecond,
	},
	"3g": {
		Name:      "3G",
		Bandwidth: ThreeG,
		Latency:   100 * time.Millisecond,
	},
	"4g": {
		Name:      "4G/LTE",
		Bandwidth: FourG,
		Latency:   50 * time.Millisecond,
	},
	"lte": {
		Name:      "LTE",
		Bandwidth: LTE,
		Latency:   30 * time.Millisecond,
	},
	"5g": {
		Name:      "5G",
		Bandwidth: FiveG,
		Latency:   10 * time.Millisecond,
	},
}

// GetProfile returns a throttle profile by name
func GetProfile(name string) (Profile, error) {
	profile, ok := Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("unknown throttle profile: %s", name)
	}
	return profile, nil
}

// ThrottledReader wraps an io.Reader to limit bandwidth
type ThrottledReader struct {
	reader    io.Reader
	bandwidth int64         // bytes per second
	lastRead  time.Time
	bucket    int64         // token bucket for smooth throttling
}

// NewThrottledReader creates a new bandwidth-limited reader
func NewThrottledReader(reader io.Reader, bandwidth int64) *ThrottledReader {
	return &ThrottledReader{
		reader:    reader,
		bandwidth: bandwidth,
		lastRead:  time.Now(),
		bucket:    bandwidth, // start with full bucket
	}
}

// Read implements io.Reader with bandwidth throttling
func (t *ThrottledReader) Read(p []byte) (n int, err error) {
	// Calculate time since last read
	now := time.Now()
	elapsed := now.Sub(t.lastRead)
	t.lastRead = now

	// Add tokens to bucket based on elapsed time
	tokensToAdd := int64(elapsed.Seconds() * float64(t.bandwidth))
	t.bucket += tokensToAdd

	// Cap bucket at bandwidth (don't accumulate too much)
	if t.bucket > t.bandwidth {
		t.bucket = t.bandwidth
	}

	// Determine how much we can read based on available tokens
	maxRead := int(t.bucket)
	if maxRead > len(p) {
		maxRead = len(p)
	}

	if maxRead == 0 {
		// No tokens available, sleep briefly
		time.Sleep(10 * time.Millisecond)
		return 0, nil
	}

	// Read up to maxRead bytes
	n, err = t.reader.Read(p[:maxRead])

	// Consume tokens
	t.bucket -= int64(n)

	return n, err
}

// ApplyLatency adds artificial latency to simulate network delay
func ApplyLatency(latency time.Duration) {
	if latency > 0 {
		time.Sleep(latency)
	}
}
