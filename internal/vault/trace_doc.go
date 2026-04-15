// Package vault — TraceClient
//
// TraceClient is a lightweight diagnostic decorator that logs every
// ReadSecrets call together with the resolved path, elapsed duration, result
// status, and the number of keys returned.
//
// It is intended for development and debugging; avoid enabling it in
// production builds where secret paths may be considered sensitive.
//
// Usage:
//
//	inner := vault.NewClient(cfg)
//	client := vault.NewTraceClient(inner, os.Stderr, "[vaultdiff] ")
//	// client now traces every ReadSecrets call to stderr
//
// Output format (one line per call):
//
//	[vaultdiff] [trace] path="secret/app" duration=1.234ms status=ok keys=5
//	[vaultdiff] [trace] path="secret/missing" duration=312µs status=error: secret not found keys=0
package vault
