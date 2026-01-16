// Example: Component Query with go-dci Library
//
// This example demonstrates how to query and analyze DCI components:
// - Listing all components across topics
// - Filtering components by topic
// - Analyzing component versions
// - Finding specific component types
//
// Usage:
//
//	export GO_DCI_ACCESSKEY="your-access-key"
//	export GO_DCI_SECRETKEY="your-secret-key"
//	go run main.go [--topic-id <topic-uuid>] [--type <component-type>]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
)

func main() {
	// Parse command line arguments
	topicID := flag.String("topic-id", "", "Filter by topic ID (optional)")
	componentType := flag.String("type", "", "Filter by component type (e.g., 'ocp', 'certsuite')")
	showVersions := flag.Bool("versions", false, "Show version breakdown")
	flag.Parse()

	// Get credentials
	accessKey := os.Getenv("GO_DCI_ACCESSKEY")
	secretKey := os.Getenv("GO_DCI_SECRETKEY")

	if accessKey == "" || secretKey == "" {
		log.Fatal("GO_DCI_ACCESSKEY and GO_DCI_SECRETKEY environment variables are required")
	}

	// Initialize client
	client := lib.NewClient(accessKey, secretKey)

	fmt.Println("=== DCI Component Analysis ===")

	// Verify authentication
	identity, err := client.GetIdentity()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	fmt.Printf("Authenticated as: %s\n\n", identity.Identity.Name)

	// Get component types
	fmt.Println("1. Available Component Types")
	fmt.Println("   " + strings.Repeat("-", 40))

	componentTypesResp, err := client.GetComponentTypes()
	if err != nil {
		log.Fatalf("Failed to get component types: %v", err)
	}

	typeMap := make(map[string]string) // name -> ID
	for _, resp := range componentTypesResp {
		for _, ct := range resp.ComponentTypes {
			typeMap[ct.Name] = ct.ID
			marker := " "
			if *componentType != "" && ct.Name == *componentType {
				marker = "*"
			}
			fmt.Printf("   %s %s (ID: %s)\n", marker, ct.Name, ct.ID[:8]+"...")
		}
	}
	fmt.Println()

	// Get topics
	fmt.Println("2. Available Topics")
	fmt.Println("   " + strings.Repeat("-", 40))

	topicsResp, err := client.GetTopics()
	if err != nil {
		log.Fatalf("Failed to get topics: %v", err)
	}

	type topicInfo struct {
		ID   string
		Name string
	}
	var topics []topicInfo
	for _, resp := range topicsResp {
		for _, topic := range resp.Topics {
			if topic.State == "active" {
				topics = append(topics, topicInfo{ID: topic.ID, Name: topic.Name})
				marker := " "
				if *topicID != "" && topic.ID == *topicID {
					marker = "*"
				}
				fmt.Printf("   %s %s (ID: %s)\n", marker, topic.Name, topic.ID[:8]+"...")
			}
		}
	}
	fmt.Println()

	// Query components
	fmt.Println("3. Components Analysis")
	fmt.Println("   " + strings.Repeat("-", 40))

	var allComponents []struct {
		Name      string
		Type      string
		Version   string
		TopicID   string
		TopicName string
	}

	// If topic ID is specified, only get components for that topic
	if *topicID != "" {
		componentsResp, err := client.GetTopicComponents(*topicID)
		if err != nil {
			log.Fatalf("Failed to get components: %v", err)
		}

		// Find topic name
		topicName := *topicID
		for _, t := range topics {
			if t.ID == *topicID {
				topicName = t.Name
				break
			}
		}

		for _, resp := range componentsResp {
			for _, comp := range resp.Components {
				allComponents = append(allComponents, struct {
					Name      string
					Type      string
					Version   string
					TopicID   string
					TopicName string
				}{
					Name:      comp.Name,
					Type:      comp.Type,
					Version:   comp.Version,
					TopicID:   *topicID,
					TopicName: topicName,
				})
			}
		}
	} else {
		// Get all components
		componentsResp, err := client.GetComponents()
		if err != nil {
			log.Fatalf("Failed to get components: %v", err)
		}

		for _, resp := range componentsResp {
			for _, comp := range resp.Components {
				allComponents = append(allComponents, struct {
					Name      string
					Type      string
					Version   string
					TopicID   string
					TopicName string
				}{
					Name:    comp.Name,
					Type:    comp.Type,
					Version: comp.Version,
					TopicID: comp.TopicID,
				})
			}
		}
	}

	// Filter by component type if specified
	if *componentType != "" {
		var filtered []struct {
			Name      string
			Type      string
			Version   string
			TopicID   string
			TopicName string
		}
		for _, comp := range allComponents {
			if comp.Type == *componentType {
				filtered = append(filtered, comp)
			}
		}
		allComponents = filtered
	}

	fmt.Printf("   Total components: %d\n\n", len(allComponents))

	// Show component breakdown by type
	typeCount := make(map[string]int)
	for _, comp := range allComponents {
		typeCount[comp.Type]++
	}

	fmt.Println("   By Type:")
	for ctype, count := range typeCount {
		fmt.Printf("   - %s: %d\n", ctype, count)
	}
	fmt.Println()

	// Show version analysis if requested
	if *showVersions {
		fmt.Println("4. Version Analysis")
		fmt.Println("   " + strings.Repeat("-", 40))

		versionCount := make(map[string]map[string]int)
		for _, comp := range allComponents {
			if versionCount[comp.Type] == nil {
				versionCount[comp.Type] = make(map[string]int)
			}
			versionCount[comp.Type][comp.Version]++
		}

		for ctype, versions := range versionCount {
			fmt.Printf("\n   %s versions:\n", ctype)

			// Sort versions
			var versionList []string
			for v := range versions {
				versionList = append(versionList, v)
			}
			sort.Strings(versionList)

			// Show last 10 versions
			start := 0
			if len(versionList) > 10 {
				start = len(versionList) - 10
				fmt.Printf("   ... %d older versions\n", start)
			}
			for i := start; i < len(versionList); i++ {
				v := versionList[i]
				fmt.Printf("   - %s: %d component(s)\n", v, versions[v])
			}
		}
		fmt.Println()
	}

	// Show sample components
	fmt.Println("5. Sample Components")
	fmt.Println("   " + strings.Repeat("-", 40))

	shown := 0
	for _, comp := range allComponents {
		if shown >= 10 {
			fmt.Printf("\n   ... and %d more components\n", len(allComponents)-10)
			break
		}
		fmt.Printf("   - %s\n", comp.Name)
		fmt.Printf("     Type: %s, Version: %s\n", comp.Type, comp.Version)
		shown++
	}

	fmt.Println("\n=== Analysis Complete ===")
}
