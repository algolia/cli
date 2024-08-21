package applications

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/applications/apikeys"
	"github.com/algolia/cli/pkg/cmd/applications/create"
	"github.com/algolia/cli/pkg/cmd/applications/list"
	"github.com/algolia/cli/pkg/cmd/applications/members"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewApplicationsCmd returns a new command for Applications.
func NewApplicationsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "applications",
		Aliases: []string{"application", "apps", "app"},
		Short:   "Manage your Algolia Applications",
	}

	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(members.NewMembersCmd(f))
	cmd.AddCommand(apikeys.NewAPIKeysCmd(f))

	return cmd
}
