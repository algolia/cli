package rules

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/compositions/rules/delete"
	"github.com/algolia/cli/pkg/cmd/compositions/rules/get"
	"github.com/algolia/cli/pkg/cmd/compositions/rules/list"
	"github.com/algolia/cli/pkg/cmd/compositions/rules/upsert"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewRulesCmd returns the compositions rules sub-group.
func NewRulesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage composition rules",
		Long:  "List, get, upsert, and delete rules for an Algolia Composition.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(upsert.NewUpsertCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	return cmd
}
