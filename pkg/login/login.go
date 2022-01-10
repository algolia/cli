package login

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/algolia/algolia-cli/pkg/config"
	"github.com/algolia/algolia-cli/pkg/iostreams"
)

// InteractiveLogin function is used to interactively ask the user for his application id / admin api key.
func InteractiveLogin(cfg *config.Config, io *iostreams.IOStreams) error {

	questions := []*survey.Question{
		{
			Name: "applicationID",
			Prompt: &survey.Input{
				Message: "Application ID:",
			},
			Validate: survey.Required,
		},
		{
			Name: "adminAPIKey",
			Prompt: &survey.Password{
				Message: "Admin API Key:",
			},
			Validate: survey.Required,
		},
	}
	survey.Ask(questions, &cfg.Profile)

	err := cfg.Profile.CreateProfile()

	if err != nil {
		return err
	}
	return nil
}
