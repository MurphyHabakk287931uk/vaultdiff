package vault

import (
	"errors"
	"fmt"
	"sync"
)

// Pool holds a named set of SecretReader clients and dispatches reads
// to the client registered under the requested name.
type Pool struct {
	mu      sync.RWMutex
	clients map[string]SecretReader
}

// NewPool returns an empty, ready-to-use Pool.
func NewPool() *Pool {
	return &Pool{
		clients: make(map[string]SecretReader),
	}
}

// Register adds a named client to the pool. Registering the same name
// twice overwrites the previous entry.
func (p *Pool) Register(name string, client SecretReader) error {
	if name == "" {
		return errors.New("vault pool: name must not be empty")
	}
	if client == nil {
		return fmt.Errorf("vault pool: client for %q must not be nil", name)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[name] = client
	return nil
}

// Remove deletes the named client from the pool. It is a no-op if the
// name is not registered.
func (p *Pool) Remove(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, name)
}

// Get returns the SecretReader registered under name, or an error if
// no such client exists.
func (p *Pool) Get(name string) (SecretReader, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	c, ok := p.clients[name]
	if !ok {
		return nil, fmt.Errorf("vault pool: no client registered for %q", name)
	}
	return c, nil
}

// ReadSecrets satisfies SecretReader by routing the request to the
// client whose name matches the first segment of path ("name/rest…").
// The prefix segment is stripped before forwarding.
func (p *Pool) ReadSecrets(path string) (map[string]string, error) {
	name, rest, err := splitPoolPath(path)
	if err != nil {
		return nil, err
	}
	client, err := p.Get(name)
	if err != nil {
		return nil, err
	}
	return client.ReadSecrets(rest)
}

// Names returns the sorted list of registered client names.
func (p *Pool) Names() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	names := make([]string, 0, len(p.clients))
	for n := range p.clients {
		names = append(names, n)
	}
	return names
}

func splitPoolPath(path string) (name, rest string, err error) {
	for i, ch := range path {
		if ch == '/' {
			return path[:i], path[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("vault pool: path %q has no name prefix", path)
}
