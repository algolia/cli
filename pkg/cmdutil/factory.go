package cmdutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	v4 "github.com/algolia/algoliasearch-client-go/v4/algolia/search"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type Factory struct {
	IOStreams       *iostreams.IOStreams
	Config          config.IConfig
	SearchClient    func() (*search.Client, error)
	V4_SearchClient func() (*v4.APIClient, error)
	CrawlerClient   func() (*crawler.Client, error)

	ExecutableName string
}

// Executable is the path to the currently invoked binary
func (f *Factory) Executable() string {
	if !strings.ContainsRune(f.ExecutableName, os.PathSeparator) {
		f.ExecutableName = executable(f.ExecutableName)
	}
	return f.ExecutableName
}

// based on https://github.com/cli/cli/blob/master/pkg/cmdutil/factory.go
func executable(fallbackName string) string {
	exe, err := os.Executable()
	if err != nil {
		return fallbackName
	}

	base := filepath.Base(exe)
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		p, err := filepath.Abs(filepath.Join(dir, base))
		if err != nil {
			continue
		}
		f, err := os.Lstat(p)
		if err != nil {
			continue
		}

		if p == exe {
			return p
		} else if f.Mode()&os.ModeSymlink != 0 {
			if t, err := os.Readlink(p); err == nil && t == exe {
				return p
			}
		}
	}

	return exe
}
