package wording

import validator "github.com/algolia/cli/pkg/cmdutil/validators"

func GetSynonymWording(synonymType validator.SynonymType) string {
	if synonymType == "" || synonymType.String() == validator.Regular {
		return "Regular synonym"
	}

	switch synonymType.String() {
	case validator.OneWay:
		return "One way synonym"
	case validator.Placeholder:
		return "Placeholder synonym"
	case validator.AltCorrection1:
		return "Alt correction 1 synonym"
	case validator.AltCorrection2:
		return "Alt correction 2 synonym"
	}

	return ""
}
