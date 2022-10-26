package synonms

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/ask"
	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
)

func ValidateSynonymFlags(flags shared.SynonymFlags) error {
	if flags.SynonymID == "" {
		return fmt.Errorf("a unique synonym id is required")
	}

	switch flags.SynonymType {
	case shared.OneWay:
		if len(flags.Synonyms) < 1 {
			return fmt.Errorf("at least 1 synonym is required")
		}
		if flags.SynonymInput == "" {
			return fmt.Errorf("a synonym input is required for one way synonyms")
		}
	case shared.AltCorrection1:
		if flags.SynonymWord == "" {
			return fmt.Errorf("synonym word is required for alt correction 1 synonyms")
		}
		if len(flags.SynonymCorrections) < 1 {
			return fmt.Errorf("synonym corrections are required for alt correction 1 synonyms")
		}
	case shared.AltCorrection2:
		if flags.SynonymWord == "" {
			return fmt.Errorf("synonym word is required for alt correction 2 synonyms")
		}
		if len(flags.SynonymCorrections) < 1 {
			return fmt.Errorf("synonym corrections are required for alt correction 2 synonyms")
		}
	case shared.Placeholder:
		if flags.SynonymPlaceholder == "" {
			return fmt.Errorf("a synonym placeholder is required for placeholder synonyms")
		}
		if len(flags.SynonymReplacements) < 1 {
			return fmt.Errorf("synonym replacements are required for placeholder synonyms")
		}
	case "", shared.Regular:
		if len(flags.Synonyms) < 1 {
			return fmt.Errorf("at least 1 synonym is required")
		}
	}

	return nil
}

type FlagsProvided struct {
	idProvided, typeProvided, synonymsProvided, inputProvided, wordProvided, placeholderProvided, correctionsProvided, replacementsProvided bool
}

func AskSynonym(flags *shared.SynonymFlags, cmd *cobra.Command) error {
	flagsProvided := FlagsProvided{
		idProvided:           cmd.Flags().Changed("id"),
		typeProvided:         cmd.Flags().Changed("type"),
		synonymsProvided:     cmd.Flags().Changed("synonyms"),
		inputProvided:        cmd.Flags().Changed("input"),
		wordProvided:         cmd.Flags().Changed("word"),
		placeholderProvided:  cmd.Flags().Changed("placeholder"),
		correctionsProvided:  cmd.Flags().Changed("corrections"),
		replacementsProvided: cmd.Flags().Changed("repalcements"),
	}

	err := AskSynonymIdQuestion(flags, flagsProvided)
	if err != nil {
		return err
	}

	err = AskSynonymTypeQuestion(flags, flagsProvided)
	if err != nil {
		return err
	}

	switch flags.SynonymType {
	case "", shared.Regular:
		return AskRegularSynonymQuestion(flags, flagsProvided)
	case shared.OneWay:
		return AskOneWaySynonymQuestions(flags, flagsProvided)
	case shared.Placeholder:
		return AskPlaceholderSynonymQuestions(flags, flagsProvided)
	case shared.AltCorrection1, shared.AltCorrection2:
		return AskAltCorrectionSynonymQuestions(flags, flagsProvided)
	default:
		return fmt.Errorf("wrong synonym type")
	}
}

func AskSynonymIdQuestion(flags *shared.SynonymFlags, flagsProvided FlagsProvided) error {
	if flagsProvided.idProvided {
		return nil
	}
	return ask.AskInputQuestion("id:", &flags.SynonymID, flags.SynonymID, survey.WithValidator(survey.Required))
}

func AskSynonymTypeQuestion(flags *shared.SynonymFlags, flagsProvided FlagsProvided) error {
	if flagsProvided.typeProvided {
		return nil
	}

	defaultType := flags.SynonymType
	if flags.SynonymType == "" {
		defaultType = shared.Regular
	}

	return ask.AskSelectQuestion(
		"type:",
		&flags.SynonymType,
		[]string{shared.Regular, shared.OneWay, shared.Placeholder, shared.AltCorrection1, shared.AltCorrection2},
		defaultType,
		survey.WithValidator(survey.Required),
	)
}

func AskRegularSynonymQuestion(flags *shared.SynonymFlags, flagsProvided FlagsProvided) error {
	if flagsProvided.synonymsProvided {
		return nil
	}

	return ask.AskCommaSeparatedInputQuestion(
		"synonyms (comma separated):",
		&flags.Synonyms,
		flags.Synonyms,
		survey.WithValidator(survey.Required),
	)
}

func AskOneWaySynonymQuestions(flags *shared.SynonymFlags, flagsProvided FlagsProvided) error {
	if !flagsProvided.inputProvided {
		err := ask.AskInputQuestion("input:", &flags.SynonymInput, flags.SynonymInput, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
	}

	return AskRegularSynonymQuestion(flags, flagsProvided)
}

func AskPlaceholderSynonymQuestions(flags *shared.SynonymFlags, flagsProvided FlagsProvided) error {
	if !flagsProvided.placeholderProvided {
		err := ask.AskInputQuestion("placeholder:", &flags.SynonymPlaceholder, flags.SynonymPlaceholder, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
	}
	if !flagsProvided.replacementsProvided {
		return ask.AskCommaSeparatedInputQuestion(
			"replacements (comma separated):",
			&flags.SynonymReplacements,
			flags.SynonymReplacements,
			survey.WithValidator(survey.Required),
		)
	}

	return nil
}

func AskAltCorrectionSynonymQuestions(flags *shared.SynonymFlags, flagsProvided FlagsProvided) error {
	if !flagsProvided.wordProvided {
		err := ask.AskInputQuestion("word:", &flags.SynonymWord, flags.SynonymWord, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
	}

	if !flagsProvided.replacementsProvided {
		return ask.AskCommaSeparatedInputQuestion(
			"corrections (comma separated):",
			&flags.SynonymCorrections,
			flags.SynonymCorrections,
			survey.WithValidator(survey.Required),
		)
	}

	return nil
}
