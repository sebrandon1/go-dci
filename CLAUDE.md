# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Golang-based CLI wrapper around the Red Hat Distributed CI (DCI) API. Provides commands for interacting with DCI topics, jobs, components, and more.

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
go test ./...
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
- **`lib/`** - DCI API client library
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

## Configuration

Create a config file at `.go-dci-config.yaml`:
```yaml
dci_client_id: "your-client-id"
dci_api_secret: "your-api-secret"
```

## Requirements

- Go 1.21+
- Valid DCI API credentials

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting
