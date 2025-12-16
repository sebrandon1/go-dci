# go-dci

![Test Incoming Changes](https://github.com/sebrandon1/go-dci/actions/workflows/pre-main.yml/badge.svg)
![Nightly](https://github.com/sebrandon1/go-dci/actions/workflows/nightly.yml/badge.svg)

## Overview

A Golang based wrapper around the Red Hat Distributed CI API:

https://doc.distributed-ci.io/dci-control-server/docs/API/

## Supported DCI API Endpoints

| Endpoint | Method | Status | CLI Command |
|----------|--------|--------|-------------|
| `/api/v1/topics` | GET | ✅ Implemented | `topics` (via lib) |
| `/api/v1/jobs` | GET | ✅ Implemented | `jobs`, `ocpcount` |
| `/api/v1/components` | GET | ✅ Implemented | `components` |
| `/api/v1/identity` | GET | ❌ Not Implemented | - |
| `/api/v1/componenttypes` | GET | ❌ Not Implemented | - |
| `/api/v1/jobs` | POST | ❌ Not Implemented | - |
| `/api/v1/jobstates` | POST | ❌ Not Implemented | - |
| `/api/v1/files` | POST | ❌ Not Implemented | - |

## CLI Usage

### Build

```bash
make build
```

### Configuration

Set up your DCI RemoteCI credentials:

```bash
./go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>
```

```
Usage:
  dci config set [flags]

Flags:
  -a, --accesskey string   The access key to set in the configuration.
  -h, --help               help for set
  -s, --secretkey string   The secret key to set in the configuration.
```

### Available Commands

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
