// Package config provides configuration loading for vaultdiff.
//
// Configuration is resolved in the following order of precedence
// (highest to lowest):
//
//  1. Environment variables (VAULT_ADDR, VAULTDIFF_AUTH_METHOD,
//     VAULTDIFF_REDACT, VAULTDIFF_OUTPUT)
//  2. YAML configuration file (path supplied by the caller)
//  3. Built-in defaults (see DefaultConfig)
//
// Supported auth methods: token, approle, kubernetes.
// Supported redact modes: none, redact, mask.
// Supported output formats: text, json.
package config
