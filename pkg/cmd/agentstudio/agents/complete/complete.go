package complete

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type CompleteOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID           string
	File              string
	Body              agentstudio.AgentCompletionRequest
	CompatibilityMode string
	Stream            bool
	Cache             bool
}

func NewCompleteCmd(f *cmdutil.Factory, runF func(*CompleteOptions) error) *cobra.Command {
	opts := &CompleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}
	cmd := &cobra.Command{
		Use:     "complete <agent_id>",
		Aliases: []string{"completion", "run"},
		Args:    cobra.ExactArgs(1),
		Short:   "Run a completion against an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Run with a JSON file holding messages
			$ algolia agentstudio agents complete a1b2 -F request.json

			# Continue an existing conversation
			$ algolia agentstudio agents complete a1b2 --id conv-1 --messages '[{"role":"user","parts":[{"type":"text","text":"hi"}]}]'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			if err := cmdutil.MergeFileAndFlagsInto(opts.IO, opts.File, cmd, cmdutil.AgentCompletionRequest, &opts.Body); err != nil {
				return err
			}
			if runF != nil {
				return runF(opts)
			}
			return runCompleteCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "Completion request JSON `file` (use \"-\" for stdin)")
	cmd.Flags().StringVar(&opts.CompatibilityMode, "compatibility-mode", "ai-sdk-5", "API compatibility mode. One of: ai-sdk-4, ai-sdk-5")
	cmd.Flags().BoolVar(&opts.Stream, "stream", false, "Stream the response as SSE bytes instead of waiting for the full JSON body")
	cmd.Flags().BoolVar(&opts.Cache, "cache", true, "Allow the API to return cached responses")
	cmdutil.AddAgentCompletionRequestFlags(cmd)
	return cmd
}

func runCompleteCmd(opts *CompleteOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Running completion")
	body, err := client.CreateCompletion(opts.AgentID, opts.Body, agentstudio.CompletionParams{
		CompatibilityMode: opts.CompatibilityMode,
		Stream:            &opts.Stream,
		Cache:             &opts.Cache,
	})
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	_, err = opts.IO.Out.Write(body)
	if err != nil {
		return err
	}
	if len(body) > 0 && body[len(body)-1] != '\n' {
		_, _ = opts.IO.Out.Write([]byte("\n"))
	}
	return nil
}
