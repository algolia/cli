package providers

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/create"
	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/delete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/get"
	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/list"
	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/models"
	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewProvidersCmd returns the providers command group.
func NewProvidersCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "Manage your Algolia Agent Studio LLM providers",
		Long:  "Create, retrieve, update, and delete the LLM provider authentications used by Algolia Agent Studio agents.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(create.NewCreateCmd(f))
	cmd.AddCommand(update.NewUpdateCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(models.NewModelsCmd(f))

	return cmd
}
