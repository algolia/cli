package ask

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/pkg/utils"
)

// https://github.com/AlecAivazis/survey#custom-types
type StringSlice struct {
	value []string
}

func (my *StringSlice) WriteAnswer(name string, value interface{}) error {
	my.value = utils.StringToSlice(value.(string))
	return nil
}

func AskCommaSeparatedInputQuestion(
	message string,
	storage *[]string,
	defaultValues []string,
	opts ...survey.AskOpt,
) error {
	stringSlice := StringSlice{}
	err := survey.AskOne(
		&survey.Input{
			Message: message,
			Default: utils.SliceToString(defaultValues),
		},
		&stringSlice,
	)
	*storage = stringSlice.value

	return err
}

func AskMultiSelectQuestion(
	message string,
	defaultValues []string,
	storage *[]string,
	options []string,
	opts ...survey.AskOpt,
) error {
	err := survey.AskOne(
		&survey.MultiSelect{
			Message: message,
			Default: defaultValues,
			Options: options,
		},
		storage,
	)

	return err
}

func AskSelectQuestion(
	message string,
	storage *string,
	options []string,
	defaultValue string,
	opts ...survey.AskOpt,
) error {
	return survey.AskOne(&survey.Select{
		Message: message,
		Options: options,
		Default: defaultValue,
	}, storage, opts...)
}

func AskInputQuestion(
	message string,
	storage *string,
	defaultValue string,
	opts ...survey.AskOpt,
) error {
	return survey.AskOne(&survey.Input{
		Message: message,
		Default: defaultValue,
	}, storage, opts...)
}

func AskInputQuestionWithSuggestion(
	message string,
	storage *string,
	defaultValue string,
	suggest func(toComplete string) []string,
	opts ...survey.AskOpt,
) error {
	return survey.AskOne(&survey.Input{
		Message: message,
		Default: defaultValue,
		Suggest: suggest,
	}, storage, opts...)
}

func AskBooleanQuestion(
	message string,
	storage *bool,
	defaultValue bool,
	opts ...survey.AskOpt,
) error {
	return survey.AskOne(&survey.Confirm{
		Message: message,
		Default: defaultValue,
	}, storage, opts...)
}
