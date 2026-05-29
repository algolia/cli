// Package planchange holds the logic shared by the "application upgrade" and
// "application downgrade" commands. Both commands change the current
// application's self-serve plan and currently behave identically, so the
// fetch / billing-check / terms / change flow lives here once.
package planchange

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// Options carries everything the shared plan-change flow needs. The upgrade and
// downgrade commands populate it the same way.
type Options struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	Plan        string // --plan (optional): target plan, e.g. "free", "grow", "grow-plus"
	DryRun      bool   // --dry-run: preview without calling the API
	AcceptTerms bool   // --accept-terms: accept ToS in non-interactive mode

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

// changeResult is the structured (-o json) payload for a successful change.
type changeResult struct {
	ApplicationID string `json:"application_id"`
	Plan          string `json:"plan"`
	PlanName      string `json:"plan_name"`
	Price         string `json:"price"`
}

// Run executes the shared plan-change flow.
func Run(opts *Options) error {
	cs := opts.IO.ColorScheme()

	appID, err := opts.Config.Profile().GetApplicationID()
	if err != nil {
		return fmt.Errorf(
			"no current application configured; configure a profile with \"algolia profile add\" or \"algolia application select\" first: %w",
			err,
		)
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	token, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	var plans []dashboard.Plan
	if err := callWithReauth(opts.IO, client, &token, "Fetching plans", func(t string) error {
		var e error
		plans, e = client.GetSelfServePlans(t)
		return e
	}); err != nil {
		return err
	}
	if len(plans) == 0 {
		return fmt.Errorf("no self-serve plans are available")
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

	target, err := selectTarget(opts, plans)
	if err != nil {
		return err
	}

	// Paid plans require a payment method that the CLI cannot collect. Only
	// block when we positively know there is none (user fetched, flag false);
	// otherwise defer to the server.
	if !target.IsFree() && user != nil && !user.HasPaymentMethod {
		return fmt.Errorf(
			"the %q plan requires a payment method, which the CLI can't collect; add one in the Algolia dashboard (Settings → Billing) and try again",
			target.Name,
		)
	}

	if opts.DryRun {
		summary := map[string]any{
			"action":      "change_application_plan",
			"application": appID,
			"plan":        target.ID,
			"dryRun":      true,
		}
		return cmdutil.PrintRunSummary(
			opts.IO,
			opts.PrintFlags,
			summary,
			fmt.Sprintf("Dry run: would change application %s to the %q plan", appID, target.Name),
		)
	}

	accepted, err := confirmToS(opts, *target)
	if err != nil {
		return err
	}
	if !accepted {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Plan change aborted; no changes were made.\n",
			cs.WarningIcon(),
		)
		return nil
	}

	if err := callWithReauth(opts.IO, client, &token, "Changing plan", func(t string) error {
		_, e := client.ChangeApplicationPlan(t, appID, target.ID)
		return e
	}); err != nil {
		return err
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, changeResult{
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
	return nil
}

// selectTarget resolves the target plan from the --plan flag or, when
// interactive and no flag is set, an interactive picker over the available
// plans.
func selectTarget(opts *Options, plans []dashboard.Plan) (*dashboard.Plan, error) {
	if opts.Plan != "" {
		return resolvePlan(plans, opts.Plan)
	}

	if !opts.IO.CanPrompt() {
		return nil, cmdutil.FlagErrorf(
			"--plan is required in non-interactive mode (one of: %s)",
			strings.Join(planChoices(plans), ", "),
		)
	}

	return pickPlan(plans)
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
		"Invalid plan %q; valid plans: %s",
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

// confirmToS displays the plan's terms and asks the user to accept them. The
// prompt defaults to yes ([Y/n]). In non-interactive mode acceptance requires
// the --accept-terms flag (chosen over silent auto-accept).
func confirmToS(opts *Options, target dashboard.Plan) (bool, error) {
	cs := opts.IO.ColorScheme()

	terms := target.AcceptTerms
	if terms == "" {
		terms = fmt.Sprintf("By proceeding, you accept the Algolia %s Plan terms.", target.Name)
	}
	fmt.Fprintf(opts.IO.Out, "\n%s\n\n", terms)

	if !opts.IO.CanPrompt() {
		if opts.AcceptTerms {
			fmt.Fprintf(opts.IO.Out, "%s Terms accepted via --accept-terms.\n", cs.SuccessIcon())
			return true, nil
		}
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
