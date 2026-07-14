package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmdExists(t *testing.T) {
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "dci", rootCmd.Use)
}

func TestRootCmdHasVersion(t *testing.T) {
	assert.NotEmpty(t, rootCmd.Version)
}

func TestRootCmdHasSubcommands(t *testing.T) {
	commands := rootCmd.Commands()
	assert.NotEmpty(t, commands)

	commandNames := make([]string, 0, len(commands))
	for _, cmd := range commands {
		commandNames = append(commandNames, cmd.Name())
	}

	expectedCommands := []string{"config", "file", "delete-file", "topics", "jobs"}
	for _, expected := range expectedCommands {
		assert.Contains(t, commandNames, expected)
	}
}

func TestInitConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		initConfig()
	})
	assert.Contains(t, configFile, "go-dci")
	assert.True(t, strings.HasSuffix(configFile, "config.yaml") || strings.HasSuffix(configFile, ".go-dci-config.yaml"))
}
