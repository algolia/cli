package agents

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/cache"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/completions"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/conversations"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/create"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/delete"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/domains"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/get"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/list"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/publish"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/unpublish"
	"github.com/algolia/cli/pkg/cmd/agentstudio/agents/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAgentsCmd returns the agents command group.
func NewAgentsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Manage your Algolia Agent Studio agents",
		Long:  "Create, retrieve, update, delete, publish, and converse with Algolia Agent Studio agents.",
	}

	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f))
	cmd.AddCommand(create.NewCreateCmd(f))
	cmd.AddCommand(update.NewUpdateCmd(f))
	cmd.AddCommand(delete.NewDeleteCmd(f))
	cmd.AddCommand(publish.NewPublishCmd(f))
	cmd.AddCommand(unpublish.NewUnpublishCmd(f))
	cmd.AddCommand(cache.NewInvalidateCacheCmd(f))
	cmd.AddCommand(completions.NewCompletionsCmd(f))
	cmd.AddCommand(domains.NewDomainsCmd(f))
	cmd.AddCommand(conversations.NewConversationsCmd(f))

	return cmd
}
