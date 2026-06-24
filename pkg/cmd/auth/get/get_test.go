package get

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

// roundTripFunc lets a test stub the HTTP transport of the dashboard client.
type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// cmdWithOpts wires runGetCmd to a custom GetOptions so tests can stub the
// auth seam (EnsureAuthenticated) or the dashboard client.
func cmdWithOpts(opts *GetOptions) *cobra.Command {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

// When there's no usable session, auth get launches the login flow (here
// stubbed) and proceeds with the resulting identity.
func TestGet_PromptsLoginWhenNoSession(t *testing.T) {
	keyring.MockInit()
	auth.ClearToken()
	t.Cleanup(auth.ClearToken)

	f, out := test.NewFactory(false, nil, nil, "")
	called := false
	opts := &GetOptions{
		IO:                 f.IOStreams,
		LoadToken:          auth.LoadToken,
		PrintFlags:         cmdutil.NewPrintFlags().WithDefaultOutput("json"),
		NewDashboardClient: func(clientID string) *dashboard.Client { return nil },
		EnsureAuthenticated: func(_ *iostreams.IOStreams, _ *dashboard.Client) (string, error) {
			called = true
			// Simulate a successful browser login persisting a session.
			require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
				AccessToken: "fresh-access",
				CreatedAt:   time.Now().Unix(),
				ExpiresIn:   3600,
				User:        &dashboard.User{ID: 7, Email: "new@example.com", Name: "New User"},
			}))
			return "fresh-access", nil
		},
	}

	out, err := test.Execute(cmdWithOpts(opts), "--output ndjson", out)
	require.NoError(t, err)
	assert.True(t, called, "expected login flow to be triggered")
	assert.Contains(t, out.String(), `"user_id":"7"`)
	assert.Contains(t, out.String(), `"email":"new@example.com"`)
}

// If the login flow fails (e.g. the user aborts), the error is propagated.
func TestGet_ReturnsErrorWhenLoginFails(t *testing.T) {
	keyring.MockInit()
	auth.ClearToken()

	f, out := test.NewFactory(false, nil, nil, "")
	opts := &GetOptions{
		IO:                 f.IOStreams,
		LoadToken:          auth.LoadToken,
		PrintFlags:         cmdutil.NewPrintFlags().WithDefaultOutput("json"),
		NewDashboardClient: func(clientID string) *dashboard.Client { return nil },
		EnsureAuthenticated: func(_ *iostreams.IOStreams, _ *dashboard.Client) (string, error) {
			return "", fmt.Errorf("authorization failed: access_denied")
		},
	}

	_, err := test.Execute(cmdWithOpts(opts), "", out)
	require.Error(t, err)
	assert.Equal(t, "authorization failed: access_denied", err.Error())
}

func TestGet_RefreshesExpiredToken(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken:  "old-access",
		RefreshToken: "valid-refresh",
		CreatedAt:    time.Now().Unix() - 7200,
		ExpiresIn:    3600,
		User: &dashboard.User{
			ID:    42,
			Email: "user@example.com",
			Name:  "Test User",
		},
	}))

	body := `{"access_token":"new-access","refresh_token":"new-refresh","expires_in":3600}`
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		}),
	}

	f, out := test.NewFactory(false, nil, nil, "")
	opts := &GetOptions{
		IO:         f.IOStreams,
		LoadToken:  auth.LoadToken,
		PrintFlags: cmdutil.NewPrintFlags().WithDefaultOutput("json"),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClientWithHTTPClient(clientID, httpClient)
		},
		// Real auth seam: GetValidToken refreshes via the stubbed client and
		// succeeds, so the browser flow is never reached.
		EnsureAuthenticated: auth.EnsureAuthenticated,
	}

	out, err := test.Execute(cmdWithOpts(opts), "--output ndjson", out)
	require.NoError(t, err)

	// Identity preserved from the pre-refresh token (refresh response has no user).
	assert.Contains(t, out.String(), `"user_id":"42"`)
	assert.Contains(t, out.String(), `"email":"user@example.com"`)
	assert.NotContains(t, out.String(), "new-access")

	// Refreshed token was persisted.
	assert.Equal(t, "new-access", auth.LoadToken().AccessToken)
}

func TestGet_PrintsIdentityWithoutTokens(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken:  "secret-access",
		RefreshToken: "secret-refresh",
		CreatedAt:    time.Now().Unix(),
		ExpiresIn:    3600,
		User: &dashboard.User{
			ID:    42,
			Email: "user@example.com",
			Name:  "Test User",
		},
	}))

	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewGetCmd(f)
	out, err := test.Execute(cmd, "--output ndjson", out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), `"user_id":"42"`)
	assert.Contains(t, out.String(), `"email":"user@example.com"`)
	assert.Contains(t, out.String(), `"name":"Test User"`)
	assert.NotContains(t, out.String(), "secret-access")
	assert.NotContains(t, out.String(), "secret-refresh")
	assert.NotContains(t, out.String(), "access_token")
	assert.NotContains(t, out.String(), "refresh_token")
}
