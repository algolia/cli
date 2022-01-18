package factory

import (
	"os"

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

	if prompt := f.Config.Profile.GetConfigField("prompt"); prompt == "disabled" {
		io.SetNeverPrompt(true)
	}

	// Pager precedence
	// 1. ALGOLIA_PAGER
	// 2. pager from config
	// 3. PAGER
	if algoliaPager, algoliaPagerExists := os.LookupEnv("ALGOLIA_PAGER"); algoliaPagerExists {
		io.SetPager(algoliaPager)
	} else if pager := f.Config.Profile.GetFieldValue("pager"); pager != "" {
		io.SetPager(pager)
	}

	return io
}

func searchClient(f *cmdutil.Factory) func() (search.ClientInterface, error) {
	return func() (search.ClientInterface, error) {
		APIKey, err := f.Config.Profile.GetAdminAPIKey()
		if err != nil {
			return nil, err
		}
		applicationID, err := f.Config.Profile.GetApplicationID()
		if err != nil {
			return nil, err
		}
		return search.NewClient(applicationID, APIKey), nil
	}
}
