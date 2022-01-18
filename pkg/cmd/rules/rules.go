package rules

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/rules/export"
	importRules "github.com/algolia/cli/pkg/cmd/rules/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewRulesCmd returns a new command for rules.
func NewRulesCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage your Algolia rules",
	}

	cmd.AddCommand(importRules.NewImportCmd(f))
	cmd.AddCommand(export.NewExportCmd(f))

	return cmd
}
