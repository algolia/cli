package compositions

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/compositions/delete"
	"github.com/algolia/cli/pkg/cmd/compositions/get"
	"github.com/algolia/cli/pkg/cmd/compositions/list"
	"github.com/algolia/cli/pkg/cmd/compositions/rules"
	compsearch "github.com/algolia/cli/pkg/cmd/compositions/search"
	"github.com/algolia/cli/pkg/cmd/compositions/upsert"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewCompositionsCmd returns the compositions command group.
func NewCompositionsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compositions",
		Short: "Manage Algolia Compositions",
		Long:  "Create, retrieve, update, delete, and search Algolia Compositions.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(upsert.NewUpsertCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(compsearch.NewSearchCmd(f))
	cmd.AddCommand(rules.NewRulesCmd(f))
	return cmd
}
