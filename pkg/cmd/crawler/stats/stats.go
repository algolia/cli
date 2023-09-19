package stats

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
)

// StatsOptions holds the options for the stats command.
type StatsOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	ID string

	PrintFlags *cmdutil.PrintFlags
}

// NewStatsCmd creates and returns a stats command for crawlers.
func NewStatsCmd(f *cmdutil.Factory, runF func(*StatsOptions) error) *cobra.Command {
	opts := &StatsOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
		PrintFlags:    cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:               "stats <crawler_id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Get statistics about a crawler",
		Long: heredoc.Doc(`
			Get a summary of the current status of crawled URLs for the specified crawler.
		`),
		Example: heredoc.Doc(`
			# Get statistics about the crawler with the ID "my-crawler"
			$ algolia crawler stats my-crawler
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}

			return runStatsCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runStatsCmd(opts *StatsOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	stats, err := client.Stats(opts.ID)
	if err != nil {
		return fmt.Errorf("cannot get stats: %w", err)
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}

		if err := p.Print(opts.IO, stats); err != nil {
			return err
		}

		return nil
	}

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("STATUS", nil, nil)
		table.AddField("CATEGORY", nil, nil)
		table.AddField("REASON", nil, nil)
		table.AddField("COUNT", nil, nil)

		table.EndRow()
	}

	status := func(s string) func(string) string {
		switch s {
		case "DONE":
			return cs.Green
		case "SKIPPED":
			return cs.Gray
		case "FAILED":
			return cs.Red
		default:
			return cs.Gray
		}
	}

	for _, stat := range stats.Data {
		table.AddField(status(stat.Status)(stat.Status), nil, nil)
		table.AddField(stat.Category, nil, nil)
		table.AddField(stat.Reason, nil, nil)
		table.AddField(fmt.Sprintf("%d", stat.Count), nil, nil)

		table.EndRow()
	}
	return table.Render()
}
