package add

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// AddOptions represents the options for the add command
type AddOptions struct {
	config *config.Config
	IO     *iostreams.IOStreams

	Interactive bool

	Name        string
	ID          string
	AdminAPIKey string
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
			$ algolia app add
			$ algolia app add --name "my-app" --app-id "my-app-id" --admin-api-key "my-admin-api-key"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			flags := cmd.Flags()
			if flags.Changed("name") || flags.Changed("app-id") || flags.Changed("admin-api-key") {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("`--name`, `--app-id` and `--admin-api-key` required when not running interactively")
				}
				opts.Interactive = true
			}

			return runAddCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.config.App.Name, "name", "n", "default", heredoc.Doc(`Name of the application.`))
	cmd.Flags().StringVar(&opts.config.App.ID, "app-id", "", heredoc.Doc(`ID of the application.`))
	cmd.Flags().StringVar(&opts.config.App.AdminAPIKey, "admin-api-key", "", heredoc.Doc(`Admin API Key of the application.`))

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
					Default: opts.config.App.Name,
				},
				Validate: survey.Required,
			},
			{
				Name: "ID",
				Prompt: &survey.Input{
					Message: "Application ID:",
					Default: opts.config.App.ID,
				},
				Validate: survey.Required,
			},
			{
				Name: "adminAPIKey",
				Prompt: &survey.Input{
					Message: "Admin API Key:",
					Default: opts.config.App.AdminAPIKey,
				},
				Validate: survey.Required,
			},
		}
		survey.Ask(questions, &opts.config.App)
	}

	// Check if the application credentials are valid
	if err := opts.config.App.Validate(); err != nil {
		return err
	}

	err := opts.config.App.AddApp()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Application `%s (%s)` successfuly added\n", cs.SuccessIcon(), opts.config.App.Name, opts.config.App.ID)
	}

	return nil
}
