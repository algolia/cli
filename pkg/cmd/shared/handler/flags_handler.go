package handler

import (
	config "github.com/algolia/cli/pkg/cmd/shared/handler/indices"
	synonms "github.com/algolia/cli/pkg/cmd/shared/handler/synonyms"
	"github.com/algolia/cli/pkg/cmd/synonyms/shared"
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

// `synonyms save`
type SynonymHandler struct {
	Flags *shared.SynonymFlags
	Cmd   *cobra.Command
}

func (handler SynonymHandler) Validate() error {
	return synonms.ValidateSynonymFlags(*handler.Flags)
}

func (handler *SynonymHandler) AskAndFill() error {
	err := synonms.AskSynonym(handler.Flags, handler.Cmd)
	if err != nil {
		return err
	}

	return synonms.ValidateSynonymFlags(*handler.Flags)
}

// `indices config export`
type IndexConfigExportHandler struct {
	Opts *config.ExportOptions
}

func (handler IndexConfigExportHandler) Validate() error {
	return config.ValidateExportConfigFlags(*handler.Opts)
}

func (handler *IndexConfigExportHandler) AskAndFill() error {
	err := config.AskExportConfig(handler.Opts)
	if err != nil {
		return err
	}

	return config.ValidateExportConfigFlags(*handler.Opts)
}
