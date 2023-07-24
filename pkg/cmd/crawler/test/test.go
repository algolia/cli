package test

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// TestOptions holds the options for the test command.
type TestOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	ID     string
	URL    string
	config *crawler.Config

	PrintFlags *cmdutil.PrintFlags
}

// NewTestCmd creates and returns a crawl command for crawlers.
func NewTestCmd(f *cmdutil.Factory, runF func(*TestOptions) error) *cobra.Command {
	opts := &TestOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
		PrintFlags:    cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	var configFile string

	cmd := &cobra.Command{
		Use:               "test <crawler_id> --url <url> [-F <file>]",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Test a URL on a crawler",
		Long: heredoc.Doc(`
			Test an URL against the given crawler's configuration and see what will be processed.
			You can also override parts of the configuration to try your changes before updating the configuration.
		`),
		Example: heredoc.Doc(`
			# Test the URL "https://www.algolia.com" against the crawler with the ID "my-crawler"
			$ algolia crawler test my-crawler --url https://www.algolia.com

			# Test the URL "https://www.algolia.com" against the crawler with the ID "my-crawler" and override the configuration with the file "config.json"
			$ algolia crawler test my-crawler --url https://www.algolia.com -F config.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]

			if cmd.Flags().Changed("config") {
				b, err := cmdutil.ReadFile(configFile, opts.IO.In)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b, &opts.config)
				if err != nil {
					return err
				}
			}

			if runF != nil {
				return runF(opts)
			}

			return runTestCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.URL, "url", "u", "", "The URL to test.")
	_ = cmd.MarkFlagRequired("url")

	cmd.Flags().StringVarP(&configFile, "config", "F", "", "The configuration file to use to override the crawler's configuration. (use \"-\" to read from standard input)")

	return cmd
}

func runTestCmd(opts *TestOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel(fmt.Sprintf("Testing URL %s on crawler %s", cs.Bold(opts.URL), cs.Bold(opts.ID)))
	res, err := client.Test(opts.ID, opts.URL, opts.config)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}
	opts.IO.StopProgressIndicator()

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	if err := p.Print(opts.IO, res); err != nil {
		return err
	}

	return nil
}
