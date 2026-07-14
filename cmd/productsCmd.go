package cmd

import (
	"encoding/json"
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting products...")

		responses, err := client.GetProducts(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get products: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printProductsJSON(responses)
		}

		printProductsStdout(responses)

		return nil
	},
}

var getProductCmd = &cobra.Command{
	Use:   "product",
	Short: "Get a specific product by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(getProductIDFlag, "product"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting product with ID: %s\n", getProductIDFlag)

		response, err := client.GetProduct(cmd.Context(), getProductIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printProductJSON(response)
		}

		printProductStdout(response)

		return nil
	},
}

func printProductsStdout(responses []lib.ProductsResponse) {
	total := 0
	for _, resp := range responses {
		total += len(resp.Products)
	}
	if total == 0 {
		fmt.Println("No products found.")
		return
	}
	fmt.Println("---")
	for _, resp := range responses {
		for _, product := range resp.Products {
			fmt.Printf("ID: %s | Name: %s | Label: %s | State: %s\n",
				product.ID, product.Name, product.Label, product.State)
		}
	}
	fmt.Printf("Total Products: %d\n", total)
}

func printProductsJSON(responses []lib.ProductsResponse) error {
	var all []lib.Product
	for _, resp := range responses {
		all = append(all, resp.Products...)
	}
	jsonBytes, err := json.Marshal(map[string]any{
		"products": all,
		"total":    len(all),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
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

func printProductJSON(response *lib.ProductResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getProductsCmd)
	rootCmd.AddCommand(getProductCmd)

	// get products flags

	// get product flags
	getProductCmd.PersistentFlags().StringVar(&getProductIDFlag, "id", "", "Product ID")
	_ = getProductCmd.MarkPersistentFlagRequired("id")
}
