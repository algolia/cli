package get

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	ObjectID string
}

// NewGetCmd creates and returns a get command for GenAI prompts.
func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
	}

	cmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get a GenAI prompt",
		Long: heredoc.Doc(`
			Get a specific GenAI prompt by ID.
		`),
		Example: heredoc.Doc(`
			# Get a prompt by ID
			$ algolia genai prompt get b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ObjectID = args[0]

			if runF != nil {
				return runF(opts)
			}

			return runGetCmd(opts)
		},
	}

	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	opts.IO.StartProgressIndicatorWithLabel("Fetching prompt")

	prompt, err := client.GetPrompt(opts.ObjectID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("ID:"), cs.Bold(prompt.ObjectID))
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Name:"), prompt.Name)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Instructions:"), prompt.Instructions)
		if prompt.Tone != "" {
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Tone:"), prompt.Tone)
		}
		fmt.Fprintf(opts.IO.Out, "%s %d\n", cs.Bold("Linked Responses:"), prompt.LinkedResponses)
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Created:"), formatTime(prompt.CreatedAt))
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Updated:"), formatTime(prompt.UpdatedAt))
	} else {
		fmt.Fprintf(opts.IO.Out, "%s\n", prompt.ObjectID)
	}

	return nil
}

// formatTime formats time to a readable format
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
