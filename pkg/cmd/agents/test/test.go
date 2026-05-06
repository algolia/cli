package test

import (
	"context"
	"os"
	"os/signal"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type TestOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	ConfigFile    string
	InputFile     string
	Message       string
	NoStream      bool
	Compatibility string
	DryRun        bool
}

func NewTestCmd(f *cmdutil.Factory, runF func(*TestOptions) error) *cobra.Command {
	opts := &TestOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:   "test --config <file> [--input <file> | --message <text>]",
		Short: "Test an Agent Studio agent configuration without persisting it",
		Long: heredoc.Doc(`
			Send a completion to /1/agents/test/completions using an
			ephemeral agent configuration. The backend's special-case
			agent_id="test" route doesn't persist anything — it
			instantiates the configuration in-memory, runs the message,
			streams the result back. The primary developer loop for
			iterating on agent prompts/tools without polluting the agent
			list.

			Streaming responses (default) are emitted as NDJSON: one
			parsed event per line as {"type":"...","data":{...}}. Pipe
			to jq to filter (e.g. select(.type=="text-delta")). Use
			--no-stream for a single buffered JSON response instead.

			Use --dry-run to print the resolved request body without
			calling the API.
		`),
		Example: heredoc.Doc(`
			# Quick one-liner with a config file
			$ algolia agents test -c cfg.json -m "Recommend a laptop under $1000"

			# Multi-turn from a messages file, streamed
			$ algolia agents test -c cfg.json -i messages.json | jq -c .

			# Just the assistant text
			$ algolia agents test -c cfg.json -m "hi" | jq -r 'select(.type=="text-delta") | .data.delta'

			# Buffered, single JSON response
			$ algolia agents test -c cfg.json -m "hi" --no-stream

			# Preview the request without sending
			$ algolia agents test -c cfg.json -m "hi" --dry-run
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runTestCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.ConfigFile, "config", "c", "", "JSON file with the agent configuration to test (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("config")
	cmd.Flags().
		StringVarP(&opts.InputFile, "input", "i", "", "JSON file with messages array (use \"-\" for stdin)")
	cmd.Flags().
		StringVarP(&opts.Message, "message", "m", "", "Single user message (convenience for one-shot prompts)")
	cmd.Flags().BoolVar(&opts.NoStream, "no-stream", false, "Request a buffered JSON response instead of SSE")
	cmd.Flags().
		StringVar(&opts.Compatibility, "compatibility", "", "Streaming protocol: v4 (ai-sdk-4) or v5 (ai-sdk-5, default)")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Print the resolved request body without calling the API")

	cmd.MarkFlagsMutuallyExclusive("input", "message")

	return cmd
}

func runTestCmd(opts *TestOptions) error {
	configuration, err := shared.ReadJSONFile(opts.IO.In, opts.ConfigFile)
	if err != nil {
		return err
	}
	messages, err := shared.BuildMessages(opts.IO.In, shared.MessageInput{
		InputFile: opts.InputFile,
		Message:   opts.Message,
	})
	if err != nil {
		return err
	}
	body, err := shared.MarshalCompletionBody(messages, configuration)
	if err != nil {
		return err
	}

	// Validate --compatibility BEFORE the dry-run short-circuit. An
	// invalid value would produce an invalid request, so previewing a
	// "valid-looking" body for it is misleading. Same rationale as
	// rejecting --input + --message together at flag-parse time.
	mode, err := shared.NormalizeCompatibility(opts.Compatibility)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, cmdutil.NewPrintFlags(), false,
			"test_completion", "POST /1/agents/test/completions", "", body, nil)
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	// Local SIGINT handling: cancels the in-flight HTTP request so the
	// transport tears down the SSE stream cleanly. Deferred stop()
	// releases the signal handler when this function returns.
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	resp, err := client.Completions(ctx, "test", body, agentstudio.CompletionOptions{
		Stream:        !opts.NoStream,
		Compatibility: mode,
	})
	if err != nil {
		return err
	}
	return shared.RenderCompletion(opts.IO, resp.Body, resp.Header.Get("Content-Type"))
}
