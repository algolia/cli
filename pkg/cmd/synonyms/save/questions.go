package save

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/utils"
)

func AskSynonym(flags *shared.SynonymFlags) error {
	err := AskSynonymIdQuestion(flags)
	if err != nil {
		fmt.Println("Error when prompting synonym id survey question:", err)
	}

	err = AskSynonymTypeQuestion(flags)
	if err != nil {
		fmt.Println("Error when prompting synonym type survey question:", err)
	}

	switch flags.SynonymType {
	case "", shared.Regular:
		err := AskRegularSynonymQuestion(flags)
		if err != nil {
			fmt.Println("Error when prompting regular synonym survey question:", err)
			return err
		}
	case shared.OneWay:
		err := AskOneWaySynonymQuestions(flags)
		if err != nil {
			fmt.Println("Error when prompting one way synonym survey questions:", err)
			return err
		}
	case shared.Placeholder:
		err := AskPlaceholderSynonymQuestions(flags)
		if err != nil {
			fmt.Println("Error when prompting placeholder synonym survey questions:", err)
			return err
		}
	case shared.AltCorrection1, shared.AltCorrection2:
		err = AskAltCorrectionSynonymQuestions(flags)
		if err != nil {
			fmt.Println("Error when prompting alt correction synonym survey questions:", err)
			return err
		}
	default:
		return nil
	}

	return nil
}

func AskSynonymIdQuestion(flags *shared.SynonymFlags) error {
	return prompt.SurveyAsk([]*survey.Question{
		{
			Name: "synonymId",
			Prompt: &survey.Input{
				Message: "id:",
				Default: flags.SynonymID,
			},
			Validate: survey.Required,
		},
	}, flags)
}

func AskSynonymTypeQuestion(flags *shared.SynonymFlags) error {
	defaultType := flags.SynonymType
	if flags.SynonymType == "" {
		defaultType = shared.Regular
	}

	return prompt.SurveyAsk([]*survey.Question{{
		Name: "SynonymType",
		Prompt: &survey.Select{
			Message: "type:",
			Options: []string{shared.Regular, shared.OneWay, shared.Placeholder, shared.AltCorrection1, shared.AltCorrection2},
			Default: defaultType,
		},
		Validate: survey.Required,
	}}, flags)
}

func AskRegularSynonymQuestion(flags *shared.SynonymFlags) error {
	var synonyms string
	err := survey.AskOne(&survey.Input{
		Message: "synonyms (comma separated):",
		Default: utils.SliceToString(flags.Synonyms),
	}, &synonyms)
	if err != nil {
		return err
	}

	flags.Synonyms = utils.StringToSlice(synonyms)
	return nil
}

func AskOneWaySynonymQuestions(flags *shared.SynonymFlags) error {
	err := prompt.SurveyAsk([]*survey.Question{
		{
			Name: "SynonymInput",
			Prompt: &survey.Input{
				Message: "input:",
				Default: flags.SynonymInput,
			},
			Validate: survey.Required,
		},
	}, flags)
	if err != nil {
		return err
	}

	err = AskRegularSynonymQuestion(flags)
	if err != nil {
		return err
	}

	return nil
}

func AskPlaceholderSynonymQuestions(flags *shared.SynonymFlags) error {
	err := prompt.SurveyAsk([]*survey.Question{
		{
			Name: "SynonymPlaceholder",
			Prompt: &survey.Input{
				Message: "placeholder:",
				Default: flags.SynonymPlaceholder,
			},
			Validate: survey.Required,
		},
	}, flags)

	if err != nil {
		return err
	}

	var replacements string
	err = survey.AskOne(&survey.Input{
		Message: "replacements (comma separated):",
		Default: utils.SliceToString(flags.SynonymReplacements),
	}, &replacements)
	if err != nil {
		return err
	}

	flags.SynonymReplacements = utils.StringToSlice(replacements)
	return nil
}

func AskAltCorrectionSynonymQuestions(flags *shared.SynonymFlags) error {
	err := prompt.SurveyAsk([]*survey.Question{
		{
			Name: "SynonymWord",
			Prompt: &survey.Input{
				Message: "word:",
				Default: flags.SynonymWord,
			},
			Validate: survey.Required,
		}}, flags)
	if err != nil {
		return err
	}

	var corrections string
	err = survey.AskOne(&survey.Input{
		Message: "corrections (comma separated):",
		Default: utils.SliceToString(flags.SynonymCorrections),
	}, &corrections)
	if err != nil {
		return err
	}

	flags.SynonymCorrections = utils.StringToSlice(corrections)
	return nil
}
