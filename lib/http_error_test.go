package lib

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatHTTPError(t *testing.T) {
	statusTests := []struct {
		name           string
		status         int
		statusText     string
		guidancePhrase string
		body           string
	}{
		{
			name:           "unauthorized",
			status:         401,
			statusText:     "Unauthorized",
			guidancePhrase: "Authentication failed.",
		},
		{
			name:           "forbidden",
			status:         403,
			statusText:     "Forbidden",
			guidancePhrase: "Permission denied.",
		},
		{
			name:           "not found",
			status:         404,
			statusText:     "Not Found",
			guidancePhrase: "Resource not found.",
		},
		{
			name:           "internal server error",
			status:         500,
			statusText:     "Internal Server Error",
			guidancePhrase: "DCI server error.",
		},
		{
			name:           "bad gateway",
			status:         502,
			statusText:     "Bad Gateway",
			guidancePhrase: "DCI server error.",
		},
		{
			name:           "service unavailable",
			status:         503,
			statusText:     "Service Unavailable",
			guidancePhrase: "DCI server error.",
		},
		{
			name:           "generic fallback",
			status:         418,
			statusText:     "I'm a teapot",
			guidancePhrase: "An unexpected error occurred",
		},
		{
			name:           "unauthorized with response body",
			status:         401,
			statusText:     "Unauthorized",
			guidancePhrase: "Authentication failed.",
			body:           "invalid access key",
		},
		{
			name:           "not found with response body",
			status:         404,
			statusText:     "Not Found",
			guidancePhrase: "Resource not found.",
			body:           "job does not exist",
		},
	}

	for _, tt := range statusTests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatHTTPError(tt.status, []byte(tt.body)).Error()

			assert.Contains(t, got, "HTTP "+strconv.Itoa(tt.status)+": "+tt.statusText)
			assert.Contains(t, got, tt.guidancePhrase)

			if tt.body != "" {
				assert.Contains(t, got, "Server response: "+tt.body)

				// Guidance must appear before the response body section.
				assert.Less(t, strings.Index(got, tt.guidancePhrase), strings.Index(got, "Server response:"))
			}
		})
	}

	bodyTests := []struct {
		name             string
		body             []byte
		wantResponseBody bool
	}{
		{name: "empty"},
		{name: "short", body: []byte("upstream detail"), wantResponseBody: true},
		{name: "499 bytes", body: []byte(strings.Repeat("x", 499)), wantResponseBody: true},
		{name: "500 bytes", body: []byte(strings.Repeat("x", 500))},
		// The cutoff is measured in bytes, not runes: 250 two-byte runes
		// reach the 500-byte limit even though the body is 250 characters.
		{name: "multi-byte under byte limit", body: []byte(strings.Repeat("é", 249)), wantResponseBody: true},
		{name: "multi-byte at byte limit", body: []byte(strings.Repeat("é", 250))},
	}

	for _, tt := range bodyTests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatHTTPError(400, tt.body).Error()

			if tt.wantResponseBody {
				assert.Contains(t, got, "Server response: "+string(tt.body))
			} else {
				assert.NotContains(t, got, "Server response")
			}
		})
	}
}
