package list

import (
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
)

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID string
	Page    int
	Limit   int

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
		Use:   "list <agent_id>",
		Args:  cobra.ExactArgs(1),
		Short: "List conversations for an agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			$ algolia agentstudio conversations list a1b2 --limit 50
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
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
	opts.IO.StartProgressIndicatorWithLabel("Fetching conversations")
	res, err := client.ListConversations(opts.AgentID, opts.Page, opts.Limit)
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
		table.AddField("TITLE", nil, nil)
		table.AddField("MESSAGES", nil, nil)
		table.AddField("UPDATED", nil, nil)
		table.EndRow()
	}
	for _, c := range res.Data {
		title := ""
		if c.Title != nil {
			title = *c.Title
		}
		count := ""
		if c.MessageCount != nil {
			count = strconv.Itoa(*c.MessageCount)
		}
		table.AddField(c.ID, nil, nil)
		table.AddField(title, nil, nil)
		table.AddField(count, nil, nil)
		table.AddField(c.UpdatedAt, nil, nil)
		table.EndRow()
	}
	return table.Render()
}
