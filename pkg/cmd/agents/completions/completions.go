package completions

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/utils"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// CompletionsOptions holds the dependencies and flags for the completions command.
type CompletionsOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)

	AgentID           string
	Message           string
	File              string
	ConversationID    string
	CompatibilityMode string
	Cache             bool
	Memory            bool
	Analytics         bool

	PrintFlags *cmdutil.PrintFlags
}

// NewCompletionsCmd returns the `agents completions` command.
func NewCompletionsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CompletionsOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "completions <agent-id>",
		Aliases:           []string{"chat", "completion"},
		Short:             "Create a completion for an agent",
		Args:              validators.ExactArgsWithMsg(1, "agents completions requires an <agent-id> argument."),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "search",
		},
		Example: heredoc.Doc(`
			# Send a one-off message to an agent
			$ algolia agents completions my-agent --message "What are your best sellers?"

			# Continue an existing conversation
			$ algolia agents completions my-agent --message "And in blue?" --conversation-id conv_123

			# Send a full completion request body (e.g. for tool approvals or multi-turn history)
			$ algolia agents completions my-agent --file request.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]

			if (opts.Message != "") == (opts.File != "") {
				return cmdutil.FlagErrorf("exactly one of `--message` or `--file` is required")
			}

			return runCompletionsCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Message to send to the agent")
	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file with a full completion request body (use - for stdin)")
	cmd.Flags().StringVar(&opts.ConversationID, "conversation-id", "", "ID of an existing conversation to continue")
	cmd.Flags().
		StringVar(&opts.CompatibilityMode, "compatibility-mode", string(agentStudio.COMPATIBILITY_MODE_AI_SDK_5), "Compatibility mode for the completion API (ai-sdk-4 or ai-sdk-5)")
	cmd.Flags().BoolVar(&opts.Cache, "cache", false, "Use cached responses if available")
	cmd.Flags().BoolVar(&opts.Memory, "memory", true, "Enable agent memory for this completion")
	cmd.Flags().BoolVar(&opts.Analytics, "analytics", true, "Enable analytics for this completion")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func buildCompletionRequest(opts *CompletionsOptions) (*agentStudio.AgentCompletionRequest, error) {
	if opts.File != "" {
		raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}
		var req agentStudio.AgentCompletionRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			return nil, fmt.Errorf("parsing completion request JSON: %w", err)
		}
		return &req, nil
	}

	messages := agentStudio.ArrayOfMessageV5AsMessagesUnion([]agentStudio.MessageV5{
		*agentStudio.UserMessageV5AsMessageV5(&agentStudio.UserMessageV5{
			Role:  "user",
			Parts: []agentStudio.TextPartV5{{Text: opts.Message}},
		}),
	})

	req := &agentStudio.AgentCompletionRequest{
		Messages: *utils.NewNullable(messages),
	}
	if opts.ConversationID != "" {
		req.Id = &opts.ConversationID
	}
	return req, nil
}

func runCompletionsCmd(opts *CompletionsOptions) error {
	completionReq, err := buildCompletionRequest(opts)
	if err != nil {
		return err
	}

	compatibilityMode, err := agentStudio.NewCompatibilityModeFromValue(opts.CompatibilityMode)
	if err != nil {
		return cmdutil.FlagErrorf("%s", err)
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	req := client.
		NewApiCreateAgentCompletionRequest(opts.AgentID, *compatibilityMode, completionReq).
		WithCache(opts.Cache).
		WithMemory(opts.Memory).
		WithAnalytics(opts.Analytics)

	opts.IO.StartProgressIndicatorWithLabel("Waiting for agent response")

	res, err := client.CreateAgentCompletion(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
