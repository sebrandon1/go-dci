package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/stretchr/testify/assert"
)

func TestPrintProductsStdout(t *testing.T) {
	response := &lib.ProductsResponse{
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
	}

	assert.NotPanics(t, func() {
		printProductsStdout(response)
	})
}

func TestPrintProductsStdout_Empty(t *testing.T) {
	response := &lib.ProductsResponse{
		Meta:     lib.Meta{Count: 0},
		Products: []lib.Product{},
	}

	assert.NotPanics(t, func() {
		printProductsStdout(response)
	})
}

func TestPrintProductsJSON(t *testing.T) {
	response := &lib.ProductsResponse{
		Meta: lib.Meta{Count: 1},
		Products: []lib.Product{
			{
				ID:    "product-abc-123",
				Name:  "Red Hat OpenShift Container Platform",
				Label: "RHOCP",
				State: "active",
			},
		},
	}

	assert.NotPanics(t, func() {
		printProductsJSON(response)
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
		printProductJSON(response)
	})
}
