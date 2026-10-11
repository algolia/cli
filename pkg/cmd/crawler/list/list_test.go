package list

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func Test_runListCmd_skipsForbiddenCrawlers(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/crawlers"):
			return jsonResponse(200, `{
				"items": [
					{"id": "ok-id", "name": "accessible"},
					{"id": "forbidden-id", "name": "not-mine"}
				],
				"page": 1, "itemsPerPage": 20, "total": 2
			}`), nil
		case strings.Contains(r.URL.Path, "forbidden-id"):
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Header:     http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
				Body:       io.NopCloser(strings.NewReader("Forbidden")),
			}, nil
		case strings.Contains(r.URL.Path, "ok-id"):
			return jsonResponse(200, `{
				"name": "accessible",
				"createdAt": "2026-01-01T00:00:00.000Z",
				"updatedAt": "2026-01-02T00:00:00.000Z",
				"running": false,
				"blocked": false,
				"config": {"appId": "APP_ID"}
			}`), nil
		default:
			t.Fatalf("unexpected request: %s", r.URL.Path)
			return nil, nil
		}
	})

	client := crawler.NewClientWithHTTPClient("user", "key", &http.Client{Transport: transport})

	ios, _, stdout, _ := iostreams.Test()
	printFlags := cmdutil.NewPrintFlags()
	printFlags.OutputFlagSpecified = func() bool { return false }
	opts := &ListOptions{
		IO: ios,
		CrawlerClient: func() (*crawler.Client, error) {
			return client, nil
		},
		PrintFlags: printFlags,
	}

	err := runListCmd(opts)
	require.NoError(t, err)

	assert.Contains(t, stdout.String(), "accessible")
	assert.NotContains(t, stdout.String(), "not-mine")
}
