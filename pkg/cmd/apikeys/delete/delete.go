package delete

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// DeleteOptions represents the options for the create command
type DeleteOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	APIKey    string
	DoConfirm bool
}

// NewDeleteCmd returns a new instance of DeleteCmd
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:           f.IOStreams,
		config:       f.Config,
		SearchClient: f.SearchClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <api-key>",
		Short: "Deletes the API key",
		Args:  validators.ExactArgs(1),
		Annotations: map[string]string{
			"acls": "admin",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.APIKey = args[0]
			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip the delete API key confirmation prompt")

	return cmd
}

// runDeleteCmd runs the delete command
func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	_, err = client.GetAPIKey(opts.APIKey)
	if err != nil {
		return fmt.Errorf("API key %q does not exist", opts.APIKey)
	}

	if opts.DoConfirm {
		var confirmed bool
		err = prompt.Confirm(fmt.Sprintf("Delete the following API key: %s?", opts.APIKey), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	_, err = client.DeleteAPIKey(opts.APIKey)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s API key successfully deleted: %s\n", cs.SuccessIcon(), opts.APIKey)
	}
	return nil
}
