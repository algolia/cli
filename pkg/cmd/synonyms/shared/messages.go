package shared

import (
	"fmt"
	"strings"

	"github.com/algolia/cli/pkg/utils"
)

func GetSynonymSuccessMessage(flags SynonymFlags, opts SaveOptions) string {
	if flags.SynonymType == "" || flags.SynonymType.String() == Regular {
		return fmt.Sprintf("%s %s '%s' successfully saved with %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Synonym",
			flags.SynonymID,
			utils.Pluralize(len(flags.Synonyms), "synonym"),
			strings.Join(flags.Synonyms, ", "),
			opts.Indice)
	}

	switch flags.SynonymType.String() {
	case OneWay:
		return fmt.Sprintf("%s %s '%s' successfully saved with input '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"One way synonym",
			flags.SynonymID,
			flags.SynonymInput,
			utils.Pluralize(len(flags.Synonyms), "synonym"),
			strings.Join(flags.Synonyms, ", "),
			opts.Indice)
	case Placeholder:
		return fmt.Sprintf("%s %s '%s' successfully saved with placeholder '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Placeholder synonym",
			flags.SynonymID,
			flags.SynonymPlaceholder,
			utils.Pluralize(len(flags.SynonymReplacements), "replacement"),
			strings.Join(flags.SynonymReplacements, ", "),
			opts.Indice)
	case AltCorrection1:
		return fmt.Sprintf("%s %s '%s' successfully saved with word '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Alt correction 1 synonym",
			flags.SynonymID,
			flags.SynonymWord,
			utils.Pluralize(len(flags.SynonymCorrections), "correction"),
			strings.Join(flags.SynonymCorrections, ", "),
			opts.Indice)
	case AltCorrection2:
		return fmt.Sprintf("%s %s '%s' successfully saved with word '%s' and %s (%s) to %s\n",
			opts.IO.ColorScheme().SuccessIcon(),
			"Alt correction 2 synonym",
			flags.SynonymID,
			flags.SynonymWord,
			utils.Pluralize(len(flags.SynonymCorrections), "correction"),
			strings.Join(flags.SynonymCorrections, ", "),
			opts.Indice)
	}

	return ""
}
