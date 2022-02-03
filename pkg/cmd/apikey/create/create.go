package create

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// CreateOptions represents the options for the create command
type CreateOptions struct {
	config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

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
		Use:   "create",
		Args:  validators.NoArgs,
		Short: "Create a new API key",
		Long:  `Create a new API key with the provided parameters.`,
		Example: heredoc.Doc(`
			$ algolia create --indices foo --acl search,browse --description "Search & Browse API Key"
			$ algolia create -i foo,bar --acl search -r "http://foo.com" --u 1h -d "Search-only API Key for foo & bar"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.ACL, "acl", nil, heredoc.Docf(`
		ACL of the API Key.

		%[1]ssearch%[1]s: allowed to perform search operations.
		%[1]sbrowse%[1]s: allowed to retrieve all index data with the browse endpoint.
		%[1]saddObject%[1]s: allowed to add or update a records in the index.
		%[1]sdeleteObject%[1]s: allowed to delete an existing record.
		%[1]slistIndexes%[1]s: allowed to get a list of all existing indices.
		%[1]sdeleteIndex%[1]s: allowed to delete an index.
		%[1]ssettings%[1]s: allowed to read all index settings.
		%[1]seditSettings%[1]s: allowed to update all index settings.
		%[1]sanalytics%[1]s: allowed to retrieve data with the Analytics API.
		%[1]srecommendation%[1]s: allowed to interact with the Recommendation API.
		%[1]susage%[1]s: allowed to retrieve data with the Usage API.
		%[1]slogs%[1]s: allowed to query the logs.
		%[1]sseeUnretrievableAttributes%[1]s: allowed to retrieve unretrievableAttributes for all operations that return records.
	`, "`"))

	cmd.Flags().StringSliceVarP(&opts.Indices, "indices", "i", nil, heredoc.Docf(`
		Specify the list of targeted indices.
		You can target all indices starting with a prefix or ending with a suffix using the %[1]s*%[1]s character.
		For example, %[1]sdev_*%[1]s matches all indices starting with %[1]sdev_%[1]s and %[1]s*_dev%[1]s matches all indices ending with %[1]s_dev%[1]s.
	`, "`"))

	cmd.Flags().DurationVarP(&opts.Validity, "validity", "u", 0, heredoc.Doc(`
		How long this API key is valid, in seconds.
		A value of 0 means the API key doesn’t expire.`,
	))

	cmd.Flags().StringSliceVarP(&opts.Referers, "referers", "r", nil, heredoc.Docf(`
		Specify the list of referrers that can perform an operation.
		You can use the %[1]s*%[1]s (asterisk) character as a wildcard to match subdomains, or all pages of a website.
	`, "`"))

	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", heredoc.Doc(`
		Specify a description of the API key.
		Used for informative purposes only. It has impact on the functionality of the API key.`,
	))

	cmd.RegisterFlagCompletionFunc("acl", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"search", "browse", "addObject", "deleteObject", "listIndexes", "deleteIndex", "settings", "editSettings", "analytics", "recommendation", "usage", "logs", "seeUnretrievableAttributes"}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

// runCreateCmd executes the create command
func runCreateCmd(opts *CreateOptions) error {
	key := search.Key{
		ACL:         opts.ACL,
		Indexes:     opts.Indices,
		Validity:    opts.Validity,
		Referers:    opts.Referers,
		Description: opts.Description,
	}

	client, err := opts.SearchClient()
	if err != nil {
		return err
	}
	res, err := client.AddAPIKey(key)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s API key created: %s\n", cs.SuccessIcon(), res.Key)
	}
	return nil
}
