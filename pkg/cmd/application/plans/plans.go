package plans

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type PlansOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	PrintFlags *cmdutil.PrintFlags

	NewDashboardClient func(clientID string) *dashboard.Client
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

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, plans)
	}

	if len(plans) == 0 {
		fmt.Fprintf(opts.IO.Out, "%s No plans available.\n", cs.WarningIcon())
		return nil
	}

	for _, plan := range plans {
		fmt.Fprintf(opts.IO.Out, "%s  %s\n", cs.Bold(plan.Name), plan.Price)
		if plan.Description != "" {
			fmt.Fprintf(opts.IO.Out, "  %s\n", plan.Description)
		}
	}
	return nil
}
