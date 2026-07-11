package conversations

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/conversations/delete"
	"github.com/algolia/cli/pkg/cmd/agents/conversations/export"
	"github.com/algolia/cli/pkg/cmd/agents/conversations/get"
	"github.com/algolia/cli/pkg/cmd/agents/conversations/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewConversationsCmd returns the `agents conversations` command group.
func NewConversationsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversations",
		Short: "Manage the conversations of an agent",
		Long:  "List, retrieve, delete, and export the conversations of an Algolia Agent Studio agent.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(export.NewExportCmd(f))

	return cmd
}
