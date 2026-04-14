package vault

import (
	"context"
	"fmt"
	"sync"
)

// BatchResult holds the result of a single path read within a batch operation.
type BatchResult struct {
	Path    string
	Secrets map[string]string
	Err     error
}

// BatchClient reads multiple paths concurrently and returns aggregated results.
type BatchClient struct {
	inner       SecretReader
	concurrency int
}

// NewBatchClient creates a BatchClient wrapping the given SecretReader.
// concurrency controls the maximum number of parallel reads; if <= 0 it defaults to 4.
func NewBatchClient(inner SecretReader, concurrency int) *BatchClient {
	if inner == nil {
		panic("vault: NewBatchClient requires a non-nil inner client")
	}
	if concurrency <= 0 {
		concurrency = 4
	}
	return &BatchClient{inner: inner, concurrency: concurrency}
}

// ReadAll reads all provided paths concurrently and returns a slice of BatchResult
// in the same order as the input paths. A failed read populates BatchResult.Err
// without aborting the remaining reads.
func (b *BatchClient) ReadAll(ctx context.Context, paths []string) []BatchResult {
	results := make([]BatchResult, len(paths))
	if len(paths) == 0 {
		return results
	}

	type work struct {
		index int
		path  string
	}

	ch := make(chan work, len(paths))
	for i, p := range paths {
		ch <- work{index: i, path: p}
	}
	close(ch)

	var wg sync.WaitGroup
	for i := 0; i < b.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range ch {
				secrets, err := b.inner.ReadSecrets(ctx, w.path)
				results[w.index] = BatchResult{
					Path:    w.path,
					Secrets: secrets,
					Err:     err,
				}
			}
		}()
	}
	wg.Wait()
	return results
}

// MergeResults merges all successful BatchResults into a single map.
// Keys are namespaced as "path:key". Returns an error if any individual read failed.
func MergeResults(results []BatchResult) (map[string]string, error) {
	merged := make(map[string]string)
	for _, r := range results {
		if r.Err != nil {
			return nil, fmt.Errorf("batch read failed for path %q: %w", r.Path, r.Err)
		}
		for k, v := range r.Secrets {
			merged[r.Path+":"+k] = v
		}
	}
	return merged, nil
}
