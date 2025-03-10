package create

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type CreateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	Name         string
	Instructions string
	Tone         string
	ObjectID     string
}

// NewCreateCmd creates and returns a create command for GenAI prompts.
func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "create <name> --instructions <instructions> [--tone <tone>] [--id <id>]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a GenAI prompt",
		Long: heredoc.Doc(`
			Create a new GenAI prompt.

			A prompt defines the instructions for how the GenAI toolkit should generate responses.
		`),
		Example: heredoc.Doc(`
			# Create a prompt with instructions
			$ algolia genai prompt create "Compare Products" --instructions "Help buyers choose products by comparing features"

			# Create a prompt with a specific tone
			$ algolia genai prompt create "Compare Products" --instructions "Help buyers choose products" --tone friendly

			# Create a prompt with a specific ID
			$ algolia genai prompt create "Compare Products" --instructions "Help buyers choose products" --id b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if opts.Instructions == "" {
				return cmdutil.FlagErrorf("--instructions is required")
			}

			if runF != nil {
				return runF(opts)
			}

			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Instructions, "instructions", "", "Instructions for the GenAI prompt")
	cmd.Flags().StringVar(&opts.Tone, "tone", "", "Tone for the prompt (natural, friendly, or professional)")
	cmd.Flags().StringVar(&opts.ObjectID, "id", "", "Optional object ID for the prompt")

	_ = cmd.MarkFlagRequired("instructions")

	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Creating prompt")

	input := genai.CreatePromptInput{
		Name:         opts.Name,
		Instructions: opts.Instructions,
	}

	if opts.Tone != "" {
		input.Tone = opts.Tone
	}

	if opts.ObjectID != "" {
		input.ObjectID = opts.ObjectID
	}

	response, err := client.CreatePrompt(input)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Prompt %s created with ID: %s\n", cs.SuccessIconWithColor(cs.Green), cs.Bold(opts.Name), cs.Bold(response.ObjectID))
	}

	return nil
}
