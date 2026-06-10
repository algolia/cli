package telemetry

import (
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

// ErrorClass returns the type of the root cause of an error, never its
// message, which could contain user data.
func ErrorClass(err error) string {
	if err == nil {
		return ""
	}
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			break
		}
		err = unwrapped
	}
	return fmt.Sprintf("%T", err)
}
