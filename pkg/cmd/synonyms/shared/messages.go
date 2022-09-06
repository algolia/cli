package shared

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/algolia/cli/pkg/utils"
)

type SuccessMessage struct {
	Icon   string
	Type   string
	Id     string
	Values string
	Indice string
}

const successTemplate = `{{ .Icon}} {{ .Type}} '{{ .Id}}' successfully saved with {{ .Values}} to {{ .Indice}}`

func GetSuccessMessage(flags SynonymFlags, opts SaveOptions) (error, string) {
	var successMessage SuccessMessage

	if flags.SynonymType == "" || flags.SynonymType.String() == Regular {
		successMessage = SuccessMessage{
			Type: "Synonym",
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("%s (%s)",
				utils.Pluralize(len(flags.Synonyms), "synonym"),
				strings.Join(flags.Synonyms, ", ")),
			Indice: opts.Indice,
		}
	}

	switch flags.SynonymType.String() {
	case OneWay:
		successMessage = SuccessMessage{
			Type: "One way synonym",
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("input '%s' and %s (%s)",
				flags.SynonymInput,
				utils.Pluralize(len(flags.Synonyms), "synonym"),
				strings.Join(flags.Synonyms, ", ")),
			Indice: opts.Indice,
		}
	case Placeholder:
		successMessage = SuccessMessage{
			Type: "Placeholder synonym",
			Id:   flags.SynonymID,
			Values: fmt.Sprintf("placeholder '%s' and %s (%s)",
				flags.SynonymPlaceholder,
				utils.Pluralize(len(flags.SynonymReplacements), "replacement"),
				strings.Join(flags.SynonymReplacements, ", ")),
			Indice: opts.Indice,
		}
	case AltCorrection1, AltCorrection2:
		altCorrectionType := "1"
		if flags.SynonymType.String() == AltCorrection2 {
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
			Indice: opts.Indice,
		}
	}

	successMessage.Icon = opts.IO.ColorScheme().SuccessIcon()
	t := template.Must(template.New("successMessage").Parse(successTemplate))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, successMessage); err != nil {
		return err, ""
	}
	return nil, tpl.String() + "\n"
}
