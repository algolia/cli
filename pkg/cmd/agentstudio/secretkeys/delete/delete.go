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
	SecretKeyID       string
	DoConfirm         bool
}

// NewDeleteCmd returns the `secret-keys delete` command.
func NewDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:               "delete <secret-key-id>",
		Short:             "Delete a secret key",
		Args:              validators.ExactArgsWithMsg(1, "secret-keys delete requires a <secret-key-id> argument."),
		ValidArgsFunction: cmdutil.SecretKeyIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Delete a secret key (with confirmation prompt)
			$ algolia secret-keys delete my-key

			# Delete without confirmation
			$ algolia secret-keys delete my-key --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SecretKeyID = args[0]

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete secret key %q?", opts.SecretKeyID),
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

	opts.IO.StartProgressIndicatorWithLabel("Deleting secret key")

	err = client.DeleteSecretKey(client.NewApiDeleteSecretKeyRequest(opts.SecretKeyID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted secret key %s\n", cs.SuccessIcon(), opts.SecretKeyID)
	}

	return nil
}
