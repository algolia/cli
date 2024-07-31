package factory

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	v4 "github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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
	f.V4_SearchClient = v4_searchClient(f, appVersion)
	f.CrawlerClient = crawlerClient(f)

	return f
}

func ioStreams(_ *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()
	return io
}

func v4_searchClient(f *cmdutil.Factory, appVersion string) func() (*v4.APIClient, error) {
	return func() (*v4.APIClient, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}
		defaultClient, err := v4.NewClient(appID, apiKey)
		if err != nil {
			return nil, err
		}
		defaultUserAgent := defaultClient.GetConfiguration().Configuration.UserAgent

		clientCfg := v4.SearchConfiguration{
			Configuration: transport.Configuration{
				AppID:     appID,
				ApiKey:    apiKey,
				UserAgent: fmt.Sprintf("Algolia CLI (%s); %s", appVersion, defaultUserAgent),
			},
		}
		return v4.NewClientWithConfig(clientCfg)
	}
}

func searchClient(f *cmdutil.Factory, appVersion string) func() (*search.Client, error) {
	return func() (*search.Client, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		APIKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		clientCfg := search.Configuration{
			AppID:          appID,
			APIKey:         APIKey,
			ExtraUserAgent: fmt.Sprintf("Algolia CLI (%s)", appVersion),
			Hosts:          f.Config.Profile().GetSearchHosts(),
		}
		return search.NewClientWithConfig(clientCfg), nil
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
