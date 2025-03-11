package generate

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GenerateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	Query                string
	DataSourceID         string
	PromptID             string
	LogRegion            string
	ObjectID             string
	NbHits               int
	AdditionalFilters    string
	WithObjectIDs        []string
	AttributesToRetrieve []string
	ConversationID       string
	Save                 bool
	UseCache             bool
}

// NewGenerateCmd creates and returns a generate command for GenAI responses.
func NewGenerateCmd(f *cmdutil.Factory, runF func(*GenerateOptions) error) *cobra.Command {
	opts := &GenerateOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
		LogRegion:   "us", // Default to US region
		NbHits:      4,    // Default number of hits
	}

	var objectIDsString string

	cmd := &cobra.Command{
		Use:   "generate --query <query> --datasource <id> --prompt <id>",
		Short: "Generate a GenAI response",
		Long: heredoc.Doc(`
			Generate a new GenAI response using a prompt and data source.
		`),
		Example: heredoc.Doc(`
			# Generate a response to a query
			$ algolia genai response generate --query "Compare iPhone 13 and Samsung S21" --datasource b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --prompt b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3

			# Generate a response with additional filters
			$ algolia genai response generate --query "Compare phones" --datasource b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --prompt b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --filters 'model:"iPhone 13" OR model:"Samsung S21"'

			# Generate a response without saving it
			$ algolia genai response generate --query "Compare phones" --datasource b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --prompt b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --no-save
			
			# Generate a response using specific object IDs instead of search
			$ algolia genai response generate --query "Compare these products" --datasource b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --prompt b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --object-ids "product1,product2,product3"
			
			# Generate a response for a conversation
			$ algolia genai response generate --query "Tell me more about the second one" --datasource b4e52d1a-2509-49ea-ba36-f6f5c3a83ba1 --prompt b4e52d1a-2509-49ea-ba36-f6f5c3a83ba3 --conversation-id conv123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Query == "" {
				return cmdutil.FlagErrorf("--query is required")
			}

			if opts.DataSourceID == "" {
				return cmdutil.FlagErrorf("--datasource is required")
			}

			if opts.PromptID == "" {
				return cmdutil.FlagErrorf("--prompt is required")
			}

			// Parse object IDs if provided
			if objectIDsString != "" {
				opts.WithObjectIDs = strings.Split(objectIDsString, ",")
			}

			if runF != nil {
				return runF(opts)
			}

			return runGenerateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Query, "query", "", "The query to generate a response for")
	cmd.Flags().StringVar(&opts.DataSourceID, "datasource", "", "The ID of the data source to use")
	cmd.Flags().StringVar(&opts.PromptID, "prompt", "", "The ID of the prompt to use")
	cmd.Flags().StringVar(&opts.LogRegion, "region", "us", "The region to use for LLM routing (us or de)")
	cmd.Flags().StringVar(&opts.ObjectID, "id", "", "Optional object ID for the response")
	cmd.Flags().IntVar(&opts.NbHits, "hits", 4, "Number of hits to retrieve as context")
	cmd.Flags().StringVar(&opts.AdditionalFilters, "filters", "", "Additional filters to apply")
	cmd.Flags().StringVar(&objectIDsString, "object-ids", "", "Specific object IDs to use instead of search (comma-separated)")
	cmd.Flags().StringSliceVar(&opts.AttributesToRetrieve, "attributes", nil, "Specific attributes to retrieve from the hits")
	cmd.Flags().StringVar(&opts.ConversationID, "conversation-id", "", "Conversation ID for follow-up queries")
	cmd.Flags().BoolVar(&opts.Save, "save", false, "Save the response")
	cmd.Flags().BoolVar(&opts.UseCache, "use-cache", false, "Use cached response if available")

	_ = cmd.MarkFlagRequired("query")
	_ = cmd.MarkFlagRequired("datasource")
	_ = cmd.MarkFlagRequired("prompt")

	return cmd
}

func runGenerateCmd(opts *GenerateOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Generating response")

	input := genai.GenerateResponseInput{
		Query:        opts.Query,
		DataSourceID: opts.DataSourceID,
		PromptID:     opts.PromptID,
		LogRegion:    opts.LogRegion,
		NbHits:       opts.NbHits,
		Save:         opts.Save,
		UseCache:     opts.UseCache,
	}

	if opts.ObjectID != "" {
		input.ObjectID = opts.ObjectID
	}

	if opts.AdditionalFilters != "" {
		input.AdditionalFilters = opts.AdditionalFilters
	}

	if len(opts.WithObjectIDs) > 0 {
		input.WithObjectIDs = opts.WithObjectIDs
	}

	if len(opts.AttributesToRetrieve) > 0 {
		input.AttributesToRetrieve = opts.AttributesToRetrieve
	}

	if opts.ConversationID != "" {
		input.ConversationID = opts.ConversationID
	}

	response, err := client.GenerateResponse(input)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Response generated with ID: %s\n\n", cs.SuccessIconWithColor(cs.Green), cs.Bold(response.ObjectID))
		fmt.Fprintf(opts.IO.Out, "%s\n", response.Response)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s", response.Response)
	}

	return nil
}
