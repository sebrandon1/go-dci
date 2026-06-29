# go-dci

[![Pre-Main Checks](https://github.com/sebrandon1/go-dci/actions/workflows/pre-main.yml/badge.svg)](https://github.com/sebrandon1/go-dci/actions/workflows/pre-main.yml)
[![DCI Verified Nightly](https://github.com/sebrandon1/go-dci/actions/workflows/nightly.yml/badge.svg)](https://github.com/sebrandon1/go-dci/actions/workflows/nightly.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sebrandon1/go-dci)](https://golang.org/)
[![License](https://img.shields.io/github/license/sebrandon1/go-dci)](https://github.com/sebrandon1/go-dci/blob/main/LICENSE)

A Go wrapper around the [Red Hat Distributed CI (DCI) API](https://doc.distributed-ci.io/dci-control-server/docs/API/). Can be used as a standalone CLI tool or imported as a Go library.

## Quick Start

```bash
go install github.com/sebrandon1/go-dci@latest
```

Configure credentials and verify:

```bash
go-dci config set --accesskey <key> --secretkey <secret>
go-dci identity
go-dci topics
```

### Library Usage

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    client := lib.NewClient(
        os.Getenv("GO_DCI_ACCESSKEY"),
        os.Getenv("GO_DCI_SECRETKEY"),
    )

    identity, err := client.GetIdentity()
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    fmt.Printf("Authenticated as: %s\n", identity.Identity.Name)

    topics, err := client.GetTopics()
    if err != nil {
        log.Fatal(err)
    }
    for _, resp := range topics {
        for _, topic := range resp.Topics {
            fmt.Printf("Topic: %s (%s)\n", topic.Name, topic.ID)
        }
    }
}
```

## Guides

| Guide | Description |
|-------|-------------|
| [Installation](docs/installation.md) | Prebuilt binaries, container image, go install, build from source |
| [Authentication](docs/authentication.md) | Config file and environment variable setup |
| [CLI Reference](docs/cli-reference.md) | All commands with flags and example output |
| [Library Guide](docs/library-guide.md) | Complete API reference for using go-dci as a library |
| [Tutorials](docs/tutorials/) | Step-by-step guides for common workflows |
| [Troubleshooting](docs/troubleshooting/) | Common issues and solutions |
| [Examples](examples/) | Runnable example programs |

## Common Usage Examples

### Query Recent Jobs
```bash
# Get jobs from last 7 days
go-dci jobs --age 7

# Get jobs in specific date range
go-dci jobs --start-date 2026-06-01 --end-date 2026-06-15

# Get OCP version statistics
go-dci ocpcount --age 30
```

### Work with Components
```bash
# List all components
go-dci components

# Filter by topic
go-dci components --topic <topic-id>

# Get specific component
go-dci component --id <component-id>

# Create a new component
go-dci create-component --name "My Component" --type ocp --topic-id <topic-id> --version 4.15.0
```

### Manage Jobs
```bash
# Get specific job
go-dci job --id <job-id>

# Create a job
go-dci create-job --topic-id <topic-id> --remoteci-id <remoteci-id>

# Update job state
go-dci update-job-state --job-id <job-id> --state success --comment "All tests passed"

# Upload a file to a job
go-dci upload-file --job-id <job-id> --file results.xml --mime-type application/xml
```

### File Operations
```bash
# List files for a job
go-dci job-files --job-id <job-id>

# Download a specific file
go-dci file --id <file-id> --output downloaded-file.xml

# Delete a file
go-dci delete-file --id <file-id>
```

## Supported DCI API Endpoints

| Endpoint | CLI Commands |
|----------|--------------|
| `/api/v1/identity` | `identity` |
| `/api/v1/topics` | `topics`, `topic`, `create-topic`, `update-topic`, `delete-topic`, `topic-components` |
| `/api/v1/jobs` | `jobs`, `job`, `ocpcount`, `create-job`, `update-job`, `delete-job`, `schedule-job`, `job-files` |
| `/api/v1/components` | `components`, `component`, `create-component`, `update-component`, `delete-component` |
| `/api/v1/componenttypes` | `componenttypes`, `componenttype`, `create-componenttype`, `update-componenttype`, `delete-componenttype` |
| `/api/v1/jobstates` | `jobstates`, `update-job-state` |
| `/api/v1/files` | `file`, `delete-file`, `upload-file`, `job-files` |
| `/api/v1/remotecis` | `remotecis`, `remoteci`, `create-remoteci`, `update-remoteci`, `delete-remoteci` |
| `/api/v1/teams` | `teams`, `team`, `create-team`, `update-team`, `delete-team` |
| `/api/v1/users` | `users`, `user`, `create-user`, `update-user`, `delete-user` |
| `/api/v1/products` | `products`, `product` |

## Development

```bash
make build    # Build binary
make test     # Run tests
make lint     # Run linters
make vet      # Run go vet
make clean    # Remove binary
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
