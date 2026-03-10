# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Golang-based CLI and library wrapper around the Red Hat Distributed CI (DCI) API. Provides commands for interacting with DCI topics, jobs, components, and more. Can be used as a standalone CLI tool or imported as a Go library.

## Common Commands

### Build
```bash
make build
```

### Run
```bash
./go-dci [command]
```

### Test
```bash
make test
```

### Lint and Vet
```bash
make lint
make vet
```

### Clean
```bash
make clean
```

### Check API Alignment
```bash
make check-swagger-alignment
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
- **`lib/`** - DCI API client library
- **`docs/`** - Documentation and tutorials
- **`examples/`** - Runnable example programs (basic-usage, certification-workflow, component-query)
- **`scripts/`** - Utility scripts (swagger alignment checker)
- **`main.go`** - Application entry point

## Supported DCI API Endpoints

| Endpoint | CLI Command |
|----------|-------------|
| `/api/v1/topics` | `topics` |
| `/api/v1/jobs` | `jobs`, `ocpcount`, `create-job` |
| `/api/v1/components` | `components` |
| `/api/v1/identity` | `identity` |
| `/api/v1/componenttypes` | `componenttypes` |
| `/api/v1/jobstates` | `update-job-state` |
| `/api/v1/files` | `upload-file` |
| `/api/v1/products` | `products` |
| `/api/v1/remotecis` | `remotecis` |
| `/api/v1/teams` | `teams` |
| `/api/v1/users` | `users` |

## Configuration

Configure credentials via CLI:
```bash
./go-dci config set --accesskey <key> --secretkey <secret>
```

Or via environment variables:
```bash
export GO_DCI_ACCESSKEY="your-access-key"
export GO_DCI_SECRETKEY="your-secret-key"
```

## Requirements

- Go 1.26+
- Valid DCI API credentials

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting
