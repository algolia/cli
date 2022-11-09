package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		"scope (comma separated):",
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

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	ImportConfig ImportConfigJson

	Indices               []string
	FilePath              string
	Scope                 []string
	ClearExistingSynonyms bool
	ClearExistingRules    bool

	ForwardSettingsToReplicas bool
	ForwardSynonymsToReplicas bool
	ForwardRulesToReplicas    bool

	DoConfirm bool
}

type ImportConfigJson struct {
	Settings *search.Settings `json:"settings,omitempty"`
	Rules    []search.Rule    `json:"rules,omitempty"`
	Synonyms []Synonym        `json:"synonyms,omitempty"`
}

type Synonym struct {
	Type                                string
	ObjectID, Word, Input, Placeholder  string
	Corrections, Synonyms, Replacements []string
}

func ValidateImportConfigFlags(opts *ImportOptions) error {
	cs := opts.IO.ColorScheme()

	if opts.FilePath == "" {
		return fmt.Errorf("%s Config file is required", cs.FailureIcon())
	}

	config, err := readConfigFromFile(cs, opts.FilePath)
	if err != nil {
		return err
	}
	opts.ImportConfig = *config

	// Required flags
	if len(opts.Scope) == 0 {
		return fmt.Errorf("%s Scope is required", cs.FailureIcon())
	}
	// Scope and replace/clear existing options
	if opts.ClearExistingRules && !utils.Contains(opts.Scope, "rules") {
		return fmt.Errorf("%s Cannot clear existing rules if rules are not in scope", cs.FailureIcon())
	}
	if opts.ClearExistingSynonyms && !utils.Contains(opts.Scope, "synonyms") {
		return fmt.Errorf("%s Cannot clear existing synonyms if synonyms are not in scope", cs.FailureIcon())
	}
	// Scope and config
	if utils.Contains(opts.Scope, "settings") && opts.ImportConfig.Settings == nil {
		return fmt.Errorf("%s No setting found in config file", cs.FailureIcon())
	}
	if utils.Contains(opts.Scope, "rules") && len(opts.ImportConfig.Rules) == 0 {
		return fmt.Errorf("%s No rule found in config file", cs.FailureIcon())
	}
	if utils.Contains(opts.Scope, "synonyms") && len(opts.ImportConfig.Synonyms) == 0 {
		return fmt.Errorf("%s No synonym found in config file", cs.FailureIcon())
	}

	return nil
}

func AskImportConfig(opts *ImportOptions) error {
	// Validate file path
	err := ask.AskInputQuestionWithSuggestion(
		"file (path of the .json config file)",
		&opts.FilePath,
		opts.FilePath,
		func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}
	config, err := readConfigFromFile(opts.IO.ColorScheme(), opts.FilePath)
	if err != nil {
		return err
	}
	opts.ImportConfig = *config

	scopeOptions := []string{}
	if len(opts.ImportConfig.Rules) > 0 {
		scopeOptions = append(scopeOptions, "rules")
	}
	if len(opts.ImportConfig.Synonyms) > 0 {
		scopeOptions = append(scopeOptions, "synonyms")
	}
	if opts.ImportConfig.Settings != nil {
		scopeOptions = append(scopeOptions, "settings")
	}
	err = ask.AskMultiSelectQuestion(
		"scope (comma separated):",
		opts.Scope,
		&opts.Scope,
		scopeOptions,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}
	if utils.Contains(opts.Scope, "synonyms") {
		err = ask.AskBooleanQuestion(
			"clearExistingSynonyms (default: false)",
			&opts.ClearExistingSynonyms,
			false,
		)
		if err != nil {
			return err
		}
		err = ask.AskBooleanQuestion(
			"forwardSynonymsToReplicas (default: false)",
			&opts.ForwardSynonymsToReplicas,
			false,
		)
		if err != nil {
			return err
		}
	}
	if utils.Contains(opts.Scope, "rules") {
		err = ask.AskBooleanQuestion(
			"clearExistingRules (default: false)",
			&opts.ClearExistingRules,
			false,
		)
		if err != nil {
			return err
		}
		err = ask.AskBooleanQuestion(
			"forwardRulesToReplicas (default: false)",
			&opts.ForwardRulesToReplicas,
			false,
		)
		if err != nil {
			return err
		}
	}
	if utils.Contains(opts.Scope, "settings") {
		err = ask.AskBooleanQuestion(
			"forwardSettingsToReplicas (default: false)",
			&opts.ForwardSettingsToReplicas,
			false,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func readConfigFromFile(cs *iostreams.ColorScheme, filePath string) (*ImportConfigJson, error) {
	var config *ImportConfigJson

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s An error occurred when opening file: %w", cs.FailureIcon(), err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("%s An error occurred when reading JSON file: %w", cs.FailureIcon(), err)
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf("%s An error occurred when parsing JSON file: %w", cs.FailureIcon(), err)
	}

	return config, nil
}
