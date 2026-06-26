package providers

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

// ListOptions collects inputs for `agents providers list`.
type ListOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags

	Page    int
	PerPage int
	Show    bool
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
		Short:   "List configured LLM providers",
		Long: heredoc.Doc(`
			List provider authentications on Agent Studio for the active
			application.

			By default, structured output (--output json) masks the
			"apiKey" field (and similar secrets) with a "***" prefix.
			Pass --show-secret to render values verbatim — useful for
			scripted exports, dangerous in shared logs.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers list
			$ algolia agents providers list --output json --show-secret
			$ algolia agents providers list --page 2 --per-page 25
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runListCmd(opts)
		},
	}

	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-indexed; 0 = backend default)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Items per page (0 = backend default, currently 10)")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields (apiKey, ...) verbatim instead of masking")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	req := client.NewApiListProvidersRequest()
	if opts.Page > 0 {
		req = req.WithPage(shared.Int32(opts.Page))
	}
	if opts.PerPage > 0 {
		req = req.WithLimit(shared.Int32(opts.PerPage))
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching providers")
	res, err := client.ListProviders(req, agentStudio.WithContext(ctx))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		// Mask the nested input.apiKey secrets in the JSON we hand to the
		// printer unless --show-secret is set.
		var out any = res
		if !opts.Show {
			out = shared.MaskSecretsInValue(res)
		}
		return opts.PrintFlags.Print(opts.IO, out)
	}

	now := nowFn()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("NAME", nil, nil)
		table.AddField("PROVIDER", nil, nil)
		table.AddField("LAST USED", nil, nil)
		table.AddField("UPDATED", nil, nil)
		table.EndRow()
	}
	for _, p := range res.Data {
		table.AddField(p.Id, nil, nil)
		table.AddField(p.Name, nil, nil)
		table.AddField(p.ProviderName, nil, nil)
		table.AddField(relTimeOrDash(p.GetLastUsedAt(), now), nil, nil)
		table.AddField(relTimeOrDash(p.UpdatedAt, now), nil, nil)
		table.EndRow()
	}
	if err := table.Render(); err != nil {
		return err
	}
	if table.IsTTY() {
		fmt.Fprintf(opts.IO.Out,
			"\n%d provider(s) — page %d of %d (total %d).\n",
			len(res.Data),
			res.Pagination.Page, res.Pagination.TotalPages, res.Pagination.TotalCount)
	}
	return nil
}
