package open

import (
	"fmt"
	"sort"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/application/selectapp"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	pkgopen "github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/printers"
)

// resourceURL is a static shortcut that does not require sign-in.
type resourceURL struct {
	// Default is the absolute URL used when no application is configured.
	Default string
	// AppPath, when set, is the dashboard path used when an application is
	// configured. It is resolved against ALGOLIA_DASHBOARD_URL as
	// {base}/apps/{appID}/{AppPath}.
	AppPath string
}

var resourceURLs = map[string]resourceURL{
	"api":       {Default: "https://www.algolia.com/doc/api-reference/rest-api/"},
	"codex":     {Default: "https://www.algolia.com/developers/code-exchange/"},
	"cli-docs":  {Default: "https://algolia.com/doc/tools/cli/get-started/overview/"},
	"cli-repo":  {Default: "https://github.com/algolia/cli"},
	"devhub":    {Default: "https://www.algolia.com/developers/"},
	"docs":      {Default: "https://algolia.com/doc/"},
	"languages": {Default: "https://alg.li/supported-languages"},
	"status": {
		Default: "https://status.algolia.com/",
		AppPath: "monitoring/status",
	},
}

// dashboardTarget is an application dashboard page. These require sign-in and
// are scoped to the current application (selecting one if none is configured).
type dashboardTarget struct {
	// path is the dashboard path for the destination.
	path string
	// accountScoped marks account-level pages. They live at
	// {base}/{path}?applicationId={appID}, whereas application pages live at
	// {base}/apps/{appID}/{path}.
	accountScoped bool
}

var dashboardTargets = map[string]dashboardTarget{
	"dashboard":  {path: "dashboard"},
	"indices":    {path: "explorer/browse"},
	"crawler":    {path: "crawler"},
	"connectors": {path: "connectors"},
	"api-keys":   {path: "account/api-keys/all", accountScoped: true},
	"usage":      {path: "account/billing/usage", accountScoped: true},
	"team":       {path: "account/teams", accountScoped: true},
	"billing":    {path: "account/billing/details", accountScoped: true},
}

// targetNames returns every supported shortcut, sorted.
func targetNames() []string {
	names := make([]string, 0, len(resourceURLs)+len(dashboardTargets))
	for name := range resourceURLs {
		names = append(names, name)
	}
	for name := range dashboardTargets {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

// pageEntry describes an open shortcut for machine-readable output.
type pageEntry struct {
	Shortcut      string `json:"shortcut"`
	URL           string `json:"url"`
	RequiresLogin bool   `json:"requiresLogin"`
}

// OpenOptions represents the options for the open command. The function fields
// are injected so the flow can be exercised without a real OAuth session or
// browser.
type OpenOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	List     bool
	Shortcut string

	PrintFlags *cmdutil.PrintFlags

	Authenticate       func(*iostreams.IOStreams, *dashboard.Client) (string, error)
	SelectApplication  func() (*dashboard.Application, error)
	NewDashboardClient func(clientID string) *dashboard.Client
	Browser            func(string) error
}

func NewOpenCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &OpenOptions{
		IO:           f.IOStreams,
		config:       f.Config,
		PrintFlags:   cmdutil.NewPrintFlags(),
		Authenticate: auth.EnsureAuthenticated,
		SelectApplication: func() (*dashboard.Application, error) {
			return selectapp.Run(f)
		},
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		Browser: pkgopen.Browser,
	}

	cmd := &cobra.Command{
		Use:       "open <shortcut>",
		ValidArgs: targetNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return targetNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Open Algolia pages in your browser",
		Long: heredoc.Doc(`
			Open Algolia pages in your browser.

			Resource shortcuts (docs, API reference, status, …) open directly.

			Application pages (dashboard, indices, crawler, connectors, api-keys,
			usage, team, billing) are scoped to the current application: they
			require you to be signed in, and prompt you to select an application
			if none is configured.

			Run 'algolia open --list' to see every shortcut.

			With an output format (--output), the resolved page links are printed
			instead of opening a browser.
		`),
		Example: heredoc.Doc(`
			# List all shortcuts
			$ algolia open --list

			# List all shortcuts as JSON
			$ algolia open --list --output json

			# Open the documentation home page
			$ algolia open docs

			# Open the dashboard for the current application
			$ algolia open dashboard

			# Open billing / payment details for the current application
			$ algolia open billing

			# Print a page link as JSON instead of opening it
			$ algolia open billing --output json
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Shortcut = args[0]
			}
			return runOpenCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.List, "list", "l", false, "List all shortcuts")
	opts.PrintFlags.AddFlags(cmd)

	auth.DisableAuthCheck(cmd)

	return cmd
}

func runOpenCmd(opts *OpenOptions) error {
	listing := opts.List || opts.Shortcut == ""

	// With an output format, emit page metadata instead of opening a browser.
	if opts.structuredOutput() {
		return printTargets(opts, listing)
	}

	if listing {
		return listTargets(opts)
	}

	// Resource shortcuts open directly, without sign-in.
	if resource, ok := resourceURLs[opts.Shortcut]; ok {
		appID, _ := opts.config.Profile().GetApplicationID()
		url := resource.Default
		if appID != "" && resource.AppPath != "" {
			baseURL := opts.NewDashboardClient(auth.OAuthClientID()).DashboardURL
			url = fmt.Sprintf("%s/apps/%s/%s", baseURL, appID, resource.AppPath)
		}
		return opts.Browser(url)
	}

	// Application pages require sign-in and an application scope.
	if target, ok := dashboardTargets[opts.Shortcut]; ok {
		return openDashboardTarget(opts, target)
	}

	return fmt.Errorf("unsupported open command, given: %s", opts.Shortcut)
}

// structuredOutput reports whether an output format was requested via --output.
func (opts *OpenOptions) structuredOutput() bool {
	return opts.PrintFlags != nil &&
		opts.PrintFlags.OutputFlagSpecified != nil &&
		opts.PrintFlags.OutputFlagSpecified()
}

// printTargets renders page metadata with the configured printer. When listing,
// every shortcut is printed; otherwise only the requested shortcut is printed.
func printTargets(opts *OpenOptions, listing bool) error {
	printer, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	if listing {
		return printer.Print(opts.IO, opts.allEntries())
	}

	baseURL, appID, displayAppID := opts.resolveScope()
	entry, ok := entryFor(opts.Shortcut, baseURL, appID, displayAppID)
	if !ok {
		return fmt.Errorf("unsupported open command, given: %s", opts.Shortcut)
	}

	return printer.Print(opts.IO, entry)
}

// resolveScope returns the dashboard base URL and the application id used to
// build dashboard links. displayAppID falls back to a placeholder so links can
// be shown even when no application is configured.
func (opts *OpenOptions) resolveScope() (baseURL, appID, displayAppID string) {
	appID, _ = opts.config.Profile().GetApplicationID()
	displayAppID = appID
	if displayAppID == "" {
		displayAppID = "<app-id>"
	}
	baseURL = opts.NewDashboardClient(auth.OAuthClientID()).DashboardURL

	return baseURL, appID, displayAppID
}

// entryFor builds the page entry for a shortcut, or returns false if the
// shortcut is unknown.
func entryFor(name, baseURL, appID, displayAppID string) (pageEntry, bool) {
	if resource, ok := resourceURLs[name]; ok {
		url := resource.Default
		if appID != "" && resource.AppPath != "" {
			url = fmt.Sprintf("%s/apps/%s/%s", baseURL, appID, resource.AppPath)
		}
		return pageEntry{Shortcut: name, URL: url}, true
	}

	if target, ok := dashboardTargets[name]; ok {
		return pageEntry{
			Shortcut:      name,
			URL:           dashboardURL(baseURL, displayAppID, target),
			RequiresLogin: true,
		}, true
	}

	return pageEntry{}, false
}

// allEntries returns every shortcut, sorted by name.
func (opts *OpenOptions) allEntries() []pageEntry {
	baseURL, appID, displayAppID := opts.resolveScope()

	entries := make([]pageEntry, 0, len(resourceURLs)+len(dashboardTargets))
	for _, name := range targetNames() {
		if entry, ok := entryFor(name, baseURL, appID, displayAppID); ok {
			entries = append(entries, entry)
		}
	}

	return entries
}

// openDashboardTarget signs the user in, resolves the current application
// (selecting one if needed), then opens the dashboard page.
func openDashboardTarget(opts *OpenOptions, target dashboardTarget) error {
	client := opts.NewDashboardClient(auth.OAuthClientID())
	if _, err := opts.Authenticate(opts.IO, client); err != nil {
		return err
	}

	appID, err := opts.config.Profile().GetApplicationID()
	if err != nil {
		app, selErr := opts.SelectApplication()
		if selErr != nil {
			return selErr
		}
		if app == nil {
			// No application is available to scope to; the selection flow has
			// already explained the situation to the user.
			return nil
		}
		appID = app.ID
	}

	// The base URL is resolved from ALGOLIA_DASHBOARD_URL by the dashboard
	// client (falling back to its compiled-in default).
	url := dashboardURL(client.DashboardURL, appID, target)

	cs := opts.IO.ColorScheme()
	fmt.Fprintf(opts.IO.Out, "Opening %s\n", cs.Bold(url))

	return opts.Browser(url)
}

// dashboardURL builds the dashboard URL for an application page. Application
// pages are scoped via the /apps/{appID} path; account pages carry the
// application in an applicationId query parameter.
func dashboardURL(baseURL, appID string, target dashboardTarget) string {
	if target.accountScoped {
		return fmt.Sprintf("%s/%s?applicationId=%s", baseURL, target.path, appID)
	}

	return fmt.Sprintf("%s/apps/%s/%s", baseURL, appID, target.path)
}

func listTargets(opts *OpenOptions) error {
	fmt.Fprintln(
		opts.IO.Out,
		"open quickly opens Algolia pages. To use, run 'algolia open <shortcut>'.",
	)
	fmt.Fprintln(opts.IO.Out, "open supports the following shortcuts:")
	fmt.Fprintln(opts.IO.Out)

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("SHORTCUT", nil, nil)
		table.AddField("URL", nil, nil)
		table.EndRow()
	}

	for _, entry := range opts.allEntries() {
		table.AddField(entry.Shortcut, nil, nil)
		table.AddField(entry.URL, nil, nil)
		table.EndRow()
	}

	return table.Render()
}
