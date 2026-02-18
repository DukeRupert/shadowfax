# shadowfax

> *"Shadowfax, lord of all horses."* — Gandalf

A small CLI tool for managing [Porkbun](https://porkbun.com) DNS records via their JSON API.

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

Shadowfax looks for credentials in `~/.dotfiles/.env`, falling back to a `.env` in the current directory.

```bash
cp .env.example ~/.dotfiles/.env
# then fill in your keys
```

Get your API keys from the [Porkbun API portal](https://porkbun.com/account/api).

```env
PORKBUN_API_KEY=pk1_your_api_key_here
PORKBUN_SECRET_KEY=sk1_your_secret_key_here
```

## Usage

### Create a record

```bash
# A record (root)
shadowfax dns create --domain example.com --type A --content 1.2.3.4

# CNAME for www
shadowfax dns create --domain example.com --type CNAME --name www --content example.com

# Custom TTL
shadowfax dns create --domain example.com --type A --content 1.2.3.4 --ttl 300
```

Default TTL is `600` seconds.

### List records

```bash
shadowfax dns list --domain example.com
```

### Delete a record

```bash
# Get the record ID from dns list first
shadowfax dns delete --domain example.com --id 123456
```

## License

MIT
