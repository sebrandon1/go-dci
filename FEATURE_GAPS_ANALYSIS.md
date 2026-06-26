# go-dci Feature Gaps Analysis

Generated: 2026-06-17

## Executive Summary

This analysis identifies 15 major feature gaps across the go-dci CLI and library wrapper for the Red Hat Distributed CI API. The codebase has good test coverage (124 tests) and implements full CRUD operations for most resources, but lacks advanced filtering, pagination controls, client configurability, and comprehensive documentation for newer resources.

---

## Feature Gaps Found

### 1. Products Resource - Missing CRUD Operations (CLI)
- **Source**: `cmd/productsCmd.go:16,46` (only getProductsCmd and getProductCmd)
- **Category**: Missing CRUD
- **Impact**: High
- **Effort**: Low
- **Description**: The Products API has full read operations in the library (`GetProducts`, `GetProduct` at `lib/client.go:1456,1479`) but the CLI is missing create, update, and delete commands. Users cannot manage Products resources through the CLI, forcing them to use the library directly or other tools.

### 2. Files Resource - Missing List/Search Operations
- **Source**: `cmd/filesCmd.go:19,67` (only getFileCmd and deleteFileCmd), `lib/client.go:713,990,1014,1032`
- **Category**: Missing CRUD
- **Impact**: Medium
- **Effort**: Medium
- **Description**: Files can only be retrieved by ID (`GetFile`), uploaded (`UploadFile`, `UploadFileContent`), deleted (`DeleteFile`), or listed per-job (`GetJobFiles`). There is no global file listing or search capability across all jobs/resources. This makes it difficult to audit uploaded files or find specific artifacts.

### 3. Jobs - Missing Filter/Query Capabilities
- **Source**: `cmd/jobs.go` (no filter flags), `cmd/get.go:563-564` (only date filtering)
- **Category**: Feature
- **Impact**: High
- **Effort**: Medium
- **Description**: Job queries support only date-based filtering (`--start-date`, `--end-date`). Missing filters include: job status/state (success/failure/running), component ID, topic ID, team ID, tags, and remote CI. The library `GetJobs` and `GetJobsByDate` methods don't expose these filtering options even though the DCI API likely supports them.

### 4. No Pagination Controls in CLI
- **Source**: Search across all `cmd/*Cmd.go` files shows zero pagination flags
- **Category**: Feature
- **Impact**: Medium
- **Effort**: Low
- **Description**: No CLI commands expose limit, offset, or page flags for list operations. The library internally uses pagination (e.g., `fetchJobs`, `fetchComponents`, `fetchTopics` with hardcoded limits), but users cannot control page size or navigate results. This impacts usability when dealing with large datasets (e.g., thousands of jobs or components).

### 5. Topics - Missing Search/Filter Capabilities
- **Source**: `cmd/topics.go` (no filter flags), `lib/client.go:352-378`
- **Category**: Feature
- **Impact**: Medium
- **Effort**: Low
- **Description**: Topic listing (`GetTopics`) returns all topics with no filtering by product, name pattern, component type, or state. Users must fetch all topics and filter client-side. Adding optional filters would improve performance and usability for organizations with many topics.

### 6. Components - Missing Advanced Query Features
- **Source**: `cmd/get.go:173` (only topic-id filter), `lib/client.go:755-872`
- **Category**: Feature
- **Impact**: Medium
- **Effort**: Medium
- **Description**: Component queries support only filtering by topic ID. Missing filters include: component type, state, version pattern matching, tags, and date ranges. This limits the ability to perform component lifecycle analysis (e.g., "find all active OCP 4.15 components").

### 7. Client Configuration - No Retry/Timeout Customization
- **Source**: `lib/client.go:60-145` (hardcoded retry logic with exponential backoff)
- **Category**: Feature
- **Impact**: Medium
- **Effort**: Medium
- **Description**: The client implements retry logic with exponential backoff but hardcodes MaxRetries (5), initial backoff (1s), and max backoff (30s). There's no way to configure retry attempts, timeouts, or disable retries for different deployment scenarios (CI/CD vs interactive). The `NewClient` function should accept an optional config struct.

### 8. Config Command - Missing List and Validate Operations
- **Source**: `cmd/config.go` (only set command exists)
- **Category**: Feature
- **Impact**: Low
- **Effort**: Low
- **Description**: The config command supports only `config set --accesskey --secretkey`. Missing operations: `config list` (show current configuration), `config validate` (test credentials against DCI API), `config unset` (remove credentials), and `config path` (show config file location). These would improve troubleshooting and credential management.

### 9. Job State Validation - Missing Transition Logic
- **Source**: `lib/structs.go:431-438` (job states defined), no validation in `lib/client.go:929`
- **Category**: Feature
- **Impact**: Low
- **Effort**: Low
- **Description**: Job states are defined (new, pre-run, running, post-run, success, failure, killed, error) but there's no validation of legal state transitions. The `UpdateJobState` method accepts any JobState without checking if the transition is valid (e.g., preventing "success" → "running"). This could lead to invalid workflows.

### 10. Batch Operations - Not Implemented
- **Source**: No batch methods exist in `lib/client.go`
- **Category**: Feature
- **Impact**: Medium
- **Effort**: High
- **Description**: No support for batch operations like creating multiple components, jobs, or resources in a single API call. Users must loop client-side which is inefficient for bulk workflows (e.g., "create jobs for all topics in product X"). The DCI API may or may not support batch operations - needs investigation.

### 11. Missing CLI Output Formats
- **Source**: All `cmd/*Cmd.go` files show only `--output json` or stdout
- **Category**: Feature
- **Impact**: Low
- **Effort**: Low
- **Description**: CLI output supports only JSON (`--output json`) or custom formatted stdout. Missing formats include: YAML (for kubectl users), CSV (for spreadsheet export), and table/tabular (for aligned columns). This limits integration with other tools and scripts.

### 12. RemoteCIs - Full CRUD Exists But Limited CLI Discoverability
- **Source**: `lib/client.go` (GetRemoteCIs, GetRemoteCI, CreateRemoteCI, UpdateRemoteCI, DeleteRemoteCI exist), `cmd/remotecisCmd.go`
- **Category**: Test Gap / Documentation
- **Impact**: Low
- **Effort**: Low
- **Description**: RemoteCI operations are fully implemented in the library and CLI with full CRUD. However, there's no CLI subcommand to list RemoteCIs filtered by team or state, making discovery difficult. Adding filter flags to `get remotecis` would improve usability.

### 13. Library Examples - Incomplete Advanced Workflows
- **Source**: `examples/` (3 examples: basic-usage, certification-workflow, component-query)
- **Category**: Documentation
- **Impact**: Medium
- **Effort**: Medium
- **Description**: Examples demonstrate basic flows but missing: topic/component lifecycle management, team/user administration, error handling patterns with retries, concurrent job monitoring, and file upload with state transitions. The certification-workflow example is partial (dry-run flag but incomplete integration).

### 14. Update Structs - Potentially Missing Fields
- **Source**: `lib/structs.go:320-598` (8 Update*Request structs)
- **Category**: Feature (Potential)
- **Impact**: Low
- **Effort**: Medium
- **Description**: Update request structs exist for all major resources but fields may not match the full DCI API schema. Without API documentation comparison, it's unclear if fields like metadata, labels, annotations, or custom properties are missing. The swagger alignment script exists (`scripts/check-swagger-alignment.go`) but endpoint coverage, not field completeness.

### 15. Missing Documentation for New Resources
- **Source**: `docs/` directory (7 markdown files)
- **Category**: Documentation
- **Impact**: Medium
- **Effort**: Medium
- **Description**: Documentation covers installation, authentication, CLI reference, library guide, and 3 tutorials. However, missing docs for: Products resource usage, RemoteCIs management, Teams/Users administration, advanced filtering patterns, error handling best practices, and rate limiting/retry strategies. The library-guide.md was last updated in March but newer CLI commands aren't documented.

---

## Summary Statistics

- **Total Gaps Identified**: 15
- **Impact Distribution**:
  - High: 2 (Jobs filtering, Products CRUD)
  - Medium: 8
  - Low: 5
- **Effort Distribution**:
  - High: 1 (Batch operations)
  - Medium: 7
  - Low: 7
- **Category Distribution**:
  - Missing CRUD: 2
  - Feature: 9
  - Documentation: 2
  - Test Gap: 1
  - Feature (Potential): 1

## High Priority Recommendations

1. **Jobs Filtering** (Gap #3) - Add status, component, topic, team filters to job queries
2. **Products CLI** (Gap #1) - Implement create/update/delete commands to match library
3. **Pagination** (Gap #4) - Expose limit/offset flags in all list commands
4. **Client Config** (Gap #7) - Add retry and timeout configuration to NewClient
5. **Documentation** (Gap #15) - Document Products, RemoteCIs, Teams/Users workflows

## Testing Coverage Notes

The library has excellent test coverage with 124 unit tests covering all major resources including:
- All CRUD operations for ComponentTypes, Topics, Jobs, Components, RemoteCIs, Teams, Users, Products, Files
- Error handling and retry logic
- Request/response serialization
- Authentication

No tests were found to be missing based on the analysis. The codebase has zero TODO/FIXME/HACK comments indicating technical debt awareness.
