package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for file command flags
var (
	getFileIDFlag    string
	getFileOutputPath string
	deleteFileIDFlag string
)

var getFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Download a file by ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getFileIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Downloading file with ID: %s\n", getFileIDFlag)
		}

		content, contentType, err := client.GetFile(getFileIDFlag)
		if err != nil {
			fmt.Printf("Failed to get file: %v\n", err)
			return
		}

		if getFileOutputPath != "" {
			if err := os.WriteFile(getFileOutputPath, content, 0644); err != nil {
				fmt.Printf("Failed to write file: %v\n", err)
				return
			}
			fmt.Printf("File saved to: %s\n", getFileOutputPath)
		} else if outputFormat == OutputFormatJSON {
			result := map[string]any{
				"id":           getFileIDFlag,
				"content_type": contentType,
				"size":         len(content),
			}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("File ID:       %s\n", getFileIDFlag)
			fmt.Printf("Content-Type:  %s\n", contentType)
			fmt.Printf("Size:          %d bytes\n", len(content))
			fmt.Println("Use --output <path> to save the file content")
		}
	},
}

var deleteFileCmd = &cobra.Command{
	Use:   "delete-file",
	Short: "Delete a file from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if deleteFileIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting file: %s\n", deleteFileIDFlag)
		}

		err = client.DeleteFile(deleteFileIDFlag)
		if err != nil {
			fmt.Printf("Failed to delete file: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteFileIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("File deleted successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(getFileCmd)
	rootCmd.AddCommand(deleteFileCmd)

	// get file flags
	getFileCmd.PersistentFlags().StringVar(&getFileIDFlag, "id", "", "File ID (required)")
	getFileCmd.PersistentFlags().StringVar(&getFileOutputPath, "output-path", "", "Path to save the file content")
	getFileCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete file flags
	deleteFileCmd.PersistentFlags().StringVar(&deleteFileIDFlag, "id", "", "File ID to delete (required)")
	deleteFileCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
