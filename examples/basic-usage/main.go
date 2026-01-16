// Example: Basic Usage of go-dci Library
//
// This example demonstrates fundamental operations with the go-dci library:
// - Initializing the client with AWS SigV4 authentication
// - Verifying authentication with identity endpoint
// - Listing topics and their components
// - Querying jobs
//
// Usage:
//
//	export GO_DCI_ACCESSKEY="your-access-key"
//	export GO_DCI_SECRETKEY="your-secret-key"
//	go run main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sebrandon1/go-dci/lib"
)

func main() {
	// Step 1: Get credentials from environment
	accessKey := os.Getenv("GO_DCI_ACCESSKEY")
	secretKey := os.Getenv("GO_DCI_SECRETKEY")

	if accessKey == "" || secretKey == "" {
		log.Fatal("GO_DCI_ACCESSKEY and GO_DCI_SECRETKEY environment variables are required")
	}

	// Step 2: Initialize the client
	client := lib.NewClient(accessKey, secretKey)

	fmt.Println("=== go-dci Basic Usage Example ===")

	// Step 3: Verify authentication by getting identity
	fmt.Println("1. Verifying authentication...")
	identity, err := client.GetIdentity()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	fmt.Printf("   Authenticated as: %s\n", identity.Identity.Name)
	fmt.Printf("   ID: %s\n", identity.Identity.ID)
	fmt.Printf("   Type: %s\n", identity.Identity.Type)
	if identity.Identity.TeamID != "" {
		fmt.Printf("   Team ID: %s\n", identity.Identity.TeamID)
	}
	fmt.Println()

	// Step 4: List available topics
	fmt.Println("2. Listing available topics...")
	topicsResp, err := client.GetTopics()
	if err != nil {
		log.Fatalf("Failed to get topics: %v", err)
	}

	totalTopics := 0
	for _, resp := range topicsResp {
		for _, topic := range resp.Topics {
			totalTopics++
			fmt.Printf("   - %s (ID: %s)\n", topic.Name, topic.ID[:8]+"...")
			fmt.Printf("     State: %s\n", topic.State)
		}
	}
	fmt.Printf("   Total: %d topics\n\n", totalTopics)

	// Step 5: List component types
	fmt.Println("3. Listing component types...")
	componentTypesResp, err := client.GetComponentTypes()
	if err != nil {
		log.Printf("   Could not get component types: %v\n", err)
	} else {
		totalTypes := 0
		for _, resp := range componentTypesResp {
			for _, ct := range resp.ComponentTypes {
				totalTypes++
				fmt.Printf("   - %s (ID: %s)\n", ct.Name, ct.ID[:8]+"...")
			}
		}
		fmt.Printf("   Total: %d component types\n", totalTypes)
	}
	fmt.Println()

	// Step 6: Get recent jobs (last 7 days)
	fmt.Println("4. Getting recent jobs (last 7 days)...")
	jobsResp, err := client.GetJobs(7)
	if err != nil {
		log.Printf("   Could not get jobs: %v\n", err)
	} else {
		totalJobs := 0
		for _, resp := range jobsResp {
			totalJobs += len(resp.Jobs)
		}
		fmt.Printf("   Found %d jobs in the last 7 days\n", totalJobs)

		// Show first 5 jobs
		shown := 0
		for _, resp := range jobsResp {
			for _, job := range resp.Jobs {
				if shown >= 5 {
					break
				}
				fmt.Printf("   - Job %s\n", job.ID[:8]+"...")
				fmt.Printf("     Status: %s\n", job.Status)
				fmt.Printf("     Created: %s\n", job.CreatedAt)
				shown++
			}
		}
		if totalJobs > 5 {
			fmt.Printf("   ... and %d more jobs\n", totalJobs-5)
		}
	}
	fmt.Println()

	// Step 7: Get products (if available)
	fmt.Println("5. Listing products...")
	products, err := client.GetProducts()
	if err != nil {
		log.Printf("   Could not get products: %v\n", err)
	} else if len(products.Products) == 0 {
		fmt.Println("   No products found")
	} else {
		for _, product := range products.Products {
			fmt.Printf("   - %s (ID: %s)\n", product.Name, product.ID[:8]+"...")
		}
	}

	fmt.Println("\n=== Example Complete ===")
}
