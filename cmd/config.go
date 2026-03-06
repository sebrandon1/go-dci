package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func UpdateConfigValue(key, value string) error {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or get configuration values",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a key value pair to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		accesskey, _ := cmd.Flags().GetString("accesskey")
		secretkey, _ := cmd.Flags().GetString("secretkey")

		if accesskey != "" {
			if err := UpdateConfigValue("accesskey", accesskey); err != nil {
				fmt.Printf("Error setting accesskey: %v\n", err)
				return
			}
			fmt.Println("Access key updated successfully")
		}

		if secretkey != "" {
			if err := UpdateConfigValue("secretkey", secretkey); err != nil {
				fmt.Printf("Error setting secretkey: %v\n", err)
				return
			}
			fmt.Println("Secret key updated successfully")
		}
	},
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset a key value pair from the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			fmt.Println("Error: --key flag is required")
			return
		}

		if !viper.IsSet(key) {
			fmt.Printf("Key '%s' does not exist in configuration\n", key)
			return
		}

		viper.Set(key, nil)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			return
		}
		fmt.Printf("Unset key '%s' from configuration\n", key)
	},
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Debug()
	},
}

func init() {
	// Add flags to the commands
	setCmd.PersistentFlags().StringP("accesskey", "a", "", "The access key to set in the configuration.")
	setCmd.PersistentFlags().StringP("secretkey", "s", "", "The secret key to set in the configuration.")

	unsetCmd.PersistentFlags().StringP("key", "k", "", "The key to unset from the configuration (e.g., accesskey, secretkey)")

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(unsetCmd)
	configCmd.AddCommand(viewCmd)

	// Add config to root command
	rootCmd.AddCommand(configCmd)
}
