package telemetry

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"gopkg.in/segmentio/analytics-go.v3"
)

// Context-related tests.
func TestEventMetadataWithGet(t *testing.T) {
	ctx := context.Background()
	event := &CLIAnalyticsEventMetadata{
		UserID:                   "user-id",
		InvocationID:             "invocation-id",
		OS:                       "os",
		CLIVersion:               "cli-version",
		CommandPath:              "command-path",
		CommandFlags:             []string{"flag1", "flag2"},
		AppID:                    "app-id",
		ConfiguredApplicationsNb: 1,
	}
	newCtx := WithEventMetadata(ctx, event)

	// Check that the event is correctly set in the context.
	require.Equal(t, event, GetEventMetadata(newCtx))
}

func TestEventMetadata_DoesNotExistsInContext(t *testing.T) {
	ctx := context.Background()
	require.Nil(t, GetEventMetadata(ctx))
}

func TestTelemetryClientWithGet(t *testing.T) {
	ctx := context.Background()

	client, err := analytics.NewWithConfig("", analytics.Config{
		Endpoint: "http://hello.com",
	})
	require.NoError(t, err)

	telemetryClient := &AnalyticsTelemetryClient{client: client}
	newCtx := WithTelemetryClient(ctx, telemetryClient)

	require.Equal(t, GetTelemetryClient(newCtx), telemetryClient)
}

func TestSetCobraCommandContext(t *testing.T) {
	event := NewEventMetadata()
	cmd := &cobra.Command{
		Use: "foo",
	}
	cmd.Flags().String("bar", "bar", "bar flag")
	cmd.SetArgs([]string{"--bar", "bar"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	event.SetCobraCommandContext(cmd)

	require.Equal(t, "foo", event.CommandPath)
	require.Equal(t, []string{"bar"}, event.CommandFlags)
}
