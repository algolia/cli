package create

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

const defaultDescription = "Created with the Algolia CLI"

// CreateOptions represents the options for the create command
type CreateOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient       func() (*search.APIClient, error)
	NewDashboardClient func(clientID string) *dashboard.Client
	LoadToken          func() *auth.StoredToken

	ACL         []string
	Description string
	Indices     []string
	Referers    []string
	Validity    time.Duration

	PrintFlags *cmdutil.PrintFlags
}

// NewCreateCmd returns a new instance of CreateCmd
func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:           f.IOStreams,
		config:       f.Config,
		SearchClient: f.SearchClient,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		LoadToken:  auth.LoadToken,
		PrintFlags: cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "n", "c"},
		Args:    validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		Short: "Create a new API key",
		Long: heredoc.Doc(`
			Create a new API key with the provided parameters.

			By default, the key is created for the current application through your
			signed-in session, so no admin API key is needed. This path supports
			--acl, --indices, --referers and --description.

			When the API key in use isn't the one the CLI provisioned for the current
			application (--api-key, ALGOLIA_API_KEY, or a key stored by a config.toml
			profile), the key is created with the Search API instead, which also
			supports --validity.
		`),
		Example: heredoc.Doc(`
			# Create a search-only API key for the current application
			$ algolia apikeys create --acl search

			# Create a new API key targeting the index "MOVIES", with the "search" and "browse" ACL and a description
			$ algolia apikeys create --indices MOVIES --acl search,browse --description "Search & Browse API Key"

			# Create a new API key targeting the indices "MOVIES" and "SERIES", with the "https://example.com" referer, with a validity of 1 hour and a description
			$ algolia apikeys create -i MOVIES,SERIES --acl search -r "https://example.com" --u 1h -d "Search-only API Key for MOVIES & SERIES" --api-key <admin-api-key>
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
		Duration (in seconds) after which the API key expires. By default (a value of 0), API keys don't expire.
		Requires an admin API key (--api-key), as expiring keys aren't supported for the signed-in session.`,
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

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runCreateCmd executes the create command
func runCreateCmd(opts *CreateOptions) error {
	if !config.ShouldUseSessionAPIKey(opts.config) {
		return runCreateWithSearchAPI(opts)
	}

	if opts.LoadToken() == nil {
		return runCreateWithSearchAPI(opts)
	}

	return runCreateWithDashboardAPI(opts)
}

func runCreateWithDashboardAPI(opts *CreateOptions) error {
	cs := opts.IO.ColorScheme()

	if len(opts.ACL) == 0 {
		return fmt.Errorf(
			"--acl is required, for example %s",
			cs.Bold("algolia apikeys create --acl search"),
		)
	}
	if opts.Validity != 0 {
		return fmt.Errorf(
			"--validity requires an admin API key: pass %s, or drop --validity",
			cs.Bold("--api-key"),
		)
	}

	appID, err := opts.config.Profile().GetApplicationID()
	if err != nil {
		return fmt.Errorf(
			"no application selected: run %s, or pass --application-id",
			cs.Bold("algolia application select"),
		)
	}

	description := opts.Description
	if description == "" {
		description = defaultDescription
	}

	params := dashboard.CreateAPIKeyRequest{
		ACL:         opts.ACL,
		Description: description,
		Indexes:     opts.Indices,
		Referers:    opts.Referers,
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	key, err := createKeyWithSession(opts, client, appID, params)
	if err != nil {
		if errors.Is(err, dashboard.ErrApplicationNotFound) {
			return fmt.Errorf(
				"application %s doesn't exist, or your account doesn't have access to it: run %s to pick one of your applications",
				cs.Bold(appID),
				cs.Bold("algolia application select"),
			)
		}
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, key)
	}

	printCreatedKey(opts.IO, key.Value)

	return nil
}

func printCreatedKey(io *iostreams.IOStreams, value string) {
	if !io.IsStdoutTTY() {
		fmt.Fprintln(io.Out, value)
		return
	}

	fmt.Fprintf(io.Out, "%s API key created: %s\n", io.ColorScheme().SuccessIcon(), value)
}

func createKeyWithSession(
	opts *CreateOptions,
	client *dashboard.Client,
	appID string,
	params dashboard.CreateAPIKeyRequest,
) (dashboard.APIKey, error) {
	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return dashboard.APIKey{}, err
	}

	opts.IO.StartProgressIndicatorWithLabel("Creating API key")
	key, err := client.CreateAPIKeyWithParams(accessToken, appID, params)
	opts.IO.StopProgressIndicator()
	if err == nil {
		return key, nil
	}

	accessToken, err = auth.ReauthenticateIfExpired(opts.IO, client, err)
	if err != nil {
		return dashboard.APIKey{}, err
	}

	opts.IO.StartProgressIndicatorWithLabel("Creating API key")
	key, err = client.CreateAPIKeyWithParams(accessToken, appID, params)
	opts.IO.StopProgressIndicator()

	return key, err
}

func runCreateWithSearchAPI(opts *CreateOptions) error {
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
		return searchAPICreateError(opts, err)
	}

	if opts.PrintFlags.HasStructuredOutput() {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, res)
	}

	printCreatedKey(opts.IO, res.Key)

	return nil
}

func searchAPICreateError(opts *CreateOptions, err error) error {
	var apiErr *search.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusForbidden {
		return err
	}

	cs := opts.IO.ColorScheme()
	if config.ShouldUseSessionAPIKey(opts.config) {
		return fmt.Errorf(
			"%w\nRun %s to create API keys without an admin key",
			err,
			cs.Bold("algolia auth login"),
		)
	}

	return fmt.Errorf(
		"%w\nThe API key in use isn't an admin key. Provide an admin key, or drop the key set through %s, %s or your profile to create the key with your signed-in session",
		err,
		cs.Bold("--api-key"),
		cs.Bold("ALGOLIA_API_KEY"),
	)
}
