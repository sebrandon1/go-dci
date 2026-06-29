# Troubleshooting Guide

This guide covers common issues and their solutions when using go-dci.

## Authentication Errors

### 401 Unauthorized

**Symptom**: Commands fail with `401 Unauthorized` error.

**Cause**: Invalid or missing credentials.

**Solution**:
1. Verify your credentials are correct:
   ```bash
   go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>
   ```

2. Test authentication:
   ```bash
   go-dci identity
   ```

3. If using environment variables, ensure they're set correctly:
   ```bash
   export GO_DCI_ACCESSKEY="your-access-key"
   export GO_DCI_SECRETKEY="your-secret-key"
   ```

4. Check that credentials haven't expired or been revoked in the DCI portal.

### 403 Forbidden

**Symptom**: Commands return `403 Forbidden` error.

**Cause**: Your credentials are valid but lack permissions for the requested operation.

**Solution**:
1. Verify your team has access to the resource (topic, job, component, etc.)
2. Contact your DCI administrator to request appropriate permissions
3. Use `go-dci teams` to check your team memberships
4. For remote CI operations, ensure your remote CI has been granted access to the topic

## Resource Not Found

### 404 Not Found

**Symptom**: Commands fail with `404 Not Found` when querying a specific resource.

**Cause**: The resource ID or name doesn't exist, or you don't have access to it.

**Solution**:
1. Verify the resource ID is correct:
   ```bash
   # List all resources to find the correct ID
   go-dci topics
   go-dci jobs
   go-dci components
   ```

2. Check for typos in the ID or name
3. Ensure the resource hasn't been deleted
4. Verify you have permission to view the resource (see 403 Forbidden above)

## Configuration Issues

### Config File Not Found

**Symptom**: Commands fail with "config file not found" or similar error.

**Cause**: No configuration file exists yet.

**Solution**:
Run the config command to create the configuration file:
```bash
go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>
```

The config file is created at `~/.go-dci/config.json` by default.

## Network Issues

### Connection Timeout

**Symptom**: Commands hang or timeout when connecting to the DCI API.

**Cause**: Network connectivity issues, firewall blocking requests, or DCI service unavailability.

**Solution**:
1. Check your internet connection
2. Verify you can reach the DCI API endpoint:
   ```bash
   curl -I https://api.distributed-ci.io/api/v1/identity
   ```

3. Check firewall rules:
   - Ensure outbound HTTPS (port 443) is allowed
   - Verify no proxy is blocking the connection
   - If behind a corporate firewall, contact your network administrator

4. Check DCI service status at [https://status.distributed-ci.io](https://status.distributed-ci.io) (if available)

5. Try increasing the client timeout if on a slow network (library usage only):
   ```go
   client := lib.NewClientWithTimeout(accessKey, secretKey, 60*time.Second)
   ```

## Query Results

### Empty Results

**Symptom**: Commands return no results when you expect data.

**Cause**: Overly restrictive filters, or genuinely no matching data.

**Solution**:
1. Remove or relax filters to see if data exists:
   ```bash
   # Instead of filtering, list all first
   go-dci topics
   go-dci jobs
   ```

2. Check pagination - you may need to increase the limit:
   ```bash
   go-dci jobs --limit 100
   ```

3. Verify your filters are correct:
   - Check date ranges are valid
   - Ensure filter values match exactly (names are case-sensitive)
   - Try broader search terms

4. Use the `--where` flag correctly for advanced filtering (see CLI Reference)

### "No components found for topic"

**Symptom**: `go-dci topic-components` returns "No components found" for a valid topic.

**Cause**: Wrong topic ID, or the topic genuinely has no components.

**Solution**:
1. Verify the topic ID is correct:
   ```bash
   go-dci topics
   ```

2. Check that components exist for this topic:
   ```bash
   # List all components and grep for the topic
   go-dci components --limit 100 | grep -i "topic-name"
   ```

3. Ensure the topic is active and has been used in jobs:
   ```bash
   go-dci jobs --where topic_id=<topic-id>
   ```

## Command-Specific Issues

### Job Creation Fails

**Symptom**: `go-dci create-job` fails or returns an error.

**Common Causes**:
- Invalid topic ID
- Missing required fields
- Remote CI not permitted for the topic
- Team lacks permissions

**Solution**:
1. Verify all required fields are provided:
   ```bash
   go-dci create-job --topic-id <id> --remoteci-id <id> --team-id <id>
   ```

2. Check that your remote CI is authorized for the topic:
   ```bash
   go-dci remotecis
   go-dci topics
   ```

3. Ensure your team has job creation permissions

### File Upload Fails

**Symptom**: `go-dci upload-file` fails with permission or validation errors.

**Common Causes**:
- Job doesn't exist or is in the wrong state
- File too large
- Invalid MIME type
- Job already finalized

**Solution**:
1. Verify the job exists and is active:
   ```bash
   go-dci job --id <job-id>
   ```

2. Check file size limits (typically 100MB max)

3. Ensure the job is in a state that accepts uploads (not finalized)

4. Provide correct MIME type if auto-detection fails:
   ```bash
   go-dci upload-file --job-id <id> --file <path> --mime-type application/json
   ```

## Library Usage Issues

### Context Timeout

**Symptom**: Operations fail with context deadline exceeded.

**Cause**: Operation takes longer than the context timeout allows.

**Solution**:
Use a longer timeout context:
```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

topics, err := client.GetTopicsWithContext(ctx)
```

### Rate Limiting

**Symptom**: Intermittent failures or 429 Too Many Requests errors.

**Cause**: Making too many API requests in a short time.

**Solution**:
1. Implement exponential backoff retry logic
2. Cache results when possible
3. Batch operations instead of individual calls
4. Add delays between bulk operations:
   ```go
   for _, item := range items {
       // Process item
       time.Sleep(100 * time.Millisecond)
   }
   ```

## Getting Help

If you're still experiencing issues:

1. Check the [DCI documentation](https://doc.distributed-ci.io/)
2. Review the [CLI Reference](cli-reference.md) for correct command usage
3. Enable debug logging (if available) to see detailed API interactions
4. [Open an issue](https://github.com/sebrandon1/go-dci/issues) with:
   - Command or code that fails
   - Full error message
   - go-dci version (`go-dci version`)
   - Expected vs actual behavior

## Common Workarounds

### Stale Cache

If you see outdated data, there may be caching at the API level. Wait a few moments and retry, or query with different parameters to bypass cache.

### Unicode/Special Characters

If resource names contain special characters, ensure your terminal encoding is UTF-8, or use the resource ID instead of the name.

### Large Result Sets

For queries that return thousands of results:
1. Use pagination with `--limit` and `--offset`
2. Apply filters with `--where` to reduce result size
3. Consider using the library API for programmatic processing
