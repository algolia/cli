package cmdutil

import (
	"github.com/algolia/algolia-cli/pkg/config"
	"github.com/algolia/algolia-cli/pkg/iostreams"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type Factory struct {
	IOStreams    *iostreams.IOStreams
	Config       *config.Config
	SearchClient func() (*search.Client, error)
}
