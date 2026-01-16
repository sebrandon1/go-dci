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

func TestFindExistingKey(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	// Set a key
	viper.Set("existingkey", "value")

	assert.True(t, findExistingKey("existingkey"))
	assert.False(t, findExistingKey("nonexistentkey"))
}

func TestValidateKeyValuePair(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		name     string
		key      string
		value    string
		expected bool
	}{
		{
			name:     "empty key",
			key:      "",
			value:    "value",
			expected: true, // returns true when validation fails
		},
		{
			name:     "empty value",
			key:      "key",
			value:    "",
			expected: true,
		},
		{
			name:     "valid new key-value pair",
			key:      "newkey",
			value:    "newvalue",
			expected: false, // returns false when validation passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateKeyValuePair(tt.key, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateKeyValuePair_ExistingKey(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	// Set an existing key
	viper.Set("existingkey", "existingvalue")

	// Should fail validation because key exists
	result := validateKeyValuePair("existingkey", "newvalue")
	assert.True(t, result) // true means validation failed
}

func TestUpdateConfigValue(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	err := UpdateConfigValue("updatekey", "updatevalue")
	assert.NoError(t, err)

	// Verify the value was set
	assert.Equal(t, "updatevalue", viper.GetString("updatekey"))
}

func TestSetConfigValue(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	err := SetConfigValue("setkey", "setvalue")
	assert.NoError(t, err)

	// Verify the value was set
	assert.Equal(t, "setvalue", viper.GetString("setkey"))
}

func TestConfigKeyValuePairAdd(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	err := ConfigKeyValuePairAdd("addkey", "addvalue")
	assert.NoError(t, err)

	// Verify the value was set
	assert.Equal(t, "addvalue", viper.GetString("addkey"))
}

func TestConfigKeyValuePairAdd_EmptyKey(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	err := ConfigKeyValuePairAdd("", "value")
	assert.Error(t, err)
}

func TestConfigKeyValuePairUpdate(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	// First add a key
	viper.Set("updateablekey", "oldvalue")
	_ = viper.WriteConfig()

	// Update it
	err := ConfigKeyValuePairUpdate("updateablekey", "newvalue")
	assert.NoError(t, err)

	// Verify
	assert.Equal(t, "newvalue", viper.GetString("updateablekey"))
}

func TestWriteKeyValuePair(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	err := writeKeyValuePair("writekey", "writevalue")
	assert.NoError(t, err)

	// Verify the value was written
	assert.Equal(t, "writevalue", viper.GetString("writekey"))
}
