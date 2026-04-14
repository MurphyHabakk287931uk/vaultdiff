package vault

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

// CounterMetrics is a simple in-process Metrics implementation that counts
// reads per path and tracks total errors. It is safe for concurrent use.
type CounterMetrics struct {
	mu      sync.Mutex
	counts  map[string]int64
	errors  map[string]int64
	totalMs int64
}

// NewCounterMetrics returns an initialised CounterMetrics.
func NewCounterMetrics() *CounterMetrics {
	return &CounterMetrics{
		counts: make(map[string]int64),
		errors: make(map[string]int64),
	}
}

// RecordRead implements Metrics.
func (c *CounterMetrics) RecordRead(path string, durationMs int64, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[path]++
	c.totalMs += durationMs
	if err != nil {
		c.errors[path]++
	}
}

// TotalCalls returns the total number of ReadSecrets calls recorded.
func (c *CounterMetrics) TotalCalls() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	var total int64
	for _, v := range c.counts {
		total += v
	}
	return total
}

// TotalErrors returns the total number of failed ReadSecrets calls.
func (c *CounterMetrics) TotalErrors() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	var total int64
	for _, v := range c.errors {
		total += v
	}
	return total
}

// WriteSummary writes a human-readable summary to w (defaults to os.Stderr).
func (c *CounterMetrics) WriteSummary(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	paths := make([]string, 0, len(c.counts))
	for p := range c.counts {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	fmt.Fprintf(w, "vault metrics summary (total_ms=%d):\n", c.totalMs)
	for _, p := range paths {
		fmt.Fprintf(w, "  %-40s calls=%-4d errors=%d\n", p, c.counts[p], c.errors[p])
	}
}
