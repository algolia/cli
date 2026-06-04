package profile

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/profile/add"
	"github.com/algolia/cli/pkg/cmd/profile/list"
	"github.com/algolia/cli/pkg/cmd/profile/remove"
	"github.com/algolia/cli/pkg/cmd/profile/setdefault"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewProfileCmd returns a new command for managing profiles.
func NewProfileCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profile",
		Aliases: []string{"profiles"},
		Short:   "(deprecated) Manage your Algolia CLI profiles",
		Long: heredoc.Doc(`
			Manage your Algolia CLI profiles.

			These commands are deprecated. Credentials now live in state.toml
			(non-secrets) and the OS keychain (secrets), managed by:

			  - algolia auth login          sign in and configure an application
			  - algolia application list    list applications, marking configured ones
			  - algolia application select  switch the active application

			Existing profiles keep working and remain resolvable as aliases via
			the deprecated --profile flag until the next major version.
		`),
		Deprecated: "use `algolia auth login`, `algolia application list`, and `algolia application select` instead. Profiles still resolve as aliases until the next major version.",
	}

	auth.DisableAuthCheck(cmd)

	cmd.AddCommand(add.NewAddCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(remove.NewRemoveCmd(f, nil))
	cmd.AddCommand(setdefault.NewSetDefaultCmd(f, nil))

	return cmd
}
