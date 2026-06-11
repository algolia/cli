package telemetry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRootError struct{}

func (testRootError) Error() string { return "boom" }

type testWrapperError struct{ inner error }

func (e testWrapperError) Error() string { return "wrap" }
func (e testWrapperError) Unwrap() error { return e.inner }

func TestErrorClass_SkipsFmtWrappers(t *testing.T) {
	wrapped := fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", testRootError{}))
	assert.Equal(t, "telemetry.testRootError", ErrorClass(wrapped))
}

func TestErrorClass_KeepsInformativeWrapperType(t *testing.T) {
	wrapped := fmt.Errorf("outer: %w", testWrapperError{inner: testRootError{}})
	assert.Equal(t, "telemetry.testWrapperError", ErrorClass(wrapped))
}

func TestErrorClass_NilError(t *testing.T) {
	assert.Equal(t, "", ErrorClass(nil))
}

func TestTrackEvent_NoClientInContextIsSafe(t *testing.T) {
	// Must not panic when the context carries no telemetry client.
	TrackEvent(context.Background(), AuthStarted(FlowLogin, false))
}

func TestAuthStarted(t *testing.T) {
	event := AuthStarted(FlowSignup, true)
	assert.Equal(t, EventAuthStarted, event.Name)
	assert.Equal(t, FlowSignup, event.Properties["flow"])
	assert.Equal(t, true, event.Properties["no_browser"])
}

func TestAuthCompleted(t *testing.T) {
	event := AuthCompleted(FlowLogin, NewFlowTracker())
	assert.Equal(t, EventAuthCompleted, event.Name)
	assert.Equal(t, FlowLogin, event.Properties["flow"])
	assert.Contains(t, event.Properties, "duration_ms")
}

func TestAuthAborted(t *testing.T) {
	tracker := NewFlowTracker()
	tracker.SetStep(StepAppSelect)

	event := AuthAborted(FlowLogin, tracker)
	assert.Equal(t, EventAuthAborted, event.Name)
	assert.Equal(t, StepAppSelect, event.Properties["step"])
}

func TestAuthFailed(t *testing.T) {
	tracker := NewFlowTracker()
	tracker.SetStep(StepAppsFetch)

	event := AuthFailed(FlowLogin, tracker, errors.New("boom"))
	assert.Equal(t, EventAuthFailed, event.Name)
	assert.Equal(t, StepAppsFetch, event.Properties["step"])
	assert.Equal(t, "*errors.errorString", event.Properties["error_class"])
	assert.Contains(t, event.Properties, "duration_ms")
}

func TestFlowTracker_NilTrackerIsSafe(t *testing.T) {
	var tracker *FlowTracker
	tracker.SetStep(StepTerms)
	assert.Equal(t, Step(""), tracker.Step())
	assert.Equal(t, int64(0), tracker.DurationMS())
}

func TestFlowTracker_TracksStepAndDuration(t *testing.T) {
	tracker := NewFlowTracker()
	tracker.SetStep(StepPlan)
	assert.Equal(t, StepPlan, tracker.Step())
	assert.GreaterOrEqual(t, tracker.DurationMS(), int64(0))
}
