// Package vault provides CompressClient and DecompressClient for transparently
// compressing and decompressing secret values using gzip and base64 encoding.
//
// # Use Case
//
// Some secrets (e.g. certificates, large JSON blobs) may approach Vault's per-value
// size constraints. Wrapping a client with CompressClient reduces stored size, while
// DecompressClient restores values on read.
//
// # Usage
//
//	writer := vault.NewCompressClient(inner)   // compress on write path
//	reader := vault.NewDecompressClient(inner) // decompress on read path
//
// # Notes
//
// CompressClient and DecompressClient must be used as a matched pair.
// Applying DecompressClient to uncompressed values will return an error.
package vault
