# API Errors

## Bad Request (400)

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
   - Check the [CLI reference](../cli-reference.md) for required flags
   - See the [library guide](../library-guide.md) for required parameters

## Not Found (404)

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

## Internal Server Error (500, 502, 503)

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
