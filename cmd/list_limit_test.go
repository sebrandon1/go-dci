package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListLimitFlagsRegisteredOnlyOnListCommands(t *testing.T) {
	listCommands := []*cobra.Command{
		getTopicsCmd,
		getComponentTypesCmd,
		getComponentsCmd,
		getJobsCmd,
		getProductsCmd,
		getTeamsCmd,
		getUsersCmd,
		getRemoteCIsCmd,
		getJobStatesCmd,
		getTopicComponentsCmd,
	}
	for _, command := range listCommands {
		assert.NotNilf(t, command.Flags().Lookup("limit"), "%s should support --limit", command.Name())
	}
	assert.Nil(t, getOcpCountCmd.Flags().Lookup("limit"), "ocpcount should not expose --limit")
}

func TestListLimitValidation(t *testing.T) {
	tests := []struct {
		value   string
		want    []int
		wantErr string
	}{
		{value: "", want: nil},
		{value: "3", want: []int{3}},
		{value: "0", wantErr: "--limit must be a positive integer"},
		{value: "-1", wantErr: "--limit must be a positive integer"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			command := &cobra.Command{Use: "list"}
			addListLimitFlag(command)
			if tt.value != "" {
				require.NoError(t, command.Flags().Set("limit", tt.value))
			}

			limits, err := listLimit(command)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, limits)
		})
	}
}

func TestTopicsCommandPassesLimitToClientAndOutput(t *testing.T) {
	requestedLimit := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedLimit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
		topics := make([]lib.Topic, 8)
		for i := range topics {
			topics[i].ID = strconv.Itoa(i)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(lib.TopicsResponse{Topics: topics})
	}))
	defer server.Close()

	client := lib.NewClient("test-key", "test-secret")
	client.BaseURL = server.URL
	oldClient := dciClient
	dciClient = client
	t.Cleanup(func() { dciClient = oldClient })

	oldOutputFormat := outputFormat
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { outputFormat = oldOutputFormat })
	oldNameFilter := nameFilter
	nameFilter = ""
	t.Cleanup(func() { nameFilter = oldNameFilter })

	limitFlag := getTopicsCmd.Flags().Lookup("limit")
	oldLimit := limitFlag.Value.String()
	oldChanged := limitFlag.Changed
	require.NoError(t, getTopicsCmd.Flags().Set("limit", "3"))
	t.Cleanup(func() {
		_ = getTopicsCmd.Flags().Set("limit", oldLimit)
		limitFlag.Changed = oldChanged
	})

	oldContext := getTopicsCmd.Context()
	getTopicsCmd.SetContext(context.Background())
	t.Cleanup(func() { getTopicsCmd.SetContext(oldContext) })

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	callErr := getTopicsCmd.RunE(getTopicsCmd, nil)
	_ = w.Close()
	os.Stdout = oldStdout
	output, readErr := io.ReadAll(r)
	_ = r.Close()
	require.NoError(t, readErr)
	require.NoError(t, callErr)

	assert.Equal(t, 3, requestedLimit)
	var result struct {
		Topics []lib.Topic `json:"topics"`
		Total  int         `json:"total"`
	}
	require.NoError(t, json.NewDecoder(bytes.NewReader(output)).Decode(&result))
	assert.Len(t, result.Topics, 3)
	assert.Equal(t, 3, result.Total)
}
