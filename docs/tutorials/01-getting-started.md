# Getting Started with go-dci

This tutorial walks you through the basics of using the go-dci library to interact with the Red Hat Distributed CI (DCI) API.

## Prerequisites

- Go 1.21 or later
- DCI RemoteCI credentials (access key and secret key)
- Access to the DCI platform

## Step 1: Install the Library

Add go-dci to your Go project:

```bash
go get github.com/sebrandon1/go-dci
```

## Step 2: Obtain DCI Credentials

DCI uses AWS Signature Version 4 for authentication. You need:

1. **Access Key** - Your DCI client ID
2. **Secret Key** - Your DCI API secret

To get these credentials:

1. Log in to the [DCI Dashboard](https://www.distributed-ci.io/)
2. Navigate to your RemoteCI configuration
3. Download or copy your credentials

## Step 3: Configure Credentials

You can provide credentials in two ways:

### Option A: Environment Variables (Recommended)

```bash
export GO_DCI_ACCESSKEY="your-access-key"
export GO_DCI_SECRETKEY="your-secret-key"
```

### Option B: Config File

Create a `.go-dci-config.yaml` file:

```yaml
dci_client_id: "your-access-key"
dci_api_secret: "your-secret-key"
```

## Step 4: Initialize the Client

Create a new Go file and initialize the client:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    // Get credentials from environment
    accessKey := os.Getenv("GO_DCI_ACCESSKEY")
    secretKey := os.Getenv("GO_DCI_SECRETKEY")

    if accessKey == "" || secretKey == "" {
        log.Fatal("Credentials not set")
    }

    // Initialize the client
    client := lib.NewClient(accessKey, secretKey)

    fmt.Println("Client initialized!")
}
```

Run your program:

```bash
go run main.go
```

## Step 5: Verify Authentication

Test your credentials by getting your identity:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    accessKey := os.Getenv("GO_DCI_ACCESSKEY")
    secretKey := os.Getenv("GO_DCI_SECRETKEY")

    client := lib.NewClient(accessKey, secretKey)

    // Verify authentication
    identity, err := client.GetIdentity()
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    fmt.Printf("Authenticated as: %s\n", identity.Identity.Name)
    fmt.Printf("ID: %s\n", identity.Identity.ID)
    fmt.Printf("Type: %s\n", identity.Identity.Type)
}
```

## Step 6: List Topics

Topics represent test suites or certification programs:

```go
// Get all topics
topicsResp, err := client.GetTopics()
if err != nil {
    log.Fatalf("Failed to get topics: %v", err)
}

for _, resp := range topicsResp {
    for _, topic := range resp.Topics {
        fmt.Printf("Topic: %s\n", topic.Name)
        fmt.Printf("  ID: %s\n", topic.ID)
        fmt.Printf("  State: %s\n", topic.State)
    }
}
```

## Step 7: Explore Components

Components represent test artifacts (e.g., OCP versions, test tools):

```go
// Get all components
componentsResp, err := client.GetComponents()
if err != nil {
    log.Fatalf("Failed to get components: %v", err)
}

for _, resp := range componentsResp {
    for _, comp := range resp.Components {
        fmt.Printf("Component: %s\n", comp.Name)
        fmt.Printf("  Type: %s\n", comp.Type)
        fmt.Printf("  Version: %s\n", comp.Version)
    }
}

// Get components for a specific topic
topicComponents, err := client.GetTopicComponents("topic-uuid-here")
```

## Understanding DCI Concepts

### Topics
Topics are certification programs or test suites (e.g., "OCP-4.14", "OCP-4.15").

### Components
Components are the artifacts being tested:
- **ocp** - OpenShift Container Platform version
- **certsuite** - Certification test suite version

### Jobs
Jobs represent test runs. Each job:
- Belongs to a topic
- Uses specific components
- Has a lifecycle (new → pre-run → running → success/failure)

### JobStates
Track the progress of a job through its lifecycle.

### Files
Attach test results and logs to jobs.

## Error Handling

```go
identity, err := client.GetIdentity()
if err != nil {
    // Check error type
    if strings.Contains(err.Error(), "401") {
        fmt.Println("Invalid credentials")
    } else if strings.Contains(err.Error(), "403") {
        fmt.Println("Access denied")
    } else {
        fmt.Printf("Unexpected error: %v\n", err)
    }
    return
}
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GO_DCI_ACCESSKEY` | Your DCI client ID / access key |
| `GO_DCI_SECRETKEY` | Your DCI API secret key |

## Complete Example

See the [basic-usage example](../../examples/basic-usage/main.go) for a complete working program.

## Next Steps

- [Certification Workflow](./02-certification-workflow.md) - Create and manage certification jobs
- [Component Analysis](./03-component-analysis.md) - Query and analyze components
