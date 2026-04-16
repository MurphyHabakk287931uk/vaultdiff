package vault

import (
	"fmt"
	"sync"
	"time"
)

// Checkpoint records a named point-in-time snapshot of secrets read from a path.
type Checkpoint struct {
	Name      string
	Path      string
	Secrets   map[string]string
	RecordedAt time.Time
}

// CheckpointStore holds named checkpoints in memory.
type CheckpointStore struct {
	mu    sync.RWMutex
	items map[string]*Checkpoint
}

// NewCheckpointStore returns an empty CheckpointStore.
func NewCheckpointStore() *CheckpointStore {
	return &CheckpointStore{items: make(map[string]*Checkpoint)}
}

// Record reads secrets from the given client at path and stores them under name.
func (s *CheckpointStore) Record(name, path string, client SecretReader) error {
	if name == "" {
		return fmt.Errorf("checkpoint name must not be empty")
	}
	secrets, err := client.ReadSecrets(path)
	if err != nil {
		return fmt.Errorf("checkpoint %q: %w", name, err)
	}
	cp := &Checkpoint{
		Name:       name,
		Path:       path,
		Secrets:    secrets,
		RecordedAt: time.Now(),
	}
	s.mu.Lock()
	s.items[name] = cp
	s.mu.Unlock()
	return nil
}

// Get returns the checkpoint stored under name, or an error if not found.
func (s *CheckpointStore) Get(name string) (*Checkpoint, error) {
	s.mu.RLock()
	cp, ok := s.items[name]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("checkpoint %q not found", name)
	}
	return cp, nil
}

// Names returns all stored checkpoint names.
func (s *CheckpointStore) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.items))
	for k := range s.items {
		names = append(names, k)
	}
	return names
}

// Delete removes a checkpoint by name.
func (s *CheckpointStore) Delete(name string) {
	s.mu.Lock()
	delete(s.items, name)
	s.mu.Unlock()
}
