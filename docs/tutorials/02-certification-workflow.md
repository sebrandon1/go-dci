# Certification Workflow with go-dci

This tutorial covers the complete certification job lifecycle using the go-dci library.

## Prerequisites

- Completed [Getting Started](./01-getting-started.md)
- Access to a DCI topic for testing

## Job Lifecycle Overview

A DCI job progresses through these states:

```
new → pre-run → running → post-run → success/failure
                                   ↘ error/killed
```

| State | Description |
|-------|-------------|
| `new` | Job created, not started |
| `pre-run` | Preparing test environment |
| `running` | Tests are executing |
| `post-run` | Tests complete, collecting results |
| `success` | All tests passed |
| `failure` | Some tests failed |
| `error` | Job encountered an error |
| `killed` | Job was manually terminated |

## Step 1: Find a Topic

First, identify the topic for your certification:

```go
package main

import (
    "context"
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

    // List available topics
    ctx := context.Background()
    topicsResp, err := client.GetTopics(ctx)
    if err != nil {
        log.Fatalf("Failed to get topics: %v", err)
    }

    for _, resp := range topicsResp {
        for _, topic := range resp.Topics {
            if topic.State == "active" {
                fmt.Printf("%s - %s\n", topic.Name, topic.ID)
            }
        }
    }
}
```

## Step 2: Get Topic Components

Find available components for testing:

```go
topicID := "your-topic-uuid"

// Get components for the topic
ctx := context.Background()
componentsResp, err := client.GetTopicComponents(ctx, topicID)
if err != nil {
    log.Fatalf("Failed to get components: %v", err)
}

var componentIDs []string
for _, resp := range componentsResp {
    for _, comp := range resp.Components {
        fmt.Printf("Component: %s (%s) - %s\n", comp.Name, comp.Type, comp.ID)
        componentIDs = append(componentIDs, comp.ID)
    }
}
```

## Step 3: Create a Job

Create a new certification job:

```go
// Create job with specific components
ctx := context.Background()
job, err := client.CreateJob(
    ctx,
    topicID,
    componentIDs,              // Components to test against
    "Certification run via API", // Comment
)
if err != nil {
    log.Fatalf("Failed to create job: %v", err)
}

fmt.Printf("Job created: %s\n", job.Job.ID)
fmt.Printf("Status: %s\n", job.Job.Status)

jobID := job.Job.ID
```

Alternatively, schedule a job (auto-selects latest components):

```go
ctx := context.Background()
job, err := client.ScheduleJob(ctx, topicID)
if err != nil {
    log.Fatalf("Failed to schedule job: %v", err)
}
```

## Step 4: Update Job State

Progress the job through its lifecycle:

```go
// Move to pre-run
ctx := context.Background()
_, err = client.UpdateJobState(ctx, jobID, lib.JobStatePreRun, "Starting pre-run setup")
if err != nil {
    log.Fatalf("Failed to update state: %v", err)
}
fmt.Println("State: pre-run")

// Perform pre-run setup here...

// Move to running
_, err = client.UpdateJobState(ctx, jobID, lib.JobStateRunning, "Executing tests")
if err != nil {
    log.Fatalf("Failed to update state: %v", err)
}
fmt.Println("State: running")

// Execute tests here...
```

### Available Job States

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

## Step 5: Upload Test Results

Attach test results to the job:

```go
// Upload from a file
ctx := context.Background()
uploadResp, err := client.UploadFile(
    ctx,
    jobID,
    "/path/to/results.xml",
    "application/junit",
)
if err != nil {
    log.Printf("Failed to upload file: %v", err)
} else {
    fmt.Printf("Uploaded: %s (ID: %s)\n", uploadResp.File.Name, uploadResp.File.ID)
}
```

### Upload from Memory

```go
// Upload content directly
ctx := context.Background()
content := []byte("<testsuite>...</testsuite>")
uploadResp, err := client.UploadFileContent(
    ctx,
    jobID,
    "results.xml",
    "application/junit",
    content,
)
```

### Common MIME Types

| Type | Description |
|------|-------------|
| `application/junit` | JUnit XML test results |
| `application/json` | JSON data |
| `text/plain` | Plain text logs |
| `application/x-tar` | Tar archives |
| `application/gzip` | Gzipped files |

## Step 6: Complete the Job

Set the final job state:

```go
// Tests passed
ctx := context.Background()
_, err = client.UpdateJobState(ctx, jobID, lib.JobStateSuccess, "All tests passed")
if err != nil {
    log.Fatalf("Failed to update state: %v", err)
}

// Or if tests failed
_, err = client.UpdateJobState(ctx, jobID, lib.JobStateFailure, "Some tests failed")
```

## Step 7: Query Job Status

Check job state at any time:

```go
// Get job details
ctx := context.Background()
job, err := client.GetJob(ctx, jobID)
if err != nil {
    log.Fatalf("Failed to get job: %v", err)
}

fmt.Printf("Job: %s\n", job.Job.ID)
fmt.Printf("Status: %s\n", job.Job.Status)
fmt.Printf("Created: %s\n", job.Job.CreatedAt)
fmt.Printf("Updated: %s\n", job.Job.UpdatedAt)
```

### Get Job State History

```go
ctx := context.Background()
states, err := client.GetJobStates(ctx, jobID)
if err != nil {
    log.Fatalf("Failed to get states: %v", err)
}

fmt.Println("State History:")
for _, state := range states.JobStates {
    fmt.Printf("  %s - %s\n", state.Status, state.CreatedAt)
    if state.Comment != "" {
        fmt.Printf("    Comment: %s\n", state.Comment)
    }
}
```

### Get Job Files

```go
ctx := context.Background()
files, err := client.GetJobFiles(ctx, jobID)
if err != nil {
    log.Fatalf("Failed to get files: %v", err)
}

fmt.Printf("Attached files: %d\n", len(files.Files))
for _, file := range files.Files {
    fmt.Printf("  - %s (ID: %s)\n", file.Name, file.ID)
}
```

## Complete Workflow Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    client := lib.NewClient(
        os.Getenv("GO_DCI_ACCESSKEY"),
        os.Getenv("GO_DCI_SECRETKEY"),
    )

    ctx := context.Background()
    topicID := "your-topic-id"

    // 1. Get components
    componentsResp, _ := client.GetTopicComponents(ctx, topicID)
    var componentIDs []string
    for _, resp := range componentsResp {
        for _, comp := range resp.Components {
            componentIDs = append(componentIDs, comp.ID)
        }
    }

    // 2. Create job
    job, err := client.CreateJob(ctx, topicID, componentIDs, "Automated certification")
    if err != nil {
        log.Fatal(err)
    }
    jobID := job.Job.ID
    fmt.Printf("Created job: %s\n", jobID)

    // 3. Pre-run
    client.UpdateJobState(ctx, jobID, lib.JobStatePreRun, "Setting up environment")
    time.Sleep(time.Second)

    // 4. Running
    client.UpdateJobState(ctx, jobID, lib.JobStateRunning, "Executing tests")

    // 5. Execute your tests here...
    testsPassed := runTests()

    // 6. Upload results
    client.UploadFile(ctx, jobID, "results.xml", "application/junit")

    // 7. Set final state
    if testsPassed {
        client.UpdateJobState(ctx, jobID, lib.JobStateSuccess, "All tests passed")
    } else {
        client.UpdateJobState(ctx, jobID, lib.JobStateFailure, "Some tests failed")
    }

    fmt.Println("Workflow complete!")
}

func runTests() bool {
    // Your test logic here
    return true
}
```

## Error Handling Best Practices

```go
// Retry on transient failures
func updateStateWithRetry(client *lib.Client, jobID string, state lib.JobState, comment string) error {
    ctx := context.Background()
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        _, err := client.UpdateJobState(ctx, jobID, state, comment)
        if err == nil {
            return nil
        }
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second)
        }
    }
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

## Managing Job Files

### Download Files

Download files attached to a job:

```go
// Get file by ID
ctx := context.Background()
fileContent, filename, err := client.GetFile(ctx, fileID)
if err != nil {
    log.Fatalf("Failed to download file: %v", err)
}

// Save to disk
err = os.WriteFile(filename, fileContent, 0644)
if err != nil {
    log.Fatalf("Failed to save file: %v", err)
}
fmt.Printf("Downloaded: %s (%d bytes)\n", filename, len(fileContent))
```

### Delete Files

Remove files from a job:

```go
ctx := context.Background()
err := client.DeleteFile(ctx, fileID)
if err != nil {
    log.Fatalf("Failed to delete file: %v", err)
}
fmt.Println("File deleted successfully")
```

### Complete File Management Example

```go
// List all files for a job
ctx := context.Background()
filesResp, err := client.GetJobFiles(ctx, jobID)
if err != nil {
    log.Fatalf("Failed to get job files: %v", err)
}

fmt.Printf("Job has %d files:\n", len(filesResp.Files))
for _, file := range filesResp.Files {
    fmt.Printf("  - %s (ID: %s, Size: %d bytes)\n", 
        file.Name, file.ID, file.Size)
    
    // Download each file
    content, filename, err := client.GetFile(ctx, file.ID)
    if err != nil {
        log.Printf("Failed to download %s: %v", file.Name, err)
        continue
    }
    
    // Save to local directory
    err = os.WriteFile(fmt.Sprintf("./downloads/%s", filename), content, 0644)
    if err != nil {
        log.Printf("Failed to save %s: %v", filename, err)
    }
}
```

## Clean Up After Testing

Delete jobs that are no longer needed:

```go
// Delete old test jobs
ctx := context.Background()
err := client.DeleteJob(ctx, jobID)
if err != nil {
    log.Fatalf("Failed to delete job: %v", err)
}
fmt.Printf("Job %s deleted\n", jobID)
```

## Complete Example

See the [certification-workflow example](../../examples/certification-workflow/main.go) for a complete working program.

## Next Steps

- [Component Analysis](./03-component-analysis.md) - Query and analyze components
