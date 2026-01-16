package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for product command flags
var (
	getProductIDFlag string
)

var getProductsCmd = &cobra.Command{
	Use:   "products",
	Short: "Get all products from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Println("Getting products...")
		}

		response, err := client.GetProducts()
		if err != nil {
			fmt.Printf("Failed to get products: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printProductsJSON(response)
		} else {
			printProductsStdout(response)
		}
	},
}

var getProductCmd = &cobra.Command{
	Use:   "product",
	Short: "Get a specific product by ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getProductIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting product with ID: %s\n", getProductIDFlag)
		}

		response, err := client.GetProduct(getProductIDFlag)
		if err != nil {
			fmt.Printf("Failed to get product: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printProductJSON(response)
		} else {
			printProductStdout(response)
		}
	},
}

func printProductsStdout(response *lib.ProductsResponse) {
	if len(response.Products) == 0 {
		fmt.Println("No products found.")
		return
	}
	fmt.Println("---")
	for _, product := range response.Products {
		fmt.Printf("ID: %s | Name: %s | Label: %s | State: %s\n",
			product.ID, product.Name, product.Label, product.State)
	}
	fmt.Printf("Total Products: %d\n", len(response.Products))
}

func printProductsJSON(response *lib.ProductsResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printProductStdout(response *lib.ProductResponse) {
	fmt.Println("---")
	fmt.Printf("ID:          %s\n", response.Product.ID)
	fmt.Printf("Name:        %s\n", response.Product.Name)
	fmt.Printf("Label:       %s\n", response.Product.Label)
	fmt.Printf("Description: %s\n", response.Product.Description)
	fmt.Printf("State:       %s\n", response.Product.State)
	fmt.Printf("Created:     %s\n", response.Product.CreatedAt)
	fmt.Printf("Updated:     %s\n", response.Product.UpdatedAt)
}

func printProductJSON(response *lib.ProductResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getProductsCmd)
	rootCmd.AddCommand(getProductCmd)

	// get products flags
	getProductsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get product flags
	getProductCmd.PersistentFlags().StringVar(&getProductIDFlag, "id", "", "Product ID (required)")
	getProductCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
