package logout

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/pkg/telemetry/telemetrytest"
)

func newLogoutOpts(srv *httptest.Server) *LogoutOptions {
	io, _, _, _ := iostreams.Test()
	return &LogoutOptions{
		IO: io,
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.DashboardURL = srv.URL
			return c
		},
	}
}

func TestLogout_EmitsAuthLogoutBeforeClearingToken(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresIn:    3600,
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/revoke", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := &telemetrytest.RecordingClient{}
	ctx := telemetry.WithTelemetryClient(context.Background(), client)

	require.NoError(t, runLogoutCmd(ctx, newLogoutOpts(srv)))

	require.Len(t, client.Events, 1)
	event := client.Events[0]
	assert.Equal(t, telemetry.EventAuthLogout, event.Name)
	assert.Equal(t, telemetry.FlowLogout, event.Properties["flow"])
	// The token is cleared after the event was tracked.
	assert.Nil(t, auth.LoadToken())
	// No Identify at logout time: Segment cannot un-identify.
	assert.Zero(t, client.Identifies)
}

func TestLogout_AlreadySignedOutEmitsNothing(t *testing.T) {
	keyring.MockInit()
	auth.ClearToken()

	client := &telemetrytest.RecordingClient{}
	ctx := telemetry.WithTelemetryClient(context.Background(), client)

	require.NoError(t, runLogoutCmd(ctx, newLogoutOpts(nil)))

	assert.Empty(t, client.Events)
}
