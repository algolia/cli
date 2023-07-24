package get

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	ID         string
	ConfigOnly bool

	PrintFlags *cmdutil.PrintFlags
}

// NewGetCmd creates and returns a get command for Crawlers.
func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
		PrintFlags:    cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:               "get <crawler_id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Get a crawler",
		Long: heredoc.Doc(`
			Get the specified crawler.
		`),
		Example: heredoc.Doc(`
			# Get the crawler with the ID "my-crawler"
			$ algolia crawler get my-crawler

			# Get the crawler with the ID "my-crawler" and display only its configuration
			$ algolia crawler get my-crawler --config-only
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}

			return runGetCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.ConfigOnly, "config-only", "c", false, "Display only the crawler configuration")

	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Fetching crawler %s", opts.ID))
	crawler, err := client.Get(opts.ID, true)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	var toPrint interface{}
	if opts.ConfigOnly {
		toPrint = crawler.Config
	} else {
		toPrint = crawler
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	if err := p.Print(opts.IO, toPrint); err != nil {
		return err
	}

	return nil
}
