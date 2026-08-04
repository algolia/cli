package login

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/pkg/telemetry/telemetrytest"
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

func TestTrackOAuthFlowOutcome(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantEvent string
		wantStep  bool
	}{
		{
			name:      "success",
			err:       nil,
			wantEvent: telemetry.EventAuthCompleted,
		},
		{
			name:      "user cancellation",
			err:       cmdutil.ErrCancel,
			wantEvent: telemetry.EventAuthAborted,
			wantStep:  true,
		},
		{
			name:      "failure",
			err:       errors.New("boom"),
			wantEvent: telemetry.EventAuthFailed,
			wantStep:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &telemetrytest.RecordingClient{}
			ctx := telemetry.WithTelemetryClient(context.Background(), client)
			tracker := telemetry.NewFlowTracker()
			tracker.SetStep(telemetry.StepAppsFetch)

			trackOAuthFlowOutcome(ctx, telemetry.FlowLogin, tracker, tt.err)

			require.Len(t, client.Events, 1)
			event := client.Events[0]
			assert.Equal(t, tt.wantEvent, event.Name)
			assert.Equal(t, telemetry.FlowLogin, event.Properties["flow"])
			if tt.wantStep {
				assert.Equal(t, telemetry.StepAppsFetch, event.Properties["step"])
			}
		})
	}
}

// `auth signup` shares LoginOptions but exposes no output flags.
func TestPrintLoginResult_NilPrintFlags(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{IO: io, NonInteractive: true}
	cmdutil.ApplyNonInteractive(opts.IO, opts.PrintFlags)

	assert.NoError(t, printLoginResult(opts, &dashboard.Application{ID: "APP1"}))
}

func TestPrintLoginResult_JSONWithoutAPIKey(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "access",
		ExpiresIn:   3600,
		User:        &dashboard.User{ID: 42, Email: "user@example.com"},
	}))

	io, _, stdout, _ := iostreams.Test()
	opts := &LoginOptions{
		IO: io,
		Config: &test.ConfigStub{
			SavedApps: map[string]test.SavedApplication{"APP1": {Alias: "my app"}},
		},
		PrintFlags:     cmdutil.NewPrintFlags(),
		NonInteractive: true,
	}
	cmdutil.ApplyNonInteractive(opts.IO, opts.PrintFlags)

	app := &dashboard.Application{
		ID:        "APP1",
		Name:      "My App",
		APIKey:    "secret-key",
		PlanLabel: "Grow",
	}
	require.NoError(t, printLoginResult(opts, app))

	var got loginResult
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &got), "stdout: %q", stdout.String())
	assert.True(t, got.Success)
	assert.Equal(t, "user@example.com", got.Email)
	require.NotNil(t, got.Application)
	assert.Equal(t, apputil.ApplicationOutput{
		ID:    "APP1",
		Alias: "my app",
		Name:  "My App",
		Plan:  "Grow",
	}, *got.Application)
	assert.NotContains(t, stdout.String(), "secret-key")
}

// Non-interactive login signs in only, so the document reports success without
// an application.
func TestPrintLoginResult_NoApplication(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(auth.ClearToken)
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "access",
		ExpiresIn:   3600,
		User:        &dashboard.User{ID: 42, Email: "user@example.com"},
	}))

	io, _, stdout, stderr := iostreams.Test()
	opts := &LoginOptions{
		IO:             io,
		Config:         &test.ConfigStub{},
		PrintFlags:     cmdutil.NewPrintFlags(),
		NonInteractive: true,
	}
	cmdutil.ApplyNonInteractive(opts.IO, opts.PrintFlags)

	require.NoError(t, printLoginResult(opts, nil))

	assert.JSONEq(t, `{"success":true,"email":"user@example.com"}`, stdout.String())
	assert.Empty(t, stderr.String())
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
