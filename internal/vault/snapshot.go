package vault

import (
	"fmt"
	"sort"
	"time"
)

// Snapshot captures the secrets at a given path at a point in time.
type Snapshot struct {
	Path      string
	Secrets   map[string]string
	CapturedAt time.Time
}

// SnapshotClient wraps a SecretReader and records snapshots of reads.
type SnapshotClient struct {
	inner     SecretReader
	snapshots []*Snapshot
}

// NewSnapshotClient wraps inner, recording every successful read as a Snapshot.
func NewSnapshotClient(inner SecretReader) *SnapshotClient {
	if inner == nil {
		panic("vault: NewSnapshotClient: inner must not be nil")
	}
	return &SnapshotClient{inner: inner}
}

// ReadSecrets delegates to the inner client and stores the result.
func (s *SnapshotClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := s.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}

	// deep-copy so mutations after the call don't corrupt the snapshot
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}

	s.snapshots = append(s.snapshots, &Snapshot{
		Path:      path,
		Secrets:   copy,
		CapturedAt: time.Now().UTC(),
	})
	return secrets, nil
}

// Snapshots returns all recorded snapshots in capture order.
func (s *SnapshotClient) Snapshots() []*Snapshot {
	out := make([]*Snapshot, len(s.snapshots))
	copy(out, s.snapshots)
	return out
}

// Latest returns the most recent snapshot for path, or an error if none exists.
func (s *SnapshotClient) Latest(path string) (*Snapshot, error) {
	for i := len(s.snapshots) - 1; i >= 0; i-- {
		if s.snapshots[i].Path == path {
			return s.snapshots[i], nil
		}
	}
	return nil, fmt.Errorf("vault: no snapshot found for path %q", path)
}

// SortedKeys returns the snapshot's secret keys in sorted order.
func (snap *Snapshot) SortedKeys() []string {
	keys := make([]string, 0, len(snap.Secrets))
	for k := range snap.Secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
