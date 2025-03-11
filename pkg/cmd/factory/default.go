package factory

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/api/genai"
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
	f.GenAIClient = genAIClient(f)

	return f
}

func ioStreams(f *cmdutil.Factory) *iostreams.IOStreams {
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

func genAIClient(f *cmdutil.Factory) func() (*genai.Client, error) {
	return func() (*genai.Client, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		APIKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		return genai.NewClient(appID, APIKey), nil
	}
}
