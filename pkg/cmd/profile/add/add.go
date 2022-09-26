package add

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// AddOptions represents the options for the add command
type AddOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	Interactive bool

	Profile config.Profile
}

// NewAddCmd returns a new instance of AddCmd
func NewAddCmd(f *cmdutil.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}
	cmd := &cobra.Command{
		Use:   "add",
		Args:  validators.NoArgs,
		Short: "Add a new profile configuration to the CLI",
		Example: heredoc.Doc(`
			# Add a new profile (interactive)
			$ algolia profile add

			# Add a new profile (non-interactive) and set it to default
			$ algolia profile add --name "my-profile" --app-id "my-app-id" --admin-api-key "my-admin-api-key" --default
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Interactive = true

			flags := cmd.Flags()
			nameProvided := flags.Changed("name")
			appIDProvided := flags.Changed("app-id")
			adminAPIKeyProvided := flags.Changed("admin-api-key")

			opts.Interactive = !(nameProvided && appIDProvided && adminAPIKeyProvided)

			if opts.Interactive && !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf("`--name`, `--app-id` and `--admin-api-key` required when not running interactively")
			}

			if !opts.Interactive {
				err := validators.ProfileNameExists(opts.config)(opts.Profile.Name)
				if err != nil {
					return err
				}

				err = validators.ApplicationIDExists(opts.config)(opts.Profile.ApplicationID)
				if err != nil {
					return err
				}
			}

			if runF != nil {
				return runF(opts)
			}

			return runAddCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Profile.Name, "name", "n", "", heredoc.Doc(`Name of the profile.`))
	cmd.Flags().StringVar(&opts.Profile.ApplicationID, "app-id", "", heredoc.Doc(`ID of the application.`))
	cmd.Flags().StringVar(&opts.Profile.AdminAPIKey, "admin-api-key", "", heredoc.Doc(`Admin API Key of the application.`))
	cmd.Flags().BoolVarP(&opts.Profile.Default, "default", "d", false, heredoc.Doc(`Set the profile as the default one.`))

	return cmd
}

// runAddCmd executes the add command
func runAddCmd(opts *AddOptions) error {
	var defaultProfile *config.Profile
	for _, profile := range opts.config.ConfiguredProfiles() {
		if profile.Default {
			defaultProfile = profile
		}
	}

	if opts.Interactive {
		questions := []*survey.Question{
			{
				Name: "Name",
				Prompt: &survey.Input{
					Message: "Name:",
					Default: opts.Profile.Name,
				},
				Validate: survey.ComposeValidators(survey.Required, validators.ProfileNameExists(opts.config)),
			},
			{
				Name: "applicationID",
				Prompt: &survey.Input{
					Message: "Application ID:",
					Default: opts.Profile.ApplicationID,
				},
				Validate: survey.ComposeValidators(survey.Required, validators.ApplicationIDExists(opts.config)),
			},
			{
				Name: "adminAPIKey",
				Prompt: &survey.Input{
					Message: "Admin API Key:",
					Default: opts.Profile.AdminAPIKey,
				},
				Validate: survey.Required,
			},
			{
				Name: "default",
				Prompt: &survey.Confirm{
					Message: "Set as default profile?",
					Default: defaultProfile == nil || opts.Profile.Default,
				},
			},
		}
		err := prompt.SurveyAsk(questions, &opts.Profile)
		if err != nil {
			return err
		}
	}

	// Check if the application credentials are valid
	_, err := search.NewClient(opts.Profile.ApplicationID, opts.Profile.AdminAPIKey).ListAPIKeys()
	if err != nil {
		return errors.New("invalid application credentials")
	}

	err = opts.Profile.Add()
	if err != nil {
		return err
	}

	if opts.Profile.Default {
		err = opts.config.SetDefaultProfile(opts.Profile.Name)
		if err != nil {
			return err
		}
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		extra := "."

		if opts.Profile.Default {
			if defaultProfile != nil {
				extra = fmt.Sprintf(". Default profile changed from '%s' to '%s'.", cs.Bold(defaultProfile.Name), cs.Bold(opts.Profile.Name))
			} else {
				extra = " and set as default."
			}
		}

		fmt.Fprintf(opts.IO.Out, "%s Profile '%s' (%s) added successfully%s\n", cs.SuccessIcon(), opts.Profile.Name, opts.Profile.ApplicationID, extra)
	}

	return nil
}
