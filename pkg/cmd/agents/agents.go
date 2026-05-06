package agents

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/cache"
	"github.com/algolia/cli/pkg/cmd/agents/config"
	"github.com/algolia/cli/pkg/cmd/agents/create"
	deletecmd "github.com/algolia/cli/pkg/cmd/agents/delete"
	"github.com/algolia/cli/pkg/cmd/agents/duplicate"
	"github.com/algolia/cli/pkg/cmd/agents/get"
	"github.com/algolia/cli/pkg/cmd/agents/list"
	"github.com/algolia/cli/pkg/cmd/agents/providers"
	"github.com/algolia/cli/pkg/cmd/agents/publish"
	"github.com/algolia/cli/pkg/cmd/agents/run"
	trycmd "github.com/algolia/cli/pkg/cmd/agents/try"
	"github.com/algolia/cli/pkg/cmd/agents/unpublish"
	"github.com/algolia/cli/pkg/cmd/agents/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewAgentsCmd returns the `algolia agents` command group, which manages
// agents on Algolia Agent Studio (https://github.com/algolia/conversational-ai).
//
// Authentication uses the same Application ID + API Key as the rest of the
// CLI; the Agent Studio host is resolved from --agent-studio-url (or the
// equivalent profile field), then the configured region, then a cluster-proxy
// fallback derived from the Application ID.
func NewAgentsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Manage Algolia Agent Studio agents",
		Long: heredoc.Doc(`
			Manage agents on Algolia Agent Studio.

			Reads credentials from your active profile. The Agent Studio host is
			resolved in this order:

			  1. ALGOLIA_AGENT_STUDIO_URL env var, or the profile's
			     "agent_studio_url" field.
			  2. The build-time default baked into the binary (ldflag-driven;
			     used by internal beta builds to point at
			     agent-studio.staging.eu.algolia.com by default).
			  3. The cluster-proxy fallback
			     https://<app-id>.algolia.net/agent-studio (recommended for
			     production end-users — the application's own cluster routes
			     the request to the right region).
		`),
		Example: heredoc.Doc(`
			# List your agents
			$ algolia agents list

			# Get one agent as JSON
			$ algolia agents get 11111111-1111-1111-1111-111111111111
		`),
	}

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

	return cmd
}
