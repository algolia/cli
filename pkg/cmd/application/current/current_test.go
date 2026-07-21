package current

import (
	"bytes"
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
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
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

func newServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/1/application/APP1", func(w http.ResponseWriter, _ *http.Request) {
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		require.NoError(t, json.NewEncoder(w).Encode(dashboard.SingleApplicationResponse{
			Data: dashboard.ApplicationResource{
				ID:   "APP1",
				Type: "application",
				Attributes: dashboard.ApplicationAttributes{
					ApplicationID: "APP1",
					Name:          "My App",
					Plan:          dashboard.ApplicationPlan{Label: "Grow Plus"},
				},
			},
		}))
	})
	return httptest.NewServer(mux)
}

func newOpts(
	t *testing.T,
	srv *httptest.Server,
	cfg *test.ConfigStub,
	output string,
	signedIn bool,
) (*CurrentOptions, *bytes.Buffer) {
	t.Helper()
	if signedIn {
		seedToken(t)
	} else {
		keyring.MockInit()
		auth.ClearToken()
	}

	io, _, stdout, _ := iostreams.Test()
	pf := cmdutil.NewPrintFlags()
	*pf.OutputFormat = output
	pf.OutputFlagSpecified = func() bool { return output != "" }

	opts := &CurrentOptions{
		IO:         io,
		Config:     cfg,
		PrintFlags: pf,
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
	return opts, stdout
}

func configWithApp(appID, alias string) *test.ConfigStub {
	cfg := &test.ConfigStub{
		CurrentProfile: config.Profile{ApplicationID: appID},
	}
	if alias != "" {
		cfg.SavedApps = map[string]test.SavedApplication{
			appID: {Alias: alias},
		}
	}
	return cfg
}

func Test_runCurrentCmd(t *testing.T) {
	srv := newServer(t, http.StatusOK)
	defer srv.Close()

	opts, out := newOpts(t, srv, configWithApp("APP1", "my-alias"), "", true)
	require.NoError(t, runCurrentCmd(opts))

	got := out.String()
	assert.Contains(t, got, "APP1")
	assert.Contains(t, got, "my-alias")
	assert.Contains(t, got, "My App")
	assert.Contains(t, got, "Grow Plus")
}

func Test_runCurrentCmd_notConfigured(t *testing.T) {
	srv := newServer(t, http.StatusOK)
	defer srv.Close()

	opts, _ := newOpts(t, srv, configWithApp("", ""), "", true)
	err := runCurrentCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no current application configured")
}

func Test_runCurrentCmd_apiFailure(t *testing.T) {
	srv := newServer(t, http.StatusInternalServerError)
	defer srv.Close()

	opts, out := newOpts(t, srv, configWithApp("APP1", "my-alias"), "", true)
	require.NoError(t, runCurrentCmd(opts))

	got := out.String()
	assert.Contains(t, got, "APP1")
	assert.Contains(t, got, "my-alias")
	assert.NotContains(t, got, "My App")
	assert.Contains(t, got, "Couldn't fetch the application name and plan")
}

func Test_runCurrentCmd_signedOut(t *testing.T) {
	hit := false
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		hit = true
	}))
	defer srv.Close()

	opts, out := newOpts(t, srv, configWithApp("APP1", "my-alias"), "", false)
	require.NoError(t, runCurrentCmd(opts))

	got := out.String()
	assert.Contains(t, got, "APP1")
	assert.Contains(t, got, "my-alias")
	assert.NotContains(t, got, "My App")
	assert.Contains(t, got, `Sign in with "algolia auth login"`)
	assert.False(t, hit, "expected no API/login call when signed out")
}

func Test_runCurrentCmd_outputJSON(t *testing.T) {
	srv := newServer(t, http.StatusOK)
	defer srv.Close()

	opts, out := newOpts(t, srv, configWithApp("APP1", "my-alias"), "json", true)
	require.NoError(t, runCurrentCmd(opts))

	got := out.String()
	assert.Contains(t, got, `"id":"APP1"`)
	assert.Contains(t, got, `"alias":"my-alias"`)
	assert.Contains(t, got, `"name":"My App"`)
	assert.Contains(t, got, `"plan":"Grow Plus"`)
}
