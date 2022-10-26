package handler

import (
	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
	"github.com/spf13/cobra"
)

type FlagsHandler interface {
	Validate() error
	AskAndFill() error
}

// Synonyms

type SynonymHandler struct {
	Flags *shared.SynonymFlags
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

func HandleFlags(handler FlagsHandler, interactive bool) error {
	err := handler.Validate()
	if interactive && err != nil {
		return handler.AskAndFill()
	}

	return err
}
