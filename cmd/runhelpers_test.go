package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// newTestClient points the package-level client at an httptest server running
// handler and restores the globals when the test finishes.
func newTestClient(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(handler)
	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"

	t.Cleanup(func() {
		server.Close()
		dciClient = nil
		resetGlobalFlags()
	})

	return server
}

// runCommand invokes a command's RunE directly with the given command as its
// own context, bypassing Cobra's flag-required and PersistentPreRunE checks.
func runCommand(command *cobra.Command, args ...string) error {
	command.SetContext(context.Background())
	return command.RunE(command, args)
}

// resetGlobalFlags restores the package-level flag variables to their defaults.
func resetGlobalFlags() {
	outputFormat = OutputFormatStdout
	yesFlag = false
	quietFlag = false
	verboseFlag = false
	dryRunFlag = false
}
