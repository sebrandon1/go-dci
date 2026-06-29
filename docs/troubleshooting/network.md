# Network Issues

## Connection Refused

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

## Connection Timeout

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

## SSL/TLS Errors

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
