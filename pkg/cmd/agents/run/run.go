package run

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

type RunOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID string
	shared.CompletionInputs
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
			persisted agent configuration. This is what downstream apps
			call in production.

			Output: TTY-attached stdout renders a flowing transcript;
			non-TTY emits NDJSON. --ndjson forces NDJSON on a TTY,
			--no-stream returns a single buffered JSON response.
		`),
		Example: heredoc.Doc(`
			$ algolia agents run <id> -m "What's new today?"

			$ cat messages.json | algolia agents run <id> -i -

			$ algolia agents run <id> -m "hi" --no-stream
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

	shared.RegisterCompletionFlags(cmd, &opts.CompletionInputs)

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
	// `agents run` uses the persisted agent's configuration; no
	// `configuration` field in the body (backend rejects it for real agents).
	body, err := shared.MarshalCompletionBody(messages, nil)
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

	ctx, stop := signal.NotifyContext(shared.OrBackground(opts.Ctx), os.Interrupt)
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
