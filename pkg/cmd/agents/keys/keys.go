package keys

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewKeysCmd is the parent for `algolia agents keys <verb>`.
// All mutating verbs require an admin API key on the active profile;
// the backend rejects non-admin keys with 403 "Admin API key required."
func NewKeysCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage Agent Studio secret keys",
		Long: heredoc.Doc(`
			Manage Agent Studio secret keys (vended per app, optionally
			scoped to specific agents). Mutating verbs require an admin
			API key.
		`),
	}
	cmd.AddCommand(newListCmd(f, nil))
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newCreateCmd(f, nil))
	cmd.AddCommand(newUpdateCmd(f, nil))
	cmd.AddCommand(newDeleteCmd(f, nil))
	return cmd
}

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// maskKey returns a copy of k with Value redacted unless show is set.
func maskKey(k agentstudio.SecretKey, show bool) agentstudio.SecretKey {
	if show {
		return k
	}
	k.Value = shared.MaskString(k.Value)
	return k
}
