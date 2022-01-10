package apikeys

import (
	"github.com/algolia/algolia-cli/pkg/cmd/apikeys/create"
	"github.com/algolia/algolia-cli/pkg/cmd/apikeys/delete"
	"github.com/algolia/algolia-cli/pkg/cmd/apikeys/list"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewAPIKeysCmd returns a new command for API Keys.
func NewAPIKeysCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apikeys",
		Short: "Manage your Algolia API keys",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(create.NewCreateCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))

	return cmd
}
