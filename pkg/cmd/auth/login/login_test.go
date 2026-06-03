package login

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/test"
)

func TestNewLoginCmd_FlagParsing(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.SetStdinTTY(false)
	io.SetStdoutTTY(false)

	f := &cmdutil.Factory{
		IOStreams: io,
		Config:    test.NewDefaultConfigStub(),
	}

	cmd := NewLoginCmd(f)
	args := []string{
		"--app-name", "My App",
		"--profile-name", "myprofile",
	}
	err := cmd.ParseFlags(args)
	require.NoError(t, err)
}

func TestSelectApplication_SingleApp(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{IO: io}

	apps := []dashboard.Application{
		{ID: "APP1", Name: "My App"},
	}

	app, err := selectApplication(opts, apps, false)
	require.NoError(t, err)
	assert.Equal(t, "APP1", app.ID)
}

func TestSelectApplication_ByName(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{IO: io, AppName: "Second App"}

	apps := []dashboard.Application{
		{ID: "APP1", Name: "First App"},
		{ID: "APP2", Name: "Second App"},
	}

	app, err := selectApplication(opts, apps, false)
	require.NoError(t, err)
	assert.Equal(t, "APP2", app.ID)
}

func TestSelectApplication_ByName_NotFound(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{IO: io, AppName: "Unknown"}

	apps := []dashboard.Application{
		{ID: "APP1", Name: "First App"},
	}

	_, err := selectApplication(opts, apps, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestApplyStoredIdentity_SetsMetadataFromToken(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)

	err := auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresIn:    3600,
		User: &dashboard.User{
			ID:    42,
			Email: "user@example.com",
			Name:  "Ada Lovelace",
		},
	})
	require.NoError(t, err)

	metadata := telemetry.NewEventMetadata()
	ctx := telemetry.WithEventMetadata(context.Background(), metadata)

	assert.True(t, applyStoredIdentity(ctx))
	assert.Equal(t, "42", metadata.UserID)
	assert.Equal(t, "user@example.com", metadata.Email)
	assert.Equal(t, "Ada Lovelace", metadata.Name)
}

func TestApplyStoredIdentity_NoTokenReturnsFalse(t *testing.T) {
	keyring.MockInit()
	auth.ClearToken()

	metadata := telemetry.NewEventMetadata()
	ctx := telemetry.WithEventMetadata(context.Background(), metadata)

	assert.False(t, applyStoredIdentity(ctx))
	assert.Empty(t, metadata.UserID)
}

func TestApplyStoredIdentity_TokenWithoutIdentityReturnsFalse(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)

	// Token persisted before identity was tracked (no user object).
	err := auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresIn:    3600,
	})
	require.NoError(t, err)

	metadata := telemetry.NewEventMetadata()
	ctx := telemetry.WithEventMetadata(context.Background(), metadata)

	assert.False(t, applyStoredIdentity(ctx))
	assert.Empty(t, metadata.UserID)
}

func TestSelectApplication_MultipleApps_NonInteractive_NoAppName(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{IO: io}

	apps := []dashboard.Application{
		{ID: "APP1", Name: "First"},
		{ID: "APP2", Name: "Second"},
	}

	_, err := selectApplication(opts, apps, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple applications found")
}
