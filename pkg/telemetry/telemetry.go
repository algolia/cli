package telemetry

import (
	"context"
	"net"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/xtgo/uuid"
	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/algolia/cli/pkg/version"
)

const AppName = "cli"
const telemetryAnalyticsURL = "https://telemetry-proxy.algolia.com/"

type telemetryMetadataKey struct{}

type telemetryClientKey struct{}

type TelemetryClient interface {
	Identify(ctx context.Context) error
	Track(ctx context.Context, event string) error
	Close() error
}

type AnalyticsTelemetryClient struct {
	client analytics.Client
}

func NewAnalyticsTelemetryClient() (TelemetryClient, error) {
	client, err := analytics.NewWithConfig("", analytics.Config{
		Endpoint: telemetryAnalyticsURL,
	})
	if err != nil {
		return nil, err
	}
	return &AnalyticsTelemetryClient{client: client}, nil
}

// userID is a unique identifier for the user of the CLI (basically the mac address of the machine)
func userID() string {
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
			return a
		}
	}
	return ""
}

type NoOpTelemetryClient struct{}

type CLIAnalyticsEventMetadata struct {
	UserId                   string // the user id is the mac address of the machine
	InvocationID             string // the invocation id is unique to each context object and represents all events coming from one command
	ConfiguredApplicationsNb int    // the number of configured applications
	AppID                    string // the app id with which the command was called
	CommandPath              string // the command path is the full path of the command
	CLIVersion               string // the version of the CLI
	OS                       string // the OS of the system
}

// NewEventMetadata initializes an instance of CLIAnalyticsEventContext
func NewEventMetadata() *CLIAnalyticsEventMetadata {
	return &CLIAnalyticsEventMetadata{
		UserId:       userID(),
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

// Identify tracks the user with the provided properties
func (a *AnalyticsTelemetryClient) Identify(ctx context.Context) error {
	metadata := GetEventMetadata(ctx)

	return a.client.Enqueue(analytics.Identify{
		AnonymousId: metadata.UserId,
		UserId:      metadata.UserId,
		Traits: map[string]interface{}{
			"configured_applications": metadata.ConfiguredApplicationsNb,
			"version":                 metadata.CLIVersion,
			"operating_system":        metadata.OS,
		},
	})
}

// Track tracks the event with the provided properties
func (a *AnalyticsTelemetryClient) Track(ctx context.Context, event string) error {
	metadata := GetEventMetadata(ctx)

	return a.client.Enqueue(analytics.Track{
		Event:       event,
		AnonymousId: metadata.UserId,
		UserId:      metadata.UserId,
		Properties: map[string]interface{}{
			"invocation_id": metadata.InvocationID,
			"app_id":        metadata.AppID,
			"command":       metadata.CommandPath,
		},
	})
}

// Close closes the client, waiting for all pending events to be sent.
func (a *AnalyticsTelemetryClient) Close() error {
	return a.client.Close()
}

func (a *NoOpTelemetryClient) Identify(ctx context.Context) error            { return nil }
func (a *NoOpTelemetryClient) Track(ctx context.Context, event string) error { return nil }
func (a *NoOpTelemetryClient) Close() error                                  { return nil }
