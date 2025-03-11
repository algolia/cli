package crawl

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

// CrawlOptions holds the options for the crawl command.
type CrawlOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	ID            string
	URLs          []string
	Save          bool
	SaveSpecified bool
}

// NewCrawlCmd creates and returns a crawl command for crawlers.
func NewCrawlCmd(f *cmdutil.Factory, runF func(*CrawlOptions) error) *cobra.Command {
	opts := &CrawlOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
	}
	cmd := &cobra.Command{
		Use:               "crawl <crawler_id> --urls <url>...",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Crawl specific URLs",
		Long: heredoc.Doc(`
			Immediately crawl the given URLs.
			The generated records are pushed to the live index if there's no ongoing reindex, and to the temporary index otherwise.
		`),
		Example: heredoc.Doc(`
			# Crawl the URLs "https://www.example.com" and "https://www.example2.com/" for the crawler with the ID "my-crawler"
			$ algolia crawler crawl my-crawler --urls https://www.example.com,https://www.example2.com/

			# Crawl the URLs "https://www.example.com" and "https://www.example2.com/" for the crawler with the ID "my-crawler" and save them in the configuration
			$ algolia crawler crawl my-crawler --urls https://www.example.com,https://www.example2.com/ --save

			# Crawl the URLs "https://www.example.com" and "https://www.example2.com/" for the crawler with the ID "my-crawler" and don't save them in the configuration
			$ algolia crawler crawl my-crawler --urls https://www.example.com,https://www.example2.com/ --save=false
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}

			// We need to know if the `save` flag was specified or not.
			if cmd.Flags().Changed("save") {
				opts.SaveSpecified = true
			}

			return runCrawlCmd(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.URLs, "urls", "u", nil, "The URLs to crawl (maximum 50).")
	_ = cmd.MarkFlagRequired("urls")

	cmd.Flags().BoolVarP(&opts.Save, "save", "s", false, heredoc.Doc(`
		When true, the given URLs are added to the extraUrls list of your configuration (unless already present in startUrls or sitemaps).
		When false, the URLs aren't saved in the configuration.
		When unspecified, the URLs are added to the extraUrls list of your configuration, but only if they haven't been indexed during the last reindex, and they aren't already present in startUrls or sitemaps.
	`))

	return cmd
}

func runCrawlCmd(opts *CrawlOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Requesting crawl for %s on crawler %s", utils.Pluralize(len(opts.URLs), "URL"), opts.ID))
	_, err = client.CrawlURLs(opts.ID, opts.URLs, opts.Save, opts.SaveSpecified)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return fmt.Errorf("%s Crawler API error: %w", cs.FailureIcon(), err)
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Successfully requested crawl for %s on crawler %s\n", cs.SuccessIconWithColor(cs.Green), utils.Pluralize(len(opts.URLs), "URL"), opts.ID)
	}

	return nil
}
