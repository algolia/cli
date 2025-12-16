package add

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type apiKeyInspector interface {
	ListAPIKeys(opts ...search.RequestOption) (*search.ListApiKeysResponse, error)
	GetAPIKey(r search.ApiGetApiKeyRequest, opts ...search.RequestOption) (*search.GetApiKeyResponse, error)
	NewAPIGetAPIKeyRequest(key string) search.ApiGetApiKeyRequest
}

// searchClientAdapter adapts the Algolia search client to our apiKeyInspector interface
type searchClientAdapter struct {
	client *search.APIClient
}

func (a *searchClientAdapter) ListAPIKeys(opts ...search.RequestOption) (*search.ListApiKeysResponse, error) {
	return a.client.ListApiKeys(opts...)
}

func (a *searchClientAdapter) GetAPIKey(r search.ApiGetApiKeyRequest, opts ...search.RequestOption) (*search.GetApiKeyResponse, error) {
	return a.client.GetApiKey(r, opts...)
}

func (a *searchClientAdapter) NewAPIGetAPIKeyRequest(key string) search.ApiGetApiKeyRequest {
	return a.client.NewApiGetApiKeyRequest(key)
}

func inspectAPIKey(client apiKeyInspector, key string) (isAdmin bool, stringACLs []string, err error) {
	// Admin API keys are special: they can list keys but aren't themselves retrievable via GET /1/keys/{key}.
	// So we use ListAPIKeys() as the admin-key check and skip GetAPIKey() in that case.
	if _, err := client.ListAPIKeys(); err == nil {
		return true, nil, nil
	}

	apiKey, err := client.GetAPIKey(client.NewAPIGetAPIKeyRequest(key))
	if err != nil {
		return false, nil, errors.New("invalid application credentials")
	}
	if len(apiKey.Acl) == 0 {
		return false, nil, errors.New("the provided API key has no ACLs")
	}

	for _, a := range apiKey.Acl {
		stringACLs = append(stringACLs, string(a))
	}

	return false, stringACLs, nil
}

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
		Args:  validators.NoArgs(),
		Short: "Add a new profile configuration to the CLI",
		Example: heredoc.Doc(`
			# Add a new profile (interactive)
			$ algolia profile add

			# Add a new profile (non-interactive) and set it to default
			$ algolia profile add --name "my-profile" --app-id "my-app-id" --api-key "my-api-key" --default
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Interactive = true

			flags := cmd.Flags()
			nameProvided := flags.Changed("name")
			appIDProvided := flags.Changed("app-id")
			APIKeyProvided := flags.Changed("api-key")

			opts.Interactive = !(nameProvided && appIDProvided && APIKeyProvided)

			if opts.Interactive && !opts.IO.CanPrompt() {
				return cmdutil.FlagErrorf(
					"`--name`, `--app-id` and `--api-key` required when not running interactively",
				)
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
	cmd.Flags().
		StringVar(&opts.Profile.ApplicationID, "app-id", "", heredoc.Doc(`ID of the application.`))
	cmd.Flags().
		StringVar(&opts.Profile.APIKey, "api-key", "", heredoc.Doc(`API Key of the application.`))
	cmd.Flags().
		BoolVarP(&opts.Profile.Default, "default", "d", false, heredoc.Doc(`Set the profile as the default one.`))

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
				Validate: survey.ComposeValidators(
					survey.Required,
					validators.ProfileNameExists(opts.config),
				),
			},
			{
				Name: "applicationID",
				Prompt: &survey.Input{
					Message: "Application ID:",
					Default: opts.Profile.ApplicationID,
				},
				Validate: survey.ComposeValidators(
					survey.Required,
					validators.ApplicationIDExists(opts.config),
				),
			},
			{
				Name: "APIKey",
				Prompt: &survey.Input{
					Message: "Write API Key:",
					Default: opts.Profile.APIKey,
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

	client, err := search.NewClient(opts.Profile.ApplicationID, opts.Profile.APIKey)
	if err != nil {
		return err
	}
	adapter := &searchClientAdapter{client: client}
	isAdminAPIKey, stringACLs, err := inspectAPIKey(adapter, opts.Profile.APIKey)
	if err != nil {
		return err
	}

	// We should have at least the ACLs for a write key, otherwise warns the user, but still allows to add the profile.
	// If it's an admin API Key, we don't need to check ACLs, but we still warn the user.
	var warning string
	if !isAdminAPIKey {
		missingACLs := utils.Differences(auth.WriteAPIKeyDefaultACLs, stringACLs)
		if len(missingACLs) > 0 {
			warning = fmt.Sprintf(
				"%s The provided API key might be missing some ACLs: %s",
				opts.IO.ColorScheme().WarningIcon(),
				missingACLs,
			)
			warning += "\n  See https://www.algolia.com/doc/guides/security/api-keys/#rights-and-restrictions for more information."
			warning += "\n  You can still add the profile, but some commands might not be available."
		}
	} else {
		warning = fmt.Sprintf("%s The provided API key is an admin API key.", opts.IO.ColorScheme().WarningIcon())
		warning += "\n  You can still add the profile, but it would be better to use a more restricted key instead."
		warning += "\n  See https://www.algolia.com/doc/guides/security/security-best-practices/#keep-your-admin-api-key-confidential for more information."
	}

	if warning != "" {
		fmt.Printf("%s\n\n", warning)
		if opts.IO.CanPrompt() {
			var confirmed bool
			err := prompt.Confirm("Do you want to continue?", &confirmed)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}
		}
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
				extra = fmt.Sprintf(
					". Default profile changed from '%s' to '%s'.",
					cs.Bold(defaultProfile.Name),
					cs.Bold(opts.Profile.Name),
				)
			} else {
				extra = " and set as default."
			}
		}

		if _, err = fmt.Fprintf(opts.IO.Out, "%s Profile '%s' (%s) added successfully%s\n", cs.SuccessIcon(), opts.Profile.Name, opts.Profile.ApplicationID, extra); err != nil {
			return err
		}
	}

	return nil
}
