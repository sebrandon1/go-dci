package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dci",
	Short: "DCI CLI",
}

var configFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile = ".go-dci-config.yaml"
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_DCI")

	// If the config file is not found, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		viper.WriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: ", err)
	}
}
