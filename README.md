# vaultdiff

A CLI tool to diff secrets between two Vault paths or environments with redacted output support.

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

```bash
# Diff secrets between two Vault paths
vaultdiff secret/prod/app secret/staging/app

# Diff with redacted values (shows keys that differ, hides values)
vaultdiff --redact secret/prod/app secret/staging/app

# Diff across environments using different Vault addresses
vaultdiff --addr-a https://vault-prod:8200 --addr-b https://vault-staging:8200 \
  secret/app secret/app
```

### Example Output

```
~ DB_PASSWORD   [redacted]
+ NEW_FEATURE_FLAG   enabled
- DEPRECATED_KEY     old-value
```

### Flags

| Flag | Description |
|------|-------------|
| `--redact` | Hide secret values in output |
| `--addr-a` | Vault address for the first path |
| `--addr-b` | Vault address for the second path |
| `--token` | Vault token (defaults to `VAULT_TOKEN` env var) |
| `--format` | Output format: `text`, `json` (default: `text`) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with a valid token and read access to compared paths

---

## License

MIT © [yourusername](https://github.com/yourusername)