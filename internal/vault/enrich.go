package vault

// EnrichClient wraps a SecretReader and injects additional static key/value
// pairs into every response. Keys provided via the enrichment map always
// override values returned by the inner client.

import "context"

// EnrichClient adds static metadata to every secret read.
type EnrichClient struct {
	inner  SecretReader
	extra  map[string]string
}

// NewEnrichClient returns an EnrichClient that merges extra into every
// response from inner. Panics if inner is nil or extra is nil.
func NewEnrichClient(inner SecretReader, extra map[string]string) *EnrichClient {
	if inner == nil {
		panic("enrich: inner client must not be nil")
	}
	if extra == nil {
		panic("enrich: extra map must not be nil")
	}
	// defensive copy
	copy := make(map[string]string, len(extra))
	for k, v := range extra {
		copy[k] = v
	}
	return &EnrichClient{inner: inner, extra: copy}
}

// ReadSecrets delegates to the inner client then merges the static extra
// keys, with extra taking precedence over inner values.
func (e *EnrichClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	secrets, err := e.inner.ReadSecrets(ctx, path)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(secrets)+len(e.extra))
	for k, v := range secrets {
		result[k] = v
	}
	for k, v := range e.extra {
		result[k] = v
	}
	return result, nil
}
