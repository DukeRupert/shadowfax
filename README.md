# shadowfax

> *"Shadowfax, lord of all horses."* — Gandalf

A CLI tool for managing [Porkbun](https://porkbun.com) DNS records, domains, and SSL certificates via their JSON API. Part of the [arnor](https://github.com/fireflysoftware) infrastructure management suite.

## Installation

```bash
go install github.com/fireflysoftware/shadowfax@latest
```

Or build from source:

```bash
git clone https://github.com/fireflysoftware/shadowfax
cd shadowfax
go build -o shadowfax .
```

## Configuration

Shadowfax loads credentials from multiple sources (highest priority first):

1. Environment variables (`PORKBUN_API_KEY` / `PORKBUN_SECRET_KEY`)
2. Config file (`~/.config/shadowfax/config.yaml`)
3. Dotfiles (`~/.dotfiles/.env`) or local `.env`

```bash
cp .env.example ~/.dotfiles/.env
# then fill in your keys
```

Get your API keys from the [Porkbun API portal](https://porkbun.com/account/api).

```env
PORKBUN_API_KEY=pk1_your_api_key_here
PORKBUN_SECRET_KEY=sk1_your_secret_key_here
```

Or use a config file at `~/.config/shadowfax/config.yaml`:

```yaml
api_key: pk1_your_api_key_here
secret_key: sk1_your_secret_key_here
```

## Usage

### Verify credentials

```bash
shadowfax ping
```

### DNS records

```bash
# Create an A record (root domain)
shadowfax dns create --domain example.com --type A --content 1.2.3.4

# Create a CNAME for www with custom TTL
shadowfax dns create --domain example.com --type CNAME --name www --content example.com --ttl 300

# List all records
shadowfax dns list --domain example.com

# List records filtered by type
shadowfax dns list-by-type --domain example.com --type A

# Edit a record by ID
shadowfax dns edit --domain example.com --id 123456 --type A --content 5.6.7.8

# Edit records by type
shadowfax dns edit-by-type --domain example.com --type A --content 5.6.7.8

# Delete a record by ID
shadowfax dns delete --domain example.com --id 123456

# Delete records by type
shadowfax dns delete-by-type --domain example.com --type A
```

Default TTL is `600` seconds.

### Domains

```bash
# List all domains in your account
shadowfax domain list

# Get pricing for a specific TLD
shadowfax domain pricing --tld com

# Get pricing for all TLDs
shadowfax domain pricing
```

### SSL certificates

```bash
# Print certificate bundle to stdout
shadowfax ssl retrieve --domain example.com

# Save certificate files to a directory
shadowfax ssl retrieve --domain example.com --output /path/to/certs
```

### Global flags

```bash
# JSON output for scripting
shadowfax dns list --domain example.com --output json

# Quiet mode (errors only)
shadowfax dns create --domain example.com --type A --content 1.2.3.4 --quiet
```

### Shell completion

```bash
# Bash
source <(shadowfax completion bash)

# Zsh
source <(shadowfax completion zsh)

# Fish
shadowfax completion fish | source
```

## License

MIT
