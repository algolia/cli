package cache

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// InvalidateCacheOptions holds the dependencies and flags for the invalidate-cache command.
type InvalidateCacheOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
}

// NewInvalidateCacheCmd returns the `agents invalidate-cache` command.
func NewInvalidateCacheCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &InvalidateCacheOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:               "invalidate-cache <agent-id>",
		Short:             "Invalidate the completion cache for an agent",
		Args:              validators.ExactArgsWithMsg(1, "agents invalidate-cache requires an <agent-id> argument."),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Invalidate the cache for the agent with ID "my-agent"
			$ algolia agents invalidate-cache my-agent
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			return runInvalidateCacheCmd(opts)
		},
	}

	return cmd
}

func runInvalidateCacheCmd(opts *InvalidateCacheOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Invalidating agent cache")

	err = client.InvalidateAgentCache(client.NewApiInvalidateAgentCacheRequest(opts.AgentID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Invalidated cache for agent %s\n", cs.SuccessIcon(), cs.Bold(opts.AgentID))
	}

	return nil
}
