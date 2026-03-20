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

func Test_runCrawlerCmd_UsesDefaultProfile(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", Default: true},
		{Name: "other"},
	})
	cfg.CurrentProfile.Name = ""

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

	assert.Equal(t, "default", cfg.CurrentProfile.Name)
	assert.Equal(t, test.CrawlerAuth{UserID: "crawler-user", APIKey: "crawler-key"}, cfg.CrawlerAuth["default"])
	assert.Contains(t, stdout.String(), "configured for profile: default")
}

func Test_runCrawlerCmd_UsesExplicitProfile(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "target"},
		{Name: "default", Default: true},
	})
	cfg.CurrentProfile.Name = "target"

	server := newCrawlerTestServer(t, "token-2", "crawler-user-2", "crawler-key-2")
	t.Cleanup(server.Close)

	err := runCrawlerCmd(&CrawlerOptions{
		IO:                 io,
		config:             cfg,
		OAuthClientID:      func() string { return "test-client-id" },
		NewDashboardClient: newDashboardTestClient(server),
		GetValidToken: func(client *dashboard.Client) (string, error) {
			return "token-2", nil
		},
	})
	require.NoError(t, err)

	assert.Equal(t, test.CrawlerAuth{UserID: "crawler-user-2", APIKey: "crawler-key-2"}, cfg.CrawlerAuth["target"])
	_, hasDefault := cfg.CrawlerAuth["default"]
	assert.False(t, hasDefault)
	assert.Contains(t, stdout.String(), "configured for profile: target")
}

func Test_runCrawlerCmd_ReturnsCrawlerAPIError(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "target"},
	})
	cfg.CurrentProfile.Name = "target"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer token-3", r.Header.Get("Authorization"))
		require.Equal(t, "/1/crawler/user", r.URL.Path)

		w.WriteHeader(http.StatusForbidden)
		_, err := fmt.Fprint(w, `{"success":false,"code":403,"message":"crawler access denied"}`)
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
	assert.Empty(t, cfg.CrawlerAuth)
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
