package run

import (
	"context"
	"fmt"
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

type RunOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID         string
	InputFile       string
	Message         string
	NoStream        bool
	Compatibility   string
	NoCache         bool
	NoMemory        bool
	NoAnalytics     bool
	SecureUserToken string
	NDJSON          bool
	DryRun          bool
}

func NewRunCmd(f *cmdutil.Factory, runF func(*RunOptions) error) *cobra.Command {
	opts := &RunOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:   "run <agent-id> [--input <file> | --message <text>]",
		Short: "Run a published Agent Studio agent and stream the response",
		Long: heredoc.Doc(`
			Send a completion to /1/agents/<id>/completions using the
			persisted agent configuration. Equivalent to "agents test"
			except it uses an existing (typically published) agent
			rather than an in-memory configuration — this is what
			downstream apps call in production.

			Streaming responses (default) are emitted as NDJSON: one
			parsed event per line as {"type":"...","data":{...}}. Use
			--no-stream for a single buffered JSON response instead.
		`),
		Example: heredoc.Doc(`
			$ algolia agents run <id> -m "What's new today?"

			$ cat messages.json | algolia agents run <id> -i -

			$ algolia agents run <id> -m "hi" --no-stream

			$ algolia agents run <id> -m "hi" --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runRunCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.InputFile, "input", "i", "", "JSON file with messages array (use \"-\" for stdin)")
	cmd.Flags().
		StringVarP(&opts.Message, "message", "m", "", "Single user message (convenience for one-shot prompts)")
	cmd.Flags().BoolVar(&opts.NoStream, "no-stream", false, "Request a buffered JSON response instead of SSE")
	cmd.Flags().
		StringVar(&opts.Compatibility, "compatibility", "", "Streaming protocol: v4 (ai-sdk-4) or v5 (ai-sdk-5, default)")
	cmd.Flags().BoolVar(&opts.NoCache, "no-cache", false, "Bypass the backend completion cache (default: cache enabled)")
	cmd.Flags().BoolVar(&opts.NoMemory, "no-memory", false, "Disable agent memory for this completion (default: memory enabled)")
	cmd.Flags().BoolVar(&opts.NoAnalytics, "no-analytics", false, "Skip Agent Studio analytics for this completion (default: analytics enabled)")
	cmd.Flags().
		StringVar(&opts.SecureUserToken, "secure-user-token", "", "Signed JWT scoping the conversation/memory/analytics partition to an end-user (X-Algolia-Secure-User-Token)")
	cmd.Flags().
		BoolVar(&opts.NDJSON, "ndjson", false, "Force NDJSON output even on a TTY (default on non-TTY; use this when you want machine-parseable events but also want to see them on screen)")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Print the resolved request body without calling the API")

	cmd.MarkFlagsMutuallyExclusive("input", "message")

	return cmd
}

func runRunCmd(opts *RunOptions) error {
	messages, err := shared.BuildMessages(opts.IO.In, shared.MessageInput{
		InputFile: opts.InputFile,
		Message:   opts.Message,
	})
	if err != nil {
		return err
	}
	// `agents run` uses the persisted agent's configuration — no
	// `configuration` field in the body (the backend would reject it
	// for a real agent, since it's only meaningful for agent_id="test").
	body, err := shared.MarshalCompletionBody(messages, nil)
	if err != nil {
		return err
	}

	// Validate --compatibility before the dry-run short-circuit; same
	// rationale as in `agents test`.
	mode, err := shared.NormalizeCompatibility(opts.Compatibility)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, cmdutil.NewPrintFlags(), false,
			"run_completion", fmt.Sprintf("POST /1/agents/%s/completions", opts.AgentID),
			"", body, map[string]any{"agentId": opts.AgentID})
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	resp, err := client.Completions(ctx, opts.AgentID, body, agentstudio.CompletionOptions{
		Stream:          !opts.NoStream,
		Compatibility:   mode,
		NoCache:         opts.NoCache,
		NoMemory:        opts.NoMemory,
		NoAnalytics:     opts.NoAnalytics,
		SecureUserToken: opts.SecureUserToken,
	})
	if err != nil {
		return err
	}
	return shared.RenderCompletion(opts.IO, resp.Body, resp.Header.Get("Content-Type"), opts.NDJSON)
}
