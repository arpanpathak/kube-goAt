package state

import "context"

// Store defines an interface to save and load serialized infrastructure DAGs.
// State representation is key for a scalable execution engine to diff resources.
type Store interface {
	Save(ctx context.Context, key string, data []byte) error
	Load(ctx context.Context, key string) ([]byte, error)
}
