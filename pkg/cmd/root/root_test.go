package root

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/telemetry"
)

// recordingTelemetryClient captures the tracked events so tests can assert on
// them without hitting the network.
type recordingTelemetryClient struct {
	events []recordedEvent
}

type recordedEvent struct {
	name  string
	props map[string]any
}

func (r *recordingTelemetryClient) Identify(ctx context.Context) error { return nil }

func (r *recordingTelemetryClient) Track(
	ctx context.Context,
	event string,
	properties map[string]any,
) error {
	r.events = append(r.events, recordedEvent{event, properties})
	return nil
}

func (r *recordingTelemetryClient) Close() {}

func newTelemetryContext(client telemetry.TelemetryClient, commandPath string) context.Context {
	metadata := telemetry.NewEventMetadata()
	metadata.SetCommandPath(commandPath)
	ctx := telemetry.WithEventMetadata(context.Background(), metadata)
	return telemetry.WithTelemetryClient(ctx, client)
}

func TestTrackCommandCompleted_SkipsWhenPreRunNeverRan(t *testing.T) {
	client := &recordingTelemetryClient{}
	// An empty command path means PersistentPreRunE never ran.
	ctx := newTelemetryContext(client, "")

	trackCommandCompleted(ctx, &cobra.Command{Use: "algolia"}, exitOK, nil, time.Second)

	if len(client.events) != 0 {
		t.Errorf("expected no event, got %d", len(client.events))
	}
}

func TestTrackCommandCompleted_ReportsSuccess(t *testing.T) {
	client := &recordingTelemetryClient{}
	ctx := newTelemetryContext(client, "algolia indices list")

	trackCommandCompleted(ctx, &cobra.Command{Use: "list"}, exitOK, nil, 1500*time.Millisecond)

	if len(client.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(client.events))
	}
	event := client.events[0]
	if event.name != telemetry.EventCommandCompleted {
		t.Errorf("event = %q, want %q", event.name, telemetry.EventCommandCompleted)
	}
	if event.props["succeeded"] != true {
		t.Errorf("succeeded = %v, want true", event.props["succeeded"])
	}
	if event.props["exit_code"] != 0 {
		t.Errorf("exit_code = %v, want 0", event.props["exit_code"])
	}
	if event.props["duration_ms"] != int64(1500) {
		t.Errorf("duration_ms = %v, want 1500", event.props["duration_ms"])
	}
	if _, ok := event.props["error_class"]; ok {
		t.Error("unexpected error_class on success")
	}
}

func TestTrackCommandCompleted_ReportsFailure(t *testing.T) {
	client := &recordingTelemetryClient{}
	ctx := newTelemetryContext(client, "algolia indices list")

	trackCommandCompleted(ctx, &cobra.Command{Use: "list"}, exitError, errors.New("boom"), time.Second)

	if len(client.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(client.events))
	}
	props := client.events[0].props
	if props["succeeded"] != false {
		t.Errorf("succeeded = %v, want false", props["succeeded"])
	}
	if props["exit_code"] != 1 {
		t.Errorf("exit_code = %v, want 1", props["exit_code"])
	}
	if props["error_class"] != "*errors.errorString" {
		t.Errorf("error_class = %v, want *errors.errorString", props["error_class"])
	}
	if props["user_cancelled"] != false {
		t.Errorf("user_cancelled = %v, want false", props["user_cancelled"])
	}
}

func TestPrintError(t *testing.T) {
	cmd := &cobra.Command{}

	type args struct {
		err   error
		cmd   *cobra.Command
		debug bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "generic error",
			args: args{
				err:   errors.New("the app exploded"),
				cmd:   nil,
				debug: false,
			},
			wantOut: "the app exploded\n",
		},
		{
			name: "DNS error",
			args: args{
				err: fmt.Errorf("DNS oopsie: %w", &net.DNSError{
					Name: "latency.algolia.net",
				}),
				cmd:   nil,
				debug: false,
			},
			wantOut: `error connecting to latency.algolia.net
check your internet connection or https://status.algolia.com
`,
		},
		{
			name: "Cobra flag error",
			args: args{
				err:   &cmdutil.FlagError{Err: errors.New("unknown flag --foo")},
				cmd:   cmd,
				debug: false,
			},
			wantOut: "unknown flag --foo\n\nUsage:\n\n",
		},
		{
			name: "unknown Cobra command error",
			args: args{
				err:   errors.New("unknown command foo"),
				cmd:   cmd,
				debug: false,
			},
			wantOut: "unknown command foo\n\nUsage:\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			printError(out, tt.args.err, tt.args.cmd, tt.args.debug)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("printError() = %q, want %q", gotOut, tt.wantOut)
			}
		})
	}
}
