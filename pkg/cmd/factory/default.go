package factory

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/call"
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

		defaultClient, _ := search.NewClient(appID, apiKey)
		defaultUserAgent := defaultClient.GetConfiguration().UserAgent

		var hosts []transport.StatefulHost
		for _, host := range f.Config.Profile().GetSearchHosts() {
			statefulHost := transport.NewStatefulHost("https", host, call.IsReadWrite)
			hosts = append(hosts, statefulHost)
		}

		clientConf := search.SearchConfiguration{
			Configuration: transport.Configuration{
				AppID:  appID,
				ApiKey: apiKey,
				UserAgent: defaultUserAgent + fmt.Sprintf(
					"Algolia CLI (%s)",
					appVersion,
				),
				Hosts:                           hosts,
				ExposeIntermediateNetworkErrors: true,
			},
		}

		return search.NewClientWithConfig(clientConf)
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
