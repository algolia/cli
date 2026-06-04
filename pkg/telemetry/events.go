package telemetry

import "time"

// Catalog of the analytics events the CLI emits, mirroring the analytics spec.
// Keeping names and property schemas here (not inline at call sites) keeps them
// auditable in one place and prevents names/keys from drifting.
//
// To add an event: add a name constant, add a constructor, then call
// telemetry.Track(ctx, telemetry.YourEvent(...)).

// Event names. Flow events follow the "CLI <Command> <Step>" convention.
const (
	EventCommandInvoked   = "Command Invoked"
	EventCommandCompleted = "Command Completed"
	EventCommandFailed    = "Command Failed"

	EventAuthStarted          = "CLI Auth Started"
	EventAuthBrowserOpened    = "CLI Auth Browser Opened"
	EventAuthBrowserFailed    = "CLI Auth Browser Failed"
	EventAuthCallbackReceived = "CLI Auth Callback Received"
	EventAuthCompleted        = "CLI Auth Completed"
	EventAuthFailed           = "CLI Auth Failed"

	EventApplicationCreateStarted   = "CLI Application Create Started"
	EventApplicationCreateCompleted = "CLI Application Create Completed"
	EventApplicationCreateFailed    = "CLI Application Create Failed"
	EventApplicationCreateAborted   = "CLI Application Create Aborted"

	EventApplicationUpgradeStarted       = "CLI Application Upgrade Started"
	EventApplicationUpgradeAcceptedTerms = "CLI Application Upgrade Accepted Terms"
	EventApplicationUpgradeDeclinedTerms = "CLI Application Upgrade Declined Terms"
	EventApplicationUpgradeFailed        = "CLI Application Upgrade Failed"
	EventApplicationUpgradeCompleted     = "CLI Application Upgrade Completed"

	EventApplicationDowngradeStarted       = "CLI Application Downgrade Started"
	EventApplicationDowngradeAcceptedTerms = "CLI Application Downgrade Accepted Terms"
	EventApplicationDowngradeDeclinedTerms = "CLI Application Downgrade Declined Terms"
	EventApplicationDowngradeFailed        = "CLI Application Downgrade Failed"
	EventApplicationDowngradeCompleted     = "CLI Application Downgrade Completed"
)

// Property values, so call sites reference a constant instead of a literal.
const (
	FlowLogin  = "login"
	FlowSignup = "signup"

	// triggered_from: auth flow vs explicit "application create" command.
	TriggeredFromAuthFlow        = "auth_flow"
	TriggeredFromExplicitCommand = "explicit_command"

	// step: where the auth flow failed (CLI Auth Failed).
	AuthStepBrowser   = "browser"
	AuthStepCallback  = "callback"
	AuthStepExchange  = "exchange"
	AuthStepAppsFetch = "apps_fetch"
)

// --- Root command events ---------------------------------------------------

// CommandInvoked is emitted when a tracked command starts.
func CommandInvoked() Event {
	return Event{Name: EventCommandInvoked}
}

// CommandCompleted is emitted (deferred) when the root command returns.
func CommandCompleted(duration time.Duration, succeeded bool, exitCode int) Event {
	return Event{
		Name: EventCommandCompleted,
		Properties: map[string]any{
			"duration_ms": duration.Milliseconds(),
			"succeeded":   succeeded,
			"exit_code":   exitCode,
		},
	}
}

// CommandFailed is emitted when the root command returns an error. httpStatus
// is omitted when zero (i.e. the error carried no HTTP status).
func CommandFailed(errorClass, errorSource string, httpStatus int) Event {
	props := map[string]any{
		"error_class":  errorClass,
		"error_source": errorSource,
	}
	if httpStatus != 0 {
		props["http_status"] = httpStatus
	}
	return Event{Name: EventCommandFailed, Properties: props}
}

// --- Auth flow events ------------------------------------------------------

// AuthStarted is emitted when the OAuth flow begins.
func AuthStarted(flow string, noBrowser bool) Event {
	return Event{
		Name: EventAuthStarted,
		Properties: map[string]any{
			"flow":       flow,
			"no_browser": noBrowser,
		},
	}
}

// AuthBrowserOpened is emitted when the default browser is launched successfully.
func AuthBrowserOpened(flow string) Event {
	return Event{
		Name:       EventAuthBrowserOpened,
		Properties: map[string]any{"flow": flow},
	}
}

// AuthBrowserFailed is emitted when launching the browser fails and the
// authorize URL is printed instead.
func AuthBrowserFailed(flow, errorClass string) Event {
	return Event{
		Name: EventAuthBrowserFailed,
		Properties: map[string]any{
			"flow":        flow,
			"error_class": errorClass,
		},
	}
}

// AuthCallbackReceived is emitted when the OAuth redirect hits the local
// callback server.
func AuthCallbackReceived(flow string, duration time.Duration) Event {
	return Event{
		Name: EventAuthCallbackReceived,
		Properties: map[string]any{
			"flow":        flow,
			"duration_ms": duration.Milliseconds(),
		},
	}
}

// AuthCompleted is emitted once the profile is fully configured at the end of
// the flow.
func AuthCompleted(
	flow string,
	duration time.Duration,
	hadExistingApps, createdAppDuringFlow bool,
) Event {
	return Event{
		Name: EventAuthCompleted,
		Properties: map[string]any{
			"flow":                    flow,
			"duration_ms":             duration.Milliseconds(),
			"had_existing_apps":       hadExistingApps,
			"created_app_during_flow": createdAppDuringFlow,
		},
	}
}

// AuthFailed is emitted on a timeout, token exchange error, or app fetch error.
// step is one of the AuthStep* constants.
func AuthFailed(flow, step, errorClass string) Event {
	return Event{
		Name: EventAuthFailed,
		Properties: map[string]any{
			"flow":        flow,
			"step":        step,
			"error_class": errorClass,
		},
	}
}

// --- Application create events ---------------------------------------------

// ApplicationCreateStarted is emitted once validation passes, before the
// Dashboard API call. triggeredFrom is one of the TriggeredFrom* constants.
func ApplicationCreateStarted(triggeredFrom string) Event {
	return Event{
		Name:       EventApplicationCreateStarted,
		Properties: map[string]any{"triggered_from": triggeredFrom},
	}
}

// ApplicationCreateCompleted is emitted when the Dashboard API returns 2xx.
func ApplicationCreateCompleted(triggeredFrom string, duration time.Duration) Event {
	return Event{
		Name: EventApplicationCreateCompleted,
		Properties: map[string]any{
			"triggered_from": triggeredFrom,
			"duration_ms":    duration.Milliseconds(),
		},
	}
}

// ApplicationCreateFailed is emitted when the Dashboard API returns an error.
// httpStatus is omitted when zero.
func ApplicationCreateFailed(triggeredFrom, errorClass string, httpStatus int) Event {
	props := map[string]any{
		"triggered_from": triggeredFrom,
		"error_class":    errorClass,
	}
	if httpStatus != 0 {
		props["http_status"] = httpStatus
	}
	return Event{Name: EventApplicationCreateFailed, Properties: props}
}

// ApplicationCreateAborted is emitted when the user declines the confirmation
// prompt.
func ApplicationCreateAborted(triggeredFrom string) Event {
	return Event{
		Name:       EventApplicationCreateAborted,
		Properties: map[string]any{"triggered_from": triggeredFrom},
	}
}

// --- Application upgrade events ---------------------------------------------

// ApplicationUpgradeStarted is emitted when the user starts the upgrade flow.
func ApplicationUpgradeStarted(plan string) Event {
	return Event{
		Name: EventApplicationUpgradeStarted,
	}
}

// ApplicationUpgradeAcceptedTerms is emitted when the user accepts the T&C.
func ApplicationUpgradeAcceptedTerms(plan string) Event {
	return Event{
		Name:       EventApplicationUpgradeAcceptedTerms,
		Properties: map[string]any{"plan": plan},
	}
}

// ApplicationUpgradeDeclinedTerms is emitted when the user declines the T&C.
func ApplicationUpgradeDeclinedTerms(plan string) Event {
	return Event{
		Name:       EventApplicationUpgradeDeclinedTerms,
		Properties: map[string]any{"plan": plan},
	}
}

// ApplicationUpgradeFailed is emitted when the Dashboard API returns an error.
// httpStatus is omitted when zero.
func ApplicationUpgradeFailed(plan, errorClass string, httpStatus int) Event {
	props := map[string]any{
		"plan":        plan,
		"error_class": errorClass,
	}
	if httpStatus != 0 {
		props["http_status"] = httpStatus
	}
	return Event{Name: EventApplicationUpgradeFailed, Properties: props}
}

// ApplicationUpgradeCompleted is emitted when the upgrade flow succeeds.
func ApplicationUpgradeCompleted(plan string) Event {
	return Event{
		Name:       EventApplicationUpgradeCompleted,
		Properties: map[string]any{"plan": plan},
	}
}

// --- Application downgrade events -------------------------------------------

// ApplicationDowngradeStarted is emitted when the user starts the downgrade flow.
func ApplicationDowngradeStarted(plan string) Event {
	return Event{
		Name: EventApplicationDowngradeStarted,
	}
}

// ApplicationDowngradeAcceptedTerms is emitted when the user accepts the T&C.
func ApplicationDowngradeAcceptedTerms(plan string) Event {
	return Event{
		Name:       EventApplicationDowngradeAcceptedTerms,
		Properties: map[string]any{"plan": plan},
	}
}

// ApplicationDowngradeDeclinedTerms is emitted when the user declines the T&C.
func ApplicationDowngradeDeclinedTerms(plan string) Event {
	return Event{
		Name:       EventApplicationDowngradeDeclinedTerms,
		Properties: map[string]any{"plan": plan},
	}
}

// ApplicationDowngradeFailed is emitted when the Dashboard API returns an error.
// httpStatus is omitted when zero.
func ApplicationDowngradeFailed(plan, errorClass string, httpStatus int) Event {
	props := map[string]any{
		"plan":        plan,
		"error_class": errorClass,
	}
	if httpStatus != 0 {
		props["http_status"] = httpStatus
	}
	return Event{Name: EventApplicationDowngradeFailed, Properties: props}
}

// ApplicationDowngradeCompleted is emitted when the downgrade flow succeeds.
func ApplicationDowngradeCompleted(plan string) Event {
	return Event{
		Name:       EventApplicationDowngradeCompleted,
		Properties: map[string]any{"plan": plan},
	}
}
