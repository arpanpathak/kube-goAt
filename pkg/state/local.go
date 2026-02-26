package state

import (
	"context"
	"os"
	"path/filepath"
)

// LocalStore implements Store using the local filesystem.
type LocalStore struct {
	dir string
}

// NewLocalStore initializes a state store in the given directory.
func NewLocalStore(dir string) *LocalStore {
	os.MkdirAll(dir, 0755)
	return &LocalStore{dir: dir}
}

// Save writes the binary gob to a local file.
func (l *LocalStore) Save(ctx context.Context, key string, data []byte) error {
	path := filepath.Join(l.dir, key+".gob")
	return os.WriteFile(path, data, 0644)
}

// Load retrieves the binary gob from a local file.
func (l *LocalStore) Load(ctx context.Context, key string) ([]byte, error) {
	path := filepath.Join(l.dir, key+".gob")
	return os.ReadFile(path)
}
