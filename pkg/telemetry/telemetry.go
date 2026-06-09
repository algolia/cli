package telemetry

import (
	"context"
	"crypto/md5" // nolint:gosec
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/segmentio/analytics-go/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/xtgo/uuid"

	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/version"
)

const (
	AppName               = "cli"
	telemetryAnalyticsURL = "https://telemetry-proxy.algolia.com/"

	// telemetryHTTPTimeout bounds the total duration of the telemetry flush
	// HTTP request (connect, TLS, and response), so closing the client at the
	// end of a command can never hang the CLI on a slow or unreachable endpoint.
	telemetryHTTPTimeout = 3 * time.Second
)

type telemetryMetadataKey struct{}

type telemetryClientKey struct{}

type TelemetryClient interface {
	Identify(ctx context.Context) error
	Track(ctx context.Context, event string, properties map[string]any) error
	Close()
}

type AnalyticsTelemetryClient struct {
	client analytics.Client
	debug  bool

	mu         sync.Mutex
	lastTS     time.Time
	eventIndex int64
}

type AnalyticsTelemetryLogger struct {
	debug  bool
	logger *log.Logger
}

func (l AnalyticsTelemetryLogger) Logf(format string, args ...interface{}) {
	if l.debug {
		fmt.Printf("INFO: "+format, args...)
	}
}

func (l AnalyticsTelemetryLogger) Errorf(format string, args ...interface{}) {
	// The telemetry should always fail silently, unless in debug mode
	if l.debug {
		fmt.Printf("ERROR: "+format, args...)
	}
}

func newTelemetryLogger(debug bool) AnalyticsTelemetryLogger {
	return AnalyticsTelemetryLogger{debug, log.New(nil, "telemetry ", log.LstdFlags)}
}

func NewAnalyticsTelemetryClient(debug bool) (TelemetryClient, error) {
	return newAnalyticsTelemetryClient(telemetryAnalyticsURL, debug)
}

func newAnalyticsTelemetryClient(endpoint string, debug bool) (TelemetryClient, error) {
	client, err := analytics.NewWithConfig("", analytics.Config{
		Endpoint: endpoint,
		Logger:   newTelemetryLogger(debug),
		// In debug mode, surface the library's own batch/flush logs.
		Verbose: debug,
		// Buffer every event into one batch flushed at Close. The default 5s
		// interval would split a long command (e.g. interactive login) across
		// requests, reordering events downstream and risking a dropped batch.
		Interval:  24 * time.Hour,
		BatchSize: 250,
		// Bound the flush request so Close() at exit can't hang the CLI.
		Transport: boundedRoundTripper{
			base:    http.DefaultTransport,
			timeout: telemetryHTTPTimeout,
		},
	})
	if err != nil {
		return nil, err
	}
	return &AnalyticsTelemetryClient{client: client, debug: debug}, nil
}

// boundedRoundTripper applies a total per-request timeout, like
// http.Client.Timeout (analytics.Config only accepts a RoundTripper).
type boundedRoundTripper struct {
	base    http.RoundTripper
	timeout time.Duration
}

func (t boundedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(req.Context(), t.timeout)
	resp, err := t.base.RoundTrip(req.WithContext(ctx))
	if err != nil {
		cancel()
		return nil, err
	}
	// Cancel when the body is closed, not now, or the response would be truncated.
	resp.Body = &cancelOnCloseBody{ReadCloser: resp.Body, cancel: cancel}
	return resp, nil
}

// cancelOnCloseBody cancels the request context when the body is closed.
type cancelOnCloseBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (b *cancelOnCloseBody) Close() error {
	err := b.ReadCloser.Close()
	b.cancel()
	return err
}

// IdentifyOnce sends a single Identify event through a short-lived client and
// flushes it before returning. It is meant for one-shot identification (for
// example, right after authentication fills the token) where the command's
// request-scoped client may already have been closed. It honors the same
// ALGOLIA_CLI_TELEMETRY and DEBUG environment variables as the root command and
// fails silently so telemetry never blocks the user.
func IdentifyOnce(ctx context.Context) {
	if os.Getenv("ALGOLIA_CLI_TELEMETRY") == "0" {
		return
	}

	client, err := NewAnalyticsTelemetryClient(os.Getenv("DEBUG") != "")
	if err != nil {
		return
	}
	defer client.Close()

	_ = client.Identify(ctx)
}

// anonymousID is a unique identifier for an anonymous user of the CLI (basically the hash of the mac address)
func anonymousID() string {
	addrs, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		a.Flags &= net.FlagUp | net.FlagLoopback
		if a.Flags == 0 {
			continue // interface down
		}
		a := a.HardwareAddr.String()
		if a != "" {
			return fmt.Sprintf("%x", md5.Sum([]byte(a))) // nolint: gosec
		}
	}
	return ""
}

type NoOpTelemetryClient struct{}

type CLIAnalyticsEventMetadata struct {
	AnonymousID              string   // the anonymous id is the hash of the mac address of the machine
	UserID                   string   // the authenticated user's id from the OAuth token; empty when logged out
	Email                    string   // the authenticated user's email, when available
	Name                     string   // the authenticated user's name, when available
	InvocationID             string   // the invocation id is unique to each context object and represents all events coming from one command
	ConfiguredApplicationsNb int      // the number of configured applications
	AppID                    string   // the app id with which the command was called
	CommandPath              string   // the command path is the full path of the command
	CommandFlags             []string // the command flags is the full list of flags passed to the command
	CLIVersion               string   // the version of the CLI
	OS                       string   // the OS of the system
}

// NewEventMetadata initializes an instance of CLIAnalyticsEventContext
func NewEventMetadata() *CLIAnalyticsEventMetadata {
	return &CLIAnalyticsEventMetadata{
		AnonymousID:  anonymousID(),
		InvocationID: uuid.NewRandom().String(),
		CLIVersion:   version.Version,
		OS:           runtime.GOOS,
	}
}

// WithEventMetadata returns a new copy of context.Context with the provided CLIAnalyticsEventMetadata
func WithEventMetadata(ctx context.Context, metadata *CLIAnalyticsEventMetadata) context.Context {
	return context.WithValue(ctx, telemetryMetadataKey{}, metadata)
}

// GetEventMetadata returns the CLIAnalyticsEventMetadata from the provided context
func GetEventMetadata(ctx context.Context) *CLIAnalyticsEventMetadata {
	metadata := ctx.Value(telemetryMetadataKey{})
	if metadata != nil {
		return metadata.(*CLIAnalyticsEventMetadata)
	}
	return nil
}

// WithTelemetryClient returns a new copy of context.Context with the provided telemetryClient
func WithTelemetryClient(ctx context.Context, client TelemetryClient) context.Context {
	return context.WithValue(ctx, telemetryClientKey{}, client)
}

// GetTelemetryClient returns the CLIAnalyticsEventMetadata from the provided context
func GetTelemetryClient(ctx context.Context) TelemetryClient {
	client := ctx.Value(telemetryClientKey{})
	if client != nil {
		return client.(TelemetryClient)
	}
	return nil
}

// SetCobraCommandContext sets the telemetry values for the command being executed.
func (e *CLIAnalyticsEventMetadata) SetCobraCommandContext(cmd *cobra.Command) {
	e.CommandPath = cmd.CommandPath()
	var flags []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flags = append(flags, f.Name)
		}
	})
	e.CommandFlags = flags
}

// SetAppID sets the AppID on the CLIAnalyticsEventContext object
func (e *CLIAnalyticsEventMetadata) SetAppID(appID string) {
	e.AppID = appID
}

// SetCommandPath sets the commandPath on the CLIAnalyticsEventContext object
func (e *CLIAnalyticsEventMetadata) SetCommandPath(commandPath string) {
	e.CommandPath = commandPath
}

// SetConfiguredApplicationsNb sets the configuredApplicationsNb on the CLIAnalyticsEventContext object
func (e *CLIAnalyticsEventMetadata) SetConfiguredApplicationsNb(nb int) {
	e.ConfiguredApplicationsNb = nb
}

// SetUser sets the authenticated user identity on the CLIAnalyticsEventContext object
func (e *CLIAnalyticsEventMetadata) SetUser(userID, email, name string) {
	e.UserID = userID
	e.Email = email
	e.Name = name
}

// Identify tracks the user with the provided properties
func (a *AnalyticsTelemetryClient) Identify(ctx context.Context) error {
	metadata := GetEventMetadata(ctx)

	var isCI int8
	if utils.IsCI() {
		isCI = 1
	}

	traits := analytics.Traits{
		"configured_applications": metadata.ConfiguredApplicationsNb,
		"version":                 metadata.CLIVersion,
		"operating_system":        metadata.OS,
		"is_ci":                   isCI,
	}

	identify := analytics.Identify{
		AnonymousId: metadata.AnonymousID,
		Traits:      traits,
		Context: &analytics.Context{
			Device: analytics.DeviceInfo{
				Id: metadata.AnonymousID,
			},
		},
	}

	if metadata.UserID != "" {
		identify.UserId = metadata.UserID
		if metadata.Email != "" {
			traits["email"] = metadata.Email
		}
		if metadata.Name != "" {
			traits["name"] = metadata.Name
		}
	}

	return a.client.Enqueue(identify)
}

// Track merges custom properties over the base properties (custom wins on collisions).
func (a *AnalyticsTelemetryClient) Track(
	ctx context.Context,
	event string,
	properties map[string]any,
) error {
	metadata := GetEventMetadata(ctx)

	props := map[string]interface{}{
		"invocation_id": metadata.InvocationID,
		"app_id":        metadata.AppID,
		"command":       metadata.CommandPath,
		"flags":         metadata.CommandFlags,
		// version lets downstream destinations split released binaries
		// (semver) from source builds ("main").
		"version": metadata.CLIVersion,
		// event_index is the emit order within one invocation: the
		// vendor-independent ground truth for sequencing, since Amplitude's
		// tie-breaking on identical timestamps is undocumented.
		"event_index": a.nextEventIndex(),
	}
	for k, v := range properties {
		props[k] = v
	}

	// In debug mode, echo each event to stderr for local observability.
	if a.debug {
		fmt.Fprintf(os.Stderr, "[telemetry] %s %v\n", event, properties)
	}

	track := analytics.Track{
		Event:       event,
		AnonymousId: metadata.AnonymousID,
		Properties:  props,
		Timestamp:   a.nextTimestamp(),
		Context: &analytics.Context{
			Device: analytics.DeviceInfo{
				Id: metadata.AnonymousID,
			},
		},
	}

	if metadata.UserID != "" {
		track.UserId = metadata.UserID
	}

	return a.client.Enqueue(track)
}

// nextEventIndex returns 1, 2, 3... in emit order for this command run.
func (a *AnalyticsTelemetryClient) nextEventIndex() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.eventIndex++
	return a.eventIndex
}

// nextTimestamp returns strictly increasing timestamps at least 1ms apart.
// Amplitude truncates event_time to milliseconds, so same-millisecond events
// are ordered non-deterministically; spacing by 1ms preserves emit order.
func (a *AnalyticsTelemetryClient) nextTimestamp() time.Time {
	a.mu.Lock()
	defer a.mu.Unlock()

	ts := time.Now()
	if !ts.After(a.lastTS) {
		ts = a.lastTS.Add(time.Millisecond)
	}
	a.lastTS = ts
	return ts
}

// Close closes the client, waiting for all pending events to be sent.
func (a *AnalyticsTelemetryClient) Close() {
	_ = a.client.Close()
}

func (a *NoOpTelemetryClient) Identify(ctx context.Context) error { return nil }

func (a *NoOpTelemetryClient) Track(
	ctx context.Context,
	event string,
	properties map[string]any,
) error {
	return nil
}
func (a *NoOpTelemetryClient) Close() {}
