// Package vault — observe.go / metrics.go
//
// ObserveClient adds non-intrusive observability to any SecretReader by
// wrapping it with timing and error recording via the Metrics interface.
//
// Usage:
//
//	inner := vault.NewClient(cfg)
//	m    := vault.NewCounterMetrics()
//	client := vault.NewObserveClient(inner, m)
//
//	// ... use client normally ...
//
//	m.WriteSummary(os.Stderr)
//
// The Metrics interface is intentionally minimal so that callers can plug in
// Prometheus counters, Datadog StatsD, or any other backend without coupling
// the vault package to a specific metrics library.
//
// CounterMetrics provides an in-process implementation suitable for CLI tools
// and integration tests where a full metrics pipeline is not available.
package vault
