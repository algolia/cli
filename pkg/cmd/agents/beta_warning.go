package agents

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/version"
)

// betaAgentsPreRunE emits a stderr banner for non-release binaries (link-time
// version.Distribution set, e.g. "beta") whenever any `agents` subtree command runs.
func betaAgentsPreRunE(f *cmdutil.Factory) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		if version.Distribution == "" {
			return nil
		}
		fmt.Fprintf(
			f.IOStreams.ErrOut,
			"warning: %s CLI build — Algolia recommends the release `algolia` binary for "+
				"production. `agents` defaults can follow your `.env` / build-time flags.\n\n",
			version.Distribution,
		)
		return nil
	}
}
