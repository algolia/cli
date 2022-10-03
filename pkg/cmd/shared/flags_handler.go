package handler

import (
	rules "github.com/algolia/cli/pkg/cmd/rules/shared"
	synonyms "github.com/algolia/cli/pkg/cmd/synonyms/shared"
	"github.com/spf13/cobra"
)

type FlagsHandler interface {
	Validate() error
	AskAndFill() error
}

func HandleFlags(handler FlagsHandler, interactive bool) error {
	err := handler.Validate()
	if interactive && err != nil {
		return handler.AskAndFill()
	}

	return err
}

// Synonyms

type SynonymHandler struct {
	Flags *synonyms.SynonymFlags
	Cmd   *cobra.Command
}

func (handler SynonymHandler) Validate() error {
	return ValidateSynonymFlags(*handler.Flags)
}

func (handler *SynonymHandler) AskAndFill() error {
	err := AskSynonym(handler.Flags, handler.Cmd)
	if err != nil {
		return err
	}

	return ValidateSynonymFlags(*handler.Flags)
}

// Rules

type RuleHandler struct {
	Flags *rules.RuleFlags
	Cmd   *cobra.Command
}

func (handler RuleHandler) Validate() error {
	return ValidateRuleFlags(*handler.Flags, GetRuleFlagsProvided(handler.Cmd))
}

func (handler *RuleHandler) AskAndFill() error {
	// TODO: interactive mode
	return nil
}
