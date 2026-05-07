package domains

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	AgentID           string
}

func newListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:     "list <agent-id>",
		Aliases: []string{"ls"},
		Short:   "List allowed domains for an agent",
		Args:    validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching allowed domains")
	res, err := client.ListAllowedDomains(shared.OrBackground(opts.Ctx), opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, res)
	}
	now := nowFnOrTime()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("DOMAIN", nil, nil)
		table.AddField("CREATED", nil, nil)
		table.EndRow()
	}
	for _, d := range res.Domains {
		table.AddField(d.ID, nil, nil)
		table.AddField(d.Domain, nil, nil)
		table.AddField(humanize.RelTime(now, d.CreatedAt, "from now", "ago"), nil, nil)
		table.EndRow()
	}
	if err := table.Render(); err != nil {
		return err
	}
	if table.IsTTY() {
		fmt.Fprintf(opts.IO.Out, "\n%d domain(s).\n", len(res.Domains))
	}
	return nil
}
