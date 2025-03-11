package list

import (
	"fmt"
	"sort"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	Name  string
	AppID string

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd creates and returns a list command for Crawlers.
func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
		PrintFlags:    cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:   "list",
		Args:  validators.NoArgs(),
		Short: "List crawlers",
		Long: heredoc.Doc(`
			List crawlers, optionally filtered by name or appID.
		`),
		Example: heredoc.Doc(`
			# List all crawlers
			$ algolia crawler list

			# List crawlers with the name "my-crawler"
			$ algolia crawler list --name my-crawler

			# List crawlers with the appID "my-app-id"
			$ algolia crawler list --app-id my-app-id
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runListCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by name")
	cmd.Flags().StringVar(&opts.AppID, "app-id", "", "Filter by appID")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching Crawlers")
	crawlersList, err := client.ListAll(opts.Name, opts.AppID)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	crawlers := make([]crawler.Crawler, 0, len(crawlersList))
	for _, item := range crawlersList {
		opts.IO.UpdateProgressIndicatorLabel(fmt.Sprintf("Fetching Crawler %s details", item.ID))
		c, err := client.Get(item.ID, true)
		if err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
		// We need to set the ID here because the API doesn't return it
		c.ID = item.ID
		crawlers = append(crawlers, *c)
	}
	opts.IO.StopProgressIndicator()

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		for _, crawler := range crawlers {
			if err := p.Print(opts.IO, crawler); err != nil {
				return err
			}
		}
		return nil
	}

	cs := opts.IO.ColorScheme()
	status := func(c crawler.Crawler) string {
		if c.Blocked {
			return cs.Red("Blocked")
		}
		// If default unix time, it means the crawler has never ran
		if c.LastReindexStartedAt.Unix() < 0 {
			return cs.Gray("Created")
		}
		if !c.Running {
			return cs.Yellow("Paused")
		}
		if c.Config.Schedule != "" {
			return cs.Blue("Scheduled")
		} else {
			// The crawler is running if `LastReindexStartedAt` > `LastReindexEndedAt`
			if c.LastReindexStartedAt.After(c.LastReindexEndedAt) {
				return cs.Green("Running")
			} else {
				return cs.Green("Finished")
			}
		}
	}

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("ID", nil, nil)
		table.AddField("NAME", nil, nil)
		table.AddField("STATUS", nil, nil)
		table.AddField("LAST REINDEX", nil, nil)

		table.EndRow()
	}

	// Sort Crawlers by updatedAt
	sort.Slice(crawlers, func(i, j int) bool {
		return crawlers[i].UpdatedAt.After(crawlers[j].UpdatedAt)
	})

	for _, crawler := range crawlers {
		table.AddField(crawler.ID, nil, nil)
		table.AddField(crawler.Name, nil, nil)
		table.AddField(status(crawler), nil, nil)

		if crawler.LastReindexStartedAt.IsZero() {
			table.AddField("N/A", nil, nil)
			table.EndRow()
			continue
		} else {
			startedAt := crawler.LastReindexStartedAt.Format("2006-01-02 15:04:05")
			if crawler.LastReindexStartedAt.Before(crawler.LastReindexEndedAt) {
				duration := crawler.LastReindexEndedAt.Sub(crawler.LastReindexStartedAt)
				table.AddField(fmt.Sprintf("%s (%s)", startedAt, duration.String()), nil, nil)
			} else {
				table.AddField(startedAt, nil, nil)
			}
		}

		table.EndRow()
	}
	return table.Render()
}
