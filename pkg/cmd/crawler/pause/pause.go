package pause

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

type PauseOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	IDs []string
}

// NewPauseCmd creates and returns a pause command for Crawlers.
func NewPauseCmd(f *cmdutil.Factory, runF func(*PauseOptions) error) *cobra.Command {
	opts := &PauseOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
	}
	cmd := &cobra.Command{
		Use:               "pause <crawler_id>...",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Pause one or multiple crawlers",
		Long: heredoc.Doc(`
			Pauses the specified crawler.
		`),
		Example: heredoc.Doc(`
			# Pause the crawler with the ID "my-crawler"
			$ algolia crawler pause my-crawler

			# Pause the crawlers with the IDs "my-crawler-1" and "my-crawler-2"
			$ algolia crawler pause my-crawler-1 my-crawler-2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.IDs = args
			if runF != nil {
				return runF(opts)
			}

			return runPauseCmd(opts)
		},
	}

	return cmd
}

func runPauseCmd(opts *PauseOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf("Pausing %s", utils.Pluralize(len(opts.IDs), "crawler")),
	)
	for _, id := range opts.IDs {
		if _, err := client.Reindex(id); err != nil {
			opts.IO.StopProgressIndicator()
			return fmt.Errorf("cannot pause crawler %s: %w", cs.Bold(id), err)
		}
	}
	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s %s\n",
			cs.SuccessIconWithColor(cs.Green),
			fmt.Sprintf("Successfully paused %s", utils.Pluralize(len(opts.IDs), "crawler")),
		)
	}

	return nil
}
