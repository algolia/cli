package secretkeys

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/secretkeys/create"
	"github.com/algolia/cli/pkg/cmd/secretkeys/delete"
	"github.com/algolia/cli/pkg/cmd/secretkeys/get"
	"github.com/algolia/cli/pkg/cmd/secretkeys/list"
	"github.com/algolia/cli/pkg/cmd/secretkeys/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewSecretKeysCmd returns the secret-keys command group.
func NewSecretKeysCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret-keys",
		Short: "Manage your Algolia Agent Studio secret keys",
		Long:  "Create, retrieve, update, and delete the secret keys used to secure Algolia Agent Studio agents.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(create.NewCreateCmd(f))
	cmd.AddCommand(update.NewUpdateCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))

	return cmd
}
