package reindex

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type ReindexOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	IDs []string
}

// NewReindexCmd creates and returns a reindex command for Crawlers.
func NewReindexCmd(f *cmdutil.Factory, runF func(*ReindexOptions) error) *cobra.Command {
	opts := &ReindexOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
	}
	cmd := &cobra.Command{
		Use:               "reindex <crawler_id>...",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Reindexs the specified crawlers",
		Long: heredoc.Doc(`
			Request the specified crawler to start (or restart) crawling.
		`),
		Example: heredoc.Doc(`
			# Reindex the crawler with the ID "my-crawler"
			$ algolia crawler reindex my-crawler

			# Reindex the crawlers with the IDs "my-crawler-1" and "my-crawler-2"
			$ algolia crawler reindex my-crawler-1 my-crawler-2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IDs = args
			if runF != nil {
				return runF(opts)
			}

			return runReindexCmd(opts)
		},
	}

	return cmd
}

func runReindexCmd(opts *ReindexOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Reindexing %s", utils.Pluralize(len(opts.IDs), "crawler")))
	for _, id := range opts.IDs {
		if _, err := client.Reindex(id); err != nil {
			opts.IO.StopProgressIndicator()
			return fmt.Errorf("cannot reindex crawler %s: %w", cs.Bold(id), err)
		}
	}
	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.SuccessIconWithColor(cs.Green), fmt.Sprintf("Successfully requested reindexing for %s", utils.Pluralize(len(opts.IDs), "crawler")))
	}

	return nil
}
