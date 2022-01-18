package cmdutil

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type Factory struct {
	IOStreams    *iostreams.IOStreams
	Config       *config.Config
	SearchClient func() (search.ClientInterface, error)
}
