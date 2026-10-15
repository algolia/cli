package export

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// ExportOptions holds the dependencies and flags for the export command.
type ExportOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID   string
	StartDate string
	EndDate   string
}

// NewExportCmd returns the `conversations export` command.
func NewExportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "export <agent-id>",
		Short: "Export Agent Studio conversations for an agent",
		Args: validators.ExactArgsWithMsg(
			1,
			"agentstudio conversations export requires a <agent-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# Export all conversations for an agent
			$ algolia agentstudio conversations export 11111111-1111-4111-8111-111111111111

			# Export a date range
			$ algolia agentstudio conversations export 11111111-1111-4111-8111-111111111111 --start-date 2026-01-01 --end-date 2026-01-31
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			return runExportCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Created after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Created before this date (YYYY-MM-DD)")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runExportCmd(opts *ExportOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	req := client.NewApiExportConversationsRequest(opts.AgentID)
	if opts.StartDate != "" {
		req = req.WithStartDate(opts.StartDate)
	}
	if opts.EndDate != "" {
		req = req.WithEndDate(opts.EndDate)
	}

	opts.IO.StartProgressIndicatorWithLabel("Exporting conversations")

	res, err := client.ExportConversations(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
