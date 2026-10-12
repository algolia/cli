package shared

import (
	"encoding/json"
	"net/http"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
)

// RawResponse adapts the (response, body, error) tuple returned by an SDK
// *WithHTTPInfo call into raw response bytes for commands that forward the
// backend payload verbatim (e.g. conversation get/export, user-data get) —
// the typed SDK models would re-serialize and reshape those payloads.
//
// A non-2xx status is converted into an *agentStudio.APIError so callers get
// the same typed error the SDK's high-level methods return.
//
// The SDK leaves the response body open for the caller; it has already been
// read into body, so it is closed here to allow connection reuse.
func RawResponse(res *http.Response, body []byte, err error) ([]byte, error) {
	if res != nil && res.Body != nil {
		_ = res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if res != nil && res.StatusCode >= 300 {
		return nil, &agentStudio.APIError{
			Message: extractAPIMessage(body),
			Status:  res.StatusCode,
		}
	}
	return body, nil
}

// extractAPIMessage pulls a human-readable message from an error body,
// preferring the Algolia {"message":...} shape, then a string {"detail":...},
// and falling back to the raw body.
func extractAPIMessage(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var m map[string]any
	if json.Unmarshal(body, &m) == nil {
		if s, ok := m["message"].(string); ok && s != "" {
			return s
		}
		if s, ok := m["detail"].(string); ok && s != "" {
			return s
		}
	}
	return string(body)
}
