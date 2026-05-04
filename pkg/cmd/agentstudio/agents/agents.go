package agents

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/complete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/create"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/delete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/duplicate"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/get"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/list"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/publish"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/unpublish"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAgentsCmd returns the parent command for agent operations.
func NewAgentsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agents",
		Aliases: []string{"agent"},
		Short:   "Manage Agent Studio agents",
		Long: heredoc.Doc(`
			Create, list, update, run, and manage Agent Studio agents.
		`),
	}

	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(get.NewGetCmd(f, nil))
	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(update.NewUpdateCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(duplicate.NewDuplicateCmd(f, nil))
	cmd.AddCommand(publish.NewPublishCmd(f, nil))
	cmd.AddCommand(unpublish.NewUnpublishCmd(f, nil))
	cmd.AddCommand(complete.NewCompleteCmd(f, nil))

	return cmd
}
