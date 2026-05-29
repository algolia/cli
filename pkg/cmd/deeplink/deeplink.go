package deeplink

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmd/application/selectapp"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// purposeTarget describes how to build the dashboard URL for a --purpose value.
type purposeTarget struct {
	// path is the dashboard path for the destination.
	path string
	// accountScoped marks account-level pages. They live at
	// {base}/{path}?applicationId={appID}, whereas application pages live at
	// {base}/apps/{appID}/{path}.
	accountScoped bool
}

// purposeTargets maps each --purpose value to its dashboard destination.
var purposeTargets = map[string]purposeTarget{
	"dashboard":  {path: "dashboard"},
	"indices":    {path: "explorer/browse"},
	"crawler":    {path: "crawler"},
	"connectors": {path: "connectors"},
	"api-keys":   {path: "account/api-keys/all", accountScoped: true},
	"usage":      {path: "account/billing/usage", accountScoped: true},
	"team":       {path: "account/teams", accountScoped: true},
	"billing":    {path: "account/billing/details", accountScoped: true},
}

// purposeOrder controls the display order for the interactive picker, the
// flag help text, and shell completion.
var purposeOrder = []string{
	"dashboard",
	"indices",
	"crawler",
	"connectors",
	"api-keys",
	"usage",
	"team",
	"billing",
}

// DeeplinkOptions holds everything the deeplink command needs. The function
// fields are injected so the flow can be exercised without a real OAuth
// session or browser.
type DeeplinkOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	Purpose string

	Authenticate       func(*iostreams.IOStreams, *dashboard.Client) (string, error)
	SelectApplication  func() (*dashboard.Application, error)
	NewDashboardClient func(clientID string) *dashboard.Client
	Browser            func(string) error
}

func NewDeeplinkCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &DeeplinkOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		Authenticate: auth.EnsureAuthenticated,
		SelectApplication: func() (*dashboard.Application, error) {
			return selectapp.Run(f)
		},
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
		Browser: open.Browser,
	}

	cmd := &cobra.Command{
		Use:   "deeplink",
		Short: "Open a dashboard page for the current application",
		Long: heredoc.Doc(`
			Open a specific page of the Algolia dashboard in your browser,
			scoped to the current application.`),
		Example: heredoc.Doc(`
			# Choose a destination from a list
			$ algolia deeplink

			# Open the API keys page for the current application
			$ algolia deeplink --purpose api-keys

			# Open billing / payment details
			$ algolia deeplink --purpose billing
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			// The command manages its own sign-in and application resolution.
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeeplinkCmd(opts)
		},
	}

	cmd.Flags().StringVar(
		&opts.Purpose,
		"purpose",
		"",
		fmt.Sprintf("Dashboard destination to open (%s)", strings.Join(purposeOrder, ", ")),
	)
	_ = cmd.RegisterFlagCompletionFunc(
		"purpose",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return purposeOrder, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func runDeeplinkCmd(opts *DeeplinkOptions) error {
	// Resolve the destination first so invalid or missing input fails before
	// any sign-in or browser side effects.
	purpose, err := opts.resolvePurpose()
	if err != nil {
		return err
	}

	// Require a valid sign-in even when an application is already configured.
	client := opts.NewDashboardClient(auth.OAuthClientID())
	if _, err := opts.Authenticate(opts.IO, client); err != nil {
		return err
	}

	appID, err := opts.Config.Profile().GetApplicationID()
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
	url := deeplinkURL(client.DashboardURL, appID, purpose)

	cs := opts.IO.ColorScheme()
	fmt.Fprintf(opts.IO.Out, "Opening %s\n", cs.Bold(url))

	return opts.Browser(url)
}

// resolvePurpose validates an explicit --purpose value or, when omitted,
// prompts for one interactively. In non-interactive mode without a value it
// returns a flag error listing the valid destinations.
func (opts *DeeplinkOptions) resolvePurpose() (string, error) {
	if opts.Purpose != "" {
		if _, ok := purposeTargets[opts.Purpose]; !ok {
			return "", cmdutil.FlagErrorf(
				"invalid purpose %q: must be one of %s",
				opts.Purpose,
				strings.Join(purposeOrder, ", "),
			)
		}
		return opts.Purpose, nil
	}

	if !opts.IO.CanPrompt() {
		return "", cmdutil.FlagErrorf(
			"--purpose is required in non-interactive mode: must be one of %s",
			strings.Join(purposeOrder, ", "),
		)
	}

	var selected int
	err := prompt.SurveyAskOne(
		&survey.Select{
			Message: "Open which dashboard page?",
			Options: purposeOrder,
		},
		&selected,
	)
	if err != nil {
		return "", err
	}

	return purposeOrder[selected], nil
}

// deeplinkURL builds the dashboard URL for an application and purpose, using
// baseURL (resolved from ALGOLIA_DASHBOARD_URL) as the host. Application pages
// are scoped via the /apps/{appID} path; account pages carry the application
// in an applicationId query parameter.
func deeplinkURL(baseURL, appID, purpose string) string {
	target := purposeTargets[purpose]
	if target.accountScoped {
		return fmt.Sprintf("%s/%s?applicationId=%s", baseURL, target.path, url.QueryEscape(appID))
	}

	return fmt.Sprintf("%s/apps/%s/%s", baseURL, appID, target.path)
}
