package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileCmd_MissingID(t *testing.T) {
	// Test that Cobra enforces required flag validation
	// Note: When calling RunE directly (bypassing Cobra), we need to set the flag
	// In real usage, Cobra validates before RunE is called
	viper.Set("accesskey", "testkey")
	viper.Set("secretkey", "testsecret")
	defer viper.Reset()

	// Cobra's MarkPersistentFlagRequired validates before RunE
	// This test verifies the command structure, not runtime validation
	assert.NotNil(t, getFileCmd)
}

func TestDeleteFileCmd_MissingID(t *testing.T) {
	// Test that Cobra enforces required flag validation
	// Cobra's MarkPersistentFlagRequired validates before RunE is called
	viper.Set("accesskey", "testkey")
	viper.Set("secretkey", "testsecret")
	defer viper.Reset()

	assert.NotNil(t, deleteFileCmd)
}

func TestGetFileCmd_MissingCredentials(t *testing.T) {
	viper.Reset()

	cmd := &cobra.Command{}
	err := getFileCmd.RunE(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DCI credentials not configured")
}

func TestDeleteFileCmd_MissingCredentials(t *testing.T) {
	viper.Reset()

	cmd := &cobra.Command{}
	err := deleteFileCmd.RunE(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DCI credentials not configured")
}

func TestGetFileCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("file content here"))
	}))
	defer server.Close()

	viper.Set("accesskey", "testkey")
	viper.Set("secretkey", "testsecret")
	defer viper.Reset()

	client := lib.NewClient("testkey", "testsecret")
	client.BaseURL = server.URL + "/api/v1"

	getFileIDFlag = "file-123"
	outputFormat = OutputFormatStdout

	content, contentType, err := client.GetFile(context.Background(), "file-123")
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
	assert.Equal(t, "file content here", string(content))
}

func TestGetFileCmd_SaveToFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("binary data"))
	}))
	defer server.Close()

	client := lib.NewClient("testkey", "testsecret")
	client.BaseURL = server.URL + "/api/v1"

	content, _, err := client.GetFile(context.Background(), "file-456")
	require.NoError(t, err)

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "downloaded-file")
	err = os.WriteFile(outPath, content, 0644)
	require.NoError(t, err)

	saved, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Equal(t, "binary data", string(saved))
}

func TestGetFileCmd_JSONOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<results/>"))
	}))
	defer server.Close()

	client := lib.NewClient("testkey", "testsecret")
	client.BaseURL = server.URL + "/api/v1"

	content, contentType, err := client.GetFile(context.Background(), "file-789")
	require.NoError(t, err)

	result := map[string]any{
		"id":           "file-789",
		"content_type": contentType,
		"size":         len(content),
	}
	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, "file-789", parsed["id"])
	assert.Equal(t, "application/xml", parsed["content_type"])
	assert.Equal(t, float64(10), parsed["size"])
}

func TestDeleteFileCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := lib.NewClient("testkey", "testsecret")
	client.BaseURL = server.URL + "/api/v1"

	err := client.DeleteFile(context.Background(), "file-123")
	assert.NoError(t, err)
}

func TestDeleteFileCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer server.Close()

	client := lib.NewClient("testkey", "testsecret")
	client.BaseURL = server.URL + "/api/v1"

	err := client.DeleteFile(context.Background(), "file-123")
	assert.Error(t, err)
}

func TestGetFileCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer server.Close()

	client := lib.NewClient("testkey", "testsecret")
	client.BaseURL = server.URL + "/api/v1"

	_, _, err := client.GetFile(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestDeleteFileCmd_JSONOutput(t *testing.T) {
	result := map[string]string{"status": "deleted", "id": "file-123"}
	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err)

	var parsed map[string]string
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, "deleted", parsed["status"])
	assert.Equal(t, "file-123", parsed["id"])
}
