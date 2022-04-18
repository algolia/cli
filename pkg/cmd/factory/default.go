package factory

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

func New(cfg *config.Config) *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config: cfg,
	}
	f.IOStreams = ioStreams(f)
	f.SearchClient = searchClient(f)

	return f
}

func ioStreams(f *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()
	return io
}

func searchClient(f *cmdutil.Factory) func() (*search.Client, error) {
	return func() (*search.Client, error) {
		appID, err := f.Config.Application.GetID()
		if err != nil {
			return nil, err
		}
		APIKey, err := f.Config.Application.GetAdminAPIKey()
		if err != nil {
			return nil, err
		}

		return search.NewClient(appID, APIKey), nil
	}
}
