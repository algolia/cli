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
	AgentID           string
	DomainID          string
	DoConfirm         bool
}

// NewDeleteCmd returns the `agents domains delete` command.
func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:   "delete <agent-id> <domain-id>",
		Short: "Delete an allowed domain from an agent",
		Args: validators.ExactArgsWithMsg(
			2,
			"agents domains delete requires an <agent-id> and a <domain-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Delete the allowed domain "domain_123" from the agent "my-agent"
			$ algolia agents domains delete my-agent domain_123 --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.DomainID = args[1]

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete allowed domain %q?", opts.DomainID),
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

	opts.IO.StartProgressIndicatorWithLabel("Deleting allowed domain")

	err = client.DeleteAllowedDomain(client.NewApiDeleteAllowedDomainRequest(opts.DomainID, opts.AgentID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted allowed domain %s\n", cs.SuccessIcon(), opts.DomainID)
	}

	return nil
}
