package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "dci",
	Short:   "CLI and library for the Red Hat Distributed CI API",
	Version: Version,
}

var (
	configFile  string
	yesFlag     bool
	quietFlag   bool
	verboseFlag bool
	dryRunFlag  bool
)

// printStatus prints a status message unless --quiet or JSON output is enabled.
// Accepts format string and args like fmt.Printf to avoid eager string formatting.
func printStatus(format string, args ...any) {
	if !quietFlag && outputFormat != OutputFormatJSON {
		if len(args) > 0 {
			fmt.Printf(format, args...)
		} else {
			fmt.Println(format)
		}
	}
}

// printVerbose prints a verbose debug message if --verbose is enabled.
// Reserved for future use (HTTP request logging, timing information).
//
//nolint:unused
func printVerbose(message string) {
	if verboseFlag {
		fmt.Printf("[VERBOSE] %s\n", message)
	}
}

// validateResourceID validates that a resource ID is non-empty and matches the UUID format used by DCI.
// Returns an error if validation fails.
func validateResourceID(id, resourceType string) error {
	if id == "" {
		return fmt.Errorf("%s ID is required", resourceType)
	}
	// DCI uses UUIDs
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	if !uuidRegex.MatchString(id) {
		return fmt.Errorf("invalid %s ID format (expected UUID): %s", resourceType, id)
	}
	return nil
}

// confirmDeletion prompts the user to confirm a deletion operation.
// Returns true if the user confirms (or --yes flag is set), false otherwise.
// Skips prompt if output format is JSON (assumes automation).
func confirmDeletion(resourceType, resourceID string) (bool, error) {
	// Skip prompt in JSON mode (automation)
	if outputFormat == OutputFormatJSON {
		return true, nil
	}

	// Skip prompt if --yes flag is set
	if yesFlag {
		return true, nil
	}

	// Interactive confirmation
	fmt.Printf("Are you sure you want to delete %s '%s'? (yes/no): ", resourceType, resourceID)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y", nil
}

func readPassword(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("reading password: %w", err)
	}
	return string(password), nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Skip confirmation prompts")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Show verbose output (HTTP requests, timing)")
	rootCmd.PersistentFlags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would happen without executing")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_DCI")

	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dciConfigDir := filepath.Join(configDir, "go-dci")
	configFile = filepath.Join(dciConfigDir, "config.yaml")

	legacyFile := ".go-dci-config.yaml"
	_, legacyErr := os.Stat(legacyFile)
	_, newErr := os.Stat(configFile)

	if legacyErr == nil && newErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: using legacy config %s\n", legacyFile)
		fmt.Fprintf(os.Stderr, "  Migrate to: %s\n", configFile)
		configFile = legacyFile
	}

	viper.SetConfigFile(configFile)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configFile), 0700); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create config dir: %v\n", err)
		}
		if err := os.WriteFile(configFile, []byte("{}\n"), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create config file: %v\n", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading config file:", err)
	}
}
