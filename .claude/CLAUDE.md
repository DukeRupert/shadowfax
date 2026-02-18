# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Shadowfax is a Go CLI tool for managing DNS records via the Porkbun API. Built with Cobra, it supports creating, listing, and deleting DNS records and verifying API credentials.

## Build & Run Commands

```bash
go build -o shadowfax .        # Build binary
go install .                    # Install locally
go test ./...                   # Run tests (none exist yet)
go mod download                 # Fetch dependencies
```

## Architecture

Two-layer design separating CLI from API:

- **`main.go`** — Entry point, calls `cmd.Execute()`
- **`cmd/root.go`** — Root Cobra command; loads credentials from `~/.dotfiles/.env` (preferred) or `./.env` (fallback) via `PersistentPreRunE`; initializes a global `client`
- **`cmd/dns.go`** — Four commands: `ping`, `dns create`, `dns list`, `dns delete`. Flags defined in `init()`, commands wired hierarchically
- **`internal/porkbun/client.go`** — HTTP client wrapping Porkbun's JSON POST API (`https://api.porkbun.com/api/json/v3`). Auth credentials embedded in every request body. All responses checked for `status == "SUCCESS"`

## Key Conventions

- Errors wrapped with `fmt.Errorf(...%w, err)`, lowercase messages (idiomatic Go)
- Success output uses `✓` prefix; tables formatted with `text/tabwriter`
- Empty `name` field = root domain; non-empty = subdomain
- Default TTL: 600 seconds
- Dependencies: `cobra`, `godotenv`, `pflag`

## Configuration

Requires `PORKBUN_API_KEY` and `PORKBUN_SECRET_KEY` env vars. Copy `.env.example` to `~/.dotfiles/.env` or `./.env`.
