package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// DeleteOptions holds the dependencies and flags for the delete command.
type DeleteOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	ProviderID        string
	DoConfirm         bool
}

// NewDeleteCmd returns the `providers delete` command.
func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:               "delete <provider-id>",
		Short:             "Delete an LLM provider authentication",
		Args:              validators.ExactArgsWithMsg(1, "providers delete requires a <provider-id> argument."),
		ValidArgsFunction: cmdutil.ProviderIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Delete a provider (with confirmation prompt)
			$ algolia providers delete my-provider

			# Delete without confirmation
			$ algolia providers delete my-provider --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete provider %q?", opts.ProviderID),
					&confirmed,
				)
				if err != nil {
					return fmt.Errorf("failed to prompt: %w", err)
				}
				if !confirmed {
					return nil
				}
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.DoConfirm, "confirm", "y", false, "Skip confirmation prompt")

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting provider")

	err = client.DeleteProvider(client.NewApiDeleteProviderRequest(opts.ProviderID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted provider %s\n", cs.SuccessIcon(), opts.ProviderID)
	}

	return nil
}
