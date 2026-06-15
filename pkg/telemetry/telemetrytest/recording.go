// Package telemetrytest provides test doubles for the telemetry package.
package telemetrytest

import (
	"context"
	"sync"

	"github.com/algolia/cli/pkg/telemetry"
)

var _ telemetry.TelemetryClient = (*RecordingClient)(nil)

// RecordedEvent is one event captured by a RecordingClient.
type RecordedEvent struct {
	Name       string
	Properties map[string]any
}

// RecordingClient is a telemetry.TelemetryClient that records the tracked
// events so tests can assert on their names, properties and order.
type RecordingClient struct {
	mu         sync.Mutex
	Events     []RecordedEvent
	Identifies int
}

func (r *RecordingClient) Identify(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Identifies++
	return nil
}

func (r *RecordingClient) Track(
	ctx context.Context,
	event string,
	properties map[string]any,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Events = append(r.Events, RecordedEvent{event, properties})
	return nil
}

func (r *RecordingClient) Close() {}
