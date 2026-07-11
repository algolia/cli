package domains

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/add"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/bulkadd"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/bulkdelete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/delete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/get"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewDomainsCmd returns the `agents domains` command group.
func NewDomainsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: "Manage the allowed domains of an agent",
		Long:  "Manage the domains allowed to embed or call an Algolia Agent Studio agent.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(add.NewAddCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(bulkadd.NewBulkAddCmd(f))
	cmd.AddCommand(bulkdelete.NewBulkDeleteCmd(f))

	return cmd
}
