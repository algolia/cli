package save

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	shared "github.com/algolia/cli/pkg/cmd/synonyms/shared"
	"github.com/algolia/cli/pkg/utils"
)

type SuccessMessage struct {
	Icon   string
	Type   string
	ID     string
	Values string
	Index  string
}

const successTemplate = `{{ .Type }} '{{ .ID }}' successfully saved with {{ .Values }} to {{ .Index }}`

func GetSuccessMessage(flags shared.SynonymFlags, index string) (string, error) {
	var successMessage SuccessMessage

	if flags.SynonymType == "" || flags.SynonymType == shared.Regular {
		successMessage = SuccessMessage{
			Type: "Synonym",
			ID:   flags.SynonymID,
			Values: fmt.Sprintf("%s (%s)",
				utils.Pluralize(len(flags.Synonyms), "synonym"),
				strings.Join(flags.Synonyms, ", ")),
			Index: index,
		}
	}

	switch flags.SynonymType {
	case shared.OneWay:
		successMessage = SuccessMessage{
			Type: "One way synonym",
			ID:   flags.SynonymID,
			Values: fmt.Sprintf("input '%s' and %s (%s)",
				flags.SynonymInput,
				utils.Pluralize(len(flags.Synonyms), "synonym"),
				strings.Join(flags.Synonyms, ", ")),
			Index: index,
		}
	case shared.Placeholder:
		successMessage = SuccessMessage{
			Type: "Placeholder synonym",
			ID:   flags.SynonymID,
			Values: fmt.Sprintf("placeholder '%s' and %s (%s)",
				flags.SynonymPlaceholder,
				utils.Pluralize(len(flags.SynonymReplacements), "replacement"),
				strings.Join(flags.SynonymReplacements, ", ")),
			Index: index,
		}
	case shared.AltCorrection1, shared.AltCorrection2:
		altCorrectionType := "1"
		if flags.SynonymType == shared.AltCorrection2 {
			altCorrectionType = "2"
		}
		altCorrectionType = "Alt correction " + altCorrectionType + " synonym"
		successMessage = SuccessMessage{
			Type: altCorrectionType,
			ID:   flags.SynonymID,
			Values: fmt.Sprintf("word '%s' and %s (%s)",
				flags.SynonymWord,
				utils.Pluralize(len(flags.SynonymCorrections), "correction"),
				strings.Join(flags.SynonymCorrections, ", ")),
			Index: index,
		}
	}

	t := template.Must(template.New("successMessage").Parse(successTemplate))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, successMessage); err != nil {
		return "", err
	}
	return tpl.String() + "\n", nil
}
