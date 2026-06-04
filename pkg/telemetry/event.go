package telemetry

import "context"

type Event struct {
	Name       string
	Properties map[string]any
}

func Track(ctx context.Context, event Event) {
	client := GetTelemetryClient(ctx)
	if client == nil {
		return
	}
	_ = client.Track(ctx, event.Name, event.Properties)
}
