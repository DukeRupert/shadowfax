# Shadowfax Roadmap

## v1.0.0 — Released ✅

| Command | Description | Status |
|---|---|---|
| `ping` | Verify credentials and connectivity | ✅ |
| `dns create` | Create a DNS record | ✅ |
| `dns list` | List all records for a domain | ✅ |
| `dns delete` | Delete a record by ID | ✅ |

---

## v1.1.0 — DNS Completeness

| Command | Description | Status |
|---|---|---|
| `dns edit` | Edit an existing record by ID | ✅ |
| `dns edit-by-type` | Edit records by domain and type | ✅ |
| `dns delete-by-type` | Delete records by domain and type | ✅ |
| `dns list-by-type` | List records filtered by type | ✅ |

---

## v1.2.0 — Domain Info

| Command | Description | Status |
|---|---|---|
| `domain list` | List all domains in your account | ✅ |
| `domain pricing` | Get pricing for a TLD | ✅ |

---

## v1.3.0 — SSL

| Command | Description | Status |
|---|---|---|
| `ssl retrieve` | Retrieve SSL certificate bundle for a domain | ✅ |

---

## v2.0.0 — Quality of Life

| Feature | Description | Status |
|---|---|---|
| Shell autocompletion | Bash/Zsh/Fish completion via Cobra | ✅ |
| `--output json` | Machine-readable JSON output for scripting | ✅ |
| `--quiet` | Suppress output except errors | ✅ |
| Config file support | `~/.config/shadowfax/config.yaml` via Viper | ✅ |