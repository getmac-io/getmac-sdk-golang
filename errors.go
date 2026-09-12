package getmac

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ErrNotFound is wrapped by errors for resources that don't exist: an API
// response with status 404, or GetByName finding no virtual machine with the
// name. Check for it with errors.Is.
var ErrNotFound = errors.New("not found")

const (
	maxErrorBodySize    = 64 << 10
	maxErrorMessageSize = 512
)

// APIError is returned when the GetMac API responds with an unexpected status
// code. Message holds the reason the API gave, if any.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("unexpected status code: %d", e.StatusCode)
	}

	return fmt.Sprintf("unexpected status code: %d: %s", e.StatusCode, e.Message)
}

// Is reports a 404 response as ErrNotFound.
func (e *APIError) Is(target error) bool {
	return target == ErrNotFound && e.StatusCode == http.StatusNotFound
}

// newAPIError reads the reason from an error response. The API sends
// {"message": "..."}; anything else, such as a proxy's error page, is kept as
// trimmed text.
func newAPIError(resp *http.Response) *APIError {
	apiErr := &APIError{StatusCode: resp.StatusCode}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))
	if err != nil || len(body) == 0 {
		return apiErr
	}

	var payload struct {
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &payload) == nil {
		apiErr.Message = payload.Message
		return apiErr
	}

	message := strings.TrimSpace(string(body))
	if len(message) > maxErrorMessageSize {
		message = message[:maxErrorMessageSize] + "..."
	}
	apiErr.Message = message

	return apiErr
}
