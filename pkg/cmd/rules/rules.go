package rules

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/rules/browse"
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

	return cmd
}
