// Package vault provides Vault client wrappers used by vaultdiff.
//
// # Batch Client
//
// BatchClient enables concurrent reads across multiple Vault paths.
// It is useful when diffing secrets spread over several paths so that
// network round-trips happen in parallel rather than sequentially.
//
// Basic usage:
//
//	inner, _ := vault.NewClient(vault.ClientConfig{Token: "root"})
//	batch := vault.NewBatchClient(inner, 4) // up to 4 concurrent reads
//
//	results := batch.ReadAll(ctx, []string{
//		"secret/app/prod",
//		"secret/app/staging",
//	})
//
//	for _, r := range results {
//		if r.Err != nil {
//			log.Printf("error reading %s: %v", r.Path, r.Err)
//			continue
//		}
//		fmt.Println(r.Path, r.Secrets)
//	}
//
// MergeResults can be used to flatten all results into a single map
// keyed by "path:secretKey", returning an error if any read failed.
package vault
