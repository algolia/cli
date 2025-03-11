package update

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type UpdateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	ObjectID     string
	Name         string
	Instructions string
	Tone         string
}

// NewUpdateCmd creates and returns an update command for GenAI prompts.
func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "update <id> [--name <name>] [--instructions <instructions>] [--tone <tone>]",
		Args:  cobra.ExactArgs(1),
		Short: "Update a GenAI prompt",
		Long: heredoc.Doc(`
			Update an existing GenAI prompt.
		`),
		Example: heredoc.Doc(`
			# Update a prompt name
			$ algolia genai prompt update b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --name "New Product Comparison"

			# Update a prompt's instructions
			$ algolia genai prompt update b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --instructions "New instructions for comparing products"

			# Update a prompt's tone
			$ algolia genai prompt update b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --tone professional
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ObjectID = args[0]

			// At least one flag should be specified
			if opts.Name == "" && opts.Instructions == "" && opts.Tone == "" {
				return cmdutil.FlagErrorf("at least one of --name, --instructions, or --tone must be specified")
			}

			if runF != nil {
				return runF(opts)
			}

			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the prompt")
	cmd.Flags().StringVar(&opts.Instructions, "instructions", "", "New instructions for the prompt")
	cmd.Flags().StringVar(&opts.Tone, "tone", "", "New tone for the prompt (natural, friendly, or professional)")

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Updating prompt")

	input := genai.UpdatePromptInput{
		ObjectID: opts.ObjectID,
	}

	if opts.Name != "" {
		input.Name = opts.Name
	}

	if opts.Instructions != "" {
		input.Instructions = opts.Instructions
	}

	if opts.Tone != "" {
		input.Tone = opts.Tone
	}

	_, err = client.UpdatePrompt(input)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Prompt %s updated\n", cs.SuccessIconWithColor(cs.Green), cs.Bold(opts.ObjectID))
	}

	return nil
}
