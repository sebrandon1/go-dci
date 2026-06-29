# Troubleshooting Guide

This guide covers common issues when using go-dci and how to resolve them.

## Authentication Errors

### Credentials Not Found

**Symptom:**
```
Error: authentication credentials not found
Error: missing access key or secret key
```

**Cause:** No credentials configured via config file or environment variables.

**Solution:**

1. Set credentials via config file:
   ```bash
   go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>
   ```
   This creates `.go-dci-config.yaml` in the current directory.

2. Or set environment variables:
   ```bash
   export GO_DCI_ACCESSKEY=<your-access-key>
   export GO_DCI_SECRETKEY=<your-secret-key>
   ```

3. Verify credentials are working:
   ```bash
   go-dci identity
   ```

**Note:** Environment variables take precedence over the config file. If both are set, the environment variables will be used.

### Permission Denied (401 Unauthorized)

**Symptom:**
```
Error: API returned 401: Unauthorized
Error: authentication failed
```

**Cause:** Invalid or expired credentials.

**Solutions:**

1. Verify your credentials are correct:
   - Check for typos in access key or secret key
   - Ensure you copied the full key (no truncation)
   - Confirm credentials haven't expired

2. Get new credentials from the [DCI dashboard](https://www.distributed-ci.io/)

3. Update your configuration:
   ```bash
   go-dci config set --accesskey <new-access-key> --secretkey <new-secret-key>
   ```

4. Test authentication:
   ```bash
   go-dci identity
   ```

### Forbidden (403)

**Symptom:**
```
Error: API returned 403: Forbidden
Error: insufficient permissions
```

**Cause:** Your credentials lack permission for the requested operation.

**Common scenarios:**

1. **Creating/updating topics** - Requires admin permissions
2. **Managing teams/users** - Requires admin permissions
3. **Accessing team resources** - Your RemoteCI must belong to the team

**Solutions:**

1. Verify your team membership and permissions in the [DCI dashboard](https://www.distributed-ci.io/)
2. Contact your DCI administrator to request necessary permissions
3. Ensure your RemoteCI is associated with the correct team

## Network Issues

### Connection Refused

**Symptom:**
```
Error: connection refused
Error: dial tcp: connection refused
```

**Cause:** Cannot reach the DCI API server.

**Solutions:**

1. Check network connectivity:
   ```bash
   ping www.distributed-ci.io
   curl -I https://api.distributed-ci.io/api/v1
   ```

2. Verify firewall settings allow outbound HTTPS (port 443)

3. If behind a corporate proxy, configure proxy settings:
   ```bash
   export https_proxy=http://proxy.example.com:8080
   export HTTPS_PROXY=http://proxy.example.com:8080
   ```

4. Check if VPN is required to access DCI

### Connection Timeout

**Symptom:**
```
Error: context deadline exceeded
Error: i/o timeout
Error: TLS handshake timeout
```

**Cause:** Request took longer than configured timeout.

**Solutions:**

**For CLI users:**

The CLI uses default timeouts (30s request, 5s TLS handshake, 10s dial). If you experience frequent timeouts on slow networks, this is a known limitation. Consider:

1. Retry the command
2. Check network conditions
3. Use the library API with custom timeouts (see below)

**For library users:**

Adjust timeout settings in your code:

```go
client := lib.NewClient(accessKey, secretKey)

// Increase timeouts for slow networks
client.RequestTimeout = 60 * time.Second  // Default: 30s
client.TLSTimeout = 10 * time.Second      // Default: 5s
client.DialTimeout = 15 * time.Second     // Default: 10s

// Increase retries for flaky connections
client.MaxRetries = 5                     // Default: 3
```

**Network troubleshooting:**

1. Test network latency:
   ```bash
   curl -w "@-" -o /dev/null -s https://api.distributed-ci.io/api/v1 <<'EOF'
   time_namelookup:  %{time_namelookup}\n
   time_connect:     %{time_connect}\n
   time_appconnect:  %{time_appconnect}\n
   time_pretransfer: %{time_pretransfer}\n
   time_starttransfer: %{time_starttransfer}\n
   time_total:       %{time_total}\n
   EOF
   ```

2. Check for network congestion or bandwidth issues
3. Verify DNS resolution is working

### SSL/TLS Errors

**Symptom:**
```
Error: x509: certificate signed by unknown authority
Error: TLS handshake failed
```

**Cause:** SSL certificate validation issues.

**Solutions:**

1. Update CA certificates:
   ```bash
   # macOS
   brew install ca-certificates
   
   # Ubuntu/Debian
   sudo apt-get update && sudo apt-get install ca-certificates
   
   # RHEL/Fedora
   sudo dnf install ca-certificates
   ```

2. Verify system time is correct (SSL validation requires accurate time)

3. If using a corporate proxy with SSL inspection, you may need to add the corporate CA certificate to your trust store

## API Errors

### Bad Request (400)

**Symptom:**
```
Error: API returned 400: Bad Request
Error: invalid request payload
```

**Cause:** Malformed request or invalid parameters.

**Common causes and solutions:**

1. **Invalid JSON payload:**
   - Check for typos in field names
   - Ensure values match expected types (string, number, boolean)
   - Verify required fields are present

2. **Invalid UUIDs:**
   ```bash
   # UUIDs must be in format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
   go-dci job --id 12345678-1234-1234-1234-123456789abc
   ```

3. **Invalid date formats:**
   ```bash
   # Use ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ
   # Not: MM/DD/YYYY or other formats
   ```

4. **Missing required fields:**
   - Check the [CLI reference](cli-reference.md) for required flags
   - See the [library guide](library-guide.md) for required parameters

### Not Found (404)

**Symptom:**
```
Error: API returned 404: Not Found
Error: resource not found
```

**Cause:** The requested resource doesn't exist.

**Solutions:**

1. Verify the ID is correct (UUIDs are case-sensitive)

2. Check if the resource was deleted

3. Ensure you have access to the resource (it may exist but be in a different team)

4. List available resources to find the correct ID:
   ```bash
   go-dci topics         # List all topics
   go-dci jobs -d 30     # List jobs from last 30 days
   go-dci components     # List all components
   ```

### Internal Server Error (500, 502, 503)

**Symptom:**
```
Error: API returned 500: Internal Server Error
Error: API returned 502: Bad Gateway
Error: API returned 503: Service Unavailable
```

**Cause:** Server-side error or maintenance.

**Solutions:**

1. **For 500 errors:**
   - This is a server bug. Check if the issue reproduces
   - Report to DCI support with the request details

2. **For 502/503 errors:**
   - The service may be temporarily unavailable
   - Wait a few minutes and retry
   - Check [DCI status page](https://www.distributed-ci.io/) for maintenance announcements

3. **Automatic retry (library only):**
   The library automatically retries 5xx errors for GET and DELETE requests up to `MaxRetries` times (default: 3).

## File Operations

### Upload Failures

**Symptom:**
```
Error: failed to upload file
Error: invalid file path
Error: file too large
```

**Solutions:**

1. **File not found:**
   ```bash
   # Verify file exists
   ls -lh /path/to/file.tar.gz
   
   # Use absolute paths
   go-dci upload-file --job-id <job-id> --file /absolute/path/to/file.tar.gz
   ```

2. **Permission denied:**
   ```bash
   # Check file permissions
   chmod 644 /path/to/file.tar.gz
   ```

3. **File too large:**
   - DCI has file size limits (typically 1GB)
   - Compress large files before uploading:
     ```bash
     tar czf results.tar.gz results/
     go-dci upload-file --job-id <job-id> --file results.tar.gz --mime-type application/gzip
     ```

4. **Wrong MIME type:**
   Specify the correct MIME type:
   ```bash
   go-dci upload-file --job-id <job-id> --file results.xml --mime-type application/xml
   go-dci upload-file --job-id <job-id> --file results.json --mime-type application/json
   go-dci upload-file --job-id <job-id> --file results.tar.gz --mime-type application/gzip
   ```

### Download Issues

**Symptom:**
```
Error: failed to download file
Error: file content empty
```

**Solutions:**

1. Verify the file ID is correct:
   ```bash
   # List files for a job
   go-dci job-files --job-id <job-id>
   ```

2. Check if you have permission to access the file (must be in the same team)

3. For library users, handle the response properly:
   ```go
   data, contentType, err := client.GetFile(fileID)
   if err != nil {
       log.Fatalf("Download failed: %v", err)
   }
   
   // Save to file
   err = os.WriteFile("output.tar.gz", data, 0644)
   if err != nil {
       log.Fatalf("Failed to save file: %v", err)
   }
   ```

## Job Management

### Invalid State Transitions

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

### Invalid Component IDs

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

### Missing Job Comment

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

## Getting Additional Help

If you continue to experience issues:

1. **Check the documentation:**
   - [CLI Reference](cli-reference.md) - All CLI commands and flags
   - [Library Guide](library-guide.md) - API reference for library usage
   - [Tutorials](tutorials/) - Step-by-step guides

2. **Enable debug output:**
   
   For CLI users, run with increased verbosity to see detailed request/response information.
   
   For library users, inspect error messages:
   ```go
   result, err := client.GetTopic(topicID)
   if err != nil {
       fmt.Printf("Full error: %v\n", err)
       // Error includes HTTP status codes and response bodies
   }
   ```

3. **Review examples:**
   Check the [examples directory](../examples/) for working code:
   - [basic-usage](../examples/basic-usage/)
   - [certification-workflow](../examples/certification-workflow/)
   - [component-query](../examples/component-query/)

4. **Contact support:**
   - Report bugs: [GitHub Issues](https://github.com/sebrandon1/go-dci/issues)
   - DCI documentation: [https://doc.distributed-ci.io/](https://doc.distributed-ci.io/)
   - DCI dashboard: [https://www.distributed-ci.io/](https://www.distributed-ci.io/)
