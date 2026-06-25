package cmdutil

import (
	"os"
	"path/filepath"
	"strings"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type Factory struct {
	IOStreams         *iostreams.IOStreams
	Config            config.IConfig
	SearchClient      func() (*search.APIClient, error)
	CrawlerClient     func() (*crawler.Client, error)
	CompositionClient func() (*composition.APIClient, error)
	// AgentStudioClient is the hand-rolled client kept for endpoints the
	// official SDK can't serve (SSE streaming for run/try, the x-hidden
	// internal endpoints, and agent duplication).
	AgentStudioClient func() (*agentstudio.Client, error)
	// AgentStudioAPIClient is the official SDK client used for the standard
	// CRUD surface. Routes to https://<appID>.algolia.net/agent-studio/1/...
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)

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
