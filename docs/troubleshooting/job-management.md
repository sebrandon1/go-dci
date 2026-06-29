# Job Management

## Invalid State Transitions

**Symptom:**
```
Error: invalid state transition
Error: cannot transition from X to Y
```

**Cause:** Invalid job state lifecycle transition.

**Valid state transitions:**

```
new → pre-run → running → post-run → success/failure/error
                    ↓
                 killed (from any state)
```

**State constants (library):**
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

**Solutions:**

1. Follow the correct state sequence:
   ```bash
   # Correct workflow
   go-dci create-job --topic-id <topic-id>
   go-dci update-job-state --job-id <job-id> --state pre-run
   go-dci update-job-state --job-id <job-id> --state running
   go-dci update-job-state --job-id <job-id> --state post-run
   go-dci update-job-state --job-id <job-id> --state success
   ```

2. Check current job state before updating:
   ```bash
   go-dci job --id <job-id>
   ```

3. You cannot transition backwards or skip states (except to "killed")

## Invalid Component IDs

**Symptom:**
```
Error: invalid component IDs
Error: component not found
Error: component not in topic
```

**Cause:** Component IDs don't exist or don't belong to the specified topic.

**Solutions:**

1. List components for a specific topic:
   ```bash
   go-dci topic-components --topic-id <topic-id>
   ```

2. Verify component IDs are correct (UUIDs are case-sensitive)

3. Use `schedule-job` to auto-select latest components:
   ```bash
   # Instead of manually specifying component IDs
   go-dci schedule-job --topic-id <topic-id>
   ```

4. For library users:
   ```go
   // Get latest components for topic
   components, err := client.GetTopicComponents(topicID)
   if err != nil {
       log.Fatalf("Failed to get components: %v", err)
   }
   
   // Extract component IDs
   var componentIDs []string
   for _, comp := range components {
       componentIDs = append(componentIDs, comp.ID)
   }
   
   // Create job with components
   job, err := client.CreateJob(topicID, componentIDs, "Test run")
   ```

## Missing Job Comment

**Symptom:**
```
Error: job comment is required
```

**Cause:** Some job operations require a comment describing the purpose.

**Solution:**

Always provide a meaningful comment when creating jobs:
```bash
go-dci create-job --topic-id <topic-id> --components <id1>,<id2> --comment "Testing OCP 4.17 upgrade"
```

Library usage:
```go
job, err := client.CreateJob(topicID, componentIDs, "Testing OCP 4.17 upgrade")
```
