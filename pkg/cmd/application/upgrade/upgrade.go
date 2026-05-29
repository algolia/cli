package upgrade

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/cmd/application/planchange"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/validators"
)

func NewUpgradeCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &planchange.Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the current application to a higher-tier plan",
		Long: heredoc.Doc(`
			Change the application associated with the current CLI profile to a
			higher-tier self-serve plan.

			Paid plans require a payment method on your account; the CLI can't
			collect card details, so add one in the Algolia dashboard first.
			You must accept the plan's terms of service before the change is applied.
		`),
		Example: heredoc.Doc(`
			# Pick a higher-tier plan interactively
			$ algolia application upgrade

			# Upgrade to a specific plan
			$ algolia application upgrade --plan grow

			# Non-interactive (accept terms up front)
			$ algolia application upgrade --plan grow-plus --accept-terms

			# Preview the change without applying it
			$ algolia application upgrade --plan grow --dry-run
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return planchange.Run(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Plan, "plan", "", "Target plan (free, grow, grow-plus)")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Preview the plan change without applying it")
	cmd.Flags().
		BoolVarP(&opts.AcceptTerms, "accept-terms", "y", false, "Accept the plan terms of service (required in non-interactive mode)")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}
