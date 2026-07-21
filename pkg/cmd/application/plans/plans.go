package plans

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/shared/apputil"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

const reasonNoPaymentMethod = "no payment method on file"

type PlansOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
}

type planOutput struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	Price             string `json:"price"`
	AcceptTerms       string `json:"accept_terms"`
	Available         bool   `json:"available"`
	UnavailableReason string `json:"unavailable_reason,omitempty"`
}

func NewPlansCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &PlansOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "plans",
		Short: "List the available self-serve plans",
		Long: heredoc.Doc(`
			List the self-serve plans you can switch to with "algolia application upgrade"
			or "algolia application downgrade", along with their pricing.
		`),
		Example: heredoc.Doc(`
			# List the available plans
			$ algolia application plans

			# Output as JSON
			$ algolia application plans --output json
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlansCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runPlansCmd(opts *PlansOptions) error {
	cs := opts.IO.ColorScheme()

	client := opts.NewDashboardClient(auth.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching plans")
	plans, err := client.GetSelfServePlans(accessToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		newToken, reAuthErr := auth.ReauthenticateIfExpired(opts.IO, client, err)
		if reAuthErr != nil {
			return reAuthErr
		}
		accessToken = newToken
		opts.IO.StartProgressIndicatorWithLabel("Fetching plans")
		plans, err = client.GetSelfServePlans(accessToken)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}
	}

	var user *dashboard.DashboardUser
	opts.IO.StartProgressIndicatorWithLabel("Checking account")
	u, userErr := client.GetUser(accessToken)
	opts.IO.StopProgressIndicator()
	if userErr == nil {
		user = u
	}

	outputs := buildPlanOutputs(plans, user)

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, outputs)
	}

	if len(outputs) == 0 {
		fmt.Fprintf(opts.IO.Out, "%s No plans available.\n", cs.WarningIcon())
		return nil
	}

	for _, plan := range outputs {
		if !plan.Available {
			fmt.Fprintf(
				opts.IO.Out,
				"%s  %s\n",
				cs.Bold(plan.Name),
				cs.Yellowf("(unavailable: %s — add billing to unlock)", plan.UnavailableReason),
			)
			continue
		}
		fmt.Fprintf(opts.IO.Out, "%s  %s\n", cs.Bold(plan.Name), plan.Price)
		if plan.Description != "" {
			fmt.Fprintf(opts.IO.Out, "  %s\n", plan.Description)
		}
	}
	return nil
}

func buildPlanOutputs(plans []dashboard.Plan, user *dashboard.DashboardUser) []planOutput {
	outputs := make([]planOutput, 0, len(plans))
	for _, p := range plans {
		outputs = append(outputs, planOutput{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Type:        p.Type,
			Price:       p.Price,
			AcceptTerms: p.AcceptTerms,
			Available:   true,
		})
	}

	if user == nil || user.HasPaymentMethod {
		return outputs
	}

	for _, p := range apputil.KnownPaidPlans() {
		if apputil.PlanAvailable(plans, p.ID) {
			continue
		}
		outputs = append(outputs, planOutput{
			ID:                p.ID,
			Name:              p.Name,
			Type:              p.Type,
			Available:         false,
			UnavailableReason: reasonNoPaymentMethod,
		})
	}
	return outputs
}
