package update

import (
	"context"
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

// newServer stubs the dashboard PATCH endpoint, echoing the requested name back
// in the response so the command's success output can be asserted.
func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		var payload dashboard.UpdateApplicationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.NoError(t, json.NewEncoder(w).Encode(dashboard.SingleApplicationResponse{
			Data: dashboard.ApplicationResource{
				ID:   "APP1",
				Type: "application",
				Attributes: dashboard.ApplicationAttributes{
					ApplicationID: "APP1",
					Name:          payload.Name,
				},
			},
		}))
	})
	return httptest.NewServer(mux)
}

func newOpts(
	t *testing.T,
	srv *httptest.Server,
	isTTY bool,
	output string,
) (*UpdateOptions, *test.CmdInOut) {
	t.Helper()
	seedToken(t)
	t.Setenv("ALGOLIA_APPLICATION_ID", "APP1")

	f, out := test.NewFactory(isTTY, nil, nil, "")
	pf := cmdutil.NewPrintFlags()
	*pf.OutputFormat = output
	pf.OutputFlagSpecified = func() bool { return output != "" }

	opts := &UpdateOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: pf,
		Name:       "Renamed App",
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
	return opts, out
}

func Test_runUpdateCmd(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	opts, out := newOpts(t, srv, true, "")
	require.NoError(t, runUpdateCmd(context.Background(), opts))

	got := out.String()
	assert.Contains(t, got, "APP1")
	assert.Contains(t, got, "Renamed App")
}

func Test_runUpdateCmd_outputJSON(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	opts, out := newOpts(t, srv, false, "json")
	require.NoError(t, runUpdateCmd(context.Background(), opts))

	got := out.String()
	assert.Contains(t, got, `"id":"APP1"`)
	assert.Contains(t, got, `"name":"Renamed App"`)
}

func TestNewUpdateCmd(t *testing.T) {
	f, _ := test.NewFactory(false, nil, nil, "")
	cmd := NewUpdateCmd(f)

	assert.Equal(t, "update", cmd.Name())
	assert.Equal(t, "true", cmd.Annotations["skipAuthCheck"])

	nameFlag := cmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag)
}
