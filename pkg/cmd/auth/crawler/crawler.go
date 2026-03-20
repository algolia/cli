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
		Short: "Load crawler auth details for the current profile",
		Args:  validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCrawlerCmd(opts)
		},
	}

	return cmd
}

func runCrawlerCmd(opts *CrawlerOptions) error {
	cs := opts.IO.ColorScheme()
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

	currentProfileName := opts.config.Profile().Name
	if currentProfileName == "" {
		defaultProfile := opts.config.Default()
		if defaultProfile != nil {
			currentProfileName = defaultProfile.Name
			opts.config.Profile().Name = currentProfileName
		}
	}
	if currentProfileName == "" {
		return fmt.Errorf("no profile selected and no default profile configured")
	}

	if err = opts.config.SetCrawlerAuth(currentProfileName, crawlerUserData.ID, crawlerUserData.APIKey); err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Crawler API auth credentials configured for profile: %s\n", cs.SuccessIcon(), currentProfileName)
	}

	return nil
}
