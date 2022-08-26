package factory

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

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
		APIKey, err := f.Config.Profile().GetAdminAPIKey()
		if err != nil {
			return nil, err
		}

		clientCfg := search.Configuration{
			AppID:          appID,
			APIKey:         APIKey,
			ExtraUserAgent: fmt.Sprintf("Algolia CLI (%s)", appVersion),
		}
		return search.NewClientWithConfig(clientCfg), nil
	}
}
