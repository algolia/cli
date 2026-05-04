package conversations

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/conversations/delete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/conversations/export"
	"github.com/algolia/cli/pkg/cmd/agentstudio/conversations/get"
	"github.com/algolia/cli/pkg/cmd/agentstudio/conversations/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

func NewConversationsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "conversations",
		Aliases: []string{"conv", "conversation"},
		Short:   "Inspect Agent Studio conversations",
		Long: heredoc.Doc(`
			List, fetch, delete, and export conversations recorded for an agent.
		`),
	}

	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(get.NewGetCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(export.NewExportCmd(f, nil))

	return cmd
}
