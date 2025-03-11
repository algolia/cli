package delete

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type DeleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	ObjectIDs []string
}

// NewDeleteCmd creates and returns a delete command for GenAI responses.
func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "delete <id>...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Delete GenAI responses",
		Long: heredoc.Doc(`
			Delete one or more GenAI responses.
		`),
		Example: heredoc.Doc(`
			# Delete a single response
			$ algolia genai response delete b4e52d1a-2509-49ea-ba36-f6f5c3a83ba9

			# Delete multiple responses
			$ algolia genai response delete b4e52d1a-2509-49ea-ba36-f6f5c3a83ba9 b4e52d1a-2509-49ea-ba36-f6f5c3a83ba8
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ObjectIDs = args

			if runF != nil {
				return runF(opts)
			}

			return runDeleteCmd(opts)
		},
	}

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Deleting response(s)")

	input := genai.DeleteResponsesInput{
		ObjectIDs: opts.ObjectIDs,
	}

	_, err = client.DeleteResponses(input)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		if len(opts.ObjectIDs) == 1 {
			fmt.Fprintf(opts.IO.Out, "%s Response %s deleted\n", cs.SuccessIconWithColor(cs.Green), cs.Bold(opts.ObjectIDs[0]))
		} else {
			fmt.Fprintf(opts.IO.Out, "%s %d responses deleted: %s\n", cs.SuccessIconWithColor(cs.Green), len(opts.ObjectIDs), cs.Bold(strings.Join(opts.ObjectIDs, ", ")))
		}
	}

	return nil
}
