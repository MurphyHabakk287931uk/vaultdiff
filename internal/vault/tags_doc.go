// Package vault — tags client
//
// # TagsClient
//
// NewTagsClient wraps any SecretReader and injects a fixed set of
// string key/value metadata tags into every successful response.
//
// Tags are written into the returned secrets map using a configurable
// prefix (default "_tag.") so they are easy to identify and strip
// downstream without conflicting with real secret keys.
//
// # Usage
//
//	client := vault.NewTagsClient(
//		baseClient,
//		map[string]string{
//			"env":    "production",
//			"region": "us-east-1",
//		},
//		"", // use default prefix "_tag."
//	)
//
// # Behaviour
//
//   - Tags are copied at construction time; later mutations to the
//     source map have no effect.
//   - On error the inner client's error is returned unchanged.
//   - Existing secret keys are never overwritten by tags.
package vault
