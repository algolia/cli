package list

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

const unroutableAPIURL = "http://127.0.0.1:1"

func freezeNow(t *testing.T) {
	t.Helper()
	oldNowFn := nowFn
	nowFn = func() time.Time { return time.Unix(1735689600, 0) } // 2025-01-01T00:00:00Z
	t.Cleanup(func() { nowFn = oldNowFn })
}

func withoutSession(t *testing.T) {
	t.Helper()
	t.Setenv("ALGOLIA_API_KEY", "")
	t.Setenv("ALGOLIA_ADMIN_API_KEY", "")
	t.Setenv("ALGOLIA_APPLICATION_ID", "")
	t.Setenv("ALGOLIA_API_URL", unroutableAPIURL)
	keyring.MockInit()
}

func managedKeyConfig() *test.ConfigStub {
	return &test.ConfigStub{
		CurrentProfile: config.Profile{ApplicationID: "APP1"},
		ActiveAppID:    "APP1",
		SavedApps: map[string]test.SavedApplication{
			"APP1": {APIKeyUUID: "uuid-1", APIKey: "cli-key"},
		},
	}
}

func explicitKeyConfig() *test.ConfigStub {
	cfg := managedKeyConfig()
	cfg.CurrentProfile.APIKey = "adm"

	return cfg
}

func unusedDashboardClient(t *testing.T) func(string) *dashboard.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected dashboard request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	return func(string) *dashboard.Client {
		t.Error("the dashboard client must not be used on the Search API path")
		c := dashboard.NewClientWithHTTPClient("test", srv.Client())
		c.APIURL = srv.URL
		return c
	}
}

func withSession(t *testing.T) {
	t.Helper()
	withoutSession(t)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "tok-1",
		ExpiresIn:   3600,
		CreatedAt:   time.Now().Unix(),
	}))
}

// listKeysServer stubs the dashboard list endpoint at wantPath, serving pages
// out of the given resource batches.
func listKeysServer(
	t *testing.T,
	wantPath string,
	pages [][]dashboard.APIKeyResource,
) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	requests := 0
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		requests++
		require.LessOrEqual(t, requests, len(pages)+1, "the pagination loop is unbounded")

		page := 1
		if r.URL.Query().Get("page") == "2" {
			page = 2
		}

		require.NoError(t, json.NewEncoder(w).Encode(dashboard.APIKeysResponse{
			Data: pages[page-1],
			Meta: dashboard.PaginationMeta{
				CurrentPage: page,
				TotalPages:  len(pages),
				TotalCount:  len(pages[page-1]),
				PerPage:     15,
			},
		}))
	})
	return httptest.NewServer(mux)
}

func newSessionOpts(
	t *testing.T,
	srv *httptest.Server,
	isTTY bool,
) (*ListOptions, *bytes.Buffer) {
	t.Helper()

	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(isTTY)

	opts := &ListOptions{
		IO:     io,
		Config: managedKeyConfig(),
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
		Reauthenticate: auth.ReauthenticateIfExpired,
		PrintFlags:     cmdutil.NewPrintFlags(),
	}
	return opts, stdout
}

func sessionKey(uuid, value, description string) dashboard.APIKeyResource {
	return dashboard.APIKeyResource{
		ID:   uuid,
		Type: "api_key",
		Attributes: dashboard.APIKeyAttributes{
			Value:       value,
			ACL:         []string{"search"},
			Description: description,
			Indexes:     []string{},
			Referers:    []string{},
			CreatedAt:   "2020-01-01T00:00:00.000Z",
		},
	}
}

func TestNewListCmd_SkipsTheAdminACLCheck(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: io}
	cmd := NewListCmd(f, nil)

	assert.Equal(t, "true", cmd.Annotations["skipAuthCheck"])
	assert.Empty(t, cmd.Annotations["acls"])
}

func Test_runListCmd(t *testing.T) {
	tests := []struct {
		name    string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "list",
			isTTY:   false,
			wantOut: "foo\ttest\t[search]\t[]\tNever expire\t0\t0\t[]\t5 years ago\n",
		},
		{
			name:    "list_tty",
			isTTY:   true,
			wantOut: "KEY  DESCRIPTION  ACL      INDICES  VALI...  MAX ...  MAX ...  REFE...  CREA...\nfoo  test         [sea...  []       Neve...  0        0        []       5 ye...\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withoutSession(t)
			freezeNow(t)

			name := "test"
			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/keys"),
				httpmock.JSONResponse(search.ListApiKeysResponse{
					Keys: []search.GetApiKeyResponse{
						{
							Value:       "foo",
							Description: &name,
							Acl:         []search.Acl{search.ACL_SEARCH},
							CreatedAt:   1577836800,
						},
					},
				}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, explicitKeyConfig(), "")
			cmd := NewListCmd(f, func(o *ListOptions) error {
				o.NewDashboardClient = unusedDashboardClient(t)
				return runListCmd(o)
			})
			out, err := test.Execute(cmd, "", out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

func Test_runListCmd_outputJSON(t *testing.T) {
	withoutSession(t)
	freezeNow(t)

	name := "test"
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/keys"),
		httpmock.JSONResponse(search.ListApiKeysResponse{
			Keys: []search.GetApiKeyResponse{
				{
					Value:       "foo",
					Description: &name,
					Acl:         []search.Acl{search.ACL_SEARCH},
					CreatedAt:   1577836800,
				},
			},
		}),
	)

	f, out := test.NewFactory(false, &r, explicitKeyConfig(), "")
	cmd := NewListCmd(f, func(o *ListOptions) error {
		o.NewDashboardClient = unusedDashboardClient(t)
		return runListCmd(o)
	})
	out, err := test.Execute(cmd, "--output json", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, out.String(), `"keys":[`)
	assert.Contains(t, out.String(), `"value":"foo"`)
}

func Test_runListCmd_ExplicitAPIKeyUsesTheSearchAPI(t *testing.T) {
	withSession(t)
	freezeNow(t)

	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/keys"),
		httpmock.JSONResponse(search.ListApiKeysResponse{
			Keys: []search.GetApiKeyResponse{{Value: "from-sapi"}},
		}),
	)

	cfg := managedKeyConfig()
	cfg.CurrentProfile.APIKey = "admin-key"

	f, out := test.NewFactory(false, &r, cfg, "")
	cmd := NewListCmd(f, func(o *ListOptions) error {
		o.NewDashboardClient = unusedDashboardClient(t)
		return runListCmd(o)
	})
	out, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "from-sapi")
}

func Test_runListCmd_EnvAPIKeyUsesTheSearchAPI(t *testing.T) {
	withSession(t)
	freezeNow(t)
	t.Setenv("ALGOLIA_API_KEY", "env-admin-key")

	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/keys"),
		httpmock.JSONResponse(search.ListApiKeysResponse{
			Keys: []search.GetApiKeyResponse{{Value: "from-sapi"}},
		}),
	)

	f, out := test.NewFactory(false, &r, managedKeyConfig(), "")
	cmd := NewListCmd(f, func(o *ListOptions) error {
		o.NewDashboardClient = unusedDashboardClient(t)
		return runListCmd(o)
	})
	out, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "from-sapi")
}

func Test_runListCmd_WithSession(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{
		{sessionKey("uuid-1", "search-key", "frontend")},
	})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)

	require.NoError(t, runListCmd(opts))

	assert.Equal(
		t,
		"search-key\tfrontend\t[search]\t[]\t-\t0\t0\t[]\t5 years ago\n",
		stdout.String(),
	)
}

func Test_runListCmd_WithSessionFollowsPagination(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{
		{sessionKey("uuid-1", "key-1", "first")},
		{sessionKey("uuid-2", "key-2", "second")},
	})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)

	require.NoError(t, runListCmd(opts))

	assert.Contains(t, stdout.String(), "key-1")
	assert.Contains(t, stdout.String(), "key-2")
}

func Test_runListCmd_WithSessionStructuredOutput(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{
		{sessionKey("uuid-1", "search-key", "frontend")},
	})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)
	format := "json"
	opts.PrintFlags.OutputFormat = &format

	require.NoError(t, runListCmd(opts))

	var keys []dashboard.APIKey
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &keys))
	require.Len(t, keys, 1)
	assert.Equal(t, "uuid-1", keys[0].UUID)
	assert.Equal(t, "search-key", keys[0].Value)
	assert.Equal(t, []string{"search"}, keys[0].ACL)
}

func Test_runListCmd_WithSessionEmpty(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{{}})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)

	require.NoError(t, runListCmd(opts))
	assert.Equal(t, "", stdout.String())
}

func Test_runListCmd_WithSessionUnknownApplication(t *testing.T) {
	withSession(t)
	freezeNow(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1/api-keys", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"status":"404","title":"Not Found"}]}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	opts, _ := newSessionOpts(t, srv, false)

	err := runListCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "APP1")
	assert.Contains(t, err.Error(), "doesn't have access to it")
}

func Test_runListCmd_SignedInWithoutAnApplication(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{{}})
	defer srv.Close()

	opts, _ := newSessionOpts(t, srv, false)
	opts.Config = &test.ConfigStub{}

	err := runListCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no application selected")
	assert.Contains(t, err.Error(), "algolia application select")
}

func Test_runListCmd_ApplicationIDFlagWithoutAStoredKeyUsesTheSession(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{
		{sessionKey("uuid-1", "search-key", "frontend")},
	})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)
	opts.Config = &test.ConfigStub{CurrentProfile: config.Profile{ApplicationID: "APP1"}}

	require.NoError(t, runListCmd(opts))

	assert.Contains(t, stdout.String(), "search-key")
}

func Test_runListCmd_SearchAPIKeyWithoutADescription(t *testing.T) {
	withoutSession(t)
	freezeNow(t)

	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/keys"),
		httpmock.JSONResponse(search.ListApiKeysResponse{
			Keys: []search.GetApiKeyResponse{
				{
					Value:     "foo",
					Acl:       []search.Acl{search.ACL_SEARCH},
					CreatedAt: 1577836800,
				},
			},
		}),
	)

	f, out := test.NewFactory(false, &r, explicitKeyConfig(), "")
	cmd := NewListCmd(f, func(o *ListOptions) error {
		o.NewDashboardClient = unusedDashboardClient(t)
		return runListCmd(o)
	})
	out, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Equal(
		t,
		"foo\t\t[search]\t[]\tNever expire\t0\t0\t[]\t5 years ago\n",
		out.String(),
	)
}

func Test_runListCmd_WithSessionEndpointNotAvailable(t *testing.T) {
	withSession(t)
	freezeNow(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1/api-keys", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("<!DOCTYPE html><html><body>Not found</body></html>"))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	opts, _ := newSessionOpts(t, srv, false)

	err := runListCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "newer Algolia API version")
	assert.NotContains(t, err.Error(), "doesn't have access to it")
}

func Test_searchAPIListError(t *testing.T) {
	forbidden := &search.APIError{Status: http.StatusForbidden, Message: "Not enough rights"}

	t.Run("403 with an admin key in play", func(t *testing.T) {
		withoutSession(t)

		io, _, _, _ := iostreams.Test()
		cfg := managedKeyConfig()
		cfg.CurrentProfile.APIKey = "weak-key"
		opts := &ListOptions{IO: io, Config: cfg}

		err := searchAPIListError(opts, forbidden)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "isn't an admin key")
		assert.Contains(t, err.Error(), "ALGOLIA_API_KEY")
		assert.NotContains(t, err.Error(), "algolia auth login")
	})

	t.Run("403 with the CLI-managed key", func(t *testing.T) {
		withoutSession(t)

		io, _, _, _ := iostreams.Test()
		opts := &ListOptions{IO: io, Config: managedKeyConfig()}

		err := searchAPIListError(opts, forbidden)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "algolia auth login")
		assert.NotContains(t, err.Error(), "--api-key")
	})

	t.Run("non-403 errors pass through", func(t *testing.T) {
		withoutSession(t)

		io, _, _, _ := iostreams.Test()
		opts := &ListOptions{IO: io, Config: managedKeyConfig()}

		other := &search.APIError{Status: http.StatusBadRequest, Message: "nope"}
		assert.Same(t, other, searchAPIListError(opts, other))

		plain := errors.New("boom")
		assert.Same(t, plain, searchAPIListError(opts, plain))
	})
}

func Test_formatCreatedAt(t *testing.T) {
	now := time.Unix(1735689600, 0)

	assert.Equal(t, "5 years ago", formatCreatedAt(now, "2020-01-01T00:00:00.000Z"))
	assert.Equal(t, "-", formatCreatedAt(now, ""))
	assert.Equal(t, "not-a-date", formatCreatedAt(now, "not-a-date"))
}

func Test_runListCmd_WithSessionMaskedKeyValue(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{
		{sessionKey("uuid-1", "", "restricted key")},
	})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)

	require.NoError(t, runListCmd(opts))

	assert.Equal(
		t,
		"-\trestricted key\t[search]\t[]\t-\t0\t0\t[]\t5 years ago\n",
		stdout.String(),
	)
}

func Test_runListCmd_WithSessionMaskedKeyValueStructuredOutput(t *testing.T) {
	withSession(t)
	freezeNow(t)

	srv := listKeysServer(t, "/1/applications/APP1/api-keys", [][]dashboard.APIKeyResource{
		{sessionKey("uuid-1", "", "restricted key")},
	})
	defer srv.Close()

	opts, stdout := newSessionOpts(t, srv, false)
	format := "json"
	opts.PrintFlags.OutputFormat = &format

	require.NoError(t, runListCmd(opts))

	var keys []map[string]any
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &keys))
	require.Len(t, keys, 1)
	assert.NotContains(t, keys[0], "value")
	assert.Equal(t, "uuid-1", keys[0]["uuid"])
	assert.Equal(t, "restricted key", keys[0]["description"])
}
