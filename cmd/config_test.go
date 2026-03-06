package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig(t *testing.T) func() {
	// Create a temporary directory for test config
	tmpDir, err := os.MkdirTemp("", "go-dci-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Save original config state
	originalConfigFile := viper.ConfigFileUsed()

	// Set up test config file
	testConfigFile := filepath.Join(tmpDir, ".go-dci-config.yaml")
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(testConfigFile)

	// Create the config file
	if err := viper.WriteConfig(); err != nil {
		// WriteConfig fails if file doesn't exist, use SafeWriteConfig
		if err := viper.SafeWriteConfig(); err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
	}

	// Return cleanup function
	return func() {
		_ = os.RemoveAll(tmpDir)
		viper.Reset()
		if originalConfigFile != "" {
			viper.SetConfigFile(originalConfigFile)
		}
	}
}

func TestGetConfigValue(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	// Set a value and retrieve it
	viper.Set("testkey", "testvalue")

	result := GetConfigValue("testkey")
	assert.Equal(t, "testvalue", result)
}

func TestGetConfigValue_NotSet(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	result := GetConfigValue("nonexistent")
	assert.Equal(t, "", result)
}

func TestUpdateConfigValue(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	err := UpdateConfigValue("updatekey", "updatevalue")
	assert.NoError(t, err)

	// Verify the value was set
	assert.Equal(t, "updatevalue", viper.GetString("updatekey"))
}

