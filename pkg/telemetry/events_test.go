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

func TestApplicationCreateCompleted(t *testing.T) {
	event := ApplicationCreateCompleted("us-east", "grow", NewFlowTracker())
	assert.Equal(t, EventApplicationCreateCompleted, event.Name)
	assert.Equal(t, "us-east", event.Properties["region"])
	assert.Equal(t, "grow", event.Properties["plan"])
	assert.Contains(t, event.Properties, "duration_ms")
}

func TestApplicationPlanChangeAborted_WithReason(t *testing.T) {
	tracker := NewFlowTracker()
	tracker.SetStep(StepPlan)

	event := ApplicationPlanChangeAborted(DirectionUpgrade, tracker, AbortReasonAlreadyOnPlan)
	assert.Equal(t, EventApplicationPlanChangeAborted, event.Name)
	assert.Equal(t, DirectionUpgrade, event.Properties["direction"])
	assert.Equal(t, StepPlan, event.Properties["step"])
	assert.Equal(t, AbortReasonAlreadyOnPlan, event.Properties["reason"])
}

func TestApplicationPlanChangeAborted_WithoutReason(t *testing.T) {
	event := ApplicationPlanChangeAborted(DirectionDowngrade, NewFlowTracker(), "")
	assert.NotContains(t, event.Properties, "reason")
}

func TestApplicationCreateAborted_WithReason(t *testing.T) {
	tracker := NewFlowTracker()
	tracker.SetStep(StepTerms)

	event := ApplicationCreateAborted(tracker, AbortReasonDeclinedTerms)
	assert.Equal(t, EventApplicationCreateAborted, event.Name)
	assert.Equal(t, StepTerms, event.Properties["step"])
	assert.Equal(t, AbortReasonDeclinedTerms, event.Properties["reason"])
}

func TestApplicationPlanChangeCompleted(t *testing.T) {
	event := ApplicationPlanChangeCompleted(DirectionUpgrade, "free", "grow", NewFlowTracker())
	assert.Equal(t, EventApplicationPlanChangeCompleted, event.Name)
	assert.Equal(t, "free", event.Properties["from_plan"])
	assert.Equal(t, "grow", event.Properties["to_plan"])
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
