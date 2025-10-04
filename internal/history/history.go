package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Timestamp time.Time     `json:"ts"`
	Method    string        `json:"method"`
	URL       string        `json:"url"`
	Status    int           `json:"status"`
	Duration  time.Duration `json:"duration"`
	BodySize  int           `json:"body_size,omitempty"`
}

func path() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mozzy", "history.json")
}

func Append(e Entry) error {
	p := path()
	_ = os.MkdirAll(filepath.Dir(p), 0o755)

	var list []Entry
	if data, err := os.ReadFile(p); err == nil {
		_ = json.Unmarshal(data, &list)
	}
	list = append(list, e)
	data, _ := json.MarshalIndent(list, "", "  ")
	return os.WriteFile(p, data, 0o644)
}

func Load() ([]Entry, error) {
	p := path()
	data, err := os.ReadFile(p)
	if err != nil { return []Entry{}, nil }
	var list []Entry
	_ = json.Unmarshal(data, &list)
	return list, nil
}
