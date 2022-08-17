package apikeys

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikeys/create"
	"github.com/algolia/cli/pkg/cmd/apikeys/delete"
	"github.com/algolia/cli/pkg/cmd/apikeys/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAPIKeysCmd returns a new command for API Keys.
func NewAPIKeysCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apikeys",
		Aliases: []string{"api-key", "api-keys", "apikey"},
		Short:   "Manage your Algolia API keys",
	}

	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))

	return cmd
}
