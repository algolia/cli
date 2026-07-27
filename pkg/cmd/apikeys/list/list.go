package list

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
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
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

// nowFn exists to make time-based output deterministic in tests.
var nowFn = time.Now

var reauthenticate = auth.ReauthenticateIfExpired

var tableHeaders = []string{
	"KEY",
	"DESCRIPTION",
	"ACL",
	"INDICES",
	"VALIDITY",
	"MAX HITS PER QUERY",
	"MAX QUERIES PER IP PER HOUR",
	"REFERERS",
	"CREATED AT",
}

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient       func() (*search.APIClient, error)
	NewDashboardClient func(clientID string) *dashboard.Client

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd creates and returns a list command for API Keys.
func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		PrintFlags: cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Args:    validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		Short: "Lists all API keys associated with your Algolia application, including their permissions and restrictions.",
		Long: heredoc.Doc(`
			Lists all API keys associated with your Algolia application, including their permissions and restrictions.

			By default, the keys of the current application are listed through your
			signed-in session, so no admin API key is needed. This only covers the keys
			created by the CLI, and doesn't report an expiry. Keys you don't have the
			rights to create are listed without their value.

			Every key of the application is listed with the Search API instead
			whenever the API key in use isn't the one the CLI provisioned for the
			current application: --api-key, ALGOLIA_API_KEY, a key stored by a
			config.toml profile, or a key kept in your keychain that the CLI didn't
			create.

			--admin-api-key and ALGOLIA_ADMIN_API_KEY are ignored while an application
			is selected.
		`),
		Example: heredoc.Doc(`
			# List the API keys the CLI created for the current application
			$ algolia apikeys list

			# List every API key of the application
			$ algolia apikeys list --api-key <admin-api-key>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runListCmd executes the list command
func runListCmd(opts *ListOptions) error {
	if !config.ShouldUseSessionAPIKey(opts.Config) {
		return runListWithSearchAPI(opts)
	}

	return runListWithSessionAPI(opts)
}

func structuredPrinter(opts *ListOptions) (printers.Printer, error) {
	if !opts.PrintFlags.HasStructuredOutput() {
		return nil, nil
	}

	return opts.PrintFlags.ToPrinter()
}

func runListWithSessionAPI(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	printer, err := structuredPrinter(opts)
	if err != nil {
		return err
	}

	appID, err := opts.Config.Profile().GetApplicationID()
	if err != nil {
		return fmt.Errorf(
			"no application selected: run %s, or pass --application-id",
			cs.Bold("algolia application select"),
		)
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	keys, err := listKeysWithSession(opts, client, appID)
	if err != nil {
		if errors.Is(err, dashboard.ErrEndpointNotAvailable) {
			return fmt.Errorf(
				"listing API keys with your signed-in session needs a newer Algolia API version than the one answering: pass %s with an admin key in the meantime",
				cs.Bold("--api-key"),
			)
		}
		if errors.Is(err, dashboard.ErrApplicationNotFound) {
			return fmt.Errorf(
				"application %s doesn't exist, or your account doesn't have access to it: run %s to pick one of your applications",
				cs.Bold(appID),
				cs.Bold("algolia application select"),
			)
		}
		return err
	}

	if printer != nil {
		return printer.Print(opts.IO, keys)
	}

	now := nowFn()
	rows := make([][]string, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, []string{
			formatKeyValue(key.Value),
			key.Description,
			fmt.Sprintf("%v", key.ACL),
			fmt.Sprintf("%v", key.Indexes),
			"-",
			formatLimit(key.MaxHitsPerQuery),
			formatLimit(key.MaxQueriesPerIPPerHour),
			fmt.Sprintf("%v", key.Referers),
			formatCreatedAt(now, key.CreatedAt),
		})
	}

	return renderTable(opts.IO, rows)
}

func listKeysWithSession(
	opts *ListOptions,
	client *dashboard.Client,
	appID string,
) ([]dashboard.APIKey, error) {
	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return nil, err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching API Keys")
	keys, err := client.ListAPIKeys(accessToken, appID)
	opts.IO.StopProgressIndicator()
	if err == nil {
		return keys, nil
	}

	accessToken, err = reauthenticate(opts.IO, client, err)
	if err != nil {
		return nil, err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching API Keys")
	keys, err = client.ListAPIKeys(accessToken, appID)
	opts.IO.StopProgressIndicator()

	return keys, err
}

func runListWithSearchAPI(opts *ListOptions) error {
	printer, err := structuredPrinter(opts)
	if err != nil {
		return err
	}

	client, err := opts.SearchClient()
	if err != nil {
		return auth.WithRemediation(err)
	}

	now := nowFn()

	opts.IO.StartProgressIndicatorWithLabel("Fetching API Keys")
	res, err := client.ListApiKeys()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return searchAPIListError(opts, err)
	}

	if printer != nil {
		return printer.Print(opts.IO, res)
	}

	// Sort API Keys by createdAt
	sort.Slice(res.Keys, func(i, j int) bool {
		return res.Keys[i].CreatedAt > res.Keys[j].CreatedAt
	})

	rows := make([][]string, 0, len(res.Keys))
	for _, key := range res.Keys {
		description := ""
		if key.Description != nil {
			description = *key.Description
		}

		rows = append(rows, []string{
			formatKeyValue(key.Value),
			description,
			fmt.Sprintf("%v", key.Acl),
			fmt.Sprintf("%v", key.Indexes),
			formatValidity(now, key.Validity),
			formatLimit(intFromInt32(key.MaxHitsPerQuery)),
			formatLimit(intFromInt32(key.MaxQueriesPerIPPerHour)),
			fmt.Sprintf("%v", key.Referers),
			humanize.RelTime(now, time.Unix(key.CreatedAt, 0), "from now", "ago"),
		})
	}

	return renderTable(opts.IO, rows)
}

func renderTable(io *iostreams.IOStreams, rows [][]string) error {
	table := printers.NewTablePrinter(io)
	if table.IsTTY() {
		for _, header := range tableHeaders {
			table.AddField(header, nil, nil)
		}
		table.EndRow()
	}

	for _, row := range rows {
		for _, field := range row {
			table.AddField(field, nil, nil)
		}
		table.EndRow()
	}

	return table.Render()
}

func formatValidity(now time.Time, validity *int32) string {
	if validity == nil || *validity == 0 {
		return "Never expire"
	}

	duration := time.Duration(*validity) * time.Second

	return humanize.RelTime(now, now.Add(duration), "from now", "ago")
}

func formatKeyValue(value string) string {
	if value == "" {
		return "-"
	}

	return value
}

func formatLimit(limit *int) string {
	if limit == nil || *limit == 0 {
		return "0"
	}

	return humanize.Comma(int64(*limit))
}

func formatCreatedAt(now time.Time, createdAt string) string {
	if createdAt == "" {
		return "-"
	}

	parsed, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return createdAt
	}

	return humanize.RelTime(now, parsed, "from now", "ago")
}

func intFromInt32(value *int32) *int {
	if value == nil {
		return nil
	}

	converted := int(*value)

	return &converted
}

func searchAPIListError(opts *ListOptions, err error) error {
	var apiErr *search.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusForbidden {
		return err
	}

	cs := opts.IO.ColorScheme()
	if config.ShouldUseSessionAPIKey(opts.Config) {
		return fmt.Errorf(
			"%w\nRun %s to list API keys without an admin key",
			err,
			cs.Bold("algolia auth login"),
		)
	}

	return fmt.Errorf(
		"%w\nThe API key in use isn't an admin key. Provide an admin key, or drop the key set through %s, %s or your profile to list the keys the CLI created for your signed-in session",
		err,
		cs.Bold("--api-key"),
		cs.Bold("ALGOLIA_API_KEY"),
	)
}
