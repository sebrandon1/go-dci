# CLI Reference

## `topics` — List Topics

Get all available topics from DCI.

```bash
# Get all topics
go-dci topics

# Output as JSON
go-dci topics --output json
```

```
Usage:
  dci topics [flags]

Flags:
  -h, --help            help for topics
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Getting all topics from DCI
ID: topic-123 | Name: OCP-4.14 | Product: OpenShift Container Platform | State: active
ID: topic-456 | Name: OCP-4.15 | Product: OpenShift Container Platform | State: active
Total Topics: 2
```

## `identity` — Verify Authentication

Verify your DCI credentials are configured correctly and display identity information.

```bash
# Check authentication
go-dci identity

# Output as JSON
go-dci identity --output json
```

```
Usage:
  dci identity [flags]

Flags:
  -h, --help            help for identity
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Authentication successful!
---
ID:       abc123-def456-ghi789
Name:     my-remoteci
Type:     remoteci
Team:     My Team
Team ID:  team-123
State:    active
```

## `jobs` — Query Jobs

Get all jobs with a specific age in days. Filters for jobs running the [certsuite](https://github.com/redhat-best-practices-for-k8s/certsuite).

```bash
# Get jobs from the last 30 days
go-dci jobs -d 30

# Output as JSON
go-dci jobs -d 30 --output json
```

```
Usage:
  dci jobs [flags]

Flags:
  -d, --age string      Age in days
  -h, --help            help for jobs
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Getting all jobs from DCI that are 30 days old
Job ID: 78ed13e1-841f-4c04-a1c6-8df9028c67cd  -  Certsuite Version: tnf-v5.1.3 (Days Since: 11.249845)
Job ID: eb491abd-ec8b-42cc-aa8b-98adf741b236  -  Certsuite Version: tnf-v5.1.3 (Days Since: 11.305408)
```

## `ocpcount` — OCP Version Statistics

Get the count of certsuite jobs for each OCP version.

```bash
# Get OCP version counts for the last 30 days
go-dci ocpcount -d 30

# Output as JSON
go-dci ocpcount -d 30 --output json
```

```
Usage:
  dci ocpcount [flags]

Flags:
  -d, --age string      Age in days
  -h, --help            help for ocpcount
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

## `componenttypes` — List Component Types

Get all available component types from DCI.

```bash
# Get all component types
go-dci componenttypes

# Output as JSON
go-dci componenttypes --output json
```

```
Usage:
  dci componenttypes [flags]

Flags:
  -h, --help            help for componenttypes
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Getting all component types from DCI
ID: ct-123 | Name: ocp | State: active
ID: ct-456 | Name: certsuite | State: active
ID: ct-789 | Name: rhel | State: active
Total Component Types: 3
```

## `components` — Query Components

Get all components, optionally filtered by topic ID.

```bash
# Get all components
go-dci components

# Get components for a specific topic
go-dci components --topic <topic-id>

# Output as JSON
go-dci components --output json
```

```
Usage:
  dci components [flags]

Flags:
  -h, --help            help for components
  -o, --output string   Output format (json) - default is stdout (default "stdout")
  -t, --topic string    Filter components by topic ID
```

Example output:

```
Getting all components from DCI
ID: abc123 | Name: OpenShift 4.14.1 | Type: ocp | Version: 4.14.1 | TopicID: topic-456
ID: def456 | Name: certsuite v5.1.3 | Type: certsuite | Version: v5.1.3 | TopicID: topic-456
Total Components: 2
```

## `create-job` — Create a New Job

Create a new job in DCI for a given topic.

```bash
# Create a job for a topic
go-dci create-job --topic-id <topic-id>

# Create a job with specific components
go-dci create-job --topic-id <topic-id> --components <comp-id-1>,<comp-id-2>

# Create a job with a comment
go-dci create-job --topic-id <topic-id> --comment "Test run for certification"

# Output as JSON
go-dci create-job --topic-id <topic-id> --output json
```

```
Usage:
  dci create-job [flags]

Flags:
      --comment string      Optional comment for the job
      --components string   Comma-separated list of component IDs
  -h, --help                help for create-job
  -o, --output string       Output format (json) - default is stdout (default "stdout")
      --topic-id string     Topic ID for the job (required)
```

Example output:

```
Creating job for topic ID: topic-123
Job created successfully!
---
Job ID:    job-456
Topic ID:  topic-123
Status:    new
State:     active
Created:   2024-01-01T00:00:00.000000
```

## `update-job-state` — Update Job State

Update the state of a job (pre-run, running, success, failure, etc.).

```bash
# Update job to running state
go-dci update-job-state --job-id <job-id> --status running

# Update job to success with comment
go-dci update-job-state --job-id <job-id> --status success --comment "All tests passed"

# Update job to failure
go-dci update-job-state --job-id <job-id> --status failure --comment "Tests failed"

# Output as JSON
go-dci update-job-state --job-id <job-id> --status success --output json
```

```
Usage:
  dci update-job-state [flags]

Flags:
      --comment string   Optional comment for the state change
  -h, --help             help for update-job-state
      --job-id string    Job ID to update (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --status string    New status (pre-run, running, success, failure, etc.) (required)
```

Valid status values: `new`, `pre-run`, `running`, `post-run`, `success`, `failure`, `killed`, `error`

Example output:

```
Updating job job-123 to status: running
Job state updated successfully!
---
JobState ID: jobstate-789
Job ID:      job-123
Status:      running
Created:     2024-01-01T00:00:00.000000
```

## `upload-file` — Upload File to Job

Upload a file (e.g., test results) to a job in DCI.

```bash
# Upload a JUnit test results file
go-dci upload-file --job-id <job-id> --file /path/to/results.xml

# Upload with custom MIME type
go-dci upload-file --job-id <job-id> --file /path/to/results.json --mime application/json

# Output as JSON
go-dci upload-file --job-id <job-id> --file /path/to/results.xml --output json
```

```
Usage:
  dci upload-file [flags]

Flags:
      --file string     Path to the file to upload (required)
  -h, --help            help for upload-file
      --job-id string   Job ID to attach the file to (required)
      --mime string     MIME type of the file (default "application/junit")
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

Example output:

```
Uploading file results.xml to job job-123
File uploaded successfully!
---
File ID:   file-456
Job ID:    job-123
Name:      results.xml
MIME Type: application/junit
Size:      1024 bytes
Created:   2024-01-01T00:00:00.000000
```

## Topics

### `topic` — Get Topic by ID

Get details for a specific topic.

```bash
# Get a topic by ID
go-dci topic --id <topic-id>

# Output as JSON
go-dci topic --id <topic-id> --output json
```

```
Usage:
  dci topic [flags]

Flags:
  -h, --help            help for topic
      --id string       Topic ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `create-topic` — Create Topic

Create a new topic in DCI.

```bash
# Create a topic
go-dci create-topic --name "OCP-4.16" --product-id <product-id>

# Create with component types
go-dci create-topic --name "OCP-4.16" --product-id <product-id> --component-types ocp,certsuite

# Output as JSON
go-dci create-topic --name "OCP-4.16" --product-id <product-id> --output json
```

```
Usage:
  dci create-topic [flags]

Flags:
      --component-types string   Comma-separated list of component type names
  -h, --help                     help for create-topic
      --name string              Topic name (required)
  -o, --output string            Output format (json) - default is stdout (default "stdout")
      --product-id string        Product ID (required)
```

### `update-topic` — Update Topic

Update an existing topic.

```bash
# Update topic name
go-dci update-topic --id <topic-id> --name "OCP-4.16.1"

# Update topic state
go-dci update-topic --id <topic-id> --state inactive

# Output as JSON
go-dci update-topic --id <topic-id> --name "New Name" --output json
```

```
Usage:
  dci update-topic [flags]

Flags:
  -h, --help             help for update-topic
      --id string        Topic ID to update (required)
      --name string      New topic name
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --state string     New state (active/inactive)
```

### `delete-topic` — Delete Topic

Delete a topic from DCI.

```bash
# Delete a topic
go-dci delete-topic --id <topic-id>

# Output as JSON
go-dci delete-topic --id <topic-id> --output json
```

```
Usage:
  dci delete-topic [flags]

Flags:
  -h, --help             help for delete-topic
      --id string        Topic ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

### `topic-components` — List Topic Components

Get all components associated with a topic.

```bash
# Get components for a topic
go-dci topic-components --topic-id <topic-id>

# Output as JSON
go-dci topic-components --topic-id <topic-id> --output json
```

```
Usage:
  dci topic-components [flags]

Flags:
  -h, --help             help for topic-components
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --topic-id string  Topic ID (required)
```

## Components

### `component` — Get Component by ID

Get details for a specific component.

```bash
# Get a component by ID
go-dci component --id <component-id>

# Output as JSON
go-dci component --id <component-id> --output json
```

```
Usage:
  dci component [flags]

Flags:
  -h, --help            help for component
      --id string       Component ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `create-component` — Create Component

Create a new component.

```bash
# Create a component
go-dci create-component --name "OCP 4.16.0" --type ocp --topic-id <topic-id> --version "4.16.0"

# Output as JSON
go-dci create-component --name "OCP 4.16.0" --type ocp --topic-id <topic-id> --version "4.16.0" --output json
```

```
Usage:
  dci create-component [flags]

Flags:
  -h, --help             help for create-component
      --name string      Component name (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --topic-id string  Topic ID (required)
      --type string      Component type (required)
      --version string   Component version
```

### `update-component` — Update Component

Update an existing component.

```bash
# Update component version
go-dci update-component --id <comp-id> --version "4.16.1"

# Update component state
go-dci update-component --id <comp-id> --state inactive

# Output as JSON
go-dci update-component --id <comp-id> --version "4.16.1" --output json
```

```
Usage:
  dci update-component [flags]

Flags:
  -h, --help             help for update-component
      --id string        Component ID to update (required)
      --name string      New component name
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --state string     New state (active/inactive)
      --tags string      Comma-separated list of tags
      --version string   New version
```

### `delete-component` — Delete Component

Delete a component from DCI.

```bash
# Delete a component
go-dci delete-component --id <component-id>

# Output as JSON
go-dci delete-component --id <component-id> --output json
```

```
Usage:
  dci delete-component [flags]

Flags:
  -h, --help             help for delete-component
      --id string        Component ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

## Component Types

### `componenttype` — Get Component Type by ID

Get details for a specific component type.

```bash
# Get a component type by ID
go-dci componenttype --id <componenttype-id>

# Output as JSON
go-dci componenttype --id <componenttype-id> --output json
```

```
Usage:
  dci componenttype [flags]

Flags:
  -h, --help            help for componenttype
      --id string       Component type ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `create-componenttype` — Create Component Type

Create a new component type.

```bash
# Create a component type
go-dci create-componenttype --name "helmchart"

# Output as JSON
go-dci create-componenttype --name "helmchart" --output json
```

```
Usage:
  dci create-componenttype [flags]

Flags:
  -h, --help            help for create-componenttype
      --name string     Component type name (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `update-componenttype` — Update Component Type

Update an existing component type.

```bash
# Update component type name
go-dci update-componenttype --id <ct-id> --name "helm-chart"

# Update component type state
go-dci update-componenttype --id <ct-id> --state inactive

# Output as JSON
go-dci update-componenttype --id <ct-id> --name "helm-chart" --output json
```

```
Usage:
  dci update-componenttype [flags]

Flags:
  -h, --help             help for update-componenttype
      --id string        Component type ID to update (required)
      --name string      New component type name
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --state string     New state (active/inactive)
```

### `delete-componenttype` — Delete Component Type

Delete a component type from DCI.

```bash
# Delete a component type
go-dci delete-componenttype --id <componenttype-id>

# Output as JSON
go-dci delete-componenttype --id <componenttype-id> --output json
```

```
Usage:
  dci delete-componenttype [flags]

Flags:
  -h, --help             help for delete-componenttype
      --id string        Component type ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

## Jobs

### `job` — Get Job by ID

Get details for a specific job.

```bash
# Get a job by ID
go-dci job --id <job-id>

# Output as JSON
go-dci job --id <job-id> --output json
```

```
Usage:
  dci job [flags]

Flags:
  -h, --help            help for job
      --id string       Job ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `update-job` — Update Job

Update an existing job.

```bash
# Update job comment
go-dci update-job --id <job-id> --comment "Updated test results"

# Update job status
go-dci update-job --id <job-id> --status success

# Output as JSON
go-dci update-job --id <job-id> --comment "Updated" --output json
```

```
Usage:
  dci update-job [flags]

Flags:
      --comment string   New job comment
  -h, --help             help for update-job
      --id string        Job ID to update (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --status string    New job status
```

### `delete-job` — Delete Job

Delete a job from DCI.

```bash
# Delete a job
go-dci delete-job --id <job-id>

# Output as JSON
go-dci delete-job --id <job-id> --output json
```

```
Usage:
  dci delete-job [flags]

Flags:
  -h, --help             help for delete-job
      --id string        Job ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

### `schedule-job` — Schedule Job with Component Selection

Schedule a job by automatically selecting the latest components for a topic.

```bash
# Schedule a job for a topic
go-dci schedule-job --topic-id <topic-id>

# Schedule with comment
go-dci schedule-job --topic-id <topic-id> --comment "Nightly certification run"

# Output as JSON
go-dci schedule-job --topic-id <topic-id> --output json
```

```
Usage:
  dci schedule-job [flags]

Flags:
      --comment string    Optional comment for the job
  -h, --help              help for schedule-job
  -o, --output string     Output format (json) - default is stdout (default "stdout")
      --topic-id string   Topic ID for the job (required)
```

### `job-files` — List Job Files

Get all files attached to a job.

```bash
# List files for a job
go-dci job-files --job-id <job-id>

# Output as JSON
go-dci job-files --job-id <job-id> --output json
```

```
Usage:
  dci job-files [flags]

Flags:
  -h, --help            help for job-files
      --job-id string   Job ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

## Job States

### `jobstates` — List Job States

Get all state transitions for a job.

```bash
# List job states
go-dci jobstates --job-id <job-id>

# Output as JSON
go-dci jobstates --job-id <job-id> --output json
```

```
Usage:
  dci jobstates [flags]

Flags:
  -h, --help            help for jobstates
      --job-id string   Job ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

## Files

### `file` — Download File

Download a file from DCI by its ID.

```bash
# Download a file
go-dci file --id <file-id>

# Output as JSON (returns metadata only)
go-dci file --id <file-id> --output json
```

```
Usage:
  dci file [flags]

Flags:
  -h, --help            help for file
      --id string       File ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `delete-file` — Delete File

Delete a file from DCI.

```bash
# Delete a file
go-dci delete-file --id <file-id>

# Output as JSON
go-dci delete-file --id <file-id> --output json
```

```
Usage:
  dci delete-file [flags]

Flags:
  -h, --help             help for delete-file
      --id string        File ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

## Remote CIs

### `remotecis` — List Remote CIs

Get all remote CI systems.

```bash
# List all remote CIs
go-dci remotecis

# Output as JSON
go-dci remotecis --output json
```

```
Usage:
  dci remotecis [flags]

Flags:
  -h, --help            help for remotecis
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `remoteci` — Get Remote CI by ID

Get details for a specific remote CI.

```bash
# Get a remote CI by ID
go-dci remoteci --id <remoteci-id>

# Output as JSON
go-dci remoteci --id <remoteci-id> --output json
```

```
Usage:
  dci remoteci [flags]

Flags:
  -h, --help            help for remoteci
      --id string       Remote CI ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `create-remoteci` — Create Remote CI

Create a new remote CI.

```bash
# Create a remote CI
go-dci create-remoteci --name "my-ci-system" --team-id <team-id>

# Output as JSON
go-dci create-remoteci --name "my-ci-system" --team-id <team-id> --output json
```

```
Usage:
  dci create-remoteci [flags]

Flags:
  -h, --help             help for create-remoteci
      --name string      Remote CI name (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --team-id string   Team ID (required)
```

### `update-remoteci` — Update Remote CI

Update an existing remote CI.

```bash
# Update remote CI name
go-dci update-remoteci --id <remoteci-id> --name "updated-ci-system"

# Update remote CI state
go-dci update-remoteci --id <remoteci-id> --state inactive

# Output as JSON
go-dci update-remoteci --id <remoteci-id> --name "updated-ci-system" --output json
```

```
Usage:
  dci update-remoteci [flags]

Flags:
  -h, --help             help for update-remoteci
      --id string        Remote CI ID to update (required)
      --name string      New remote CI name
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --state string     New state (active/inactive)
```

### `delete-remoteci` — Delete Remote CI

Delete a remote CI from DCI.

```bash
# Delete a remote CI
go-dci delete-remoteci --id <remoteci-id>

# Output as JSON
go-dci delete-remoteci --id <remoteci-id> --output json
```

```
Usage:
  dci delete-remoteci [flags]

Flags:
  -h, --help             help for delete-remoteci
      --id string        Remote CI ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

## Teams

### `teams` — List Teams

Get all teams.

```bash
# List all teams
go-dci teams

# Output as JSON
go-dci teams --output json
```

```
Usage:
  dci teams [flags]

Flags:
  -h, --help            help for teams
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `team` — Get Team by ID

Get details for a specific team.

```bash
# Get a team by ID
go-dci team --id <team-id>

# Output as JSON
go-dci team --id <team-id> --output json
```

```
Usage:
  dci team [flags]

Flags:
  -h, --help            help for team
      --id string       Team ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `create-team` — Create Team

Create a new team.

```bash
# Create a team
go-dci create-team --name "my-team"

# Output as JSON
go-dci create-team --name "my-team" --output json
```

```
Usage:
  dci create-team [flags]

Flags:
  -h, --help            help for create-team
      --name string     Team name (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `update-team` — Update Team

Update an existing team.

```bash
# Update team name
go-dci update-team --id <team-id> --name "updated-team-name"

# Update team state
go-dci update-team --id <team-id> --state inactive

# Output as JSON
go-dci update-team --id <team-id> --name "updated-team-name" --output json
```

```
Usage:
  dci update-team [flags]

Flags:
  -h, --help             help for update-team
      --id string        Team ID to update (required)
      --name string      New team name
  -o, --output string    Output format (json) - default is stdout (default "stdout")
      --state string     New state (active/inactive)
```

### `delete-team` — Delete Team

Delete a team from DCI.

```bash
# Delete a team
go-dci delete-team --id <team-id>

# Output as JSON
go-dci delete-team --id <team-id> --output json
```

```
Usage:
  dci delete-team [flags]

Flags:
  -h, --help             help for delete-team
      --id string        Team ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

## Users

### `users` — List Users

Get all users.

```bash
# List all users
go-dci users

# Output as JSON
go-dci users --output json
```

```
Usage:
  dci users [flags]

Flags:
  -h, --help            help for users
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `user` — Get User by ID

Get details for a specific user.

```bash
# Get a user by ID
go-dci user --id <user-id>

# Output as JSON
go-dci user --id <user-id> --output json
```

```
Usage:
  dci user [flags]

Flags:
  -h, --help            help for user
      --id string       User ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `create-user` — Create User

Create a new user.

```bash
# Create a user
go-dci create-user --name "jdoe" --email "jdoe@example.com" --fullname "John Doe" --team-id <team-id> --password "secure123"

# Output as JSON
go-dci create-user --name "jdoe" --email "jdoe@example.com" --fullname "John Doe" --team-id <team-id> --password "secure123" --output json
```

```
Usage:
  dci create-user [flags]

Flags:
      --email string      User email address (required)
      --fullname string   User full name (required)
  -h, --help              help for create-user
      --name string       Username (required)
  -o, --output string     Output format (json) - default is stdout (default "stdout")
      --password string   User password (required)
      --team-id string    Team ID (required)
```

### `update-user` — Update User

Update an existing user.

```bash
# Update user email
go-dci update-user --id <user-id> --email "newemail@example.com"

# Update user full name
go-dci update-user --id <user-id> --fullname "Jane Doe"

# Output as JSON
go-dci update-user --id <user-id> --email "newemail@example.com" --output json
```

```
Usage:
  dci update-user [flags]

Flags:
      --email string      New email address
      --fullname string   New full name
  -h, --help              help for update-user
      --id string         User ID to update (required)
      --name string       New username
  -o, --output string     Output format (json) - default is stdout (default "stdout")
      --password string   New password
```

### `delete-user` — Delete User

Delete a user from DCI.

```bash
# Delete a user
go-dci delete-user --id <user-id>

# Output as JSON
go-dci delete-user --id <user-id> --output json
```

```
Usage:
  dci delete-user [flags]

Flags:
  -h, --help             help for delete-user
      --id string        User ID to delete (required)
  -o, --output string    Output format (json) - default is stdout (default "stdout")
```

## Products

### `products` — List Products

Get all products.

```bash
# List all products
go-dci products

# Output as JSON
go-dci products --output json
```

```
Usage:
  dci products [flags]

Flags:
  -h, --help            help for products
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```

### `product` — Get Product by ID

Get details for a specific product.

```bash
# Get a product by ID
go-dci product --id <product-id>

# Output as JSON
go-dci product --id <product-id> --output json
```

```
Usage:
  dci product [flags]

Flags:
  -h, --help            help for product
      --id string       Product ID (required)
  -o, --output string   Output format (json) - default is stdout (default "stdout")
```
