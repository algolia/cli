package ask

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/algolia/cli/pkg/utils"
)

func AskInputQuestion(message string, storage *string, defaultValue string, opts ...survey.AskOpt) error {
	return survey.AskOne(&survey.Input{
		Message: message,
		Default: defaultValue,
	}, storage, opts...)
}

func AskCommaSeparatedInputQuestion(message string, storage *[]string, defaultValues []string, opts ...survey.AskOpt) error {
	values := ""

	err := survey.AskOne(&survey.Input{
		Message: message,
		Default: utils.SliceToString(defaultValues),
	}, &values, opts...)
	if err != nil {
		return err
	}

	*storage = utils.StringToSlice(values)
	return nil
}

func AskSelectQuestion(message string, storage *string, options []string, defaultValue string, opts ...survey.AskOpt) error {
	return survey.AskOne(&survey.Select{
		Message: message,
		Options: options,
		Default: defaultValue,
	}, storage, opts...)
}
