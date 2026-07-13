# go-dci Library Guide

This guide covers how to use go-dci as a Go library in your applications.

## Installation

```bash
go get github.com/sebrandon1/go-dci
```

## Client Initialization

```go
package main

import (
    "context"
    "os"
    "time"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    // Get credentials from environment (recommended)
    accessKey := os.Getenv("GO_DCI_ACCESSKEY")
    secretKey := os.Getenv("GO_DCI_SECRETKEY")

    // Create client with default timeouts
    client := lib.NewClient(accessKey, secretKey)

    // Optional: Customize timeouts for specific network conditions
    client.RequestTimeout = 60 * time.Second  // Overall request timeout (default: 30s)
    client.TLSTimeout = 10 * time.Second      // TLS handshake timeout (default: 5s)
    client.DialTimeout = 15 * time.Second     // Connection dial timeout (default: 10s)
    client.MaxRetries = 5                     // Max retry attempts (default: 3)

    // Use client...
}
```

### Default Timeouts

The client is configured with sensible default timeouts to prevent requests from hanging indefinitely:

- **Request Timeout**: 30 seconds - Maximum time for the entire request/response cycle
- **TLS Handshake Timeout**: 5 seconds - Maximum time to establish TLS connection
- **Dial Timeout**: 10 seconds - Maximum time to establish TCP connection
- **Max Retries**: 3 - Number of retry attempts for retriable errors (5xx responses on GET/DELETE)

These defaults work well for most scenarios. Adjust them if you experience:
- Frequent timeouts on slow networks (increase timeouts)
- Need faster failure detection (decrease timeouts)
- Flaky network connections (increase `MaxRetries`)

## Authentication

DCI uses AWS Signature Version 4 for authentication. The library handles this automatically.

### Verify Authentication

```go
identity, err := client.GetIdentity(context.Background())
if err != nil {
    log.Fatalf("Authentication failed: %v", err)
}
fmt.Printf("Authenticated as: %s\n", identity.Identity.Name)
```

## API Categories

### Identity Operations

```go
// Get current identity
identity, err := client.GetIdentity(context.Background())
```

### Topic Operations

```go
// List all topics
topics, err := client.GetTopics(context.Background())

// Get specific topic
topic, err := client.GetTopic(context.Background(), topicID)

// Create topic (admin)
topic, err := client.CreateTopic(context.Background(), name, productID, componentTypes)

// Update topic
topic, err := client.UpdateTopic(context.Background(), topicID, lib.UpdateTopicRequest{Name: "new-name"})

// Delete topic (admin)
err := client.DeleteTopic(context.Background(), topicID)

// Get topic components
components, err := client.GetTopicComponents(context.Background(), topicID)
```

### Component Type Operations

```go
// List all component types
types, err := client.GetComponentTypes(context.Background())

// Get specific type
compType, err := client.GetComponentType(context.Background(), typeID)

// Create type (admin)
compType, err := client.CreateComponentType(context.Background(), name)

// Update type
compType, err := client.UpdateComponentType(context.Background(), typeID, lib.UpdateComponentTypeRequest{Name: "new"})

// Delete type (admin)
err := client.DeleteComponentType(context.Background(), typeID)
```

### Component Operations

```go
// List all components
components, err := client.GetComponents(context.Background())

// Get components by topic
components, err := client.GetComponentsByTopicID(context.Background(), topicID)

// Get specific component
component, err := client.GetComponent(context.Background(), componentID)

// Create component
component, err := client.CreateComponent(context.Background(), name, componentType, topicID, version)

// Update component
component, err := client.UpdateComponent(context.Background(), componentID, lib.UpdateComponentRequest{
    Name:    "updated-name",
    Version: "1.2.3",
})

// Delete component
err := client.DeleteComponent(context.Background(), componentID)
```

### Job Operations

```go
// List jobs (with day filter)
jobs, err := client.GetJobs(context.Background(), 30) // Last 30 days

// List jobs by date range
jobs, err := client.GetJobsByDate(context.Background(), startTime, endTime)

// Get specific job
job, err := client.GetJob(context.Background(), jobID)

// Create job with components
job, err := client.CreateJob(context.Background(), topicID, componentIDs, comment)

// Schedule job (auto-select components)
job, err := client.ScheduleJob(context.Background(), topicID)

// Update job
job, err := client.UpdateJob(context.Background(), jobID, lib.UpdateJobRequest{Comment: "updated"})

// Delete job
err := client.DeleteJob(context.Background(), jobID)

// Get job files
files, err := client.GetJobFiles(context.Background(), jobID)
```

### Job State Operations

```go
// Update job state
state, err := client.UpdateJobState(context.Background(), jobID, lib.JobStateRunning, "Starting tests")

// Mark job as successful
state, err = client.UpdateJobState(context.Background(), jobID, lib.JobStateSuccess, "Tests passed")

// Get job states history
states, err := client.GetJobStates(context.Background(), jobID)
```

### Job States

```go
lib.JobStateNew       // "new"
lib.JobStatePreRun    // "pre-run"
lib.JobStateRunning   // "running"
lib.JobStatePostRun   // "post-run"
lib.JobStateSuccess   // "success"
lib.JobStateFailure   // "failure"
lib.JobStateError     // "error"
lib.JobStateKilled    // "killed"
```

### File Operations

```go
// Upload file from path
upload, err := client.UploadFile(context.Background(), jobID, filePath, mimeType)

// Upload file content
upload, err := client.UploadFileContent(context.Background(), jobID, fileName, mimeType, content)

// Get file (download)
data, contentType, err := client.GetFile(context.Background(), fileID)

// Delete file
err := client.DeleteFile(context.Background(), fileID)
```

### RemoteCI Operations

```go
// List RemoteCIs
remotecis, err := client.GetRemoteCIs(context.Background())

// Get specific RemoteCI
remoteci, err := client.GetRemoteCI(context.Background(), remoteciID)

// Create RemoteCI (admin)
remoteci, err := client.CreateRemoteCI(context.Background(), name, teamID)

// Update RemoteCI
remoteci, err := client.UpdateRemoteCI(context.Background(), remoteciID, lib.UpdateRemoteCIRequest{Name: "new"})

// Delete RemoteCI (admin)
err := client.DeleteRemoteCI(context.Background(), remoteciID)
```

### Team Operations

```go
// List teams
teams, err := client.GetTeams(context.Background())

// Get specific team
team, err := client.GetTeam(context.Background(), teamID)

// Create team (admin)
team, err := client.CreateTeam(context.Background(), name)

// Update team
team, err := client.UpdateTeam(context.Background(), teamID, lib.UpdateTeamRequest{Name: "new"})

// Delete team (admin)
err := client.DeleteTeam(context.Background(), teamID)
```

### User Operations

```go
// List users
users, err := client.GetUsers(context.Background())

// Get specific user
user, err := client.GetUser(context.Background(), userID)

// Create user (admin)
user, err := client.CreateUser(context.Background(), name, email, fullname, teamID, password)

// Update user
user, err := client.UpdateUser(context.Background(), userID, lib.UpdateUserRequest{Name: "new"})

// Delete user (admin)
err := client.DeleteUser(context.Background(), userID)
```

### Product Operations

```go
// List products
products, err := client.GetProducts(context.Background())

// Get specific product
product, err := client.GetProduct(context.Background(), productID)
```

## Error Handling

```go
result, err := client.GetTopic(context.Background(), topicID)
if err != nil {
    errStr := err.Error()

    switch {
    case strings.Contains(errStr, "401"):
        // Invalid credentials
    case strings.Contains(errStr, "403"):
        // Permission denied
    case strings.Contains(errStr, "404"):
        // Not found
    default:
        // Other error
    }
}
```

## Retry Logic

Implement retries for transient failures:

```go
func withRetry(fn func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return err
}
```

## Pagination

The library handles pagination internally for list operations. Methods that return slices (like `GetTopics()`) automatically fetch all pages.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GO_DCI_ACCESSKEY` | Your DCI client ID / access key |
| `GO_DCI_SECRETKEY` | Your DCI API secret key |

## Best Practices

1. **Environment Variables** - Store credentials in environment variables, not code
2. **Error Handling** - Always check errors and handle appropriately
3. **Pagination** - The library handles pagination automatically
4. **State Transitions** - Follow the correct job state lifecycle
5. **File Uploads** - Use appropriate MIME types for test results

## Examples

See the [examples directory](../examples/) for complete working programs:

- [basic-usage](../examples/basic-usage/) - Getting started
- [certification-workflow](../examples/certification-workflow/) - Complete job lifecycle
- [component-query](../examples/component-query/) - Component analysis

## Tutorials

Step-by-step guides in the [tutorials directory](./tutorials/):

1. [Getting Started](./tutorials/01-getting-started.md)
2. [Certification Workflow](./tutorials/02-certification-workflow.md)
3. [Component Analysis](./tutorials/03-component-analysis.md)
