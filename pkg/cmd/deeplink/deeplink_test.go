package deeplink

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestDeeplinkURL(t *testing.T) {
	base := "https://dashboard.algolia.com"
	tests := []struct {
		purpose string
		want    string
	}{
		{"dashboard", "https://dashboard.algolia.com/apps/APP123/dashboard"},
		{"indices", "https://dashboard.algolia.com/apps/APP123/explorer/browse"},
		{"crawler", "https://dashboard.algolia.com/apps/APP123/crawler"},
		{"connectors", "https://dashboard.algolia.com/apps/APP123/connectors"},
		{"api-keys", "https://dashboard.algolia.com/account/api-keys/all?applicationId=APP123"},
		{"usage", "https://dashboard.algolia.com/account/billing/usage?applicationId=APP123"},
		{"team", "https://dashboard.algolia.com/account/teams?applicationId=APP123"},
		{"billing", "https://dashboard.algolia.com/account/billing/details?applicationId=APP123"},
	}

	for _, tt := range tests {
		t.Run(tt.purpose, func(t *testing.T) {
			assert.Equal(t, tt.want, deeplinkURL(base, "APP123", tt.purpose))
		})
	}
}

// TestPurposeOrderMatchesTargets guards against the ordered list and the
// target map drifting apart.
func TestPurposeOrderMatchesTargets(t *testing.T) {
	assert.Len(t, purposeOrder, len(purposeTargets))
	for _, p := range purposeOrder {
		_, ok := purposeTargets[p]
		assert.Truef(t, ok, "purpose %q listed in order but has no target", p)
	}
}

func newTestOptions(
	io *iostreams.IOStreams,
	cfg config.IConfig,
) (*DeeplinkOptions, *string, *bool) {
	opened := new(string)
	authed := new(bool)

	opts := &DeeplinkOptions{
		IO:     io,
		Config: cfg,
		Authenticate: func(_ *iostreams.IOStreams, _ *dashboard.Client) (string, error) {
			*authed = true
			return "test-token", nil
		},
		SelectApplication: func() (*dashboard.Application, error) {
			return nil, errors.New("SelectApplication should not be called")
		},
		NewDashboardClient: func(string) *dashboard.Client {
			return &dashboard.Client{DashboardURL: "https://dashboard.algolia.com"}
		},
		Browser: func(u string) error {
			*opened = u
			return nil
		},
	}

	return opts, opened, authed
}

func TestRunDeeplinkCmd_ConfiguredApp(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, authed := newTestOptions(io, cfg)
	opts.Purpose = "api-keys"

	err := runDeeplinkCmd(opts)
	require.NoError(t, err)
	assert.True(t, *authed, "expected sign-in to be required")
	assert.Equal(
		t,
		"https://dashboard.algolia.com/account/api-keys/all?applicationId=APP123",
		*opened,
	)
}

func TestRunDeeplinkCmd_UsesConfiguredDashboardURL(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, _ := newTestOptions(io, cfg)
	opts.Purpose = "usage"
	opts.NewDashboardClient = func(string) *dashboard.Client {
		return &dashboard.Client{DashboardURL: "https://staging.algolia.test"}
	}

	err := runDeeplinkCmd(opts)
	require.NoError(t, err)
	assert.Equal(
		t,
		"https://staging.algolia.test/account/billing/usage?applicationId=APP123",
		*opened,
	)
}

func TestRunDeeplinkCmd_SelectsAppWhenNoneConfigured(t *testing.T) {
	t.Setenv("ALGOLIA_APPLICATION_ID", "")

	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "", Default: true},
	})

	opts, opened, _ := newTestOptions(io, cfg)
	opts.Purpose = "dashboard"
	opts.SelectApplication = func() (*dashboard.Application, error) {
		return &dashboard.Application{ID: "SELECTED", Name: "Picked"}, nil
	}

	err := runDeeplinkCmd(opts)
	require.NoError(t, err)
	assert.Equal(t, "https://dashboard.algolia.com/apps/SELECTED/dashboard", *opened)
}

func TestRunDeeplinkCmd_InvalidPurpose(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	opts, opened, authed := newTestOptions(io, test.NewDefaultConfigStub())
	opts.Purpose = "bogus"

	err := runDeeplinkCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid purpose")
	assert.False(t, *authed, "should fail before sign-in")
	assert.Empty(t, *opened, "browser should not be opened")
}

func TestRunDeeplinkCmd_NoPurposeNonInteractive(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(false)
	io.SetStdinTTY(false)

	opts, opened, _ := newTestOptions(io, test.NewDefaultConfigStub())
	opts.Purpose = ""

	err := runDeeplinkCmd(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--purpose is required")
	assert.Empty(t, *opened, "browser should not be opened")
}
