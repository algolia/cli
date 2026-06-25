package list

import (
	"context"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

// nowFn is overridable for deterministic time-based output in tests.
var nowFn = time.Now

type ListOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig
	Ctx    context.Context

	AgentStudioAPIClient func() (*agentStudio.APIClient, error)

	PrintFlags *cmdutil.PrintFlags

	Page       int
	PerPage    int
	ProviderID string
}

func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:                   f.IOStreams,
		Config:               f.Config,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List Agent Studio agents on the active application",
		Long: heredoc.Doc(`
			List agents on Agent Studio for the active application.

			Pagination follows the backend defaults (10 per page) unless
			--page or --per-page is provided.
		`),
		Example: heredoc.Doc(`
			# List with backend defaults (page 1, 10 per page)
			$ algolia agents list

			# Second page, 25 items
			$ algolia agents list --page 2 --per-page 25

			# Filter by LLM provider
			$ algolia agents list --provider-id prov_123
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
		StringVar(&opts.ProviderID, "provider-id", "", "Filter by LLM provider authentication ID")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	req := client.NewApiListAgentsRequest()
	if opts.Page > 0 {
		req = req.WithPage(int32(opts.Page))
	}
	if opts.PerPage > 0 {
		req = req.WithLimit(int32(opts.PerPage))
	}
	if opts.ProviderID != "" {
		req = req.WithProviderId(opts.ProviderID)
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching agents")
	res, err := client.ListAgents(req, agentStudio.WithContext(ctx))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, res)
	}

	now := nowFn()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("NAME", nil, nil)
		table.AddField("STATUS", nil, nil)
		table.AddField("PROVIDER", nil, nil)
		table.AddField("UPDATED", nil, nil)
		table.EndRow()
	}

	for _, a := range res.Data {
		table.AddField(a.Id, nil, nil)
		table.AddField(a.Name, nil, nil)
		table.AddField(string(a.Status), nil, nil)
		table.AddField(stringOrDash(a.GetProviderId()), nil, nil)
		table.AddField(relTimeOrDash(a.GetUpdatedAt(), now), nil, nil)
		table.EndRow()
	}

	if err := table.Render(); err != nil {
		return err
	}

	if table.IsTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"\n%d agent(s) — page %d of %d (total %d).\n",
			len(res.Data),
			res.Pagination.Page,
			res.Pagination.TotalPages,
			res.Pagination.TotalCount,
		)
	}

	return nil
}

func stringOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func relTimeOrDash(ts string, now time.Time) string {
	if ts == "" {
		return "-"
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil || t.IsZero() {
		return "-"
	}
	// Order matches existing CLI usage (see pkg/cmd/apikeys/list): label
	// for "future" first, label for "past" second.
	return humanize.RelTime(now, t, "from now", "ago")
}
