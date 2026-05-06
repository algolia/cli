package keys

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
)

type ListOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	Page, Limit       int
	ShowSecret        bool
}

func newListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
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
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching secret keys")
	res, err := client.ListSecretKeys(ctxOrBackground(opts.Ctx),
		agentstudio.ListSecretKeysParams{Page: opts.Page, Limit: opts.Limit})
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
		table.AddField(k.ID, nil, nil)
		table.AddField(k.Name, nil, nil)
		table.AddField(k.Value, nil, nil)
		table.AddField(boolWord(k.IsDefault), nil, nil)
		table.AddField(fmt.Sprintf("%d", len(k.AgentIDs)), nil, nil)
		if k.LastUsedAt != nil {
			table.AddField(humanize.RelTime(now, *k.LastUsedAt, "from now", "ago"), nil, nil)
		} else {
			table.AddField("-", nil, nil)
		}
		table.AddField(humanize.RelTime(now, k.CreatedAt, "from now", "ago"), nil, nil)
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
