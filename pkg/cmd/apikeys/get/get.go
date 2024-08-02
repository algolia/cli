package get

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikeys/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// GetOptions represents the options for the get command
type GetOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	APIKey string

	PrintFlags *cmdutil.PrintFlags
}

// NewGetCmd returns a new instance of DeleteCmd
func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:           f.IOStreams,
		config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get <api-key>",
		Short: "Get API key",
		Long: heredoc.Doc(`
			Get the details of a given API Key (ACLs, description, indexes, and other attributes).
		`),
		Example: heredoc.Doc(`
			# Get an API key
			$ algolia --application-id app-id apikeys get abcdef1234567890
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.APIKey = args[0]

			if runF != nil {
				return runF(opts)
			}

			return runGetCmd(opts)
		},
	}

	return cmd
}

// runGetCmd runs the get command
func runGetCmd(opts *GetOptions) error {
	opts.config.Profile().APIKey = opts.APIKey
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	key, err := client.GetApiKey(client.NewApiGetApiKeyRequest(opts.APIKey))
	if err != nil {
		return fmt.Errorf("API key %q does not exist", opts.APIKey)
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	keyResult := shared.JSONKey{
		ACL:                    key.Acl,
		CreatedAt:              key.CreatedAt,
		Description:            *key.Description,
		Indexes:                key.Indexes,
		MaxQueriesPerIPPerHour: key.MaxQueriesPerIPPerHour,
		MaxHitsPerQuery:        key.MaxHitsPerQuery,
		Referers:               key.Referers,
		QueryParameters:        key.QueryParameters,
		Validity:               key.Validity,
		Value:                  *key.Value,
	}

	if err := p.Print(opts.IO, keyResult); err != nil {
		return err
	}

	return nil
}
