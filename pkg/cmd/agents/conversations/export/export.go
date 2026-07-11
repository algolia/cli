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

	AgentID   string
	StartDate string
	EndDate   string

	PrintFlags *cmdutil.PrintFlags
}

// NewExportCmd returns the `agents conversations export` command.
func NewExportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ExportOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "export <agent-id>",
		Short: "Export the conversations of an agent",
		Args: validators.ExactArgsWithMsg(
			1,
			"agents conversations export requires an <agent-id> argument.",
		),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# Export all conversations of the agent "my-agent"
			$ algolia agents conversations export my-agent

			# Export conversations created between two dates
			$ algolia agents conversations export my-agent --start-date 2026-01-01 --end-date 2026-01-31
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			return runExportCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Filter conversations created after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Filter conversations created before this date (YYYY-MM-DD)")

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
