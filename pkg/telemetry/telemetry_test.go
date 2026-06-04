package telemetry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

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

func TestTrack_TimestampsAreStrictlyIncreasing(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	// Back-to-back events land in the same millisecond but must still get
	// strictly increasing timestamps so Amplitude preserves emit order.
	const n = 5
	for i := 0; i < n; i++ {
		require.NoError(t, client.Track(ctx, "Event", nil))
	}
	require.Len(t, fake.messages, n)

	var prev time.Time
	for i, m := range fake.messages {
		ts := m.(analytics.Track).Timestamp
		if i > 0 {
			assert.True(
				t,
				ts.After(prev),
				"timestamp %d (%v) must be strictly after previous (%v)",
				i, ts, prev,
			)
		}
		prev = ts
	}
}

// collectBatches starts a server recording each batch POST's event names in order.
func collectBatches(t *testing.T) (url string, batches *[][]string) {
	t.Helper()

	var (
		mu  sync.Mutex
		got [][]string
	)
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload struct {
				Batch []struct {
					Event string `json:"event"`
				} `json:"batch"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))

			names := make([]string, 0, len(payload.Batch))
			for _, m := range payload.Batch {
				names = append(names, m.Event)
			}

			mu.Lock()
			got = append(got, names)
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
		}),
	)
	t.Cleanup(srv.Close)

	return srv.URL, &got
}

// TestAnalyticsClient_SendsAllEventsInOneOrderedBatch is the regression test for
// the "events out of order / sometimes missing" bug: every event must reach the
// backend in one ordered batch, not split across the library's periodic flushes.
func TestAnalyticsClient_SendsAllEventsInOneOrderedBatch(t *testing.T) {
	url, batches := collectBatches(t)

	client, err := newAnalyticsTelemetryClient(url, false)
	require.NoError(t, err)

	metadata := NewEventMetadata()
	metadata.AnonymousID = "anon-test" // ensure messages validate without a MAC
	ctx := WithEventMetadata(context.Background(), metadata)

	want := []string{
		EventCommandInvoked,
		EventAuthStarted,
		EventAuthBrowserOpened,
		EventAuthCallbackReceived,
		EventAuthCompleted,
		EventCommandCompleted,
	}
	for _, name := range want {
		require.NoError(t, client.Track(ctx, name, nil))
		// A real command emits these over time; the gap must not trigger a flush.
		time.Sleep(2 * time.Millisecond)
	}

	// Close is the single flush point and must block until the batch is sent.
	client.Close()

	require.Len(t, *batches, 1, "all events must be delivered in exactly one batch")
	assert.Equal(t, want, (*batches)[0], "events must arrive in emit order")
}

// TestBoundedRoundTripper_TimesOut verifies the flush request is bounded, so
// closing the client at exit can never hang the CLI on a stalled endpoint.
func TestBoundedRoundTripper_TimesOut(t *testing.T) {
	blocked := make(chan struct{})
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-blocked // never respond until the test tears the server down
		}),
	)
	t.Cleanup(func() {
		close(blocked)
		srv.Close()
	})

	rt := boundedRoundTripper{base: http.DefaultTransport, timeout: 50 * time.Millisecond}
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := rt.RoundTrip(req)
	elapsed := time.Since(start)

	require.Error(t, err, "a stalled endpoint must surface as an error, not a hang")
	if resp != nil {
		_ = resp.Body.Close()
	}
	assert.Less(t, elapsed, time.Second, "request must be abandoned near the configured timeout")
}

func TestTrack_MergesCustomProperties(t *testing.T) {
	fake := &fakeAnalyticsClient{}
	client := &AnalyticsTelemetryClient{client: fake}

	metadata := NewEventMetadata()
	ctx := WithEventMetadata(context.Background(), metadata)

	props := map[string]any{"flow": "signup", "duration_ms": int64(1200)}
	require.NoError(t, client.Track(ctx, EventAuthCompleted, props))
	require.Len(t, fake.messages, 1)

	track, ok := fake.messages[0].(analytics.Track)
	require.True(t, ok)
	assert.Equal(t, EventAuthCompleted, track.Event)
	assert.Equal(t, "signup", track.Properties["flow"])
	assert.Equal(t, int64(1200), track.Properties["duration_ms"])
	assert.Contains(t, track.Properties, "invocation_id")
	assert.Contains(t, track.Properties, "command")
}
