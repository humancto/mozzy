package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Request struct {
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        string            `json:"body,omitempty"`
	Description string            `json:"description,omitempty"`
}

type Collection struct {
	Requests map[string]Request `json:"requests"`
}

func collectionPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mozzy", "collections.json")
}

func Load() (*Collection, error) {
	p := collectionPath()
	data, err := os.ReadFile(p)
	if err != nil {
		// Return empty collection if file doesn't exist
		if os.IsNotExist(err) {
			return &Collection{Requests: make(map[string]Request)}, nil
		}
		return nil, err
	}

	var coll Collection
	if err := json.Unmarshal(data, &coll); err != nil {
		return nil, err
	}
	if coll.Requests == nil {
		coll.Requests = make(map[string]Request)
	}
	return &coll, nil
}

func (c *Collection) Save() error {
	p := collectionPath()
	_ = os.MkdirAll(filepath.Dir(p), 0o755)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func (c *Collection) Add(req Request) error {
	if req.Name == "" {
		return fmt.Errorf("request name is required")
	}
	c.Requests[req.Name] = req
	return c.Save()
}

func (c *Collection) Get(name string) (Request, error) {
	req, ok := c.Requests[name]
	if !ok {
		return Request{}, fmt.Errorf("request %q not found in collection", name)
	}
	return req, nil
}

func (c *Collection) Delete(name string) error {
	delete(c.Requests, name)
	return c.Save()
}

func (c *Collection) List() []Request {
	var reqs []Request
	for _, req := range c.Requests {
		reqs = append(reqs, req)
	}
	return reqs
}
