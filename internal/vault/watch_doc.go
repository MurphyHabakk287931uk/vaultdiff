// Package vault – watch.go
//
// WatchClient provides a lightweight polling mechanism that monitors a Vault
// secret path for changes and notifies callers via callbacks.
//
// # Usage
//
//	stop := vault.NewWatchClient(ctx, client, "secret/myapp", vault.WatchConfig{
//		Interval: 30 * time.Second,
//		OnChange: func(path string, secrets map[string]string) {
//			log.Printf("secrets changed at %s", path)
//		},
//		OnError: func(path string, err error) {
//			log.Printf("watch error for %s: %v", path, err)
//		},
//	})
//	defer stop()
//
// # Behaviour
//
//   - Polling starts immediately on the first tick of the interval ticker.
//   - OnChange is only invoked when the secret map actually differs from the
//     previous read (deep-equality check on key/value pairs).
//   - OnError is optional; if nil, read errors are silently discarded and the
//     poller continues on the next tick.
//   - Calling the returned CancelFunc or cancelling the parent context stops
//     the background goroutine cleanly.
package vault
