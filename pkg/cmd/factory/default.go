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
	f.V4SearchClient = v4searchClient(f, appVersion)
	f.CrawlerClient = crawlerClient(f)

	return f
}

func ioStreams(_ *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()
	return io
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

func v4searchClient(f *cmdutil.Factory, appVersion string) func() (*v4.APIClient, error) {
	return func() (*v4.APIClient, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		defaultClient, _ := v4.NewClient(appID, apiKey)
		defaultUserAgent := defaultClient.GetConfiguration().UserAgent

		// TODO: Doesn't support custom `search_hosts` yet.
		// To support it, it's best to transform the GetSearchHosts() function
		clientConf := v4.SearchConfiguration{
			Configuration: transport.Configuration{
				AppID:     appID,
				ApiKey:    apiKey,
				UserAgent: defaultUserAgent + fmt.Sprintf("Algolia CLI (%s)", appVersion),
			},
		}

		return v4.NewClientWithConfig(clientConf)
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
