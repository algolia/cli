package try

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

// TryOptions configures `algolia agents try`.
//
// Note: no DryRun field. The whole command IS the dry-run in the
// conversational-ai sense — it sends the configuration to the
// backend's `agent_id="test"` route, which runs a real completion
// against an in-memory configuration without persisting anything.
// Adding a `--dry-run` flag on top would mean "dry-run a dry-run",
// the awkward case that motivated the rename from `agents test`.
// See AGENTS.md → "Agent Studio" → "On `--dry-run`".
type TryOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	ConfigFile      string
	InputFile       string
	Message         string
	NoStream        bool
	Compatibility   string
	NoCache         bool
	NoMemory        bool
	NoAnalytics     bool
	SecureUserToken string
	NDJSON          bool
}

func NewTryCmd(f *cmdutil.Factory, runF func(*TryOptions) error) *cobra.Command {
	opts := &TryOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:   "try --config <file> [--input <file> | --message <text>]",
		Short: "Try an Agent Studio agent configuration without persisting it",
		Long: heredoc.Doc(`
			Send a completion to /1/agents/test/completions using an
			ephemeral agent configuration. The backend's special-case
			agent_id="test" route doesn't persist anything — it
			instantiates the configuration in-memory, runs the message,
			streams the result back. The primary developer loop for
			iterating on agent prompts/tools without polluting the
			agents list, and the conversational-ai-side equivalent of
			"dry-run" semantics for an agentic configuration.

			Streaming responses (default) are emitted as NDJSON: one
			parsed event per line as {"type":"...","data":{...}}. Pipe
			to jq to filter (e.g. select(.type=="text-delta")). Use
			--no-stream for a single buffered JSON response instead.

			There is no --dry-run flag here on purpose — the whole
			command is the dry-run. To preview the request body
			without calling the API, redirect stdout: --no-stream
			against an unreachable backend, or build the body
			yourself. To preview a *create* or *update* request
			without sending it, use those commands' --dry-run flags.
		`),
		Example: heredoc.Doc(`
			# Quick one-liner with a config file
			$ algolia agents try -c cfg.json -m "Recommend a laptop under $1000"

			# Multi-turn from a messages file, streamed
			$ algolia agents try -c cfg.json -i messages.json | jq -c .

			# Just the assistant text
			$ algolia agents try -c cfg.json -m "hi" | jq -r 'select(.type=="text-delta") | .data.delta'

			# Buffered, single JSON response
			$ algolia agents try -c cfg.json -m "hi" --no-stream
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runTryCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.ConfigFile, "config", "c", "", "JSON file with the agent configuration to try (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("config")
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

	cmd.MarkFlagsMutuallyExclusive("input", "message")

	return cmd
}

func runTryCmd(opts *TryOptions) error {
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
	mode, err := shared.NormalizeCompatibility(opts.Compatibility)
	if err != nil {
		return err
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
