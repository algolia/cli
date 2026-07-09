package status

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type StatusOptions struct {
	IO        *iostreams.IOStreams
	Config    config.IConfig
	LoadToken func() *auth.StoredToken
}

func NewStatusCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &StatusOptions{
		IO:        f.IOStreams,
		Config:    f.Config,
		LoadToken: auth.LoadToken,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show authentication and configuration status",
		Long: heredoc.Doc(`
			Show whether you're signed in, which application is current, and
			whether API credentials are available.

			This command is read-only: it never prompts, opens a browser, or
			modifies stored credentials. It exits with a non-zero status when
			no usable credentials are found.
		`),
		Example: heredoc.Doc(`
			# Check the authentication status
			$ algolia auth status
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusCmd(opts)
		},
	}

	return cmd
}

func runStatusCmd(opts *StatusOptions) error {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	sessionOK := false
	token := opts.LoadToken()
	switch {
	case token == nil:
		fmt.Fprintf(out, "%s Not signed in\n", cs.FailureIcon())
	case token.IsExpired() && token.RefreshToken == "":
		fmt.Fprintf(
			out,
			"%s Session expired for %s — run `algolia auth login` to sign in again\n",
			cs.FailureIcon(),
			cs.Bold(token.Email),
		)
	default:
		sessionOK = true
		who := token.Email
		if who == "" {
			who = token.UserID
		}
		fmt.Fprintf(out, "%s Signed in as %s\n", cs.SuccessIcon(), cs.Bold(who))
		if token.IsExpired() {
			fmt.Fprintf(
				out,
				"%s Access token expired; it refreshes automatically on the next command\n",
				cs.WarningIcon(),
			)
		}
	}

	appID, appErr := opts.Config.Profile().GetApplicationID()
	credentialsOK := false
	if appErr != nil {
		fmt.Fprintf(
			out,
			"%s No application selected — run `algolia application select`, or set ALGOLIA_APPLICATION_ID\n",
			cs.WarningIcon(),
		)
	} else {
		fmt.Fprintf(out, "%s Current application: %s\n", cs.SuccessIcon(), cs.Bold(appID))
		if _, keyErr := opts.Config.Profile().GetAPIKey(); keyErr != nil {
			fmt.Fprintf(out, "%s No API key available: %s\n", cs.WarningIcon(), keyErr)
		} else {
			credentialsOK = true
			fmt.Fprintf(out, "%s API key: available\n", cs.SuccessIcon())
		}
	}

	printEnvOverrides(opts.IO)

	if !sessionOK && !credentialsOK {
		fmt.Fprintf(
			out,
			"\nRun `algolia auth login` to sign in, or set ALGOLIA_APPLICATION_ID and ALGOLIA_API_KEY.\n",
		)
		return cmdutil.ErrSilent
	}
	return nil
}

func printEnvOverrides(io *iostreams.IOStreams) {
	cs := io.ColorScheme()

	if v := os.Getenv("ALGOLIA_APPLICATION_ID"); v != "" {
		fmt.Fprintf(
			io.Out,
			"%s ALGOLIA_APPLICATION_ID is set (%s) — it takes precedence over the selected application\n",
			cs.WarningIcon(),
			v,
		)
	}
	if os.Getenv("ALGOLIA_API_KEY") != "" {
		fmt.Fprintf(
			io.Out,
			"%s ALGOLIA_API_KEY is set — it takes precedence over the stored API key\n",
			cs.WarningIcon(),
		)
	}
	if os.Getenv("ALGOLIA_ADMIN_API_KEY") != "" {
		fmt.Fprintf(
			io.Out,
			"%s ALGOLIA_ADMIN_API_KEY is set (deprecated) — use ALGOLIA_API_KEY instead\n",
			cs.WarningIcon(),
		)
	}
	if v := os.Getenv("ALGOLIA_SEARCH_HOSTS"); v != "" {
		fmt.Fprintf(
			io.Out,
			"%s ALGOLIA_SEARCH_HOSTS is set — API requests go to: %s\n",
			cs.WarningIcon(),
			v,
		)
	}
}
