package selectapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/keychain"
	"github.com/algolia/cli/test"
)

func seedToken(t *testing.T) {
	t.Helper()
	keyring.MockInit()
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "test-token",
		ExpiresIn:   3600,
		CreatedAt:   time.Now().Unix(),
	}))
}

// selectServer stubs the dashboard endpoints select uses: listing applications
// and creating an API key. createHit records whether the key-creation endpoint
// was called.
func selectServer(t *testing.T, createHit *bool) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, _ *http.Request) {
		require.NoError(t, json.NewEncoder(w).Encode(dashboard.ApplicationsResponse{
			Data: []dashboard.ApplicationResource{{
				ID:   "APP1",
				Type: "application",
				Attributes: dashboard.ApplicationAttributes{
					ApplicationID: "APP1",
					Name:          "My App",
				},
			}},
			Meta: dashboard.PaginationMeta{CurrentPage: 1, TotalPages: 1},
		}))
	})
	mux.HandleFunc(
		"/1/applications/APP1/api-keys",
		func(w http.ResponseWriter, _ *http.Request) {
			*createHit = true
			w.WriteHeader(http.StatusCreated)
			require.NoError(t, json.NewEncoder(w).Encode(dashboard.CreateAPIKeyResponse{
				Data: dashboard.APIKeyResource{
					ID:         "new-uuid",
					Attributes: dashboard.APIKeyAttributes{Value: "new-key"},
				},
			}))
		},
	)
	return httptest.NewServer(mux)
}

func newSelectOpts(t *testing.T, srv *httptest.Server, cfg *test.ConfigStub) *SelectOptions {
	t.Helper()
	seedToken(t)
	io, _, _, _ := iostreams.Test()
	return &SelectOptions{
		IO:      io,
		Config:  cfg,
		AppName: "My App", // bypasses the interactive picker
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
}

func Test_runSelectCmd_RegeneratesKeyWhenNoUUID(t *testing.T) {
	createHit := false
	srv := selectServer(t, &createHit)
	defer srv.Close()

	// Migrated application: present in state with an alias, but no UUID.
	cfg := &test.ConfigStub{
		SavedApps: map[string]test.SavedApplication{
			"APP1": {Alias: "my app", APIKey: "old-key"},
		},
	}
	opts := newSelectOpts(t, srv, cfg)

	app, err := runSelectCmd(opts)
	require.NoError(t, err)
	require.NotNil(t, app)

	assert.True(t, createHit, "expected a fresh key to be generated when no UUID is stored")
	assert.Equal(t, "new-uuid", cfg.SavedApps["APP1"].APIKeyUUID)
	assert.Equal(t, "new-key", cfg.SavedApps["APP1"].APIKey)
}

func Test_runSelectCmd_ReusesKeyWhenUUIDPresent(t *testing.T) {
	createHit := false
	srv := selectServer(t, &createHit)
	defer srv.Close()

	cfg := &test.ConfigStub{
		SavedApps: map[string]test.SavedApplication{
			"APP1": {Alias: "my app", APIKeyUUID: "existing-uuid", APIKey: "old-key"},
		},
	}
	opts := newSelectOpts(t, srv, cfg)
	// A key in the keychain lets ReuseExistingAPIKey succeed.
	require.NoError(t, keychain.SaveAppSecrets("APP1", keychain.AppSecrets{APIKey: "kc-key"}))

	app, err := runSelectCmd(opts)
	require.NoError(t, err)
	require.NotNil(t, app)

	assert.False(t, createHit, "expected no new key when a UUID is already stored")
	assert.Equal(t, "existing-uuid", cfg.SavedApps["APP1"].APIKeyUUID)
}
