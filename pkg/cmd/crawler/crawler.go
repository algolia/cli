package crawler

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/crawler/crawl"
	"github.com/algolia/cli/pkg/cmd/crawler/create"
	"github.com/algolia/cli/pkg/cmd/crawler/get"
	"github.com/algolia/cli/pkg/cmd/crawler/list"
	"github.com/algolia/cli/pkg/cmd/crawler/pause"
	"github.com/algolia/cli/pkg/cmd/crawler/reindex"
	"github.com/algolia/cli/pkg/cmd/crawler/run"
	"github.com/algolia/cli/pkg/cmd/crawler/stats"
	"github.com/algolia/cli/pkg/cmd/crawler/test"
	"github.com/algolia/cli/pkg/cmd/crawler/unblock"
	"github.com/algolia/cli/pkg/cmdutil"
)

const (
	AuthMethodHelpMsg = `In order to use the 'crawler' commands, you will need to authenticate with the Algolia Crawler API. You can do so by either:
  - Export your Algolia Crawler username and API Key as ALGOLIA_CRAWLER_USER_ID and ALGOLIA_CRAWLER_API_KEY environment variables.
  - Add your Algolia Crawler 'crawler_user_id' and 'crawler_api_key' credentials to your profile file (~/.config/algolia/config.tml).`
)

// NewCrawlersCmd returns a new command to manage your Algolia Crawlers.
func NewCrawlersCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crawler",
		Aliases: []string{"crawlers"},
		Short:   "Manage your Algolia crawlers",
		Long: heredoc.Docf(`
			Manage your Algolia crawlers.

			%s
		`, AuthMethodHelpMsg),
		// Check Crawler specific Authentication
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := f.CrawlerClient()
			if err != nil {
				authError := errors.New("authError")
				stderr := f.IOStreams.ErrOut
				fmt.Fprintf(stderr, "Crawler authentication error: %s\n", err)
				fmt.Fprintln(stderr, "")
				fmt.Fprintln(stderr, AuthMethodHelpMsg)
				return authError
			}
			return nil
		},
	}

	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(reindex.NewReindexCmd(f, nil))
	cmd.AddCommand(stats.NewStatsCmd(f, nil))
	cmd.AddCommand(unblock.NewUnblockCmd(f, nil))
	cmd.AddCommand(run.NewRunCmd(f, nil))
	cmd.AddCommand(pause.NewPauseCmd(f, nil))
	cmd.AddCommand(crawl.NewCrawlCmd(f, nil))
	cmd.AddCommand(test.NewTestCmd(f, nil))
	cmd.AddCommand(get.NewGetCmd(f, nil))
	cmd.AddCommand(create.NewCreateCmd(f, nil))

	return cmd
}
