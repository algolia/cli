package create

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

type CreateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	Name   string
	config crawler.Config
}

// NewCreateCmd creates and returns a create command for Crawlers.
func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
	}

	var configFile string

	cmd := &cobra.Command{
		Use:   "create <name> -F <file>",
		Args:  cobra.ExactArgs(1),
		Short: "Create a crawler",
		Long: heredoc.Doc(`
			Create a new crawler from the given configuration.
		`),
		Example: heredoc.Doc(`
			# Create a crawler named "my-crawler" with the configuration in the file "config.json"
			$ algolia crawler create my-crawler -F config.json

			# Create a crawler from another crawler's configuration
			$ algolia crawler get another-crawler --config-only | algolia crawler create my-crawler -F -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			b, err := cmdutil.ReadFile(configFile, opts.IO.In)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, &opts.config)
			if err != nil {
				return err
			}

			if runF != nil {
				return runF(opts)
			}

			return runCreateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&configFile, "file", "F", "", "Path to the configuration file (use \"-\" to read from standard input)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Creating crawler")
	id, err := client.Create(opts.Name, opts.config)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Crawler %s created: %s\n",
			cs.SuccessIconWithColor(cs.Green),
			cs.Bold(opts.Name),
			cs.Bold(id),
		)
	}

	return nil
}
