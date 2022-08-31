package wording

import (
	"fmt"
	"strings"

	validator "github.com/algolia/cli/pkg/cmdutil/validators"
	"github.com/algolia/cli/pkg/utils"
)

func GetSynonymSuccessWording(opts validator.SaveOptions, saveWording string) string {
	if opts.SynonymType == "" || opts.SynonymType.String() == validator.Regular {
		return fmt.Sprintf("%s %s '%s' successfully %s with %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Synonym",
			opts.SynonymID,
			saveWording,
			utils.Pluralize(len(opts.Synonyms), "synonym"),
			strings.Join(opts.Synonyms, ", "),
			opts.Indice)
	}

	switch opts.SynonymType.String() {
	case validator.OneWay:
		return fmt.Sprintf("%s %s '%s' successfully %s with input '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"One way synonym",
			opts.SynonymID,
			saveWording,
			opts.SynonymInput,
			utils.Pluralize(len(opts.Synonyms), "synonym"),
			strings.Join(opts.Synonyms, ", "),
			opts.Indice)
	case validator.Placeholder:
		return fmt.Sprintf("%s %s '%s' successfully %s with placeholder '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Placeholder synonym",
			opts.SynonymID,
			saveWording,
			opts.SynonymPlaceholder,
			utils.Pluralize(len(opts.SynonymReplacements), "replacement"),
			strings.Join(opts.SynonymReplacements, ", "),
			opts.Indice)
	case validator.AltCorrection1:
		return fmt.Sprintf("%s %s '%s' successfully %s with word '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Alt correction 1 synonym",
			opts.SynonymID,
			saveWording,
			opts.SynonymWord,
			utils.Pluralize(len(opts.SynonymCorrections), "correction"),
			strings.Join(opts.SynonymCorrections, ", "),
			opts.Indice)
	case validator.AltCorrection2:
		return fmt.Sprintf("%s %s '%s' successfully %s with word '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Alt correction 2 synonym",
			opts.SynonymID,
			saveWording,
			opts.SynonymWord,
			utils.Pluralize(len(opts.SynonymCorrections), "correction"),
			strings.Join(opts.SynonymCorrections, ", "),
			opts.Indice)
	}

	return ""
}
