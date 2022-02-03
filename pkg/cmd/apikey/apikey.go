package apikey

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikey/create"
	"github.com/algolia/cli/pkg/cmd/apikey/delete"
	"github.com/algolia/cli/pkg/cmd/apikey/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAPIKeyCmd returns a new command for API Keys.
func NewAPIKeyCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apikey",
		Aliases: []string{"api-key", "api-keys", "apikeys"},
		Short:   "Manage your Algolia API keys",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))

	return cmd
}
