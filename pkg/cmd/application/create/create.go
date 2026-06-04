package create

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	pkgopen "github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/pkg/validators"
)

type CreateOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	Name        string
	Region      string
	ProfileName string
	Plan        string
	Default     bool
	DryRun      bool
	AcceptTerms bool

	nameProvided bool

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
	Browser            func(string) error
}

func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		Browser: pkgopen.Browser,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Algolia application",
		Long: heredoc.Doc(`
			Create a new Algolia application and optionally configure it as a CLI profile.
			Requires an active session (run "algolia auth login" first).`),
		Example: heredoc.Doc(`
			# Create an application interactively (prompts for name, plan, and terms)
			$ algolia application create

			# Create a Free application non-interactively
			$ algolia application create --name "My App" --region CA --accept-terms

			# Create on a paid plan (requires a payment method on file)
			$ algolia application create --name "My App" --region CA --plan grow --accept-terms

			# Create and set the new profile as the default
			$ algolia application create --name "My App" --region CA --accept-terms --default

			# Preview what would be created without actually creating it
			$ algolia application create --name "My App" --region CA --plan grow --dry-run
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.nameProvided = cmd.Flags().Changed("name")
			return runCreateCmd(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "My First Application", "Name for the application")
	cmd.Flags().StringVar(&opts.Region, "region", "", "Region code (e.g. EU, UK, USC, USE, USW)")
	cmd.Flags().
		StringVar(&opts.ProfileName, "profile-name", "", "Name for the CLI profile (defaults to app name)")
	cmd.Flags().
		StringVar(&opts.Plan, "plan", "", "Self-serve plan to create the application on (free, grow, grow-plus)")
	cmd.Flags().BoolVar(&opts.Default, "default", false, "Set the new profile as the default")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Preview the create request without sending it")
	cmd.Flags().
		BoolVarP(&opts.AcceptTerms, "accept-terms", "y", false, "Accept the selected plan's terms of service (required in non-interactive mode)")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runCreateCmd(ctx context.Context, opts *CreateOptions) error {
	cs := opts.IO.ColorScheme()

	name, err := resolveName(opts)
	if err != nil {
		return err
	}

	if opts.DryRun {
		planLabel := opts.Plan
		if planLabel == "" {
			planLabel = dashboard.PlanTypeFree
		}
		summary := map[string]any{
			"action":  "create_application",
			"name":    name,
			"region":  opts.Region,
			"plan":    planLabel,
			"default": opts.Default,
			"dryRun":  true,
		}
		return cmdutil.PrintRunSummary(
			opts.IO,
			opts.PrintFlags,
			summary,
			fmt.Sprintf(
				"Dry run: would create application %q in region %q on the %q plan",
				name,
				opts.Region,
				planLabel,
			),
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

	// Best-effort: continue without billing status if /1/user fails.
	var user *dashboard.DashboardUser
	if err := callWithReauth(opts.IO, client, &token, "Checking account", func(t string) error {
		var e error
		user, e = client.GetUser(t)
		return e
	}); err != nil {
		user = nil
	}

	target, err := selectPlan(opts, plans, user)
	if err != nil {
		return err
	}

	if !target.IsFree() {
		billingMissing := !apputil.PlanAvailable(plans, target.ID) ||
			(user != nil && !user.HasPaymentMethod)
		if billingMissing {
			return offerBilling(opts, client, *target)
		}
	}

	accepted, err := confirmToS(opts, *target)
	if err != nil {
		return err
	}
	if !accepted {
		telemetry.Track(
			ctx,
			telemetry.ApplicationCreateAborted(telemetry.TriggeredFromExplicitCommand),
		)
		fmt.Fprintf(opts.IO.Out, "%s Aborted; no application was created.\n", cs.WarningIcon())
		return nil
	}

	appDetails, err := apputil.CreateAndFetchApplication(
		ctx,
		opts.IO,
		client,
		token,
		opts.Region,
		name,
		telemetry.TriggeredFromExplicitCommand,
	)
	if err != nil {
		return err
	}

	if !target.IsFree() {
		if err := callWithReauth(opts.IO, client, &token, "Applying plan", func(t string) error {
			_, e := client.ChangeApplicationPlan(t, appDetails.ID, target.ID)
			return e
		}); err != nil {
			fmt.Fprintf(
				opts.IO.ErrOut,
				"%s Application %s was created on the Free plan, but applying the %s plan failed: %v\n",
				cs.WarningIcon(),
				cs.Bold(appDetails.ID),
				cs.Bold(target.Name),
				err,
			)
			fmt.Fprintf(
				opts.IO.ErrOut,
				"  Add a payment method if needed, then retry with: algolia application upgrade --plan %s\n",
				target.ID,
			)
			if !opts.structuredOutput() {
				_ = apputil.ConfigureProfile(
					opts.IO,
					opts.Config,
					appDetails,
					opts.ProfileName,
					opts.Default,
				)
			}
			return fmt.Errorf(
				"failed to apply the %q plan to application %s: %w",
				target.Name,
				appDetails.ID,
				err,
			)
		}
		fmt.Fprintf(
			opts.IO.Out,
			"%s Application %s created on the %s plan.\n",
			cs.SuccessIcon(),
			cs.Bold(appDetails.ID),
			cs.Bold(target.Name),
		)
	}

	if opts.structuredOutput() {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, appDetails)
	}

	return apputil.ConfigureProfile(
		opts.IO,
		opts.Config,
		appDetails,
		opts.ProfileName,
		opts.Default,
	)
}

func (opts *CreateOptions) structuredOutput() bool {
	return opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil
}

// resolveName returns the application name, prompting when interactive and --name is unset.
func resolveName(opts *CreateOptions) (string, error) {
	if opts.nameProvided || !opts.IO.CanPrompt() {
		return opts.Name, nil
	}

	var name string
	if err := prompt.SurveyAskOne(
		&survey.Input{
			Message: "Name:",
			Default: opts.Name,
		},
		&name,
	); err != nil {
		return "", err
	}
	if name == "" {
		name = opts.Name
	}
	return name, nil
}

// selectPlan resolves the target plan from --plan or an interactive picker.
func selectPlan(
	opts *CreateOptions,
	plans []dashboard.Plan,
	user *dashboard.DashboardUser,
) (*dashboard.Plan, error) {
	if opts.Plan != "" {
		target, err := apputil.ResolvePlan(plans, opts.Plan)
		if err == nil {
			return target, nil
		}
		if paid := apputil.KnownPaidPlan(opts.Plan); paid != nil {
			return paid, nil
		}
		return nil, err
	}

	if !opts.IO.CanPrompt() {
		free := apputil.FindFreePlan(plans)
		if free == nil {
			return nil, fmt.Errorf(
				"no free plan is available; pass --plan to choose one of: %s",
				strings.Join(apputil.PlanChoices(plans), ", "),
			)
		}
		return free, nil
	}

	hideNonFree := user != nil && !user.HasPaymentMethod
	candidates := apputil.SelectablePlans(plans, hideNonFree)
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no self-serve plans are available")
	}
	if hideNonFree {
		fmt.Fprintln(
			opts.IO.Out,
			"No payment method on file — only the Free plan is available. Add billing in the Algolia dashboard to unlock paid plans.",
		)
	}
	if len(candidates) == 1 {
		return &candidates[0], nil
	}
	return apputil.PickPlan(candidates)
}

// confirmToS shows the plan's terms and returns whether they were accepted.
func confirmToS(opts *CreateOptions, plan dashboard.Plan) (bool, error) {
	cs := opts.IO.ColorScheme()

	terms := plan.AcceptTerms
	if terms == "" {
		terms = fmt.Sprintf("By proceeding, you accept the Algolia %s Plan terms.", plan.Name)
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
	if err := prompt.Confirm("Do you accept these terms and want to create the application?", &accepted); err != nil {
		return false, err
	}
	return accepted, nil
}

// offerBilling tells the user a paid plan needs billing and offers the billing page.
func offerBilling(opts *CreateOptions, client *dashboard.Client, plan dashboard.Plan) error {
	cs := opts.IO.ColorScheme()
	url := client.DashboardURL + "/account/billing/details"

	fmt.Fprintf(
		opts.IO.Out,
		"\nThe %s plan requires a payment method on file before a paid application can be provisioned.\nThe CLI can't collect card details.\n",
		cs.Bold(plan.Name),
	)

	if opts.IO.CanPrompt() && opts.IO.IsStdoutTTY() {
		browser := opts.Browser
		if browser == nil {
			browser = pkgopen.Browser
		}
		open := true
		if err := prompt.Confirm("Open the billing page to add a payment method?", &open); err != nil {
			return err
		}
		if !open {
			return nil
		}
		fmt.Fprintf(opts.IO.Out, "Opening %s\n", cs.Bold(url))
		return browser(url)
	}

	fmt.Fprintf(
		opts.IO.Out,
		"Add a payment method here, then re-run with --plan %s:\n%s\n",
		plan.ID,
		url,
	)
	return fmt.Errorf("the %q plan requires a payment method; none is on file", plan.Name)
}

// callWithReauth runs fn, re-authenticating once and retrying on an expired session.
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
