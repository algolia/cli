package telemetry

import (
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
