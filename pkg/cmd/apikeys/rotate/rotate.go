package rotate

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// RotateOptions represents the options for the rotate command.
type RotateOptions struct {
	IO     *iostreams.IOStreams
	Config config.IConfig

	NewDashboardClient func(clientID string) *dashboard.Client
}

// NewRotateCmd returns a new instance of the rotate command.
func NewRotateCmd(f *cmdutil.Factory, runF func(*RotateOptions) error) *cobra.Command {
	opts := &RotateOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		NewDashboardClient: func(clientID string) *dashboard.Client {
			return dashboard.NewClient(clientID)
		},
	}

	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate the CLI-managed API key for the current application",
		Long: heredoc.Doc(`
			Rotate (regenerate) the CLI-managed API key for the current application.

			The previous key is invalidated and replaced by a new one, which is then
			stored for the current application.
		`),
		Example: heredoc.Doc(`
			# Rotate the current application's CLI-managed key
			$ algolia apikeys rotate
		`),
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"skipAuthCheck": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runRotateCmd(opts)
		},
	}

	return cmd
}

// runRotateCmd executes the rotate command.
func runRotateCmd(opts *RotateOptions) error {
	cs := opts.IO.ColorScheme()

	appID := opts.Config.ActiveApplicationID()
	if appID == "" {
		return fmt.Errorf(
			"no current application selected; run %s first",
			cs.Bold("algolia application select"),
		)
	}

	keyUUID, hasStored := opts.Config.APIKeyUUID(appID)
	if !hasStored {
		return fmt.Errorf(
			"no CLI-managed API key found for application %s; run %s to regenerate one",
			cs.Bold(appID),
			cs.Bold("algolia application select"),
		)
	}

	client := opts.NewDashboardClient(auth.OAuthClientID())

	accessToken, err := auth.EnsureAuthenticated(opts.IO, client)
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Rotating API key")
	created, err := client.RotateAPIKey(accessToken, appID, keyUUID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		newToken, reAuthErr := auth.ReauthenticateIfExpired(opts.IO, client, err)
		if reAuthErr != nil {
			return reAuthErr
		}
		accessToken = newToken
		opts.IO.StartProgressIndicatorWithLabel("Rotating API key")
		created, err = client.RotateAPIKey(accessToken, appID, keyUUID)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}
	}

	// We rotated the application's own CLI-managed key, so the previous value is
	// now invalid: persist the new one for the current application.
	if err := opts.Config.SaveApplication(appID, "", keyUUID, created.Value, false); err != nil {
		return fmt.Errorf("API key rotated but could not be saved locally: %w", err)
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s API key rotated: %s\n", cs.SuccessIcon(), created.Value)
	}

	return nil
}
