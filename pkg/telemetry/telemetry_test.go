package telemetry

import (
	"context"
	"testing"

	"github.com/segmentio/analytics-go/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestSetUser(t *testing.T) {
	event := NewEventMetadata()
	event.SetUser("user-42", "user@test.com", "Test User")

	assert.Equal(t, "user-42", event.UserID)
	assert.Equal(t, "user@test.com", event.Email)
	assert.Equal(t, "Test User", event.Name)
}

// fakeAnalyticsClient captures the messages enqueued by the telemetry client so
// tests can assert on the payload without hitting the network.
type fakeAnalyticsClient struct {
	messages []analytics.Message
}

func (f *fakeAnalyticsClient) Enqueue(msg analytics.Message) error {
	f.messages = append(f.messages, msg)
	return nil
}

func (f *fakeAnalyticsClient) Close() error { return nil }

func TestIdentify_IncludesUserWhenAuthenticated(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	metadata.SetUser("user-42", "user@test.com", "Test User")
	ctx := WithEventMetadata(context.Background(), metadata)

	require.NoError(t, client.Identify(ctx))
	require.Len(t, fake.messages, 1)

	identify, ok := fake.messages[0].(analytics.Identify)
	require.True(t, ok)
	assert.Equal(t, "user-42", identify.UserId)
	assert.Equal(t, metadata.AnonymousID, identify.AnonymousId)
	assert.Equal(t, "user@test.com", identify.Traits["email"])
	assert.Equal(t, "Test User", identify.Traits["name"])
}

func TestIdentify_OmitsUserWhenAnonymous(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	require.NoError(t, client.Identify(ctx))
	require.Len(t, fake.messages, 1)

	identify, ok := fake.messages[0].(analytics.Identify)
	require.True(t, ok)
	assert.Empty(t, identify.UserId)
	assert.NotContains(t, identify.Traits, "email")
	assert.NotContains(t, identify.Traits, "name")
}

func TestTrack_IncludesUserWhenAuthenticated(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	metadata.SetUser("user-42", "user@test.com", "Test User")
	ctx := WithEventMetadata(context.Background(), metadata)

	require.NoError(t, client.Track(ctx, "Command Invoked"))
	require.Len(t, fake.messages, 1)

	track, ok := fake.messages[0].(analytics.Track)
	require.True(t, ok)
	assert.Equal(t, "user-42", track.UserId)
	assert.Equal(t, "Command Invoked", track.Event)
}

func TestTrack_OmitsUserWhenAnonymous(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	require.NoError(t, client.Track(ctx, "Command Invoked"))
	require.Len(t, fake.messages, 1)

	track, ok := fake.messages[0].(analytics.Track)
	require.True(t, ok)
	assert.Empty(t, track.UserId)
}
