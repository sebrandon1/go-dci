package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// captureStdout runs fn with os.Stdout redirected and returns what it wrote.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	old := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = old })

	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close write end of pipe: %v", err)
	}
	<-done
	_ = r.Close()

	return buf.String()
}

// withStdin replaces os.Stdin with a file containing contents for the test.
func withStdin(t *testing.T, contents string) {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatalf("create stdin fixture: %v", err)
	}
	if _, err := f.WriteString(contents); err != nil {
		t.Fatalf("write stdin fixture: %v", err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("rewind stdin fixture: %v", err)
	}

	old := os.Stdin
	os.Stdin = f
	t.Cleanup(func() {
		os.Stdin = old
		_ = f.Close()
	})
}

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

func TestValidateResourceID(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		resourceType string
		wantErr      string
	}{
		{"valid uuid", "550e8400-e29b-41d4-a716-446655440000", "job", ""},
		{"valid uuid another type", "550e8400-e29b-41d4-a716-446655440000", "topic", ""},
		{"empty id", "", "job", "job ID is required"},
		{"uppercase hex", "550E8400-E29B-41D4-A716-446655440000", "topic", "invalid topic ID format"},
		{"missing hyphens", "550e8400e29b41d4a716446655440000", "component", "invalid component ID format"},
		{"trailing junk", "550e8400-e29b-41d4-a716-446655440000-extra", "team", "invalid team ID format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceID(tt.id, tt.resourceType)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestConfirmDeletion(t *testing.T) {
	tests := []struct {
		name         string
		outputFormat string
		yes          bool
		stdin        string
		want         bool
	}{
		{"json mode skips prompt", OutputFormatJSON, false, "", true},
		{"yes flag skips prompt", OutputFormatStdout, true, "", true},
		{"interactive yes", OutputFormatStdout, false, "yes\n", true},
		{"interactive y", OutputFormatStdout, false, "y\n", true},
		{"interactive no", OutputFormatStdout, false, "no\n", false},
		{"interactive empty", OutputFormatStdout, false, "\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFormat = tt.outputFormat
			yesFlag = tt.yes
			t.Cleanup(func() { outputFormat = OutputFormatStdout; yesFlag = false })

			if tt.stdin != "" {
				withStdin(t, tt.stdin)
			}

			got, err := confirmDeletion("topic", "550e8400-e29b-41d4-a716-446655440000")
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInitConfigLegacyMigration(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)

	legacy := filepath.Join(workDir, ".go-dci-config.yaml")
	if err := os.WriteFile(legacy, []byte("accesskey: legacy\n"), 0600); err != nil {
		t.Fatalf("write legacy config: %v", err)
	}
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	previousConfigFile := configFile
	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
		configFile = previousConfigFile
	})

	initConfig()

	assert.Equal(t, ".go-dci-config.yaml", configFile)
	assert.Equal(t, "legacy", viper.GetString("accesskey"))
}

func TestInitConfigCreatesFreshConfig(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)

	home := t.TempDir()
	configHome := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", configHome)

	previousConfigFile := configFile
	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
		configFile = previousConfigFile
	})

	initConfig()

	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("config file was not created at %s: %v", configFile, err)
	}
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	dirInfo, err := os.Stat(filepath.Dir(configFile))
	if err != nil {
		t.Fatalf("config directory was not created: %v", err)
	}
	assert.Equal(t, os.FileMode(0700), dirInfo.Mode().Perm())
}

func TestConfirmDeletion_StdinError(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatalf("create stdin fixture: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close stdin fixture: %v", err)
	}

	old := os.Stdin
	os.Stdin = f
	t.Cleanup(func() { os.Stdin = old })

	outputFormat = OutputFormatStdout
	yesFlag = false
	t.Cleanup(func() { outputFormat = OutputFormatStdout; yesFlag = false })

	confirmed, err := confirmDeletion("topic", "550e8400-e29b-41d4-a716-446655440000")
	assert.Error(t, err)
	assert.False(t, confirmed)
}

func TestInitConfigUnwritableConfigDir(t *testing.T) {
	t.Chdir(t.TempDir())

	configRoot := t.TempDir()
	t.Setenv("HOME", configRoot)
	t.Setenv("XDG_CONFIG_HOME", configRoot)

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("resolve user config dir: %v", err)
	}
	if err := os.MkdirAll(userConfigDir, 0700); err != nil {
		t.Fatalf("create user config dir: %v", err)
	}
	// A read-only parent makes MkdirAll and WriteFile fail, exercising the
	// warning branches.
	// #nosec G302 -- this is a directory mode, and only for the test duration.
	if err := os.Chmod(userConfigDir, 0500); err != nil {
		t.Fatalf("make user config dir read-only: %v", err)
	}
	// #nosec G302 -- restore the directory mode so t.TempDir cleanup can run.
	t.Cleanup(func() { _ = os.Chmod(userConfigDir, 0700) })

	previousConfigFile := configFile
	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
		configFile = previousConfigFile
	})

	assert.NotPanics(t, initConfig)
}
