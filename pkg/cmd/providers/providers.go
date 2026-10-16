package providers

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/providers/create"
	"github.com/algolia/cli/pkg/cmd/providers/delete"
	"github.com/algolia/cli/pkg/cmd/providers/get"
	"github.com/algolia/cli/pkg/cmd/providers/list"
	"github.com/algolia/cli/pkg/cmd/providers/models"
	"github.com/algolia/cli/pkg/cmd/providers/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewProvidersCmd returns the providers command group.
func NewProvidersCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "Manage Agent Studio LLM providers",
		Long:  "List, get, create, update, delete, and inspect models for Agent Studio providers.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(create.NewCreateCmd(f))
	cmd.AddCommand(update.NewUpdateCmd(f))
	cmd.AddCommand(models.NewModelsCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	return cmd
}
