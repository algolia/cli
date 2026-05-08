package agents

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/version"
)

const betaWarningLine = "[BETA] WARNING: This version should not be used in production."

// betaAgentsPreRunE emits a stderr banner for non-release binaries (link-time
// version.Distribution set, e.g. "beta") whenever any `agents` subtree command runs.
func betaAgentsPreRunE(f *cmdutil.Factory) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		if version.Distribution == "" {
			return nil
		}

		w := f.IOStreams.ErrOut
		line := betaWarningLine
		if f.IOStreams.ColorEnabled() {
			cs := f.IOStreams.ColorScheme()
			line = cs.Bold(cs.Yellow(betaWarningLine))
		}
		fmt.Fprintf(w, "%s\n\n", line)
		return nil
	}
}
