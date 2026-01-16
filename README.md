# go-dci

![Test Incoming Changes](https://github.com/sebrandon1/go-dci/actions/workflows/pre-main.yml/badge.svg)
![Nightly](https://github.com/sebrandon1/go-dci/actions/workflows/nightly.yml/badge.svg)

## Overview

A Golang based wrapper around the Red Hat Distributed CI API:

https://doc.distributed-ci.io/dci-control-server/docs/API/

## Quick Start

### As a CLI Tool

```bash
# Build the CLI
make build

# Configure credentials
./go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>

# Verify authentication
./go-dci identity

# List topics
./go-dci topics

# Get recent jobs
./go-dci jobs -d 30
```

### As a Go Library

```bash
go get github.com/sebrandon1/go-dci
```

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    // Initialize client with AWS SigV4 authentication
    client := lib.NewClient(
        os.Getenv("GO_DCI_ACCESSKEY"),
        os.Getenv("GO_DCI_SECRETKEY"),
    )

    // Verify authentication
    identity, err := client.GetIdentity()
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }
    fmt.Printf("Authenticated as: %s\n", identity.Identity.Name)

    // List topics
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

## Documentation

### Tutorials

Step-by-step guides for common workflows:

- [Getting Started](docs/tutorials/01-getting-started.md) - Installation, authentication, first API calls
- [Certification Workflow](docs/tutorials/02-certification-workflow.md) - Complete job lifecycle
- [Component Analysis](docs/tutorials/03-component-analysis.md) - Query and analyze components

### Examples

Runnable example programs in the [examples/](examples/) directory:

- [basic-usage](examples/basic-usage/) - Authentication, topics, components, jobs
- [certification-workflow](examples/certification-workflow/) - Complete certification job lifecycle
- [component-query](examples/component-query/) - Component filtering and version analysis

### Library Reference

- [Library Guide](docs/library-guide.md) - Complete API reference for Go library usage

## Supported DCI API Endpoints

| Endpoint | Method | Status | CLI Command |
|----------|--------|--------|-------------|
| `/api/v1/topics` | GET | ✅ Implemented | `topics` |
| `/api/v1/jobs` | GET | ✅ Implemented | `jobs`, `ocpcount` |
| `/api/v1/components` | GET | ✅ Implemented | `components` |
| `/api/v1/identity` | GET | ✅ Implemented | `identity` |
| `/api/v1/componenttypes` | GET | ✅ Implemented | `componenttypes` |
| `/api/v1/jobs` | POST | ✅ Implemented | `create-job` |
| `/api/v1/jobstates` | POST | ✅ Implemented | `update-job-state` |
| `/api/v1/files` | POST | ✅ Implemented | `upload-file` |

## CLI Usage

### Build

```bash
make build
```

### Configuration

Set up your DCI RemoteCI credentials using one of the following methods:

#### Option 1: Config File

```bash
./go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>
```

This creates a `.go-dci-config.yaml` file in the current directory.

```
Usage:
  dci config set [flags]

Flags:
  -a, --accesskey string   The access key to set in the configuration.
  -h, --help               help for set
  -s, --secretkey string   The secret key to set in the configuration.
```

#### Option 2: Environment Variables

You can also set credentials via environment variables (useful for CI/CD):

```bash
export GO_DCI_ACCESSKEY=<your-access-key>
export GO_DCI_SECRETKEY=<your-secret-key>
```

| Variable | Description |
|----------|-------------|
| `GO_DCI_ACCESSKEY` | Your DCI client ID / access key |
| `GO_DCI_SECRETKEY` | Your DCI API secret key |

Environment variables take precedence over values in the config file.

### Available Commands

#### `topics` - List Topics

Get all available topics from DCI.

```bash
# Get all topics
./go-dci topics

# Output as JSON
./go-dci topics --output json
```

```
Usage:
  dci topics [flags]

Flags:
  -h, --help            help for topics
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Getting all topics from DCI
ID: topic-123 | Name: OCP-4.14 | Product: OpenShift Container Platform | State: active
ID: topic-456 | Name: OCP-4.15 | Product: OpenShift Container Platform | State: active
Total Topics: 2
```

#### `identity` - Verify Authentication

Verify your DCI credentials are configured correctly and display identity information.

```bash
# Check authentication
./go-dci identity

# Output as JSON
./go-dci identity --output json
```

```
Usage:
  dci identity [flags]

Flags:
  -h, --help            help for identity
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Authentication successful!
---
ID:       abc123-def456-ghi789
Name:     my-remoteci
Type:     remoteci
Team:     My Team
Team ID:  team-123
State:    active
```

#### `jobs` - Query Jobs

Get all jobs with a specific age in days. Filters for jobs running the [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite).

```bash
# Get jobs from the last 30 days
./go-dci jobs -d 30

# Output as JSON
./go-dci jobs -d 30 --output json
```

```
Usage:
  dci jobs [flags]

Flags:
  -d, --age string      Age in days
  -h, --help            help for jobs
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Getting all jobs from DCI that are 30 days old
Job ID: 78ed13e1-841f-4c04-a1c6-8df9028c67cd  -  Certsuite Version: tnf-v5.1.3 (Days Since: 11.249845)
Job ID: eb491abd-ec8b-42cc-aa8b-98adf741b236  -  Certsuite Version: tnf-v5.1.3 (Days Since: 11.305408)
```

#### `ocpcount` - OCP Version Statistics

Get the count of certsuite jobs for each OCP version.

```bash
# Get OCP version counts for the last 30 days
./go-dci ocpcount -d 30

# Output as JSON
./go-dci ocpcount -d 30 --output json
```

```
Usage:
  dci ocpcount [flags]

Flags:
  -d, --age string      Age in days
  -h, --help            help for ocpcount
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

#### `componenttypes` - List Component Types

Get all available component types from DCI.

```bash
# Get all component types
./go-dci componenttypes

# Output as JSON
./go-dci componenttypes --output json
```

```
Usage:
  dci componenttypes [flags]

Flags:
  -h, --help            help for componenttypes
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Getting all component types from DCI
ID: ct-123 | Name: ocp | State: active
ID: ct-456 | Name: certsuite | State: active
ID: ct-789 | Name: rhel | State: active
Total Component Types: 3
```

#### `components` - Query Components

Get all components, optionally filtered by topic ID.

```bash
# Get all components
./go-dci components

# Get components for a specific topic
./go-dci components --topic <topic-id>

# Output as JSON
./go-dci components --output json
```

```
Usage:
  dci components [flags]

Flags:
  -h, --help            help for components
  -o, --output string   Output format (json) - default is stdout (default "stdout")
  -t, --topic string    Filter components by topic ID
```

Example output:

```
Getting all components from DCI
ID: abc123 | Name: OpenShift 4.14.1 | Type: ocp | Version: 4.14.1 | TopicID: topic-456
ID: def456 | Name: certsuite v5.1.3 | Type: certsuite | Version: v5.1.3 | TopicID: topic-456
Total Components: 2
```

#### `create-job` - Create a New Job

Create a new job in DCI for a given topic.

```bash
# Create a job for a topic
./go-dci create-job --topic-id <topic-id>

# Create a job with specific components
./go-dci create-job --topic-id <topic-id> --components <comp-id-1>,<comp-id-2>

# Create a job with a comment
./go-dci create-job --topic-id <topic-id> --comment "Test run for certification"

# Output as JSON
./go-dci create-job --topic-id <topic-id> --output json
```

```
Usage:
  dci create-job [flags]

Flags:
      --comment string      Optional comment for the job
      --components string   Comma-separated list of component IDs
  -h, --help                help for create-job
  -o, --output string       Output format (json) - default is stdout (default "stdout")
      --topic-id string     Topic ID for the job (required)
```

Example output:

```
Creating job for topic ID: topic-123
Job created successfully!
---
Job ID:    job-456
Topic ID:  topic-123
Status:    new
State:     active
Created:   2024-01-01T00:00:00.000000
```

#### `update-job-state` - Update Job State

Update the state of a job (pre-run, running, success, failure, etc.).

```bash
# Update job to running state
./go-dci update-job-state --job-id <job-id> --status running

# Update job to success with comment
./go-dci update-job-state --job-id <job-id> --status success --comment "All tests passed"

# Update job to failure
./go-dci update-job-state --job-id <job-id> --status failure --comment "Tests failed"

# Output as JSON
./go-dci update-job-state --job-id <job-id> --status success --output json
```

```
Usage:
  dci update-job-state [flags]

Flags:
      --comment string   Optional comment for the state change
  -h, --help             help for update-job-state
      --job-id string    Job ID to update (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --status string    New status (pre-run, running, success, failure, etc.) (required)
```

Valid status values: `new`, `pre-run`, `running`, `post-run`, `success`, `failure`, `killed`, `error`

Example output:

```
Updating job job-123 to status: running
Job state updated successfully!
---
JobState ID: jobstate-789
Job ID:      job-123
Status:      running
Created:     2024-01-01T00:00:00.000000
```

#### `upload-file` - Upload File to Job

Upload a file (e.g., test results) to a job in DCI.

```bash
# Upload a JUnit test results file
./go-dci upload-file --job-id <job-id> --file /path/to/results.xml

# Upload with custom MIME type
./go-dci upload-file --job-id <job-id> --file /path/to/results.json --mime application/json

# Output as JSON
./go-dci upload-file --job-id <job-id> --file /path/to/results.xml --output json
```

```
Usage:
  dci upload-file [flags]

Flags:
      --file string     Path to the file to upload (required)
  -h, --help            help for upload-file
      --job-id string   Job ID to attach the file to (required)
      --mime string     MIME type of the file (default "application/junit")
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Uploading file results.xml to job job-123
File uploaded successfully!
---
File ID:   file-456
Job ID:    job-123
Name:      results.xml
MIME Type: application/junit
Size:      1024 bytes
Created:   2024-01-01T00:00:00.000000
```

## Development

### Run Tests

```bash
make test
```

### Run Linter

```bash
make lint
```

### Build

```bash
make build
```
