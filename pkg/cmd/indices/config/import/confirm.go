package indeximport

import (
	"fmt"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

func GetConfirmMessage(cs *iostreams.ColorScheme, scope []string, clearExistingRules, clearExistingSynonyms bool) string {
	scopeToClear := []string{}
	scopeToUpdate := []string{}
	message := ""

	if utils.Contains(scope, "settings") {
		scopeToClear = append(scopeToClear, "settings")
	}
	if utils.Contains(scope, "rules") {
		if clearExistingRules {
			scopeToClear = append(scopeToClear, "rules")
		} else {
			scopeToUpdate = append(scopeToUpdate, "rules")
		}
	}
	if utils.Contains(scope, "synonyms") {
		if clearExistingSynonyms {
			scopeToClear = append(scopeToClear, "synonyms")
		} else {
			scopeToUpdate = append(scopeToUpdate, "synonyms")
		}
	}
	if len(scopeToClear) > 0 {
		message = fmt.Sprintf("%s Your %s will be %s\n",
			cs.WarningIcon(), utils.SliceToReadableString(scopeToClear), cs.Bold("CLEARED and REPLACED."))
	}
	if len(scopeToUpdate) > 0 {
		message = fmt.Sprintf("%s%s Your %s will be %s\n",
			message, cs.WarningIcon(), utils.SliceToReadableString(scopeToUpdate), cs.Bold("UPDATED"))
	}

	return message
}
