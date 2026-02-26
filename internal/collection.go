package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SavedRequest is a serialized request that can be stored and replayed.
type SavedRequest struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Auth    string            `json:"auth,omitempty"`
	Timeout time.Duration     `json:"timeout_ns,omitempty"`
}

// Collection is the top-level JSON structure for the collections file.
type Collection struct {
	Requests []SavedRequest `json:"requests"`
}

// collectionFilePath returns the path to the collections JSON file.
// It lives at ~/.apitester/collections.json.
func collectionFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	dir := filepath.Join(home, ".apitester")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(dir, "collections.json"), nil
}

// loadCollection reads the existing collections file, returning an empty
// Collection if none exists yet.
func loadCollection() (Collection, error) {
	path, err := collectionFilePath()
	if err != nil {
		return Collection{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Collection{}, nil
	}
	if err != nil {
		return Collection{}, fmt.Errorf("could not read collections file: %w", err)
	}

	var col Collection
	if err := json.Unmarshal(data, &col); err != nil {
		return Collection{}, fmt.Errorf("invalid collections file: %w", err)
	}
	return col, nil
}

// saveCollection writes the collection back to disk as JSON.
func saveCollection(col Collection) error {
	path, err := collectionFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(col, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize collections: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// SaveRequest saves (or overwrites) a named request into the collection.
func SaveRequest(req SavedRequest) error {
	col, err := loadCollection()
	if err != nil {
		return err
	}

	// Overwrite if a request with the same name exists.
	for i, r := range col.Requests {
		if r.Name == req.Name {
			col.Requests[i] = req
			if err := saveCollection(col); err != nil {
				return err
			}
			fmt.Printf("✅ Updated request %q in collection.\n", req.Name)
			return nil
		}
	}

	col.Requests = append(col.Requests, req)
	if err := saveCollection(col); err != nil {
		return err
	}
	fmt.Printf("✅ Saved request %q to collection.\n", req.Name)
	return nil
}

// GetRequest retrieves a saved request by name.
func GetRequest(name string) (SavedRequest, error) {
	col, err := loadCollection()
	if err != nil {
		return SavedRequest{}, err
	}

	for _, r := range col.Requests {
		if r.Name == name {
			return r, nil
		}
	}
	return SavedRequest{}, fmt.Errorf("no saved request named %q found", name)
}

// DeleteRequest removes a saved request by name.
func DeleteRequest(name string) error {
	col, err := loadCollection()
	if err != nil {
		return err
	}

	newRequests := make([]SavedRequest, 0, len(col.Requests))
	found := false
	for _, r := range col.Requests {
		if r.Name == name {
			found = true
			continue
		}
		newRequests = append(newRequests, r)
	}

	if !found {
		return fmt.Errorf("no saved request named %q found", name)
	}

	col.Requests = newRequests
	return saveCollection(col)
}

// ListRequests returns all saved requests.
func ListRequests() ([]SavedRequest, error) {
	col, err := loadCollection()
	if err != nil {
		return nil, err
	}
	return col.Requests, nil
}
