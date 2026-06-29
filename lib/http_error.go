package lib

import (
	"fmt"
	"net/http"
)

// formatHTTPError returns a user-friendly error message with actionable guidance
// based on the HTTP status code and response body.
func formatHTTPError(statusCode int, body []byte) error {
	baseMsg := fmt.Sprintf("HTTP %d: %s", statusCode, http.StatusText(statusCode))

	var guidance string
	switch statusCode {
	case http.StatusUnauthorized:
		guidance = "Authentication failed. Please check your DCI credentials:\n" +
			"  - Verify your access key and secret key are correct\n" +
			"  - Run: go-dci config set --accesskey <key> --secretkey <secret>\n" +
			"  - Or set environment variables: GO_DCI_ACCESSKEY and GO_DCI_SECRETKEY"
	case http.StatusForbidden:
		guidance = "Permission denied. Your credentials are valid but lack permission for this resource:\n" +
			"  - Contact your DCI administrator to request access\n" +
			"  - Verify you're using the correct team/project credentials"
	case http.StatusNotFound:
		guidance = "Resource not found. The requested item does not exist:\n" +
			"  - Check the resource ID/name is correct\n" +
			"  - Verify the resource hasn't been deleted\n" +
			"  - Use list commands to see available resources"
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		guidance = "DCI server error. This is a temporary issue:\n" +
			"  - Try again in a few moments\n" +
			"  - Check DCI service status\n" +
			"  - Contact support if the issue persists"
	default:
		guidance = "An unexpected error occurred"
	}

	msg := fmt.Sprintf("%s\n%s", baseMsg, guidance)

	// Include response body if present and not too large
	if len(body) > 0 && len(body) < 500 {
		msg += fmt.Sprintf("\n\nServer response: %s", string(body))
	}

	return fmt.Errorf("%s", msg)
}
