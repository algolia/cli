package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// memoryCall is the shared verb signature for memorize/ponder/consolidate.
type memoryCall func(ctx context.Context, agentID string, body json.RawMessage) (json.RawMessage, error)

type memoryOptions struct {
	io                *iostreams.IOStreams
	ctx               context.Context
	agentStudioClient func() (*agentstudio.Client, error)
	printFlags        *cmdutil.PrintFlags
	verb              string
	agentID           string
	body              string
	file              string
	dryRun            bool
}

func newMemoryCmd(
	f *cmdutil.Factory,
	verb, summary, longText string,
	pickCall func(*agentstudio.Client) memoryCall,
	runF func(*memoryOptions) error,
) *cobra.Command {
	opts := &memoryOptions{
		io:                f.IOStreams,
		agentStudioClient: f.AgentStudioClient,
		printFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
		verb:              verb,
	}
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <agent-id> [--body <json> | -F file]", verb),
		Short: summary,
		Long:  longText,
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.agentID = args[0]
			opts.ctx = cmd.Context()
			if opts.agentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.body != "" && opts.file != "" {
				return cmdutil.FlagErrorf("--body and --file are mutually exclusive")
			}
			if opts.body == "" && opts.file == "" {
				return cmdutil.FlagErrorf("provide --body or --file (\"-\" for stdin)")
			}
			if runF != nil {
				return runF(opts)
			}
			return runMemoryCmd(opts, pickCall)
		},
	}
	cmd.Flags().StringVar(&opts.body, "body", "", "Inline JSON body")
	cmd.Flags().StringVarP(&opts.file, "file", "F", "", "JSON file path (\"-\" for stdin)")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Print what would be sent without calling the API")
	opts.printFlags.AddFlags(cmd)
	return cmd
}

func runMemoryCmd(opts *memoryOptions, pickCall func(*agentstudio.Client) memoryCall) error {
	var body []byte
	if opts.body != "" {
		if !json.Valid([]byte(opts.body)) {
			return cmdutil.FlagErrorf("--body is not valid JSON")
		}
		body = []byte(opts.body)
	} else {
		var err error
		body, err = shared.ReadJSONFile(opts.io.In, opts.file)
		if err != nil {
			return err
		}
	}

	if opts.dryRun {
		var pretty bytes.Buffer
		_ = json.Indent(&pretty, body, "  ", "  ")
		fmt.Fprintf(opts.io.Out,
			"Dry run: would POST /1/agents/agents/%s/%s\n  body: %s\n",
			opts.agentID, opts.verb, pretty.String())
		return nil
	}

	client, err := opts.agentStudioClient()
	if err != nil {
		return err
	}
	call := pickCall(client)

	opts.io.StartProgressIndicatorWithLabel(fmt.Sprintf("Calling %s", opts.verb))
	out, err := call(shared.OrBackground(opts.ctx), opts.agentID, json.RawMessage(body))
	opts.io.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.printFlags.HasStructuredOutput() {
		var anyV any
		if err := json.Unmarshal(out, &anyV); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		return opts.printFlags.Print(opts.io, anyV)
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, out, "", "  "); err != nil {
		_, _ = opts.io.Out.Write(out)
		return nil
	}
	_, _ = opts.io.Out.Write(pretty.Bytes())
	_, _ = opts.io.Out.Write([]byte("\n"))
	return nil
}

func newMemorizeCmd(f *cmdutil.Factory, runF func(*memoryOptions) error) *cobra.Command {
	return newMemoryCmd(f, "memorize",
		"Extract semantic memories from a conversation (unstable)",
		heredoc.Doc(`
			POST /1/agents/agents/{id}/memorize. Body shape is the
			AgentMemorizeRequest schema in the OpenAPI spec
			({providerID, model, messages[], targetMemories?, ...}).
			Pass --body for inline JSON or -F file (\"-\" for stdin).
		`),
		func(c *agentstudio.Client) memoryCall { return c.AgentMemorize }, runF)
}

func newPonderCmd(f *cmdutil.Factory, runF func(*memoryOptions) error) *cobra.Command {
	return newMemoryCmd(f, "ponder",
		"Extract episodic memories (OTAR episodes) from a conversation (unstable)",
		heredoc.Doc(`
			POST /1/agents/agents/{id}/ponder. Body shape is the
			AgentPonderRequest schema in the OpenAPI spec.
		`),
		func(c *agentstudio.Client) memoryCall { return c.AgentPonder }, runF)
}

func newConsolidateCmd(f *cmdutil.Factory, runF func(*memoryOptions) error) *cobra.Command {
	return newMemoryCmd(f, "consolidate",
		"Consolidate an agent's memories by loading existing ones (unstable)",
		heredoc.Doc(`
			POST /1/agents/agents/{id}/consolidate. Body shape is the
			AgentConsolidateRequest schema in the OpenAPI spec
			(includes maxExisting + memoryType).
		`),
		func(c *agentstudio.Client) memoryCall { return c.AgentConsolidate }, runF)
}
