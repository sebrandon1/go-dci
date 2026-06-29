package lib

import "fmt"

// formatHTTPError formats HTTP error responses into user-friendly messages
func formatHTTPError(statusCode int, body []byte) error {
	switch statusCode {
	case 401:
		return fmt.Errorf("authentication failed (401 Unauthorized)\nCheck your credentials with: go-dci identity")
	case 403:
		return fmt.Errorf("access denied (403 Forbidden)\nYour credentials lack permission for this resource")
	case 404:
		return fmt.Errorf("resource not found (404 Not Found)\nCheck the ID and try again")
	case 500:
		return fmt.Errorf("DCI API server error (500 Internal Server Error)\nTry again later or check status at https://www.distributed-ci.io/")
	default:
		// For other errors, show first 200 chars of body
		bodyStr := string(body)
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return fmt.Errorf("HTTP %d: %s", statusCode, bodyStr)
	}
}
