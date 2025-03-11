package run

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// RunOptions holds the options for the stats command.
type RunOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	ID string
}

// NewRunCmd creates and returns a run command for crawlers.
func NewRunCmd(f *cmdutil.Factory, runF func(*RunOptions) error) *cobra.Command {
	opts := &RunOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
	}
	cmd := &cobra.Command{
		Use:               "run <crawler_id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Start or resume a crawler",
		Long: heredoc.Doc(`
			Unpause the specified crawler.
			Previously ongoing crawls will be resumed. Otherwise, the crawler waits for its next scheduled run.
		`),
		Example: heredoc.Doc(`
			# Run the crawler with the ID "my-crawler"
			$ algolia crawler run my-crawler
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}

			return runRunCmd(opts)
		},
	}

	return cmd
}

func runRunCmd(opts *RunOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	_, err = client.Run(opts.ID)
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Crawler %s started\n", cs.SuccessIconWithColor(cs.Green), opts.ID)
	}

	return nil
}
