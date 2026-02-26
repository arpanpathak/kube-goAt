package state

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStore(t *testing.T) {
	dir, err := os.MkdirTemp("", "kube-goat-state")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	store := NewLocalStore(dir)

	ctx := context.Background()
	key := "test-state"
	data := []byte("hello world binary payload")

	// Test Save
	if err := store.Save(ctx, key, data); err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Verify file exists
	path := filepath.Join(dir, key+".gob")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file at %s to exist", path)
	}

	// Test Load
	loaded, err := store.Load(ctx, key)
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if string(loaded) != string(data) {
		t.Errorf("Expected %s, got %s", data, loaded)
	}

	// Test Load Non-Existent
	_, err = store.Load(ctx, "does-not-exist")
	if err == nil {
		t.Error("Expected error loading non-existent state")
	}
}
