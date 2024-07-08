package cmd

import (
	"fmt"
	"log"

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

func ConfigKeyValuePairAdd(key string, value string) {
	if validateKeyValuePair(key, value) {
		log.Printf("Validation not met for %s.", key)
	} else {
		writeKeyValuePair(key, value)
	}
}

func validateKeyValuePair(key string, value string) bool {
	if len(key) == 0 || len(value) == 0 {
		fmt.Println("The key and value must both contain contents to write to the configuration file.")
		return true
	}
	// Determine if an existing key, if so notify.
	if findExistingKey(key) {
		fmt.Println("This key already exists. Create a key value pair with a different key, or if this is an update use the update command.")
		return true
	}
	return false
}

func findExistingKey(theKey string) bool {
	existingKey := false
	for i := 0; i < len(viper.AllKeys()); i++ {
		if viper.AllKeys()[i] == theKey {
			existingKey = true
		}
	}
	return existingKey
}

func ConfigKeyValuePairUpdate(key string, value string) {
	writeKeyValuePair(key, value)
}

func writeKeyValuePair(key string, value interface{}) {
	viper.Set(key, value)
	viper.WriteConfig()
	fmt.Printf("Wrote the %s pair.\n", key)
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
		accesskey, _ := cmd.Flags().GetString("accesskey")
		secretkey, _ := cmd.Flags().GetString("secretkey")

		if accesskey != "" {
			UpdateConfigValue("accesskey", accesskey)
		}

		if secretkey != "" {
			UpdateConfigValue("secretkey", secretkey)
		}
	},
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset a key value pair from the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("")
		viper.Set(key, nil)
		err := viper.WriteConfig()
		if err != nil {
			panic(err)
		}
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

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(unsetCmd)
	configCmd.AddCommand(viewCmd)

	// Add config to root command
	rootCmd.AddCommand(configCmd)
}
