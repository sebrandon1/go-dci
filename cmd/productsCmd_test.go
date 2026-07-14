package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPrintProductsStdout(t *testing.T) {
	responses := []lib.ProductsResponse{
		{
			Meta: lib.Meta{Count: 2},
			Products: []lib.Product{
				{
					ID:          "product-abc-123",
					Name:        "Red Hat OpenShift Container Platform",
					Label:       "RHOCP",
					Description: "Enterprise Kubernetes platform",
					State:       "active",
					CreatedAt:   "2024-01-15T10:30:00.000000",
					UpdatedAt:   "2024-06-20T14:45:00.000000",
				},
				{
					ID:          "product-def-456",
					Name:        "Red Hat Enterprise Linux",
					Label:       "RHEL",
					Description: "Enterprise Linux distribution",
					State:       "active",
					CreatedAt:   "2024-02-01T08:00:00.000000",
					UpdatedAt:   "2024-07-10T12:00:00.000000",
				},
			},
		},
	}

	assert.NotPanics(t, func() {
		printProductsStdout(responses)
	})
}

func TestPrintProductsStdout_Empty(t *testing.T) {
	responses := []lib.ProductsResponse{
		{
			Meta:     lib.Meta{Count: 0},
			Products: []lib.Product{},
		},
	}

	assert.NotPanics(t, func() {
		printProductsStdout(responses)
	})
}

func TestPrintProductsJSON(t *testing.T) {
	responses := []lib.ProductsResponse{
		{
			Meta: lib.Meta{Count: 1},
			Products: []lib.Product{
				{
					ID:    "product-abc-123",
					Name:  "Red Hat OpenShift Container Platform",
					Label: "RHOCP",
					State: "active",
				},
			},
		},
	}

	assert.NotPanics(t, func() {
		err := printProductsJSON(responses)
		assert.NoError(t, err)
	})
}

func TestPrintProductStdout(t *testing.T) {
	response := &lib.ProductResponse{
		Product: lib.Product{
			ID:          "product-abc-123",
			Name:        "Red Hat OpenShift Container Platform",
			Label:       "RHOCP",
			Description: "Enterprise Kubernetes platform",
			State:       "active",
			CreatedAt:   "2024-01-15T10:30:00.000000",
			UpdatedAt:   "2024-06-20T14:45:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printProductStdout(response)
	})
}

func TestPrintProductJSON(t *testing.T) {
	response := &lib.ProductResponse{
		Product: lib.Product{
			ID:          "product-abc-123",
			Name:        "Red Hat OpenShift Container Platform",
			Label:       "RHOCP",
			Description: "Enterprise Kubernetes platform",
			State:       "active",
			CreatedAt:   "2024-01-15T10:30:00.000000",
			UpdatedAt:   "2024-06-20T14:45:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		err := printProductJSON(response)
		assert.NoError(t, err)
	})
}

func TestGetProductCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.ProductResponse{Product: lib.Product{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "RHOCP"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getProductIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getProductIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getProductCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetProductCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getProductIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getProductIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getProductCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}
