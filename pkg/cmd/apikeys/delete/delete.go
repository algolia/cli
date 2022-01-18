package delete

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// DeleteOptions represents the options for the create command
type DeleteOptions struct {
	config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (search.ClientInterface, error)

	APIKeys   []string
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
		Use:   "delete <api_key>, <api_key>...",
		Short: "Delete API key(s)",
		Long:  `Delete the given API key(s).`,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.APIKeys = args
			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when passing a single argument")
				}
				opts.DoConfirm = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompt")

	return cmd
}

// runDeleteCmd runs the delete command
func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	// Check that all the API keys exists
	for _, apiKey := range opts.APIKeys {
		_, err := client.GetAPIKey(apiKey)
		if err != nil {
			return fmt.Errorf("API key %q does not exist", apiKey)
		}
	}

	if opts.DoConfirm {
		var confirmed bool
		p := &survey.Confirm{
			Message: fmt.Sprintf("Delete the following API Key(s) %v?", opts.APIKeys),
			Default: false,
		}
		err = prompt.SurveyAskOne(p, &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	// Delete all the API keys
	for _, apiKey := range opts.APIKeys {
		_, err = client.DeleteAPIKey(apiKey)
		if err != nil {
			return err
		}
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s API key(s) successfully deleted: %v\n", cs.SuccessIcon(), opts.APIKeys)
	}
	return nil
}
