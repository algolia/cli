// Package planchange holds the logic shared by the "application upgrade" and
// "application downgrade" commands. Both commands change the current
// application's self-serve plan and currently behave identically, so the
// fetch / billing-check / terms / change flow lives here once.
package planchange

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	pkgopen "github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/telemetry"
)

// Direction selects whether the flow offers higher-tier (upgrade) or
// lower-tier (downgrade) plans.
type Direction int

const (
	DirectionUpgrade Direction = iota
	DirectionDowngrade
)

// Options carries everything the shared plan-change flow needs. The upgrade and
// downgrade commands populate it the same way.
type Options struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	Direction   Direction // upgrade or downgrade
	Plan        string    // --plan (optional): target plan, e.g. "free", "grow", "grow-plus"
	DryRun      bool      // --dry-run: preview without calling the API
	AcceptTerms bool      // --accept-terms: accept ToS in non-interactive mode

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
	Browser            func(string) error
}

// changeResult is the structured (-o json) payload for a successful change.
type changeResult struct {
	ApplicationID string `json:"application_id"`
	Plan          string `json:"plan"`
	PlanName      string `json:"plan_name"`
	Price         string `json:"price"`
}

// telemetryDirection maps the flow direction to its telemetry value.
func (opts *Options) telemetryDirection() telemetry.Direction {
	if opts.Direction == DirectionDowngrade {
		return telemetry.DirectionDowngrade
	}
	return telemetry.DirectionUpgrade
}

// planChangeResult carries what the plan change flow produced, for telemetry.
type planChangeResult struct {
	changed     bool
	fromPlan    string
	toPlan      string
	abortReason telemetry.AbortReason
}

// Run executes the shared plan-change flow.
func Run(ctx context.Context, opts *Options) error {
	if opts.DryRun {
		// A dry run is not a funnel: no events, no tracker.
		_, err := changePlan(ctx, opts, nil)
		return err
	}

	direction := opts.telemetryDirection()
	tracker := telemetry.NewFlowTracker()
	telemetry.TrackEvent(ctx, telemetry.ApplicationPlanChangeStarted(direction))

	result, err := changePlan(ctx, opts, tracker)
	trackPlanChangeOutcome(ctx, direction, tracker, result, err)
	return err
}

// trackPlanChangeOutcome reports how the plan change flow ended: completed,
// aborted (with the reason why), or failed.
func trackPlanChangeOutcome(
	ctx context.Context,
	direction telemetry.Direction,
	tracker *telemetry.FlowTracker,
	result planChangeResult,
	err error,
) {
	switch {
	case result.changed:
		// The plan was changed: report Completed even when a post-success
		// step (courtesy prompt, output printing) failed.
		telemetry.TrackEvent(
			ctx,
			telemetry.ApplicationPlanChangeCompleted(direction, result.fromPlan, result.toPlan, tracker),
		)
	case err == nil || result.abortReason != "" || cmdutil.IsUserCancellation(err):
		// Stopped without changing anything: declined terms, already on the
		// plan, nothing to change to, billing wall, or user cancellation.
		reason := result.abortReason
		if reason == "" && cmdutil.IsUserCancellation(err) {
			reason = telemetry.AbortReasonCancelled
		}
		telemetry.TrackEvent(
			ctx,
			telemetry.ApplicationPlanChangeAborted(direction, tracker, reason),
		)
	default:
		telemetry.TrackEvent(ctx, telemetry.ApplicationPlanChangeFailed(direction, tracker, err))
	}
}

func changePlan(
	ctx context.Context,
	opts *Options,
	tracker *telemetry.FlowTracker,
) (planChangeResult, error) {
	var result planChangeResult
	cs := opts.IO.ColorScheme()

	tracker.SetStep(telemetry.StepAuth)
	appID, err := opts.Config.Profile().GetApplicationID()
	if err != nil {
		return result, fmt.Errorf(
			"no current application configured; run \"algolia auth login\" or \"algolia application select\" first: %w",
			err,
		)
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	token, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return result, err
	}

	tracker.SetStep(telemetry.StepPlan)
	var plans []dashboard.Plan
	if err := callWithReauth(opts.IO, client, &token, "Fetching plans", func(t string) error {
		var e error
		plans, e = client.GetSelfServePlans(t)
		return e
	}); err != nil {
		return result, err
	}
	if len(plans) == 0 {
		return result, fmt.Errorf("no self-serve plans are available")
	}

	// Billing status is best-effort. If /1/user is unavailable we continue
	// without it and let the server enforce billing validity.
	var user *dashboard.DashboardUser
	if err := callWithReauth(opts.IO, client, &token, "Checking account", func(t string) error {
		var e error
		user, e = client.GetUser(t)
		return e
	}); err != nil {
		user = nil
	}

	app := fetchApplication(opts, client, &token, appID)
	if app != nil {
		result.fromPlan = currentPlanTelemetryID(plans, app)
	}

	target, err := resolveTarget(opts, appID, app, plans, user)
	if err != nil {
		return result, err
	}
	if target == nil {
		result.abortReason = telemetry.AbortReasonNoCandidates
		return result, nil
	}
	result.toPlan = apputil.PlanTelemetryID(*target)

	if isCurrentPlan(app, *target) {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Application %s is already on the %s plan; no change needed.\n",
			cs.WarningIcon(),
			cs.Bold(appID),
			cs.Bold(target.Name),
		)
		result.abortReason = telemetry.AbortReasonAlreadyOnPlan
		return result, nil
	}

	// Paid plans require a payment method that the CLI cannot collect. Block
	// when the server hides the plan (paid plans appear only once billing is
	// on file) or when the user record says no payment method is on file.
	if !target.IsFree() &&
		(!apputil.PlanAvailable(plans, target.ID) || (user != nil && !user.HasPaymentMethod)) {
		result.abortReason = telemetry.AbortReasonBillingRequired
		return result, apputil.OfferBilling(opts.IO, opts.Browser, client.DashboardURL, *target)
	}

	if opts.DryRun {
		summary := map[string]any{
			"action":      "change_application_plan",
			"application": appID,
			"plan":        target.ID,
			"dryRun":      true,
		}
		return result, cmdutil.PrintRunSummary(
			opts.IO,
			opts.PrintFlags,
			summary,
			fmt.Sprintf("Dry run: would change application %s to the %q plan", appID, target.Name),
		)
	}

	// With --plan the interactive picker (which shows the current application)
	// is skipped; show which application the terms apply to before asking.
	if opts.Plan != "" {
		printCurrentApplication(opts, appID, app)
	}

	tracker.SetStep(telemetry.StepTerms)
	accepted, err := confirmToS(opts, *target)
	if err != nil {
		return result, err
	}
	if !accepted {
		telemetry.TrackEvent(
			ctx,
			telemetry.ApplicationPlanChangeDeclinedTerms(
				opts.telemetryDirection(),
				apputil.PlanTelemetryID(*target),
			),
		)
		fmt.Fprintf(
			opts.IO.Out,
			"%s Plan change aborted; no changes were made.\n",
			cs.WarningIcon(),
		)
		result.abortReason = telemetry.AbortReasonDeclinedTerms
		return result, nil
	}
	telemetry.TrackEvent(
		ctx,
		telemetry.ApplicationPlanChangeAcceptedTerms(
			opts.telemetryDirection(),
			apputil.PlanTelemetryID(*target),
		),
	)

	tracker.SetStep(telemetry.StepAPICall)
	if err := callWithReauth(opts.IO, client, &token, "Changing plan", func(t string) error {
		_, e := client.ChangeApplicationPlan(t, appID, target.ID)
		return e
	}); err != nil {
		return result, err
	}
	result.changed = true

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return result, err
		}
		return result, p.Print(opts.IO, changeResult{
			ApplicationID: appID,
			Plan:          target.ID,
			PlanName:      target.Name,
			Price:         target.Price,
		})
	}

	fmt.Fprintf(
		opts.IO.Out,
		"%s Application %s changed to the %s plan.\n",
		cs.SuccessIcon(),
		cs.Bold(appID),
		cs.Bold(target.Name),
	)

	if target.IsFree() {
		return result, nil
	}

	return result, offerCostManagementBudget(opts, client.DashboardURL, appID)
}

// offerCostManagementBudget tells the user they can create a budget and, when
// confirmed, opens the cost management page in the browser.
func offerCostManagementBudget(opts *Options, dashboardURL, appID string) error {
	if !opts.IO.CanPrompt() || !opts.IO.IsStdoutTTY() {
		return nil
	}

	browser := opts.Browser
	if browser == nil {
		browser = pkgopen.Browser
	}

	cs := opts.IO.ColorScheme()
	fmt.Fprintf(
		opts.IO.Out,
		"\nWant to create a budget and monitor your costs?\n",
	)

	openPage := true
	if err := prompt.Confirm("Open cost management?", &openPage); err != nil {
		return err
	}
	if !openPage {
		return nil
	}

	url := fmt.Sprintf("%s/account/billing/cost-management?applicationId=%s", dashboardURL, appID)
	fmt.Fprintf(opts.IO.Out, "Opening %s\n", cs.Bold(url))

	return browser(url)
}

// resolveTarget picks the target plan: --plan overrides the direction filter,
// otherwise candidates are filtered by direction and chosen interactively. A
// nil plan with a nil error means there is nothing to switch to.
func resolveTarget(
	opts *Options,
	appID string,
	app *dashboard.Application,
	plans []dashboard.Plan,
	user *dashboard.DashboardUser,
) (*dashboard.Plan, error) {
	if opts.Plan != "" {
		target, err := resolvePlan(plans, opts.Plan)
		if err == nil {
			return target, nil
		}
		if paid := apputil.KnownPaidPlan(opts.Plan); paid != nil {
			return paid, nil
		}
		return nil, err
	}

	candidates := filterByDirection(withKnownPaidPlans(plans, user), app, opts.Direction)

	if len(candidates) == 0 {
		reportNoCandidates(opts, appID, app, opts.Direction)
		return nil, nil
	}

	if !opts.IO.CanPrompt() {
		return nil, cmdutil.FlagErrorf(
			"--plan is required in non-interactive mode (one of: %s)",
			strings.Join(planChoices(candidates), ", "),
		)
	}

	printCurrentApplication(opts, appID, app)

	return pickPlan(candidates)
}

// fetchApplication returns the current application, or nil if it can't be fetched.
func fetchApplication(
	opts *Options,
	client *dashboard.Client,
	token *string,
	appID string,
) *dashboard.Application {
	var app *dashboard.Application
	if err := callWithReauth(opts.IO, client, token, "Fetching application", func(t string) error {
		var e error
		app, e = client.GetApplication(t, appID)
		return e
	}); err != nil {
		return nil
	}
	return app
}

// filterByDirection returns plans above (upgrade) or below (downgrade) the
// current plan in the API's tier order, or all plans when it isn't found.
func filterByDirection(
	plans []dashboard.Plan,
	app *dashboard.Application,
	dir Direction,
) []dashboard.Plan {
	idx := currentPlanIndex(plans, app)
	if idx < 0 {
		return plans
	}
	if dir == DirectionDowngrade {
		return plans[:idx]
	}
	return plans[idx+1:]
}

func withKnownPaidPlans(plans []dashboard.Plan, user *dashboard.DashboardUser) []dashboard.Plan {
	if user == nil || user.HasPaymentMethod {
		return plans
	}
	augmented := append([]dashboard.Plan(nil), plans...)
	for _, p := range apputil.KnownPaidPlans() {
		if !apputil.PlanAvailable(augmented, p.ID) {
			augmented = append(augmented, p)
		}
	}
	return augmented
}

// normalizePlanKey normalizes a plan label/name for comparison; the CLI joins
// the current plan to the self-serve list by matching label against Name.
func normalizePlanKey(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// currentPlanTelemetryID maps the app's current plan to the identifier used in
// telemetry properties, falling back to the normalized label when the plan is
// not in the self-serve list.
func currentPlanTelemetryID(plans []dashboard.Plan, app *dashboard.Application) string {
	if idx := currentPlanIndex(plans, app); idx >= 0 {
		return apputil.PlanTelemetryID(plans[idx])
	}
	return normalizePlanKey(app.PlanLabel)
}

// currentPlanIndex returns the index of the app's current plan in plans, or -1.
func currentPlanIndex(plans []dashboard.Plan, app *dashboard.Application) int {
	if app == nil {
		return -1
	}
	label := normalizePlanKey(app.PlanLabel)
	if label == "" {
		return -1
	}
	for i := range plans {
		if normalizePlanKey(plans[i].Name) == label {
			return i
		}
	}
	return -1
}

// isCurrentPlan reports whether target is the plan the application is already on.
func isCurrentPlan(app *dashboard.Application, target dashboard.Plan) bool {
	if app == nil {
		return false
	}
	label := normalizePlanKey(app.PlanLabel)
	return label != "" && label == normalizePlanKey(target.Name)
}

// reportNoCandidates tells the user they're already at the highest/lowest plan.
func reportNoCandidates(
	opts *Options,
	appID string,
	app *dashboard.Application,
	dir Direction,
) {
	cs := opts.IO.ColorScheme()
	current := ""
	if app != nil && app.PlanLabel != "" {
		current = fmt.Sprintf(" (%s)", app.PlanLabel)
	}
	tier, verb := "highest", "upgrade"
	if dir == DirectionDowngrade {
		tier, verb = "lowest", "downgrade"
	}
	fmt.Fprintf(
		opts.IO.Out,
		"%s Application %s is already on the %s self-serve plan%s; nothing to %s to.\n",
		cs.WarningIcon(),
		cs.Bold(appID),
		tier,
		current,
		verb,
	)
}

// printCurrentApplication prints the current app and plan before the picker
// or, with --plan, before the terms confirmation.
func printCurrentApplication(opts *Options, appID string, app *dashboard.Application) {
	cs := opts.IO.ColorScheme()
	label := cs.Bold(appID)
	if app != nil && app.Name != "" {
		label = fmt.Sprintf("%s (%s)", cs.Bold(appID), app.Name)
	}
	if app != nil && app.PlanLabel != "" {
		fmt.Fprintf(
			opts.IO.Out,
			"Current application: %s — current plan: %s\n\n",
			label,
			cs.Bold(app.PlanLabel),
		)
		return
	}
	fmt.Fprintf(opts.IO.Out, "Current application: %s\n\n", label)
}

// resolvePlan maps a --plan value to one of the fetched plans.
func resolvePlan(plans []dashboard.Plan, value string) (*dashboard.Plan, error) {
	// Exact match on the plan id (configuration.plan).
	for i := range plans {
		if plans[i].ID == value {
			return &plans[i], nil
		}
	}
	// The user-facing "free" choice maps to the free-type template, whose id is
	// not fixed (it can be "build"); match on type rather than a hard-coded id.
	if value == dashboard.PlanTypeFree {
		for i := range plans {
			if plans[i].IsFree() {
				return &plans[i], nil
			}
		}
	}
	return nil, cmdutil.FlagErrorf(
		"Invalid plan %q; available plans: %s",
		value,
		strings.Join(planChoices(plans), ", "),
	)
}

// pickPlan shows an interactive selector over the candidate plans.
func pickPlan(candidates []dashboard.Plan) (*dashboard.Plan, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no plans are available")
	}
	labels := make([]string, len(candidates))
	for i, p := range candidates {
		if p.Price == "" {
			p.Price = "Pay as you go"
		}
		labels[i] = fmt.Sprintf("%s — %s", p.Name, p.Price)
	}
	var selected int
	if err := prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select a plan:",
			Options: labels,
		},
		&selected,
	); err != nil {
		return nil, err
	}
	return &candidates[selected], nil
}

// confirmToS shows the plan's terms and returns whether they were accepted.
func confirmToS(opts *Options, target dashboard.Plan) (bool, error) {
	cs := opts.IO.ColorScheme()

	terms := target.AcceptTerms
	if terms == "" {
		terms = fmt.Sprintf("By proceeding, you accept the Algolia %s Plan terms.", target.Name)
	}
	fmt.Fprintf(opts.IO.Out, "\n%s\n\n", terms)

	if opts.AcceptTerms {
		fmt.Fprintf(opts.IO.Out, "%s Terms accepted via --accept-terms.\n", cs.SuccessIcon())
		return true, nil
	}

	if !opts.IO.CanPrompt() {
		return false, cmdutil.FlagErrorf(
			"the plan terms must be accepted in non-interactive mode; pass --accept-terms to confirm",
		)
	}

	accepted := true
	if err := prompt.Confirm("Do you accept these terms and want to change the plan?", &accepted); err != nil {
		return false, err
	}
	return accepted, nil
}

// planChoices returns the user-facing plan identifiers (the free plan is shown
// as "free" regardless of its underlying id).
func planChoices(plans []dashboard.Plan) []string {
	choices := make([]string, 0, len(plans))
	for _, p := range plans {
		if p.IsFree() {
			choices = append(choices, dashboard.PlanTypeFree)
		} else {
			choices = append(choices, p.ID)
		}
	}
	return choices
}

// callWithReauth runs fn with the current token; if it fails with a
// session-expired error it re-authenticates once and retries, mirroring the
// retry pattern used by the other application commands.
func callWithReauth(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	token *string,
	label string,
	fn func(token string) error,
) error {
	io.StartProgressIndicatorWithLabel(label)
	err := fn(*token)
	io.StopProgressIndicator()
	if err == nil {
		return nil
	}

	newToken, reAuthErr := auth.ReauthenticateIfExpired(io, client, err)
	if reAuthErr != nil {
		return reAuthErr
	}
	*token = newToken

	io.StartProgressIndicatorWithLabel(label)
	err = fn(*token)
	io.StopProgressIndicator()
	return err
}
