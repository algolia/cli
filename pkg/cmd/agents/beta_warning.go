package agents

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/version"
)

const betaWarningLine = "[BETA] WARNING: This version should not be used in production."

// PrintBetaWarning writes the beta banner to stderr for non-release binaries.
// The root command calls it before running any `agents` subtree command,
// because a PersistentPreRunE on `agents` would replace the root hook (auth,
// ACL checks and telemetry).
func PrintBetaWarning(f *cmdutil.Factory) {
	_ = betaAgentsPreRunE(f)(nil, nil)
}

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
