// Package config provides loading and validation of vaultdiff configuration
// files. Configuration is expressed as YAML and merged over built-in defaults,
// allowing users to override only the fields they care about.
//
// A minimal config file looks like:
//
//	vault:
//	  address: https://vault.example.com
//	  token: s.mytoken
//	diff:
//	  redact_mode: redact
//	output:
//	  format: json
//
// If no config file is provided, DefaultConfig() values are used.
package config
