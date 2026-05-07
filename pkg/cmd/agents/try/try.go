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

// TryOptions configures `algolia agents try`. No DryRun field — the
// command IS the dry-run (see docs/agents.md "On --dry-run").
type TryOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	ConfigFile string
	shared.CompletionInputs
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
			ephemeral agent configuration. The backend's agent_id="test"
			route runs the message in-memory and streams the result —
			nothing is persisted. The primary developer loop for
			iterating on agent prompts/tools without polluting the
			agents list.

			Output: TTY-attached stdout renders a flowing transcript
			(text inline, tool calls/results dim, errors red); non-TTY
			emits NDJSON one event per line. Pass --ndjson to force
			NDJSON on a TTY, or --no-stream for a single buffered
			JSON response.

			There is no --dry-run flag — the whole command is the
			dry-run. See docs/agents.md.
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
	shared.RegisterCompletionFlags(cmd, &opts.CompletionInputs)

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

	// SIGINT cancels the in-flight request so the SSE stream tears down cleanly.
	ctx, stop := signal.NotifyContext(shared.OrBackground(opts.Ctx), os.Interrupt)
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
