package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func Test_runCrawlerCmd_StoresCrawlerKeyForActiveApp(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(true)

	cfg := test.NewDefaultConfigStub()
	cfg.ActiveAppID = "APP1"

	server := newCrawlerTestServer(t, "token-1", "crawler-user", "crawler-key")
	t.Cleanup(server.Close)

	err := runCrawlerCmd(&CrawlerOptions{
		IO:                 io,
		config:             cfg,
		OAuthClientID:      func() string { return "test-client-id" },
		NewDashboardClient: newDashboardTestClient(server),
		GetValidToken: func(client *dashboard.Client) (string, error) {
			return "token-1", nil
		},
	})
	require.NoError(t, err)

	assert.Equal(t, "crawler-key", cfg.CrawlerKeys["APP1"])
	assert.Contains(t, stdout.String(), "configured for application: APP1")
}

func Test_runCrawlerCmd_ErrorsWithoutActiveApplication(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()

	cfg := test.NewDefaultConfigStub() // ActiveAppID left empty

	err := runCrawlerCmd(&CrawlerOptions{
		IO:                 io,
		config:             cfg,
		OAuthClientID:      func() string { return "test-client-id" },
		NewDashboardClient: func(string) *dashboard.Client { return nil },
		GetValidToken: func(client *dashboard.Client) (string, error) {
			return "", nil
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no application configured")
	assert.Empty(t, cfg.CrawlerKeys)
	assert.Empty(t, stdout.String())
}

func Test_runCrawlerCmd_ReturnsCrawlerAPIError(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "target"},
	})
	cfg.ActiveAppID = "APP1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer token-3", r.Header.Get("Authorization"))
		require.Equal(t, "/1/crawler/user", r.URL.Path)

		w.WriteHeader(http.StatusForbidden)
		_, err := fmt.Fprint(w, `{"errors":[{"status":"Forbidden","title":"Forbidden","detail":"crawler access denied"}]}`)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	err := runCrawlerCmd(&CrawlerOptions{
		IO:                 io,
		config:             cfg,
		OAuthClientID:      func() string { return "test-client-id" },
		NewDashboardClient: newDashboardTestClient(server),
		GetValidToken: func(client *dashboard.Client) (string, error) {
			return "token-3", nil
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get crawler user data: crawler access denied")
	assert.Empty(t, cfg.CrawlerKeys)
	assert.Empty(t, stdout.String())
}

func newCrawlerTestServer(t *testing.T, token, userID, apiKey string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer "+token, r.Header.Get("Authorization"))

		switch r.URL.Path {
		case "/1/crawler/user":
			_, err := fmt.Fprintf(w, `{"data":{"id":%q,"email":"crawler@example.com","name":"Crawler User","apiKey":%q}}`, userID, apiKey)
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
}

func newDashboardTestClient(server *httptest.Server) func(string) *dashboard.Client {
	return func(clientID string) *dashboard.Client {
		client := dashboard.NewClientWithHTTPClient(clientID, server.Client())
		client.APIURL = server.URL
		return client
	}
}
