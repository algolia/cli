package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type addSearchIndexOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID          string
	ToolName         string
	Index            string
	Description      string
	SearchParameters string
	OutputChanged    bool
}

func newAddSearchIndexCmd(f *cmdutil.Factory, runF func(*addSearchIndexOptions) error) *cobra.Command {
	opts := &addSearchIndexOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "add-search-index <agent-id>",
		Short: "Add an Algolia Search index tool (algolia_search_index) to an agent",
		Long: heredoc.Doc(`
			Fetches the agent, merges an Algolia Search built-in tool entry
			(type algolia_search_index) into its tools array, and PATCHes
			the agent.

			Each Algolia search tool must include a short "name" for the tool
			(3–32 characters) per the Agent Studio OpenAPI schema, in addition
			to "type" and "indices":
			https://www.algolia.com/doc/rest-api/agent-studio/agents/update-agent

			If the agent already has an algolia_search_index tool, the new
			index is appended to its indices list (unless that index name
			is already present — then the command fails).

			For narrative context on search tools, see:
			https://www.algolia.com/doc/guides/algolia-ai/agent-studio/how-to/tools
		`),
		Example: heredoc.Doc(`
			$ algolia agents tools add-search-index $AGENT_ID --tool-name product_search --index PRODUCTS --description "Product catalog"

			$ algolia agents tools add-search-index $AGENT_ID --index PRODUCTS --description "Catalog" \
			  --search-parameters '{"filters":"inStock:true","attributesToRetrieve":["title","price"]}'
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.Index == "" {
				return cmdutil.FlagErrorf("--index is required")
			}
			if opts.Description == "" {
				opts.Description = fmt.Sprintf("Algolia index %s", opts.Index)
			}
			opts.ToolName = strings.TrimSpace(opts.ToolName)
			if opts.ToolName == "" {
				opts.ToolName = deriveDefaultToolName(opts.Index)
			}
			if err := validateToolName(opts.ToolName); err != nil {
				return cmdutil.FlagErrorf(
					"%v; set --tool-name to a value between 3 and 32 characters",
					err,
				)
			}
			if runF != nil {
				return runF(opts)
			}
			return runAddSearchIndexCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.ToolName, "tool-name", "", `OpenAPI tool "name" field (3–32 chars; default derived from --index)`)
	cmd.Flags().StringVar(&opts.Index, "index", "", "Algolia index name (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Human-facing description for the LLM (default: \"Algolia index <name>\")")
	cmd.Flags().
		StringVar(&opts.SearchParameters, "search-parameters", "", "JSON object merged into the index entry as \"searchParameters\" (optional)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runAddSearchIndexCmd(opts *addSearchIndexOptions) error {
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching agent")
	agent, err := client.GetAgent(ctx, opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	mergedTools, err := mergeAlgoliaSearchIndexTool(
		agent.Tools,
		opts.ToolName,
		opts.Index,
		opts.Description,
		[]byte(opts.SearchParameters),
	)
	if err != nil {
		return err
	}

	patch, err := json.Marshal(map[string]json.RawMessage{
		"tools": mergedTools,
	})
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating agent tools")
	updated, err := client.UpdateAgent(ctx, opts.AgentID, patch)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, updated)
}

// deriveDefaultToolName builds a 3–32 character identifier from the index
// name when --tool-name is omitted.
func deriveDefaultToolName(index string) string {
	const maxRunes = 32
	s := sanitizeIndexForToolName(index)
	if s == "" {
		return "search_tool"
	}
	runes := []rune(s)
	if len(runes) > maxRunes {
		s = string(runes[:maxRunes])
		runes = []rune(s)
	}
	if len(runes) < 3 {
		s = "srch_" + s
		runes = []rune(s)
		if len(runes) > maxRunes {
			s = string(runes[:maxRunes])
		}
	}
	if len([]rune(s)) < 3 {
		return "search_tool"
	}
	return s
}

func sanitizeIndexForToolName(index string) string {
	index = strings.ToLower(strings.TrimSpace(index))
	var b strings.Builder
	for _, r := range index {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_', r == '-', r == ' ', r == '.':
			b.WriteByte('_')
		}
	}
	s := strings.Trim(b.String(), "_")
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	return s
}
