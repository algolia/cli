package domains

import (
	"context"
	"fmt"
	"time"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags
	AgentID              string
}

func newListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags(),
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
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching allowed domains")
	res, err := client.ListAgentAllowedDomains(
		client.NewApiListAgentAllowedDomainsRequest(opts.AgentID),
		agentStudio.WithContext(shared.OrBackground(opts.Ctx)),
	)
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
		created := "-"
		if t, err := time.Parse(time.RFC3339, d.CreatedAt); err == nil {
			created = humanize.RelTime(now, t, "from now", "ago")
		}
		table.AddField(d.Id, nil, nil)
		table.AddField(d.Domain, nil, nil)
		table.AddField(created, nil, nil)
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
