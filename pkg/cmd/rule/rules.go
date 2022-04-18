package rule

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/rule/browse"
	importRules "github.com/algolia/cli/pkg/cmd/rule/import"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewRuleCmd returns a new command for rules.
func NewRuleCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Manage your Algolia rules",
	}

	cmd.AddCommand(importRules.NewImportCmd(f))
	cmd.AddCommand(browse.NewBrowseCmd(f))

	return cmd
}
