package apikeys

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/applications/apikeys/create"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAPIKeysCmd returns a new command for Application API Keys.
func NewAPIKeysCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apikeys",
		Aliases: []string{"apikey", "key", "keys"},
		Short:   "Manage your Algolia Applications API Keys",
	}

	cmd.AddCommand(create.NewCreateCmd(f, nil))

	return cmd
}
