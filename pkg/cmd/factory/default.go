package factory

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

func New(appVersion string, cfg config.IConfig) *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config:         cfg,
		ExecutableName: "gh",
	}
	f.IOStreams = ioStreams(f)
	f.SearchClient = searchClient(f, appVersion)
	f.CrawlerClient = crawlerClient(f)

	return f
}

func ioStreams(_ *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()
	return io
}

func searchClient(f *cmdutil.Factory, appVersion string) func() (*search.APIClient, error) {
	return func() (*search.APIClient, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}
		defaultClient, err := search.NewClient(appID, apiKey)
		if err != nil {
			return nil, err
		}
		defaultUserAgent := defaultClient.GetConfiguration().Configuration.UserAgent

		clientCfg := search.SearchConfiguration{
			Configuration: transport.Configuration{
				AppID:     appID,
				ApiKey:    apiKey,
				UserAgent: fmt.Sprintf("Algolia CLI (%s); %s", appVersion, defaultUserAgent),
			},
		}
		return search.NewClientWithConfig(clientCfg)
	}
}

func crawlerClient(f *cmdutil.Factory) func() (*crawler.Client, error) {
	return func() (*crawler.Client, error) {
		userID, err := f.Config.Profile().GetCrawlerUserID()
		if err != nil {
			return nil, err
		}
		APIKey, err := f.Config.Profile().GetCrawlerAPIKey()
		if err != nil {
			return nil, err
		}

		return crawler.NewClient(userID, APIKey), nil
	}
}
