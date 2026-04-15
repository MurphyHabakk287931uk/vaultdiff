// Package vault — hedge.go
//
// HedgeClient implements the "hedged request" pattern: a second, redundant
// request is issued after a configurable delay if the first has not yet
// responded.  Whichever response arrives first is returned to the caller.
//
// This trades a small increase in backend load for a significant reduction in
// tail latency — particularly useful when reading secrets over a high-latency
// or occasionally-slow Vault cluster.
//
// Usage:
//
//	client := vault.NewHedgeClient(
//		baseClient,
//		vault.HedgeConfig{Delay: 150 * time.Millisecond},
//	)
//
// If Delay is zero or negative, NewHedgeClient returns the inner client
// unwrapped so there is no overhead in the hot path.
package vault
