# Component Analysis with go-dci

This tutorial covers how to query, filter, and analyze DCI components.

## Prerequisites

- Completed [Getting Started](./01-getting-started.md)

## Understanding Components

Components represent test artifacts in DCI:

| Type | Description |
|------|-------------|
| `ocp` | OpenShift Container Platform version |
| `certsuite` | Red Hat Certification Test Suite |
| `rhel` | Red Hat Enterprise Linux version |

Each component has:
- **ID** - Unique identifier
- **Name** - Human-readable name
- **Type** - Component type (ocp, certsuite, etc.)
- **Version** - Version string
- **TopicID** - Parent topic

## Listing All Components

```go
package main

import (
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

    // Get all components
    componentsResp, err := client.GetComponents()
    if err != nil {
        log.Fatalf("Failed to get components: %v", err)
    }

    totalComponents := 0
    for _, resp := range componentsResp {
        totalComponents += len(resp.Components)
        for _, comp := range resp.Components {
            fmt.Printf("%s - %s (%s)\n", comp.Name, comp.Version, comp.Type)
        }
    }

    fmt.Printf("\nTotal: %d components\n", totalComponents)
}
```

## Filtering by Topic

Get components for a specific certification topic:

```go
topicID := "your-topic-uuid"

componentsResp, err := client.GetTopicComponents(topicID)
if err != nil {
    log.Fatalf("Failed to get components: %v", err)
}

for _, resp := range componentsResp {
    for _, comp := range resp.Components {
        fmt.Printf("%s v%s\n", comp.Name, comp.Version)
    }
}
```

## Filtering by Topic ID (All Components)

```go
topicID := "your-topic-uuid"

componentsResp, err := client.GetComponentsByTopicID(topicID)
if err != nil {
    log.Fatalf("Failed to get components: %v", err)
}

for _, resp := range componentsResp {
    for _, comp := range resp.Components {
        fmt.Printf("%s (%s)\n", comp.Name, comp.Type)
    }
}
```

## Getting a Specific Component

```go
componentID := "component-uuid"

component, err := client.GetComponent(componentID)
if err != nil {
    log.Fatalf("Failed to get component: %v", err)
}

fmt.Printf("Name: %s\n", component.Component.Name)
fmt.Printf("Type: %s\n", component.Component.Type)
fmt.Printf("Version: %s\n", component.Component.Version)
fmt.Printf("Topic ID: %s\n", component.Component.TopicID)
fmt.Printf("Created: %s\n", component.Component.CreatedAt)
```

## Component Types

List and understand available component types:

```go
typesResp, err := client.GetComponentTypes()
if err != nil {
    log.Fatalf("Failed to get types: %v", err)
}

for _, resp := range typesResp {
    for _, ct := range resp.ComponentTypes {
        fmt.Printf("Type: %s\n", ct.Name)
        fmt.Printf("  ID: %s\n", ct.ID)
        fmt.Printf("  State: %s\n", ct.State)
    }
}
```

## Version Analysis

Analyze component versions across topics:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "sort"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    client := lib.NewClient(
        os.Getenv("GO_DCI_ACCESSKEY"),
        os.Getenv("GO_DCI_SECRETKEY"),
    )

    // Get all components
    componentsResp, err := client.GetComponents()
    if err != nil {
        log.Fatal(err)
    }

    // Group by type and version
    versionsByType := make(map[string]map[string]int)

    for _, resp := range componentsResp {
        for _, comp := range resp.Components {
            if versionsByType[comp.Type] == nil {
                versionsByType[comp.Type] = make(map[string]int)
            }
            versionsByType[comp.Type][comp.Version]++
        }
    }

    // Print analysis
    for compType, versions := range versionsByType {
        fmt.Printf("\n%s versions:\n", compType)

        // Sort versions
        var versionList []string
        for v := range versions {
            versionList = append(versionList, v)
        }
        sort.Strings(versionList)

        for _, v := range versionList {
            fmt.Printf("  %s: %d\n", v, versions[v])
        }
    }
}
```

## OCP Version Statistics

Analyze OpenShift versions in certification jobs:

```go
package main

import (
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

    // Get recent jobs
    jobsResp, err := client.GetJobs(30) // Last 30 days
    if err != nil {
        log.Fatal(err)
    }

    // Count OCP versions
    ocpCount := make(map[string]int)

    for _, resp := range jobsResp {
        for _, job := range resp.Jobs {
            // Each job may have OCP version in components
            for _, compRef := range job.Components {
                // Get component details
                comp, err := client.GetComponent(compRef.ID)
                if err != nil {
                    continue
                }
                if comp.Component.Type == "ocp" {
                    ocpCount[comp.Component.Version]++
                }
            }
        }
    }

    fmt.Println("OCP Version Distribution (Last 30 Days):")
    for version, count := range ocpCount {
        fmt.Printf("  %s: %d jobs\n", version, count)
    }
}
```

## Creating Components

Create a new component (admin only):

```go
component, err := client.CreateComponent(
    "My Component",      // name
    "ocp",              // type
    "topic-uuid",       // topic ID
    "4.15.0",           // version
)
if err != nil {
    log.Fatalf("Failed to create component: %v", err)
}

fmt.Printf("Created component: %s\n", component.Component.ID)
```

## Updating Components

Update component metadata:

```go
updated, err := client.UpdateComponent("component-uuid", lib.UpdateComponentRequest{
    Name:    "Updated Name",
    Version: "4.15.1",
})
if err != nil {
    log.Fatalf("Failed to update: %v", err)
}

fmt.Printf("Updated: %s v%s\n", updated.Component.Name, updated.Component.Version)
```

## Deleting Components

Remove a component (admin only):

```go
err := client.DeleteComponent("component-uuid")
if err != nil {
    log.Fatalf("Failed to delete: %v", err)
}
fmt.Println("Component deleted")
```

## Component Type Management

### Get a Specific Type

```go
compType, err := client.GetComponentType("type-uuid")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Type: %s\n", compType.ComponentType.Name)
```

### Create a Type (Admin)

```go
newType, err := client.CreateComponentType("my-custom-type")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created type: %s\n", newType.ComponentType.Name)
```

### Update a Type

```go
updated, err := client.UpdateComponentType("type-uuid", lib.UpdateComponentTypeRequest{
    Name: "updated-name",
})
```

### Delete a Type

```go
err := client.DeleteComponentType("type-uuid")
```

## Practical Example: Find Latest Versions

```go
package main

import (
    "fmt"
    "log"
    "os"
    "sort"

    "github.com/sebrandon1/go-dci/lib"
)

func main() {
    client := lib.NewClient(
        os.Getenv("GO_DCI_ACCESSKEY"),
        os.Getenv("GO_DCI_SECRETKEY"),
    )

    // Get components for a topic
    topicID := "your-topic-uuid"
    componentsResp, err := client.GetTopicComponents(topicID)
    if err != nil {
        log.Fatal(err)
    }

    // Group by type, find latest
    latestByType := make(map[string]struct {
        Version string
        ID      string
    })

    for _, resp := range componentsResp {
        for _, comp := range resp.Components {
            current := latestByType[comp.Type]
            if current.Version == "" || comp.Version > current.Version {
                latestByType[comp.Type] = struct {
                    Version string
                    ID      string
                }{
                    Version: comp.Version,
                    ID:      comp.ID,
                }
            }
        }
    }

    fmt.Println("Latest components for topic:")
    for compType, info := range latestByType {
        fmt.Printf("  %s: %s (ID: %s)\n", compType, info.Version, info.ID[:8])
    }
}
```

## Complete Example

See the [component-query example](../../examples/component-query/main.go) for a complete working program.

## Best Practices

1. **Cache Component Data** - Component lists are relatively static; cache when appropriate
2. **Use Topic Filtering** - Filter by topic ID to reduce API calls
3. **Handle Pagination** - The library handles pagination internally for large result sets
4. **Version Comparison** - Use semantic versioning libraries for accurate comparisons
