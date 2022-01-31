package rule

import (
	"github.com/spf13/cobra"

	importRules "github.com/algolia/cli/pkg/cmd/rule/import"
	"github.com/algolia/cli/pkg/cmd/rule/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewRuleCmd returns a new command for rules.
func NewRuleCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Manage your Algolia rules",
	}

	cmd.AddCommand(importRules.NewImportCmd(f))
	cmd.AddCommand(list.NewListCmd(f))

	return cmd
}
