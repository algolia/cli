package rotate

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

// rotateServer stubs the dashboard rotate endpoint at wantPath, returning
// newValue as the rotated key's value.
func rotateServer(t *testing.T, wantPath, newValue string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, json.NewEncoder(w).Encode(dashboard.CreateAPIKeyResponse{
			Data: dashboard.APIKeyResource{
				ID:         "key-uuid-123",
				Type:       "api_key",
				Attributes: dashboard.APIKeyAttributes{Value: newValue},
			},
		}))
	})
	return httptest.NewServer(mux)
}

func newRotateOpts(
	t *testing.T,
	srv *httptest.Server,
	cfg config.IConfig,
	isTTY bool,
) (*RotateOptions, *bytes.Buffer) {
	t.Helper()
	seedToken(t)

	io, _, stdout, _ := iostreams.Test()
	io.SetStdoutTTY(isTTY)

	opts := &RotateOptions{
		IO:     io,
		Config: cfg,
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
	return opts, stdout
}

func TestNewRotateCmd(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: io}
	cmd := NewRotateCmd(f, nil)

	assert.Equal(t, "rotate", cmd.Name())
	assert.Equal(t, "true", cmd.Annotations["skipAuthCheck"])
	assert.Nil(t, cmd.Flags().Lookup("key-id"))
}

func Test_runRotateCmd_ResolutionErrors(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *test.ConfigStub
		wantErr string
	}{
		{
			name:    "no active application",
			cfg:     &test.ConfigStub{ActiveAppID: ""},
			wantErr: "no current application selected",
		},
		{
			name: "active application without a stored UUID",
			cfg: &test.ConfigStub{
				ActiveAppID: "APP1",
				SavedApps: map[string]test.SavedApplication{
					"APP1": {Alias: "prod", APIKey: "old-key"},
				},
			},
			wantErr: "no CLI-managed API key found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			opts := &RotateOptions{IO: io, Config: tt.cfg}

			err := runRotateCmd(opts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func Test_runRotateCmd_PersistsCLIManagedKey(t *testing.T) {
	srv := rotateServer(t, "/1/applications/APP1/api-keys/key-uuid-123/rotate", "rotated-key")
	defer srv.Close()

	cfg := &test.ConfigStub{
		ActiveAppID: "APP1",
		SavedApps: map[string]test.SavedApplication{
			"APP1": {Alias: "prod", APIKeyUUID: "key-uuid-123", APIKey: "old-key"},
		},
	}
	opts, stdout := newRotateOpts(t, srv, cfg, true)

	require.NoError(t, runRotateCmd(opts))

	assert.Contains(t, stdout.String(), "rotated-key")
	// The CLI-managed key was rotated, so the new value replaces the stored one.
	assert.Equal(t, "rotated-key", cfg.SavedApps["APP1"].APIKey)
}
