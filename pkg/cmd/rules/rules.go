package rules

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/rules/browse"
	"github.com/algolia/cli/pkg/cmd/rules/delete"
	importRules "github.com/algolia/cli/pkg/cmd/rules/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewRulesCmd returns a new command for rules.
func NewRulesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rules",
		Aliases: []string{"rule"},
		Short:   "Manage your Algolia rules",
	}

	cmd.AddCommand(importRules.NewImportCmd(f, nil))
	cmd.AddCommand(browse.NewBrowseCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))

	return cmd
}
