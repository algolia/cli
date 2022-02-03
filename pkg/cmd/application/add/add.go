package add

import (
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
	config *config.Config
	IO     *iostreams.IOStreams

	Interactive bool

	Application config.Application
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
		Short: "Add a new application",
		Long:  `Add a new application configuration to the CLI.`,
		Example: heredoc.Doc(`
			# Add a new application (interactive)
			$ algolia application add

			# Add a new application (non-interactive)
			$ algolia application add --name "my-app" --app-id "my-app-id" --admin-api-key "my-admin-api-key"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			opts.Interactive = true

			flags := cmd.Flags()
			nameProvided := flags.Changed("name")
			appIDProvided := flags.Changed("app-id")
			adminAPIKeyProvided := flags.Changed("admin-api-key")

			opts.Interactive = !(nameProvided && appIDProvided && adminAPIKeyProvided)

			if opts.Interactive && !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf("`--name`, `--app-id` and `--admin-api-key` required when not running interactively")
			}

			return runAddCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Application.Name, "name", "n", "default", heredoc.Doc(`Name of the application.`))
	cmd.Flags().StringVar(&opts.Application.ID, "app-id", "", heredoc.Doc(`ID of the application.`))
	cmd.Flags().StringVar(&opts.Application.AdminAPIKey, "admin-api-key", "", heredoc.Doc(`Admin API Key of the application.`))

	return cmd
}

// runAddCmd executes the add command
func runAddCmd(opts *AddOptions) error {
	if opts.Interactive {
		questions := []*survey.Question{
			{
				Name: "Name",
				Prompt: &survey.Input{
					Message: "Name:",
					Default: opts.Application.Name,
				},
				Validate: survey.Required,
			},
			{
				Name: "ID",
				Prompt: &survey.Input{
					Message: "Application ID:",
					Default: opts.Application.ID,
				},
				Validate: survey.Required,
			},
			{
				Name: "adminAPIKey",
				Prompt: &survey.Input{
					Message: "Admin API Key:",
					Default: opts.Application.AdminAPIKey,
				},
				Validate: survey.Required,
			},
		}
		err := prompt.SurveyAsk(questions, &opts.Application)
		if err != nil {
			return err
		}
	}

	// Check if the application credentials are valid
	_, err := search.NewClient(opts.Application.ID, opts.Application.AdminAPIKey).ListAPIKeys()
	if err != nil {
		return fmt.Errorf("invalid application credentials: %s", err)
	}

	err = opts.config.AddApplication(&opts.Application)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Application '%s' (%s) successfuly added\n", cs.SuccessIcon(), opts.Application.Name, opts.Application.ID)
	}

	return nil
}
