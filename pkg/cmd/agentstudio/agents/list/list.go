package list

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	Page  int
	Limit int

	PrintFlags *cmdutil.PrintFlags
}

func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Args:    validators.NoArgs(),
		Short:   "List Agent Studio agents",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# List the first page of agents
			$ algolia agentstudio agents list

			# Page through results
			$ algolia agentstudio agents list --page 2 --limit 50
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}

	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-based; default API behavior when omitted)")
	cmd.Flags().IntVar(&opts.Limit, "limit", 0, "Items per page (default API behavior when omitted)")
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching agents")
	res, err := client.ListAgents(opts.Page, opts.Limit)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, res)
	}

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("NAME", nil, nil)
		table.AddField("STATUS", nil, nil)
		table.AddField("UPDATED", nil, nil)
		table.EndRow()
	}
	for _, a := range res.Data {
		table.AddField(a.ID, nil, nil)
		table.AddField(a.Name, nil, nil)
		table.AddField(a.Status, nil, nil)
		table.AddField(a.UpdatedAt, nil, nil)
		table.EndRow()
	}
	return table.Render()
}
