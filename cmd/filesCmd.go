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
	getFileIDFlag     string
	getFileOutputPath string
	deleteFileIDFlag  string
)

var getFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Download a file by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Downloading file with ID: %s\n", getFileIDFlag)

		content, contentType, err := client.GetFile(cmd.Context(), getFileIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get file: %v", err)
		}

		if getFileOutputPath != "" {
			if err := os.WriteFile(getFileOutputPath, content, 0644); err != nil {
				return fmt.Errorf("failed to write file: %v", err)
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

		return nil
	},
}

var deleteFileCmd = &cobra.Command{
	Use:   "delete-file",
	Short: "Delete a file from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete file: id=%s\n", deleteFileIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("file", deleteFileIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Deleting file: %s\n", deleteFileIDFlag)

		err = client.DeleteFile(cmd.Context(), deleteFileIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete file: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteFileIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("File deleted successfully!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getFileCmd)
	rootCmd.AddCommand(deleteFileCmd)

	// get file flags
	getFileCmd.PersistentFlags().StringVar(&getFileIDFlag, "id", "", "File ID")
	_ = getFileCmd.MarkPersistentFlagRequired("id")
	getFileCmd.PersistentFlags().StringVar(&getFileOutputPath, "output-path", "", "Path to save the file content")
	getFileCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete file flags
	deleteFileCmd.PersistentFlags().StringVar(&deleteFileIDFlag, "id", "", "File ID to delete")
	_ = deleteFileCmd.MarkPersistentFlagRequired("id")
	deleteFileCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
