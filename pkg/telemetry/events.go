package telemetry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Event names. New flow events follow the `CLI <Command> <Step>` convention;
// the command lifecycle events stay unprefixed for consistency with the
// historical "Command Invoked".
const (
	EventCommandInvoked   = "Command Invoked"
	EventCommandCompleted = "Command Completed"

	EventAuthStarted   = "CLI Auth Started"
	EventAuthCompleted = "CLI Auth Completed"
	EventAuthFailed    = "CLI Auth Failed"
	EventAuthAborted   = "CLI Auth Aborted"
	EventAuthLogout    = "CLI Auth Logout"

	EventApplicationCreateStarted       = "CLI Application Create Started"
	EventApplicationCreateAcceptedTerms = "CLI Application Create Accepted Terms"
	EventApplicationCreateDeclinedTerms = "CLI Application Create Declined Terms"
	EventApplicationCreateCompleted     = "CLI Application Create Completed"
	EventApplicationCreateFailed        = "CLI Application Create Failed"
	EventApplicationCreateAborted       = "CLI Application Create Aborted"

	EventApplicationPlanChangeStarted       = "CLI Application Plan Change Started"
	EventApplicationPlanChangeAcceptedTerms = "CLI Application Plan Change Accepted Terms"
	EventApplicationPlanChangeDeclinedTerms = "CLI Application Plan Change Declined Terms"
	EventApplicationPlanChangeCompleted     = "CLI Application Plan Change Completed"
	EventApplicationPlanChangeFailed        = "CLI Application Plan Change Failed"
	EventApplicationPlanChangeAborted       = "CLI Application Plan Change Aborted"
)

// Flow is the kind of auth flow the user is going through.
type Flow string

const (
	FlowLogin  Flow = "login"
	FlowSignup Flow = "signup"
	FlowLogout Flow = "logout"
)

// Step locates where the user is inside an interactive flow, so aborts and
// failures can tell where the user stopped.
type Step string

const (
	// Auth flow steps.
	StepBrowserWait      Step = "browser_wait"
	StepCodeExchange     Step = "code_exchange"
	StepAppsFetch        Step = "apps_fetch"
	StepAppSelect        Step = "app_select"
	StepAppCreate        Step = "app_create"
	StepProfileConfigure Step = "profile_configure"

	// Application create and plan change flow steps.
	StepName      Step = "name"
	StepPlan      Step = "plan"
	StepTerms     Step = "terms"
	StepRegion    Step = "region"
	StepAPICall   Step = "api_call"
	StepApplyPlan Step = "apply_plan"
)

// Direction is the direction of a plan change.
type Direction string

const (
	DirectionUpgrade   Direction = "upgrade"
	DirectionDowngrade Direction = "downgrade"
)

// FlowTracker carries the state of one interactive flow: the step the user is
// currently in and the flow start time, to compute durations. All its methods
// are safe on a nil tracker, so helpers shared by several flows can take an
// optional tracker.
type FlowTracker struct {
	start time.Time
	step  Step
}

func NewFlowTracker() *FlowTracker {
	return &FlowTracker{start: time.Now()}
}

// SetStep records the step the flow is entering.
func (f *FlowTracker) SetStep(step Step) {
	if f == nil {
		return
	}
	f.step = step
}

// Step returns the step the flow is currently in.
func (f *FlowTracker) Step() Step {
	if f == nil {
		return ""
	}
	return f.step
}

// DurationMS returns the time elapsed since the flow started, in milliseconds.
func (f *FlowTracker) DurationMS() int64 {
	if f == nil {
		return 0
	}
	return time.Since(f.start).Milliseconds()
}

// Event is a fully assembled telemetry event, ready to be tracked.
type Event struct {
	Name       string
	Properties map[string]any
}

// TrackEvent sends the event through the context's telemetry client, silently
// doing nothing when no client is present.
func TrackEvent(ctx context.Context, event Event) {
	client := GetTelemetryClient(ctx)
	if client == nil {
		return
	}
	_ = client.Track(ctx, event.Name, event.Properties)
}

// AuthStarted is emitted when the browser-based OAuth flow begins.
func AuthStarted(flow Flow, noBrowser bool) Event {
	return Event{EventAuthStarted, map[string]any{
		"flow":       flow,
		"no_browser": noBrowser,
	}}
}

// AuthCompleted is emitted when the profile is fully configured at the end of
// the auth flow.
func AuthCompleted(flow Flow, tracker *FlowTracker) Event {
	return Event{EventAuthCompleted, map[string]any{
		"flow":        flow,
		"duration_ms": tracker.DurationMS(),
	}}
}

// AuthAborted is emitted when the user cancelled the auth flow, with the step
// they stopped at.
func AuthAborted(flow Flow, tracker *FlowTracker) Event {
	return Event{EventAuthAborted, map[string]any{
		"flow": flow,
		"step": tracker.Step(),
	}}
}

// AuthLogout is emitted when the user signs out. It must be tracked before
// the local state is cleared, while the user identifier is still attached to
// the telemetry metadata; no Identify follows, since Segment's identity graph
// has no concept of un-identifying.
func AuthLogout() Event {
	return Event{EventAuthLogout, map[string]any{
		"flow": FlowLogout,
	}}
}

// AuthFailed is emitted when the auth flow failed, with the step it failed at.
func AuthFailed(flow Flow, tracker *FlowTracker, err error) Event {
	return Event{EventAuthFailed, map[string]any{
		"flow":        flow,
		"step":        tracker.Step(),
		"duration_ms": tracker.DurationMS(),
		"error_class": ErrorClass(err),
	}}
}

// ErrorClass returns the type of the first informative error of the chain,
// skipping the anonymous wrappers created by fmt.Errorf. It never returns an
// error message, which could contain user data.
func ErrorClass(err error) string {
	for err != nil {
		class := fmt.Sprintf("%T", err)
		switch class {
		case "*fmt.wrapError", "*fmt.wrapErrors", "*errors.joinError":
			if unwrapped := errors.Unwrap(err); unwrapped != nil {
				err = unwrapped
				continue
			}
		}
		return class
	}
	return ""
}
