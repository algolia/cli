package keys

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
)

type ListOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags
	Page, Limit          int
	ShowSecret           bool
}

func newListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List secret keys",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (default 1)")
	cmd.Flags().IntVar(&opts.Limit, "limit", 0, "Page size")
	cmd.Flags().BoolVar(&opts.ShowSecret, "show-secret", false, "Reveal raw key values (default redacted as ***)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	req := client.NewApiListSecretKeysRequest()
	if opts.Page > 0 {
		req = req.WithPage(int32(opts.Page))
	}
	if opts.Limit > 0 {
		req = req.WithLimit(int32(opts.Limit))
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching secret keys")
	res, err := client.ListSecretKeys(req, agentStudio.WithContext(shared.OrBackground(opts.Ctx)))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	for i := range res.Data {
		res.Data[i] = maskKey(res.Data[i], opts.ShowSecret)
	}
	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, res)
	}
	now := nowFnOrTime()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("NAME", nil, nil)
		table.AddField("VALUE", nil, nil)
		table.AddField("DEFAULT", nil, nil)
		table.AddField("AGENTS", nil, nil)
		table.AddField("LAST USED", nil, nil)
		table.AddField("CREATED", nil, nil)
		table.EndRow()
	}
	for _, k := range res.Data {
		table.AddField(k.Id, nil, nil)
		table.AddField(k.Name, nil, nil)
		table.AddField(k.Value, nil, nil)
		table.AddField(boolWord(k.GetIsDefault()), nil, nil)
		table.AddField(fmt.Sprintf("%d", len(k.AgentIds)), nil, nil)
		table.AddField(relTimeOrDash(k.GetLastUsedAt(), now), nil, nil)
		table.AddField(relTimeOrDash(k.CreatedAt, now), nil, nil)
		table.EndRow()
	}
	if err := table.Render(); err != nil {
		return err
	}
	if table.IsTTY() {
		fmt.Fprintf(opts.IO.Out, "\n%d / %d key(s) (page %d).\n",
			len(res.Data), res.Pagination.TotalCount, res.Pagination.Page)
	}
	return nil
}

func boolWord(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// relTimeOrDash formats an RFC3339 timestamp as a humanized relative time, or
// "-" when the value is empty or unparseable.
func relTimeOrDash(ts string, now time.Time) string {
	if ts == "" {
		return "-"
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil || t.IsZero() {
		return "-"
	}
	return humanize.RelTime(now, t, "from now", "ago")
}
