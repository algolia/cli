package crawler

import (
	"fmt"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
	"github.com/spf13/cobra"
)

type CrawlerOptions struct {
	IO                 *iostreams.IOStreams
	config             config.IConfig
	OAuthClientID      func() string
	NewDashboardClient func(clientID string) *dashboard.Client
	GetValidToken      func(client *dashboard.Client) (string, error)
}

func NewCrawlerCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CrawlerOptions{
		IO:            f.IOStreams,
		config:        f.Config,
		OAuthClientID: auth.OAuthClientID,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		GetValidToken: auth.GetValidToken,
	}

	cmd := &cobra.Command{
		Use:   "crawler",
		Short: "Configure the crawler API key for the current application",
		Args:  validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCrawlerCmd(opts)
		},
	}

	return cmd
}

func runCrawlerCmd(opts *CrawlerOptions) error {
	cs := opts.IO.ColorScheme()

	appID := opts.config.ActiveApplicationID()
	if appID == "" {
		return fmt.Errorf(
			"no application configured: run `algolia auth login` or `algolia application select` first",
		)
	}

	dashboardClient := opts.NewDashboardClient(opts.OAuthClientID())

	accessToken, err := opts.GetValidToken(dashboardClient)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching crawler information")
	crawlerUserData, err := dashboardClient.GetCrawlerUser(accessToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if err := opts.config.SetCrawlerAPIKey(appID, crawlerUserData.APIKey); err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Crawler API key configured for application: %s\n",
			cs.SuccessIcon(), appID)
	}

	return nil
}
