package members

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/applications/members/add"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewMembersCmd returns a new command for Applications Members.
func NewMembersCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "members",
		Aliases: []string{"member", "users", "user"},
		Short:   "Manage your Algolia Applications Members",
	}

	cmd.AddCommand(add.NewAddCmd(f, nil))

	return cmd
}
