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
    "os"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    // Get credentials from environment (recommended)
    accessKey := os.Getenv("GO_DCI_ACCESSKEY")
    secretKey := os.Getenv("GO_DCI_SECRETKEY")

    // Create client
    client := lib.NewClient(accessKey, secretKey)

    // Use client...
}
```

## Authentication

DCI uses AWS Signature Version 4 for authentication. The library handles this automatically.

### Verify Authentication

```go
identity, err := client.GetIdentity()
if err != nil {
    log.Fatalf("Authentication failed: %v", err)
}
fmt.Printf("Authenticated as: %s\n", identity.Identity.Name)
```

## API Categories

### Identity Operations

```go
// Get current identity
identity, err := client.GetIdentity()
```

### Topic Operations

```go
// List all topics
topics, err := client.GetTopics()

// Get specific topic
topic, err := client.GetTopic(topicID)

// Create topic (admin)
topic, err := client.CreateTopic(name, productID, componentTypes)

// Update topic
topic, err := client.UpdateTopic(topicID, lib.UpdateTopicRequest{Name: "new-name"})

// Delete topic (admin)
err := client.DeleteTopic(topicID)

// Get topic components
components, err := client.GetTopicComponents(topicID)
```

### Component Type Operations

```go
// List all component types
types, err := client.GetComponentTypes()

// Get specific type
compType, err := client.GetComponentType(typeID)

// Create type (admin)
compType, err := client.CreateComponentType(name)

// Update type
compType, err := client.UpdateComponentType(typeID, lib.UpdateComponentTypeRequest{Name: "new"})

// Delete type (admin)
err := client.DeleteComponentType(typeID)
```

### Component Operations

```go
// List all components
components, err := client.GetComponents()

// Get components by topic
components, err := client.GetComponentsByTopicID(topicID)

// Get specific component
component, err := client.GetComponent(componentID)

// Create component
component, err := client.CreateComponent(name, componentType, topicID, version)

// Update component
component, err := client.UpdateComponent(componentID, lib.UpdateComponentRequest{
    Name:    "updated-name",
    Version: "1.2.3",
})

// Delete component
err := client.DeleteComponent(componentID)
```

### Job Operations

```go
// List jobs (with day filter)
jobs, err := client.GetJobs(30) // Last 30 days

// List jobs by date range
jobs, err := client.GetJobsByDate(startTime, endTime)

// Get specific job
job, err := client.GetJob(jobID)

// Create job with components
job, err := client.CreateJob(topicID, componentIDs, comment)

// Schedule job (auto-select components)
job, err := client.ScheduleJob(topicID)

// Update job
job, err := client.UpdateJob(jobID, lib.UpdateJobRequest{Comment: "updated"})

// Delete job
err := client.DeleteJob(jobID)

// Get job files
files, err := client.GetJobFiles(jobID)
```

### Job State Operations

```go
// Update job state
state, err := client.UpdateJobState(jobID, lib.JobStateRunning, "Starting tests")

// Create job state entry
state, err := client.CreateJobState(jobID, lib.JobStateSuccess, "Tests passed")

// Get job states history
states, err := client.GetJobStates(jobID)
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
upload, err := client.UploadFile(jobID, filePath, mimeType)

// Upload file content
upload, err := client.UploadFileContent(jobID, fileName, mimeType, content)

// Get file (download)
data, contentType, err := client.GetFile(fileID)

// Delete file
err := client.DeleteFile(fileID)
```

### RemoteCI Operations

```go
// List RemoteCIs
remotecis, err := client.GetRemoteCIs()

// Get specific RemoteCI
remoteci, err := client.GetRemoteCI(remoteciID)

// Create RemoteCI (admin)
remoteci, err := client.CreateRemoteCI(name, teamID)

// Update RemoteCI
remoteci, err := client.UpdateRemoteCI(remoteciID, lib.UpdateRemoteCIRequest{Name: "new"})

// Delete RemoteCI (admin)
err := client.DeleteRemoteCI(remoteciID)
```

### Team Operations

```go
// List teams
teams, err := client.GetTeams()

// Get specific team
team, err := client.GetTeam(teamID)

// Create team (admin)
team, err := client.CreateTeam(name)

// Update team
team, err := client.UpdateTeam(teamID, lib.UpdateTeamRequest{Name: "new"})

// Delete team (admin)
err := client.DeleteTeam(teamID)
```

### User Operations

```go
// List users
users, err := client.GetUsers()

// Get specific user
user, err := client.GetUser(userID)

// Create user (admin)
user, err := client.CreateUser(name, email, fullname, teamID, password)

// Update user
user, err := client.UpdateUser(userID, lib.UpdateUserRequest{Name: "new"})

// Delete user (admin)
err := client.DeleteUser(userID)
```

### Product Operations

```go
// List products
products, err := client.GetProducts()

// Get specific product
product, err := client.GetProduct(productID)
```

## Error Handling

```go
result, err := client.GetTopic(topicID)
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
