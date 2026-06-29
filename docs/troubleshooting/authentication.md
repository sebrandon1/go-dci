# Authentication Errors

## Credentials Not Found

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

## Permission Denied (401 Unauthorized)

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

## Forbidden (403)

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
