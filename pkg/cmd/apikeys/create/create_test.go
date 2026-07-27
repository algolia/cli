package create

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
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

const (
	createdValue     = "new-value"
	createdUUID      = "uuid-123"
	unroutableAPIURL = "http://127.0.0.1:1"
)

type ttys struct {
	stdin  bool
	stdout bool
	stderr bool
}

func withoutSession(t *testing.T) {
	t.Helper()
	t.Setenv("ALGOLIA_API_KEY", "")
	t.Setenv("ALGOLIA_ADMIN_API_KEY", "")
	t.Setenv("ALGOLIA_APPLICATION_ID", "")
	t.Setenv("ALGOLIA_API_URL", unroutableAPIURL)
	keyring.MockInit()
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

func createdKeyResponse() dashboard.CreateAPIKeyResponse {
	maxHits := 42
	maxQueries := 1000
	queryParameters := "filters=visible:true"

	return dashboard.CreateAPIKeyResponse{
		Data: dashboard.APIKeyResource{
			ID:   createdUUID,
			Type: "api_key",
			Attributes: dashboard.APIKeyAttributes{
				Value:                  createdValue,
				ApplicationID:          "APP1",
				ACL:                    []string{"browse"},
				Description:            "stored description",
				Indexes:                []string{"SERIES"},
				Referers:               []string{"https://stored.example.com"},
				MaxHitsPerQuery:        &maxHits,
				MaxQueriesPerIPPerHour: &maxQueries,
				QueryParameters:        &queryParameters,
			},
		},
	}
}

// createKeyServer stubs the dashboard create endpoint at wantPath, recording the
// received payload.
func createKeyServer(
	t *testing.T,
	wantPath string,
	got *dashboard.CreateAPIKeyRequest,
) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, json.NewDecoder(r.Body).Decode(got))
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(createdKeyResponse()))
	})
	return httptest.NewServer(mux)
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

func newDashboardOpts(
	t *testing.T,
	srv *httptest.Server,
	tty ttys,
) (*CreateOptions, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	io, _, stdout, stderr := iostreams.Test()
	io.SetStdinTTY(tty.stdin)
	io.SetStdoutTTY(tty.stdout)
	io.SetStderrTTY(tty.stderr)

	opts := &CreateOptions{
		IO:     io,
		config: managedKeyConfig(),
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
		PrintFlags: cmdutil.NewPrintFlags(),
	}
	return opts, stdout, stderr
}

func TestNewCreateCmd(t *testing.T) {
	oneHour, _ := time.ParseDuration("1h")

	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts CreateOptions
	}{
		{
			name:     "all the flags",
			cli:      "-i foo,bar --acl search,browse -r \"http://foo.com\" -u 1h -d \"description\"",
			tty:      false,
			wantsErr: false,
			wantsOpts: CreateOptions{
				ACL:         []string{"search", "browse"},
				Indices:     []string{"foo", "bar"},
				Description: "description",
				Referers:    []string{"http://foo.com"},
				Validity:    oneHour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()
			if tt.tty {
				io.SetStdinTTY(tt.tty)
				io.SetStdoutTTY(tt.tty)
			}

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			var opts *CreateOptions
			cmd := NewCreateCmd(f, func(o *CreateOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.ACL, opts.ACL)
			assert.Equal(t, tt.wantsOpts.Indices, opts.Indices)
			assert.Equal(t, tt.wantsOpts.Description, opts.Description)
			assert.Equal(t, tt.wantsOpts.Referers, opts.Referers)
			assert.Equal(t, tt.wantsOpts.Validity, opts.Validity)
		})
	}
}

func TestNewCreateCmd_SkipsTheAdminACLCheck(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: io}
	cmd := NewCreateCmd(f, nil)

	assert.Equal(t, "true", cmd.Annotations["skipAuthCheck"])
	assert.Empty(t, cmd.Annotations["acls"])
}

func Test_runCreateCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "no TTY prints the bare key",
			cli:     "--acl search",
			isTTY:   false,
			wantOut: "foo\n",
		},
		{
			name:    "TTY",
			cli:     "--acl search",
			isTTY:   true,
			wantOut: "✓ API key created: foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withoutSession(t)

			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/keys"),
				httpmock.JSONResponse(search.AddApiKeyResponse{Key: "foo"}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, explicitKeyConfig(), "")
			cmd := NewCreateCmd(f, func(o *CreateOptions) error {
				o.NewDashboardClient = unusedDashboardClient(t)
				return runCreateCmd(o)
			})
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}

func Test_runCreateCmd_ExplicitAPIKeyUsesTheSearchAPI(t *testing.T) {
	withSession(t)

	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "1/keys"),
		httpmock.JSONResponse(search.AddApiKeyResponse{Key: "from-sapi"}),
	)

	cfg := managedKeyConfig()
	cfg.CurrentProfile.APIKey = "admin-key"

	f, out := test.NewFactory(true, &r, cfg, "")
	cmd := NewCreateCmd(f, func(o *CreateOptions) error {
		o.NewDashboardClient = unusedDashboardClient(t)
		return runCreateCmd(o)
	})
	out, err := test.Execute(cmd, "--acl search", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "from-sapi")
}

func Test_runCreateCmd_EnvAPIKeyUsesTheSearchAPI(t *testing.T) {
	withSession(t)
	t.Setenv("ALGOLIA_API_KEY", "env-admin-key")

	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("POST", "1/keys"),
		httpmock.JSONResponse(search.AddApiKeyResponse{Key: "from-sapi"}),
	)

	f, out := test.NewFactory(true, &r, managedKeyConfig(), "")
	cmd := NewCreateCmd(f, func(o *CreateOptions) error {
		o.NewDashboardClient = unusedDashboardClient(t)
		return runCreateCmd(o)
	})
	out, err := test.Execute(cmd, "--acl search", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "from-sapi")
}

func Test_runCreateCmd_WithSession(t *testing.T) {
	withSession(t)

	var got dashboard.CreateAPIKeyRequest
	srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
	defer srv.Close()

	opts, stdout, _ := newDashboardOpts(t, srv, ttys{stdout: true, stderr: true})
	opts.ACL = []string{"search", "browse"}
	opts.Indices = []string{"MOVIES"}
	opts.Referers = []string{"https://example.com"}

	require.NoError(t, runCreateCmd(opts))

	assert.Equal(t, []string{"search", "browse"}, got.ACL)
	assert.Equal(t, []string{"MOVIES"}, got.Indexes)
	assert.Equal(t, []string{"https://example.com"}, got.Referers)
	assert.Equal(t, defaultDescription, got.Description)
	assert.Equal(t, "✓ API key created: "+createdValue+"\n", stdout.String())
}

func Test_runCreateCmd_WithSessionNoTTYPrintsTheBareKey(t *testing.T) {
	withSession(t)

	var got dashboard.CreateAPIKeyRequest
	srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
	defer srv.Close()

	opts, stdout, _ := newDashboardOpts(t, srv, ttys{})
	opts.ACL = []string{"search"}

	require.NoError(t, runCreateCmd(opts))

	assert.Equal(t, createdValue+"\n", stdout.String())
}

func Test_runCreateCmd_WithSessionKeepsTheGivenDescription(t *testing.T) {
	withSession(t)

	var got dashboard.CreateAPIKeyRequest
	srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
	defer srv.Close()

	opts, _, _ := newDashboardOpts(t, srv, ttys{stdout: true, stderr: true})
	opts.ACL = []string{"search"}
	opts.Description = "frontend search key"

	require.NoError(t, runCreateCmd(opts))

	assert.Equal(t, "frontend search key", got.Description)
}

func Test_runCreateCmd_WithSessionStructuredOutput(t *testing.T) {
	withSession(t)

	var got dashboard.CreateAPIKeyRequest
	srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
	defer srv.Close()

	opts, stdout, _ := newDashboardOpts(t, srv, ttys{})
	opts.ACL = []string{"search"}
	format := "json"
	opts.PrintFlags.OutputFormat = &format

	require.NoError(t, runCreateCmd(opts))

	var key dashboard.APIKey
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &key))
	assert.Equal(t, createdValue, key.Value)
	assert.Equal(t, createdUUID, key.UUID)
	assert.Equal(t, "APP1", key.ApplicationID)
	assert.Equal(t, []string{"browse"}, key.ACL)
	assert.Equal(t, "stored description", key.Description)
	assert.Equal(t, []string{"SERIES"}, key.Indexes)
	assert.Equal(t, []string{"https://stored.example.com"}, key.Referers)
	require.NotNil(t, key.MaxHitsPerQuery)
	assert.Equal(t, 42, *key.MaxHitsPerQuery)
	require.NotNil(t, key.MaxQueriesPerIPPerHour)
	assert.Equal(t, 1000, *key.MaxQueriesPerIPPerHour)
	require.NotNil(t, key.QueryParameters)
	assert.Equal(t, "filters=visible:true", *key.QueryParameters)
}

func Test_runCreateCmd_WithSessionUnknownApplication(t *testing.T) {
	withSession(t)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/applications/APP1/api-keys",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"errors":[{"status":"404","title":"Not Found"}]}`))
		},
	)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	opts, _, _ := newDashboardOpts(t, srv, ttys{stdout: true, stderr: true})
	opts.ACL = []string{"search"}

	err := runCreateCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "APP1")
	assert.Contains(t, err.Error(), "doesn't have access to it")
	assert.Contains(t, err.Error(), "algolia application select")
}

func Test_runCreateCmd_WithSessionValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		opts    func(*CreateOptions)
		wantErr string
	}{
		{
			name:    "no ACL",
			opts:    func(o *CreateOptions) {},
			wantErr: "--acl is required",
		},
		{
			name: "validity",
			opts: func(o *CreateOptions) {
				o.ACL = []string{"search"}
				o.Validity = time.Hour
			},
			wantErr: "--validity isn't supported with your signed-in session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withSession(t)

			var got dashboard.CreateAPIKeyRequest
			srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
			defer srv.Close()

			opts, _, _ := newDashboardOpts(t, srv, ttys{stdout: true, stderr: true})
			tt.opts(opts)

			err := runCreateCmd(opts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func Test_runCreateCmd_SignedInWithoutAnApplication(t *testing.T) {
	withSession(t)

	var got dashboard.CreateAPIKeyRequest
	srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
	defer srv.Close()

	opts, _, _ := newDashboardOpts(t, srv, ttys{stdout: true, stderr: true})
	opts.ACL = []string{"search"}
	opts.config = &test.ConfigStub{}

	err := runCreateCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no application selected")
	assert.Contains(t, err.Error(), "algolia application select")
}

func Test_runCreateCmd_ApplicationIDFlagWithoutAStoredKeyUsesTheSession(t *testing.T) {
	withSession(t)

	var got dashboard.CreateAPIKeyRequest
	srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
	defer srv.Close()

	opts, stdout, _ := newDashboardOpts(t, srv, ttys{})
	opts.ACL = []string{"search"}
	opts.config = &test.ConfigStub{CurrentProfile: config.Profile{ApplicationID: "APP1"}}

	require.NoError(t, runCreateCmd(opts))

	assert.Equal(t, []string{"search"}, got.ACL)
	assert.Equal(t, createdValue+"\n", stdout.String())
}

func Test_runCreateCmd_ACLIsRequiredOnTheSearchAPIPath(t *testing.T) {
	withoutSession(t)

	r := httpmock.Registry{}
	r.Register(httpmock.REST("POST", "1/keys"), func(*http.Request) (*http.Response, error) {
		t.Error("no key must be created without --acl")
		return nil, errors.New("unexpected request")
	})

	f, out := test.NewFactory(false, &r, explicitKeyConfig(), "")
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--acl is required")
}

func Test_runCreateCmd_SessionStructuredOutputExposesTheKey(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		template string
		want     string
	}{
		{
			name:   "json",
			format: "json",
			want:   `"key":"` + createdValue + `"`,
		},
		{
			name:     "jsonpath",
			format:   "jsonpath",
			template: "{.key}",
			want:     createdValue + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withSession(t)

			var got dashboard.CreateAPIKeyRequest
			srv := createKeyServer(t, "/1/applications/APP1/api-keys", &got)
			defer srv.Close()

			opts, stdout, _ := newDashboardOpts(t, srv, ttys{})
			opts.ACL = []string{"search"}
			format := tt.format
			opts.PrintFlags.OutputFormat = &format
			if tt.template != "" {
				template := tt.template
				opts.PrintFlags.JSONPathPrintFlags.TemplateArgument = &template
			}

			require.NoError(t, runCreateCmd(opts))
			assert.Contains(t, stdout.String(), tt.want)
		})
	}
}

func Test_runCreateCmd_SearchAPIStructuredOutputExposesTheKey(t *testing.T) {
	tests := []struct {
		name string
		cli  string
		want string
	}{
		{
			name: "json",
			cli:  "--acl search -o json",
			want: `"key":"foo"`,
		},
		{
			name: "jsonpath",
			cli:  "--acl search -o jsonpath --template {.key}",
			want: "foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withoutSession(t)

			r := httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/keys"),
				httpmock.JSONResponse(search.AddApiKeyResponse{Key: "foo"}),
			)

			f, out := test.NewFactory(false, &r, explicitKeyConfig(), "")
			cmd := NewCreateCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			require.NoError(t, err)

			assert.Contains(t, out.String(), tt.want)
		})
	}
}

func Test_runCreateCmd_UnsupportedOutputFormatFailsBeforeCreating(t *testing.T) {
	t.Run("session path", func(t *testing.T) {
		withSession(t)

		srv := httptest.NewServer(
			http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				t.Errorf("no key must be created: %s %s", r.Method, r.URL.Path)
			}),
		)
		defer srv.Close()

		opts, stdout, _ := newDashboardOpts(t, srv, ttys{})
		opts.ACL = []string{"search"}
		format := "xml"
		opts.PrintFlags.OutputFormat = &format

		err := runCreateCmd(opts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to match a printer")
		assert.Empty(t, stdout.String())
	})

	t.Run("search API path", func(t *testing.T) {
		withoutSession(t)

		r := httpmock.Registry{}
		r.Register(httpmock.REST("POST", "1/keys"), func(*http.Request) (*http.Response, error) {
			t.Error("no key must be created")
			return nil, errors.New("unexpected request")
		})

		f, out := test.NewFactory(false, &r, explicitKeyConfig(), "")
		cmd := NewCreateCmd(f, nil)
		_, err := test.Execute(cmd, "--acl search -o xml", out)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to match a printer")
		assert.Empty(t, out.String())
	})
}

func Test_runCreateCmd_WithoutASessionOrAnApplication(t *testing.T) {
	withoutSession(t)

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		t.Errorf("no request must be made: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	opts, stdout, _ := newDashboardOpts(t, srv, ttys{})
	opts.ACL = []string{"search"}
	opts.config = &test.ConfigStub{}

	err := runCreateCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no application selected")
	assert.Contains(t, err.Error(), "algolia application select")
	assert.Empty(t, stdout.String())
}

func Test_runCreateCmd_WithoutASessionStdoutPipedStderrTTY(t *testing.T) {
	withoutSession(t)

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		t.Errorf("no request must be made without a session: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	opts, stdout, stderr := newDashboardOpts(t, srv, ttys{stdin: true, stderr: true})
	opts.ACL = []string{"search"}
	require.False(t, opts.IO.CanPrompt())

	err := runCreateCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a terminal")
	assert.Contains(t, err.Error(), "algolia auth login")
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "not logged in")

	prompting, _, _ := newDashboardOpts(t, srv, ttys{stdin: true, stdout: true, stderr: true})
	assert.True(t, prompting.IO.CanPrompt())
}

func Test_runCreateCmd_SearchAPIClientErrorSurfacesTheRemediation(t *testing.T) {
	tests := []struct {
		name    string
		session bool
		want    string
	}{
		{
			name: "signed out",
			want: "algolia auth login",
		},
		{
			name:    "signed in",
			session: true,
			want:    "algolia application select",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.session {
				withSession(t)
			} else {
				withoutSession(t)
			}
			t.Setenv("ALGOLIA_API_KEY", "adm")

			io, _, stdout, _ := iostreams.Test()
			opts := &CreateOptions{
				IO:     io,
				config: &test.ConfigStub{},
				ACL:    []string{"search"},
				SearchClient: func() (*search.APIClient, error) {
					return nil, config.ErrApplicationIDNotConfigured
				},
				NewDashboardClient: unusedDashboardClient(t),
				PrintFlags:         cmdutil.NewPrintFlags(),
			}

			err := runCreateCmd(opts)
			require.Error(t, err)
			assert.ErrorIs(t, err, config.ErrApplicationIDNotConfigured)
			assert.Contains(t, err.Error(), tt.want)
			assert.Empty(t, stdout.String())
		})
	}
}

func Test_searchAPICreateError(t *testing.T) {
	forbidden := &search.APIError{Status: http.StatusForbidden, Message: "Not enough rights"}

	t.Run("403 with an admin key in play", func(t *testing.T) {
		withoutSession(t)

		io, _, _, _ := iostreams.Test()
		cfg := managedKeyConfig()
		cfg.CurrentProfile.APIKey = "weak-key"
		opts := &CreateOptions{IO: io, config: cfg}

		err := searchAPICreateError(opts, forbidden)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "isn't an admin key")
		assert.Contains(t, err.Error(), "ALGOLIA_API_KEY")
		assert.NotContains(t, err.Error(), "algolia auth login")
	})

	t.Run("403 with the CLI-managed key", func(t *testing.T) {
		withoutSession(t)

		io, _, _, _ := iostreams.Test()
		opts := &CreateOptions{IO: io, config: managedKeyConfig()}

		err := searchAPICreateError(opts, forbidden)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "algolia auth login")
		assert.NotContains(t, err.Error(), "--api-key")
	})

	t.Run("non-403 errors pass through", func(t *testing.T) {
		withoutSession(t)

		io, _, _, _ := iostreams.Test()
		opts := &CreateOptions{IO: io, config: managedKeyConfig()}

		other := &search.APIError{Status: http.StatusBadRequest, Message: "nope"}
		assert.Same(t, other, searchAPICreateError(opts, other))

		plain := errors.New("boom")
		assert.Same(t, plain, searchAPICreateError(opts, plain))
	})
}
