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

### Coverage
```bash
make coverage
```

### Clean
```bash
make clean
```

### Run (build + execute)
```bash
make run
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

| Endpoint | CLI Commands |
|----------|-------------|
| `/api/v1/topics` | `topics`, `topic`, `create-topic`, `update-topic`, `delete-topic`, `topic-components` |
| `/api/v1/jobs` | `jobs`, `job`, `create-job`, `update-job`, `delete-job`, `schedule-job`, `job-files`, `ocpcount` |
| `/api/v1/components` | `components`, `component`, `create-component`, `update-component`, `delete-component` |
| `/api/v1/componenttypes` | `componenttypes`, `componenttype`, `create-componenttype`, `update-componenttype`, `delete-componenttype` |
| `/api/v1/remotecis` | `remotecis`, `remoteci`, `create-remoteci`, `update-remoteci`, `delete-remoteci` |
| `/api/v1/teams` | `teams`, `team`, `create-team`, `update-team`, `delete-team` |
| `/api/v1/users` | `users`, `user`, `create-user`, `update-user`, `delete-user` |
| `/api/v1/products` | `products`, `product` |
| `/api/v1/files` | `file`, `delete-file`, `upload-file` |
| `/api/v1/jobstates` | `jobstates`, `update-job-state` |
| `/api/v1/identity` | `identity` |
| Config | `config set`, `config unset`, `config view` |

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
