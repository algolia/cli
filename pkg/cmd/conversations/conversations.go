package conversations

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/conversations/delete"
	"github.com/algolia/cli/pkg/cmd/conversations/export"
	"github.com/algolia/cli/pkg/cmd/conversations/get"
	"github.com/algolia/cli/pkg/cmd/conversations/list"
	"github.com/algolia/cli/pkg/cmd/conversations/purge"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewConversationsCmd returns the conversations command group.
func NewConversationsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversations",
		Short: "Manage Agent Studio conversations",
		Long:  "List, get, export, delete, and purge Agent Studio conversations.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(export.NewExportCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(purge.NewPurgeCmd(f))
	return cmd
}
