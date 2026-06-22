# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultdiff/releases).

---

## Usage

Compare two versions of a secret within the same Vault path:

```bash
vaultdiff --path secret/data/myapp --v1 3 --v2 5
```

Diff secrets across environments:

```bash
vaultdiff --path secret/data/myapp \
  --env-a staging \
  --env-b production \
  --addr-a https://vault-staging:8200 \
  --addr-b https://vault-prod:8200
```

Output highlights added, removed, and changed keys while **redacting sensitive values** by default. Use `--show-values` to display plaintext diffs.

### Common Flags

| Flag | Description |
|------|-------------|
| `--path` | Vault KV path to compare |
| `--v1` | First version number |
| `--v2` | Second version number |
| `--show-values` | Display secret values in diff output |
| `--output` | Output format: `text`, `json`, `yaml` |
| `--token` | Vault token (overrides `VAULT_TOKEN`) |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `VAULT_ADDR` | Vault server address |
| `VAULT_TOKEN` | Authentication token |
| `VAULT_NAMESPACE` | Vault namespace (Enterprise) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine

---

## License

[MIT](LICENSE)