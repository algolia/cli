package agents

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/cache"
	"github.com/algolia/cli/pkg/cmd/agents/config"
	"github.com/algolia/cli/pkg/cmd/agents/conversations"
	"github.com/algolia/cli/pkg/cmd/agents/create"
	deletecmd "github.com/algolia/cli/pkg/cmd/agents/delete"
	"github.com/algolia/cli/pkg/cmd/agents/domains"
	"github.com/algolia/cli/pkg/cmd/agents/duplicate"
	"github.com/algolia/cli/pkg/cmd/agents/feedback"
	"github.com/algolia/cli/pkg/cmd/agents/get"
	internalcmd "github.com/algolia/cli/pkg/cmd/agents/internal"
	"github.com/algolia/cli/pkg/cmd/agents/keys"
	"github.com/algolia/cli/pkg/cmd/agents/list"
	"github.com/algolia/cli/pkg/cmd/agents/providers"
	"github.com/algolia/cli/pkg/cmd/agents/publish"
	"github.com/algolia/cli/pkg/cmd/agents/run"
	trycmd "github.com/algolia/cli/pkg/cmd/agents/try"
	"github.com/algolia/cli/pkg/cmd/agents/unpublish"
	"github.com/algolia/cli/pkg/cmd/agents/update"
	"github.com/algolia/cli/pkg/cmd/agents/userdata"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAgentsCmd returns the `algolia agents` command group for Algolia
// Agent Studio (github.com/algolia/conversational-ai). See docs/agents.md.
func NewAgentsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Manage Algolia Agent Studio agents",
		Long: heredoc.Doc(`
			Manage agents on Algolia Agent Studio. Credentials come from
			the active profile. Host resolution: --agent-studio-url /
			ALGOLIA_AGENT_STUDIO_URL → profile region → ldflag default →
			cluster-proxy fallback.
		`),
		Example: heredoc.Doc(`
			# List your agents
			$ algolia agents list

			# Get one agent as JSON
			$ algolia agents get 11111111-1111-1111-1111-111111111111
		`),
	}

	cmd.PersistentPreRunE = betaAgentsPreRunE(f)

	cmd.AddCommand(list.NewListCmd(f, nil))
	cmd.AddCommand(get.NewGetCmd(f, nil))
	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(update.NewUpdateCmd(f, nil))
	cmd.AddCommand(deletecmd.NewDeleteCmd(f, nil))
	cmd.AddCommand(publish.NewPublishCmd(f, nil))
	cmd.AddCommand(unpublish.NewUnpublishCmd(f, nil))
	cmd.AddCommand(duplicate.NewDuplicateCmd(f, nil))
	cmd.AddCommand(trycmd.NewTryCmd(f, nil))
	cmd.AddCommand(run.NewRunCmd(f, nil))
	cmd.AddCommand(cache.NewCacheCmd(f))
	cmd.AddCommand(providers.NewProvidersCmd(f))
	cmd.AddCommand(config.NewConfigCmd(f))
	cmd.AddCommand(conversations.NewConversationsCmd(f))
	cmd.AddCommand(domains.NewDomainsCmd(f))
	cmd.AddCommand(keys.NewKeysCmd(f))
	cmd.AddCommand(feedback.NewFeedbackCmd(f))
	cmd.AddCommand(userdata.NewUserDataCmd(f))
	cmd.AddCommand(internalcmd.NewInternalCmd(f))

	h := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		_ = betaAgentsPreRunE(f)(nil, nil)
		if h != nil {
			h(c, args)
		}
	})

	return cmd
}
