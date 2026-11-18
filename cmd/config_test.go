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

func TestMaskCredential(t *testing.T) {
	assert.Equal(t, "", maskCredential(""))
	assert.Equal(t, "****", maskCredential("abc"))
	assert.Equal(t, "****", maskCredential("abcd"))
	assert.Equal(t, "****fghi", maskCredential("abcdefghi"))
	assert.Equal(t, "****cdef", maskCredential("abcdef"))
}

func TestSetConfigCmd_Success(t *testing.T) {
	t.Cleanup(setupTestConfig(t))

	if err := setCmd.ParseFlags([]string{"--accesskey", "myaccesskey", "--secretkey", "mysecretkey"}); err != nil {
		t.Fatalf("parse set flags: %v", err)
	}
	t.Cleanup(func() {
		_ = setCmd.Flags().Set(configKeyAccessKey, "")
		_ = setCmd.Flags().Set(configKeySecretKey, "")
	})

	assert.NoError(t, runCommand(setCmd))
	assert.Equal(t, "myaccesskey", viper.GetString(configKeyAccessKey))
	assert.Equal(t, "mysecretkey", viper.GetString(configKeySecretKey))
}

func TestSetConfigCmd_WriteError(t *testing.T) {
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(t.TempDir(), "absent", "config.yaml"))
	t.Cleanup(viper.Reset)

	if err := setCmd.ParseFlags([]string{"--accesskey", "myaccesskey"}); err != nil {
		t.Fatalf("parse set flags: %v", err)
	}
	t.Cleanup(func() { _ = setCmd.Flags().Set(configKeyAccessKey, "") })

	err := runCommand(setCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "setting accesskey")
}

func TestUnsetConfigCmd_Success(t *testing.T) {
	t.Cleanup(setupTestConfig(t))
	viper.Set(configKeyAccessKey, "myaccesskey")

	if err := unsetCmd.ParseFlags([]string{"--key", configKeyAccessKey}); err != nil {
		t.Fatalf("parse unset flags: %v", err)
	}
	t.Cleanup(func() { _ = unsetCmd.Flags().Set("key", "") })

	assert.NoError(t, runCommand(unsetCmd))
	assert.False(t, viper.IsSet(configKeyAccessKey))
}

func TestUnsetConfigCmd_MissingKey(t *testing.T) {
	t.Cleanup(setupTestConfig(t))

	if err := unsetCmd.ParseFlags([]string{"--key", "notpresent"}); err != nil {
		t.Fatalf("parse unset flags: %v", err)
	}
	t.Cleanup(func() { _ = unsetCmd.Flags().Set("key", "") })

	assert.NoError(t, runCommand(unsetCmd))
}

func TestViewConfigCmd_MasksCredentials(t *testing.T) {
	t.Cleanup(setupTestConfig(t))
	viper.Set(configKeyAccessKey, "abcdefghij")
	viper.Set(configKeySecretKey, "secretvalue")

	assert.NoError(t, runCommand(viewCmd))
}

func TestFormatConfigValue(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		value    string
		fallback string
		want     string
	}{
		{"value present", "Access Key", "mykey", "(not set)", "  Access Key:  mykey\n"},
		{"empty value uses fallback", "Secret Key", "", "(not set)", "  Secret Key:  (not set)\n"},
		{"short label", "ID", "abc", "-", "  ID:          abc\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatConfigValue(tt.label, tt.value, tt.fallback))
		})
	}
}

func TestWriteConfigSecure_SetsFileMode(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("{}\n"), 0600); err != nil {
		t.Fatalf("write config fixture: %v", err)
	}

	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	viper.Set(configKeyAccessKey, "myaccesskey")
	t.Cleanup(viper.Reset)

	if err := writeConfigSecure(); err != nil {
		t.Fatalf("writeConfigSecure returned an error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat config file: %v", err)
	}
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	assert.Equal(t, "myaccesskey", viper.GetString(configKeyAccessKey))
}

func TestWriteConfigSecure_WriteError(t *testing.T) {
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(t.TempDir(), "absent", "config.yaml"))
	t.Cleanup(viper.Reset)

	assert.Error(t, writeConfigSecure())
}
