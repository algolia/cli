package crawler

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestClient(fn roundTripFunc) *Client {
	return NewClientWithHTTPClient("user", "key", &http.Client{Transport: fn})
}

func TestRequest_PlainTextError(t *testing.T) {
	client := newTestClient(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
			Body:       io.NopCloser(strings.NewReader("Forbidden")),
		}, nil
	})

	_, err := client.Get("some-id", true)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", apiErr.StatusCode)
	}
	if apiErr.Error() != "Forbidden" {
		t.Errorf("expected message %q, got %q", "Forbidden", apiErr.Error())
	}
}
