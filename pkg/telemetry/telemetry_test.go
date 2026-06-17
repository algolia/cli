package telemetry

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/segmentio/analytics-go/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/version"
)

// captureTransport records the request it receives without hitting the network.
type captureTransport struct {
	req *http.Request
}

func (c *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.req = req
	return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
}

func TestEnvHeaderTransport_SetsHeaderWithoutMutatingRequest(t *testing.T) {
	capture := &captureTransport{}
	transport := &envHeaderTransport{base: capture, env: "prod"}

	req, err := http.NewRequest(http.MethodPost, "https://example.com/v1/batch", nil)
	require.NoError(t, err)

	_, err = transport.RoundTrip(req)
	require.NoError(t, err)

	assert.Equal(t, "prod", capture.req.Header.Get(envHeader))
	// RoundTrippers must not mutate the caller's request.
	assert.Empty(t, req.Header.Get(envHeader))
}

func TestTelemetryEnv(t *testing.T) {
	orig := version.Version
	t.Cleanup(func() { version.Version = orig })

	version.Version = "main"
	assert.Equal(t, "dev", telemetryEnv())

	version.Version = "1.20.0"
	assert.Equal(t, "prod", telemetryEnv())
}

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
	mu       sync.Mutex
	messages []analytics.Message
}

func (f *fakeAnalyticsClient) Enqueue(msg analytics.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
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

	require.NoError(t, client.Track(ctx, "Command Invoked", nil))
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

	require.NoError(t, client.Track(ctx, "Command Invoked", nil))
	require.Len(t, fake.messages, 1)

	track, ok := fake.messages[0].(analytics.Track)
	require.True(t, ok)
	assert.Empty(t, track.UserId)
}

func TestTrack_MergesCustomProperties(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	metadata.SetAppID("app-id")
	ctx := WithEventMetadata(context.Background(), metadata)

	require.NoError(t, client.Track(ctx, "CLI Auth Started", map[string]any{"flow": "login"}))
	require.Len(t, fake.messages, 1)

	track, ok := fake.messages[0].(analytics.Track)
	require.True(t, ok)
	assert.Equal(t, "login", track.Properties["flow"])
	assert.Equal(t, metadata.InvocationID, track.Properties["invocation_id"])
	assert.Equal(t, "app-id", track.Properties["app_id"])
}

func TestTrack_SequenceIsMonotonic(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	for i := 0; i < 3; i++ {
		require.NoError(t, client.Track(ctx, "Command Invoked", nil))
	}
	require.Len(t, fake.messages, 3)

	for i, msg := range fake.messages {
		track, ok := msg.(analytics.Track)
		require.True(t, ok)
		assert.Equal(t, int64(i+1), track.Properties["sequence"])
	}
}

func TestTrack_SequenceIsUniqueUnderConcurrency(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	const n = 100
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = client.Track(ctx, "Command Invoked", nil)
		}()
	}
	wg.Wait()

	require.Len(t, fake.messages, n)
	seen := make(map[int64]bool, n)
	for _, msg := range fake.messages {
		track, ok := msg.(analytics.Track)
		require.True(t, ok)
		seq, ok := track.Properties["sequence"].(int64)
		require.True(t, ok)
		assert.False(t, seen[seq], "duplicate sequence %d", seq)
		seen[seq] = true
	}
}

func TestTrack_CustomPropertiesCannotOverrideBase(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	require.NoError(t, client.Track(ctx, "Command Invoked", map[string]any{
		"invocation_id": "spoofed",
		"sequence":      int64(999),
	}))
	require.Len(t, fake.messages, 1)

	track, ok := fake.messages[0].(analytics.Track)
	require.True(t, ok)
	assert.Equal(t, metadata.InvocationID, track.Properties["invocation_id"])
	assert.Equal(t, int64(1), track.Properties["sequence"])
}
