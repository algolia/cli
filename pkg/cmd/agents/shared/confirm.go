package shared

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// AddConfirmFlag registers the standard `-y/--confirm` boolean on cmd
// and binds it to dst. Pair with ResolveConfirm in RunE and Confirm in
// the runner.
func AddConfirmFlag(cmd *cobra.Command, dst *bool) {
	cmd.Flags().BoolVarP(dst, "confirm", "y", false, "Skip confirmation prompt")
}

// ResolveConfirm enforces the "non-interactive shells require -y" rule
// and reports whether the runner should prompt.
func ResolveConfirm(ios *iostreams.IOStreams, confirm bool) (bool, error) {
	if confirm {
		return false, nil
	}
	if !ios.CanPrompt() {
		return false, cmdutil.FlagErrorf(
			"--confirm required when non-interactive shell is detected",
		)
	}
	return true, nil
}

// targetLabel describes the application a destructive command will act on. It
// is set once by the agents command group (see SetTargetLabel) so every
// confirmation prompt names the target, which can come from --application-id,
// the environment, --profile or the current application.
var targetLabel func() string

// SetTargetLabel registers the function Confirm uses to name the target
// application in its prompt.
func SetTargetLabel(f func() string) { targetLabel = f }

// Confirm prompts with msg, followed by the target application when known. Returns (true, nil) on yes, (false, nil)
// on no, or (false, err) on prompt failure. Callers use the bool to
// short-circuit silently — matches existing agents-CLI behaviour.
func Confirm(msg string) (bool, error) {
	if targetLabel != nil {
		if t := targetLabel(); t != "" {
			msg = fmt.Sprintf("%s [%s]", msg, t)
		}
	}
	var ok bool
	if err := prompt.Confirm(msg, &ok); err != nil {
		return false, fmt.Errorf("failed to prompt: %w", err)
	}
	return ok, nil
}
