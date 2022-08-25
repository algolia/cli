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
	Id     string
	Values string
	Indice string
}

const successTemplate = `{{ .Type}} '{{ .Id}}' successfully saved with {{ .Values}} to {{ .Indice}}`

func GetSuccessMessage(flags shared.SynonymFlags, indice string) (error, string) {
	var successMessage SuccessMessage

	if flags.SynonymType == "" || flags.SynonymType.String() == shared.Regular {
		successMessage = SuccessMessage{
			Type: "Synonym",
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("%s (%s)",
				utils.Pluralize(len(flags.Synonyms), "synonym"),
				strings.Join(flags.Synonyms, ", ")),
			Indice: indice,
		}
	}

	switch flags.SynonymType.String() {
	case shared.OneWay:
		successMessage = SuccessMessage{
			Type: "One way synonym",
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("input '%s' and %s (%s)",
				flags.SynonymInput,
				utils.Pluralize(len(flags.Synonyms), "synonym"),
				strings.Join(flags.Synonyms, ", ")),
			Indice: indice,
		}
	case shared.Placeholder:
		successMessage = SuccessMessage{
			Type: "Placeholder synonym",
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("placeholder '%s' and %s (%s)",
				flags.SynonymPlaceholder,
				utils.Pluralize(len(flags.SynonymReplacements), "replacement"),
				strings.Join(flags.SynonymReplacements, ", ")),
			Indice: indice,
		}
	case shared.AltCorrection1, shared.AltCorrection2:
		altCorrectionType := "1"
		if flags.SynonymType.String() == shared.AltCorrection2 {
			altCorrectionType = "2"
		}
		altCorrectionType = "Alt correction " + altCorrectionType + " synonym"
		successMessage = SuccessMessage{
			Type: altCorrectionType,
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("word '%s' and %s (%s)",
				flags.SynonymWord,
				utils.Pluralize(len(flags.SynonymCorrections), "correction"),
				strings.Join(flags.SynonymCorrections, ", ")),
			Indice: indice,
		}
	}

	t := template.Must(template.New("successMessage").Parse(successTemplate))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, successMessage); err != nil {
		return err, ""
	}
	return nil, tpl.String() + "\n"
}
