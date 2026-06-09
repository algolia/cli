package open

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestDashboardURL(t *testing.T) {
	base := "https://dashboard.algolia.com"
	tests := []struct {
		shortcut string
		want     string
	}{
		{"dashboard", "https://dashboard.algolia.com/apps/APP123/dashboard"},
		{"indices", "https://dashboard.algolia.com/apps/APP123/explorer/browse"},
		{"crawler", "https://dashboard.algolia.com/apps/APP123/crawler"},
		{"connectors", "https://dashboard.algolia.com/apps/APP123/connectors"},
		{"api-keys", "https://dashboard.algolia.com/account/api-keys/all?applicationId=APP123"},
		{"usage", "https://dashboard.algolia.com/account/billing/usage?applicationId=APP123"},
		{"team", "https://dashboard.algolia.com/account/teams?applicationId=APP123"},
		{"billing", "https://dashboard.algolia.com/account/billing/details?applicationId=APP123"},
		{"cost-management", "https://dashboard.algolia.com/account/billing/cost-management?applicationId=APP123"},
	}

	for _, tt := range tests {
		t.Run(tt.shortcut, func(t *testing.T) {
			assert.Equal(t, tt.want, dashboardURL(base, "APP123", dashboardTargets[tt.shortcut]))
		})
	}
}

// TestTargetNamesNoCollision guards against a shortcut existing in both maps,
// which would make its behavior ambiguous.
func TestTargetNamesNoCollision(t *testing.T) {
	for name := range dashboardTargets {
		_, dup := resourceURLs[name]
		assert.Falsef(t, dup, "shortcut %q exists in both resourceURLs and dashboardTargets", name)
	}
}

func newTestOptions(
	io *iostreams.IOStreams,
	cfg config.IConfig,
) (*OpenOptions, *string, *bool) {
	opened := new(string)
	authed := new(bool)

	opts := &OpenOptions{
		IO:         io,
		config:     cfg,
		PrintFlags: cmdutil.NewPrintFlags(),
		Authenticate: func(_ context.Context, _ *iostreams.IOStreams, _ *dashboard.Client) (string, error) {
			*authed = true
			return "test-token", nil
		},
		SelectApplication: func(context.Context) (*dashboard.Application, error) {
			return nil, errors.New("SelectApplication should not be called")
		},
		NewDashboardClient: func(string) *dashboard.Client {
			return &dashboard.Client{DashboardURL: "https://dashboard.algolia.com"}
		},
		ListApplications: func(_ *dashboard.Client, _ string) ([]dashboard.Application, error) {
			return []dashboard.Application{{ID: "APP123", Name: "Test"}}, nil
		},
		Browser: func(u string) error {
			*opened = u
			return nil
		},
	}

	return opts, opened, authed
}

// withOutputFormat configures opts to emit structured output, as the --output
// flag would when parsed by cobra.
func withOutputFormat(opts *OpenOptions, format string) {
	*opts.PrintFlags.OutputFormat = format
	opts.PrintFlags.OutputFlagSpecified = func() bool { return true }
}

func TestRunOpenCmd_ResourceShortcutNoAuth(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts, opened, authed := newTestOptions(io, test.NewDefaultConfigStub())
	opts.Shortcut = "docs"

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.False(t, *authed, "resource shortcuts must not require sign-in")
	assert.Equal(t, "https://algolia.com/doc/", *opened)
}

func TestRunOpenCmd_ResourceShortcutWithAppID(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, authed := newTestOptions(io, cfg)
	opts.Shortcut = "status"
	opts.NewDashboardClient = func(string) *dashboard.Client {
		return &dashboard.Client{DashboardURL: "https://staging.algolia.test"}
	}

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.True(t, *authed, "app-scoped status validates the profile application against the signed-in account")
	assert.Equal(t, "https://staging.algolia.test/apps/APP123/monitoring/status", *opened)
}

func TestRunOpenCmd_ResourceShortcutNoAppUsesDefault(t *testing.T) {
	t.Setenv("ALGOLIA_APPLICATION_ID", "")

	io, _, _, _ := iostreams.Test()

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "", Default: true},
	})

	opts, opened, _ := newTestOptions(io, cfg)
	opts.Shortcut = "status"

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, "https://status.algolia.com/", *opened)
}

func TestRunOpenCmd_DashboardTargetConfiguredApp(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, authed := newTestOptions(io, cfg)
	opts.Shortcut = "billing"

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.True(t, *authed, "application pages require sign-in")
	assert.Equal(
		t,
		"https://dashboard.algolia.com/account/billing/details?applicationId=APP123",
		*opened,
	)
}

func TestRunOpenCmd_DashboardTargetAppScoped(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, _ := newTestOptions(io, cfg)
	opts.Shortcut = "dashboard"

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, "https://dashboard.algolia.com/apps/APP123/dashboard", *opened)
}

func TestRunOpenCmd_DashboardTargetUsesConfiguredDashboardURL(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, _ := newTestOptions(io, cfg)
	opts.Shortcut = "usage"
	opts.NewDashboardClient = func(string) *dashboard.Client {
		return &dashboard.Client{DashboardURL: "https://staging.algolia.test"}
	}

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(
		t,
		"https://staging.algolia.test/account/billing/usage?applicationId=APP123",
		*opened,
	)
}

func TestRunOpenCmd_DashboardTargetIgnoresStaleProfileApp(t *testing.T) {
	t.Setenv("ALGOLIA_APPLICATION_ID", "")

	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "USER_A_APP", APIKey: "key", Default: true},
	})

	selectCalled := false
	opts, opened, _ := newTestOptions(io, cfg)
	opts.Shortcut = "billing"
	opts.ListApplications = func(_ *dashboard.Client, _ string) ([]dashboard.Application, error) {
		return []dashboard.Application{{ID: "USER_B_APP", Name: "User B App"}}, nil
	}
	opts.SelectApplication = func(context.Context) (*dashboard.Application, error) {
		selectCalled = true
		return &dashboard.Application{ID: "USER_B_APP", Name: "User B App"}, nil
	}

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.True(t, selectCalled, "stale profile application must trigger selection")
	assert.Equal(
		t,
		"https://dashboard.algolia.com/account/billing/details?applicationId=USER_B_APP",
		*opened,
	)
}

func TestRunOpenCmd_DashboardTargetSelectsAppWhenNoneConfigured(t *testing.T) {
	t.Setenv("ALGOLIA_APPLICATION_ID", "")

	io, _, _, _ := iostreams.Test()
	io.SetStdoutTTY(true)
	io.SetStdinTTY(true)

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "", Default: true},
	})

	opts, opened, _ := newTestOptions(io, cfg)
	opts.Shortcut = "dashboard"
	opts.SelectApplication = func(context.Context) (*dashboard.Application, error) {
		return &dashboard.Application{ID: "SELECTED", Name: "Picked"}, nil
	}

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)
	assert.Equal(t, "https://dashboard.algolia.com/apps/SELECTED/dashboard", *opened)
}

func TestRunOpenCmd_Unsupported(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts, opened, authed := newTestOptions(io, test.NewDefaultConfigStub())
	opts.Shortcut = "bogus"

	err := runOpenCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported open command, given: bogus")
	assert.Contains(t, err.Error(), "Available shortcuts:")
	assert.Contains(t, err.Error(), "billing")
	assert.Contains(t, err.Error(), "docs")
	assert.False(t, *authed)
	assert.Empty(t, *opened)
}

func TestRunOpenCmd_ListIncludesBothKinds(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()

	opts, _, _ := newTestOptions(io, test.NewDefaultConfigStub())
	opts.List = true

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)

	out := stdout.String()
	assert.Contains(t, out, "docs")
	assert.Contains(t, out, "billing")
}

func TestRunOpenCmd_ListJSONOutput(t *testing.T) {
	t.Setenv("ALGOLIA_APPLICATION_ID", "")

	io, _, stdout, _ := iostreams.Test()

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, authed := newTestOptions(io, cfg)
	opts.List = true
	withOutputFormat(opts, "json")

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)

	var entries []pageEntry
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &entries))
	assert.Len(t, entries, len(resourceURLs)+len(dashboardTargets))

	byName := make(map[string]pageEntry, len(entries))
	for _, e := range entries {
		byName[e.Shortcut] = e
	}

	assert.Equal(
		t,
		pageEntry{
			Shortcut:      "billing",
			URL:           "https://dashboard.algolia.com/account/billing/details?applicationId=APP123",
			RequiresLogin: true,
		},
		byName["billing"],
	)
	assert.Equal(
		t,
		pageEntry{Shortcut: "docs", URL: "https://algolia.com/doc/"},
		byName["docs"],
	)

	// Structured output never opens a browser or signs in.
	assert.False(t, *authed)
	assert.Empty(t, *opened)
}

func TestRunOpenCmd_SingleShortcutJSONOutput(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()

	cfg := test.NewConfigStubWithProfiles([]*config.Profile{
		{Name: "default", ApplicationID: "APP123", APIKey: "key", Default: true},
	})

	opts, opened, authed := newTestOptions(io, cfg)
	opts.Shortcut = "billing"
	withOutputFormat(opts, "json")

	err := runOpenCmd(context.Background(), opts)
	require.NoError(t, err)

	var entry pageEntry
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &entry))
	assert.Equal(
		t,
		pageEntry{
			Shortcut:      "billing",
			URL:           "https://dashboard.algolia.com/account/billing/details?applicationId=APP123",
			RequiresLogin: true,
		},
		entry,
	)

	// A dashboard target with --output does not sign in or open a browser.
	assert.False(t, *authed)
	assert.Empty(t, *opened)
}

func TestRunOpenCmd_UnsupportedWithJSONOutput(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts, _, _ := newTestOptions(io, test.NewDefaultConfigStub())
	opts.Shortcut = "bogus"
	withOutputFormat(opts, "json")

	err := runOpenCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported open command, given: bogus")
	assert.Contains(t, err.Error(), "Available shortcuts:")
}

func TestRunOpenCmd_InvalidOutputFormat(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts, _, _ := newTestOptions(io, test.NewDefaultConfigStub())
	opts.List = true
	withOutputFormat(opts, "yaml")

	err := runOpenCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to match a printer")
}
