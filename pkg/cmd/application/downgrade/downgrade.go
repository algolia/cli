package downgrade

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/cmd/application/planchange"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/validators"
)

func NewDowngradeCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &planchange.Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		Direction:  planchange.DirectionDowngrade,
		PrintFlags: cmdutil.NewPrintFlags(),
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "downgrade",
		Short: "Downgrade the current application to a lower-tier plan",
		Long: heredoc.Doc(`
			Change the application associated with the current CLI profile to a
			lower-tier self-serve plan.

			You must accept the target plan's terms of service before the change
			is applied.
		`),
		Example: heredoc.Doc(`
			# Pick a lower-tier plan interactively
			$ algolia application downgrade

			# Downgrade to a specific plan
			$ algolia application downgrade --plan free

			# Non-interactive (accept terms up front)
			$ algolia application downgrade --plan grow --accept-terms

			# Preview the change without applying it
			$ algolia application downgrade --plan free --dry-run
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return planchange.Run(cmd.Context(), opts)
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
