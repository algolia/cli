package selectapp

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
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
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
	// --app-name bypasses the interactive picker.
	return newSelectOptsWithSelector(t, srv, cfg, func(opts *SelectOptions) {
		opts.AppName = "My App"
	})
}

func newSelectOptsWithSelector(
	t *testing.T,
	srv *httptest.Server,
	cfg *test.ConfigStub,
	selector func(*SelectOptions),
) *SelectOptions {
	t.Helper()
	seedToken(t)
	io, _, _, _ := iostreams.Test()
	opts := &SelectOptions{
		IO:     io,
		Config: cfg,
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
	selector(opts)
	return opts
}

func Test_runSelectCmd_NonInteractiveWritesJSONOnlyToStdout(t *testing.T) {
	createHit := false
	srv := selectServer(t, &createHit)
	defer srv.Close()

	cfg := &test.ConfigStub{}
	opts := newSelectOpts(t, srv, cfg)
	opts.NonInteractive = true
	opts.PrintFlags = cmdutil.NewPrintFlags()
	// Mirrors what the command does before running.
	applyNonInteractive(opts)

	stdout, stderr := captureOutput(t, opts.IO)

	app, err := runSelectCmd(opts)
	require.NoError(t, err)
	require.NoError(t, printSelection(opts, app))

	// The JSON document is the whole output: progress narration is dropped, not
	// moved to stderr, and key material never appears.
	var got apputil.ApplicationOutput
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &got), "stdout: %q", stdout.String())
	assert.Equal(t, apputil.ApplicationOutput{
		ID:    "APP1",
		Alias: "my app",
		Name:  "My App",
	}, got)
	assert.NotContains(t, stdout.String(), "API key")
	assert.NotContains(t, stdout.String(), "new-key")
	assert.Empty(t, stderr.String())
}

func Test_applyNonInteractive_DefaultsToJSON(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &SelectOptions{IO: io, PrintFlags: cmdutil.NewPrintFlags()}

	// Flag off: nothing changes.
	applyNonInteractive(opts)
	assert.False(t, opts.PrintFlags.HasStructuredOutput())
	assert.False(t, io.GetNeverPrompt())

	opts.NonInteractive = true
	applyNonInteractive(opts)
	assert.True(t, io.GetNeverPrompt())
	assert.Equal(t, "json", *opts.PrintFlags.OutputFormat)
}

// The spinner writes to stderr even on a TTY, so a non-interactive run has to
// silence it too.
func Test_applyNonInteractive_DisablesProgressIndicator(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetProgressIndicatorEnabled(true)
	opts := &SelectOptions{IO: io, PrintFlags: cmdutil.NewPrintFlags(), NonInteractive: true}

	applyNonInteractive(opts)

	assert.False(t, io.GetProgressIndicatorEnabled())
}

func Test_applyNonInteractive_KeepsExplicitOutput(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &SelectOptions{IO: io, PrintFlags: cmdutil.NewPrintFlags(), NonInteractive: true}
	*opts.PrintFlags.OutputFormat = "jsonpath={.id}"

	applyNonInteractive(opts)
	assert.Equal(t, "jsonpath={.id}", *opts.PrintFlags.OutputFormat)
}

func TestNewSelectCmd_NonInteractiveRequiresSelector(t *testing.T) {
	f, inOut := test.NewFactory(true, nil, &test.ConfigStub{}, "")

	_, err := test.Execute(NewSelectCmd(f), "--non-interactive", inOut)
	assert.ErrorContains(t, err, "--app-id or --app-name is required in non-interactive mode")
	assert.Empty(t, inOut.OutBuf.String())
}

func Test_printSelection_NoApplication(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &SelectOptions{IO: io, PrintFlags: cmdutil.NewPrintFlags(), NonInteractive: true}
	applyNonInteractive(opts)

	assert.ErrorContains(t, printSelection(opts, nil), "no applications found")
}

// captureOutput swaps the test streams for buffers we can assert on separately.
func captureOutput(t *testing.T, io *iostreams.IOStreams) (*bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	io.Out = stdout
	io.ErrOut = stderr

	return stdout, stderr
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

func Test_runSelectCmd_SelectsByAppID(t *testing.T) {
	createHit := false
	srv := selectServer(t, &createHit)
	defer srv.Close()

	cfg := &test.ConfigStub{}
	opts := newSelectOptsWithSelector(t, srv, cfg, func(o *SelectOptions) {
		o.AppID = "APP1"
	})

	app, err := runSelectCmd(opts)
	require.NoError(t, err)
	require.NotNil(t, app)

	assert.Equal(t, "APP1", app.ID)
	assert.Equal(t, "My App", app.Name)
	assert.Equal(t, "new-key", cfg.SavedApps["APP1"].APIKey)
}

func Test_runSelectCmd_UnknownAppID(t *testing.T) {
	createHit := false
	srv := selectServer(t, &createHit)
	defer srv.Close()

	opts := newSelectOptsWithSelector(t, srv, &test.ConfigStub{}, func(o *SelectOptions) {
		o.AppID = "NOPE"
	})

	_, err := runSelectCmd(opts)
	require.Error(t, err)
	assert.Equal(t, `application with ID "NOPE" not found`, err.Error())
	assert.False(t, createHit)
}

func Test_runSelectCmd_RequiresSelectorWhenNonInteractive(t *testing.T) {
	createHit := false
	srv := selectServer(t, &createHit)
	defer srv.Close()

	// iostreams.Test() is not a TTY, so the picker is unavailable.
	opts := newSelectOptsWithSelector(t, srv, &test.ConfigStub{}, func(*SelectOptions) {})

	_, err := runSelectCmd(opts)
	require.Error(t, err)
	assert.Equal(t, "--app-id or --app-name is required in non-interactive mode", err.Error())
}
