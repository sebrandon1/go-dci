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
	RunE: func(cmd *cobra.Command, args []string) error {
		accesskey, _ := cmd.Flags().GetString("accesskey")
		secretkey, _ := cmd.Flags().GetString("secretkey")

		if accesskey != "" {
			if err := UpdateConfigValue("accesskey", accesskey); err != nil {
				return fmt.Errorf("setting accesskey: %v", err)
			}
			printStatus("Access key updated successfully")
		}

		if secretkey != "" {
			if err := UpdateConfigValue("secretkey", secretkey); err != nil {
				return fmt.Errorf("setting secretkey: %v", err)
			}
			printStatus("Secret key updated successfully")
		}

		return nil
	},
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset a key value pair from the configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, _ := cmd.Flags().GetString("key")

		if !viper.IsSet(key) {
			fmt.Printf("Key '%s' does not exist in configuration\n", key)
			return nil
		}

		viper.Set(key, nil)
		err := viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("writing config: %v", err)
		}
		fmt.Printf("Unset key '%s' from configuration\n", key)

		return nil
	},
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Debug()

		return nil
	},
}

func init() {
	// Add flags to the commands
	setCmd.PersistentFlags().StringP("accesskey", "a", "", "The access key to set in the configuration.")
	setCmd.PersistentFlags().StringP("secretkey", "s", "", "The secret key to set in the configuration.")

	unsetCmd.PersistentFlags().StringP("key", "k", "", "The key to unset from the configuration (e.g., accesskey, secretkey)")
	_ = unsetCmd.MarkPersistentFlagRequired("key")

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(unsetCmd)
	configCmd.AddCommand(viewCmd)

	// Add config to root command
	rootCmd.AddCommand(configCmd)
}
