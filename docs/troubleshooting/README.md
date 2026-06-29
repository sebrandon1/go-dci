# Troubleshooting Guide

This guide covers common issues when using go-dci and how to resolve them.

## Categories

- [Authentication Errors](authentication.md) - Credential issues, permission denied, forbidden access
- [Network Issues](network.md) - Connection problems, timeouts, SSL/TLS errors
- [API Errors](api-errors.md) - Bad requests, not found, server errors
- [File Operations](file-operations.md) - Upload/download failures
- [Job Management](job-management.md) - State transitions, component IDs, job comments

## Getting Additional Help

If you continue to experience issues:

1. **Check the documentation:**
   - [CLI Reference](../cli-reference.md) - All CLI commands and flags
   - [Library Guide](../library-guide.md) - API reference for library usage
   - [Tutorials](../tutorials/) - Step-by-step guides

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
   Check the [examples directory](../../examples/) for working code:
   - [basic-usage](../../examples/basic-usage/)
   - [certification-workflow](../../examples/certification-workflow/)
   - [component-query](../../examples/component-query/)

4. **Contact support:**
   - Report bugs: [GitHub Issues](https://github.com/sebrandon1/go-dci/issues)
   - DCI documentation: [https://doc.distributed-ci.io/](https://doc.distributed-ci.io/)
   - DCI dashboard: [https://www.distributed-ci.io/](https://www.distributed-ci.io/)
