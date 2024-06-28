package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetConfigValue(key, value string) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func UpdateConfigValue(key, value string) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or get configuration values",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a key value pair to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")
		SetConfigValue(key, value)
	},
}

var accessKeyCmd = &cobra.Command{
	Use:   "accesskey",
	Short: "Set the access key for the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		value, _ := cmd.Flags().GetString("value")
		UpdateConfigValue(key, value)
	},
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset a key value pair from the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		viper.Set(key, nil)
		err := viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	// Add flags to the access key command
	accessKeyCmd.Flags()

	// Add the commands to the config command
	setCmd.AddCommand(accessKeyCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(unsetCmd)

	// Add config to root command
	rootCmd.AddCommand(configCmd)

	// configCmd.PersistentFlags().StringP("key", "k", "", "The key for the key value set to add to the configuration.")
	// configCmd.PersistentFlags().StringP("value", "v", "", "The value for the key value set to add to the configuration.")
}
