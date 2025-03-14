package create

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// CreateOptions represents the options for the create command
type CreateOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	ACL         []string
	Description string
	Indices     []string
	Referers    []string
	Validity    time.Duration
}

// NewCreateCmd returns a new instance of CreateCmd
func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:           f.IOStreams,
		config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "n", "c"},
		Args:    validators.NoArgs(),
		Annotations: map[string]string{
			"acls": "admin",
		},
		Short: "Create a new API key",
		Long:  `Create a new API key with the provided parameters.`,
		Example: heredoc.Doc(`
			# Create a new API key targeting the index "MOVIES", with the "search" and "browse" ACL and a description
			$ algolia apikeys create --indices MOVIES --acl search,browse --description "Search & Browse API Key"

			# Create a new API key targeting the indices "MOVIES" and "SERIES", with the "https://example.com" referer, with a validity of 1 hour and a description
			$ algolia apikeys create -i MOVIES,SERIES --acl search -r "https://example.com" --u 1h -d "Search-only API Key for MOVIES & SERIES"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.ACL, "acl", nil, heredoc.Docf(`
		API key's ACL.

			%[1]ssearch%[1]s: can perform search operations.
			%[1]sbrowse%[1]s: can retrieve all index data with the browse endpoint.
			%[1]saddObject%[1]s: can add or update records in the index.
			%[1]sdeleteObject%[1]s: can delete an existing record.
			%[1]slistIndexes%[1]s: can get a list of all indices.
			%[1]sdeleteIndex%[1]s: can delete an index.
			%[1]ssettings%[1]s: can read all index settings.
			%[1]seditSettings%[1]s: can update all index settings.
			%[1]sanalytics%[1]s: can retrieve data with the Analytics API.
			%[1]srecommendation%[1]s: can interact with the Recommendation API.
			%[1]susage%[1]s: can retrieve data with the Usage API.
			%[1]slogs%[1]s: can query the logs.
			%[1]sseeUnretrievableAttributes%[1]s: can retrieve unretrievableAttributes for all operations that return records.
	`, "`"))

	cmd.Flags().StringSliceVarP(&opts.Indices, "indices", "i", nil, heredoc.Docf(`
		Index names or patterns that this API key can access. By default, an API key can access all indices in the same application.

		You can use leading and trailing wildcard characters (%[1]s*%[1]s).
		For example, %[1]sdev_*%[1]s matches all indices starting with %[1]sdev_%[1]s. %[1]s*_dev%[1]s matches all indices ending with %[1]s_dev%[1]s. %[1]s*_products_*%[1]s matches all indices containing %[1]sproducts%[1]s.
	`, "`"))

	cmd.Flags().DurationVarP(&opts.Validity, "validity", "u", 0, heredoc.Doc(`
		Duration (in seconds) after which the API key expires. By default (a value of 0), API keys don't expire.`,
	))

	cmd.Flags().StringSliceVarP(&opts.Referers, "referers", "r", nil, heredoc.Docf(`
		Specify the list of referrers that can perform an operation.
		You can use the wildcard character (%[1]s*%[1]s) to match subdomains or entire websites.
	`, "`"))

	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", heredoc.Doc(`
		Describe an API key to help you identify its uses.`,
	))

	_ = cmd.RegisterFlagCompletionFunc(
		"indices",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			client, err := f.SearchClient()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			indicesRes, err := client.ListIndices(client.NewApiListIndicesRequest())
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			allowedIndices := make([]string, 0, len(indicesRes.Items))
			for _, index := range indicesRes.Items {
				allowedIndices = append(
					allowedIndices,
					fmt.Sprintf("%s\t%s records", index.Name, humanize.Comma(int64(index.Entries))),
				)
			}
			return allowedIndices, cobra.ShellCompDirectiveNoFileComp
		},
	)

	_ = cmd.RegisterFlagCompletionFunc("acl",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"search":                     "perform search operations",
			"browse":                     "retrieve all index data with the browse endpoint",
			"addObject":                  "add or update a records in the index",
			"deleteObject":               "delete an existing record",
			"listIndexes":                "get a list of all existing indices",
			"deleteIndex":                "delete an index",
			"settings":                   "read all index settings",
			"editSettings":               "update all index settings",
			"analytics":                  "retrieve data with the Analytics API",
			"recommendation":             "interact with the Recommendation API",
			"usage":                      "retrieve data with the Usage API",
			"logs":                       "query the logs",
			"seeUnretrievableAttributes": "retrieve unretrievableAttributes for all operations that return records",
		}, "can"))

	return cmd
}

// runCreateCmd executes the create command
func runCreateCmd(opts *CreateOptions) error {
	var acls []search.Acl
	for _, a := range opts.ACL {
		acls = append(acls, search.Acl(a))
	}
	validity := int32(opts.Validity.Seconds())
	key := search.ApiKey{
		Acl:         acls,
		Indexes:     opts.Indices,
		Validity:    &validity,
		Referers:    opts.Referers,
		Description: &opts.Description,
	}

	client, err := opts.SearchClient()
	if err != nil {
		return err
	}
	res, err := client.AddApiKey(client.NewApiAddApiKeyRequest(&key))
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s API key created: %s\n", cs.SuccessIcon(), res.Key)
	}
	return nil
}
