package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/ask"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type ExportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	ExistingIndices []string
	Indices         []string
	Scope           []string
	Directory       string

	SearchClient func() (*search.Client, error)
}

func ValidateExportConfigFlags(opts ExportOptions) error {
	cs := opts.IO.ColorScheme()

	for _, indexToCheck := range opts.Indices {
		if !utils.Contains(opts.ExistingIndices, indexToCheck) {
			return fmt.Errorf("%s Indice '%s' doesn't exist", cs.FailureIcon(), indexToCheck)
		}
	}

	return nil
}

func AskExportConfig(opts *ExportOptions) error {
	err := ask.AskMultiSelectQuestion(
		"replacements (comma separated):",
		opts.Scope,
		&opts.Scope,
		[]string{"settings", "synonyms", "rules"},
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = ask.AskInputQuestion(
		"directory (default to current folder)", &opts.Directory, opts.Directory)
	if err != nil {
		return err
	}

	return nil
}

// Matching Algolia Dashboard file naming
// https://github.com/algolia/AlgoliaWeb/blob/develop/_client/src/routes/explorer/components/Explorer/IndexExportSettingsModal.tsx#L88
func GetConfigFileName(path string, indiceName string, appId string) string {
	rootPath := ""
	if path != "" {
		rootPath = path + "/"
	}

	return fmt.Sprintf("%sexport-%s-%s-%s.json", rootPath, indiceName, appId, strconv.FormatInt(time.Now().UTC().Unix(), 10))
}
