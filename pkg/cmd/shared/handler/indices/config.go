package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"

	"github.com/algolia/cli/pkg/ask"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type ExportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	Indices   []string
	Index     string
	Scope     []string
	Directory string

	SearchClient func() (*search.APIClient, error)
}

func ValidateExportConfigFlags(opts ExportOptions) error {
	cs := opts.IO.ColorScheme()

	if !utils.Contains(opts.Indices, opts.Index) {
		return fmt.Errorf("%s Index '%s' doesn't exist", cs.FailureIcon(), opts.Index)
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
func GetConfigFileName(path string, indexName string, appID string) string {
	rootPath := ""
	if path != "" {
		rootPath = path + "/"
	}

	return fmt.Sprintf(
		"%sexport-%s-%s-%s.json",
		rootPath,
		indexName,
		appID,
		strconv.FormatInt(time.Now().UTC().Unix(), 10),
	)
}

type ImportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	ImportConfig ImportConfigJSON

	Index                 string
	FilePath              string
	Scope                 []string
	ClearExistingSynonyms bool
	ClearExistingRules    bool

	ForwardSettingsToReplicas bool
	ForwardSynonymsToReplicas bool
	ForwardRulesToReplicas    bool

	DoConfirm bool
}

type ImportConfigJSON struct {
	Settings *search.IndexSettings `json:"settings,omitempty"`
	Rules    []search.Rule         `json:"rules,omitempty"`
	Synonyms []search.SynonymHit   `json:"synonyms,omitempty"`
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
		return fmt.Errorf(
			"%s Cannot clear existing rules if rules are not in scope",
			cs.FailureIcon(),
		)
	}
	if opts.ClearExistingSynonyms && !utils.Contains(opts.Scope, "synonyms") {
		return fmt.Errorf(
			"%s Cannot clear existing synonyms if synonyms are not in scope",
			cs.FailureIcon(),
		)
	}
	// Scope and config
	if (utils.Contains(opts.Scope, "settings") && opts.ImportConfig.Settings != nil) ||
		(utils.Contains(opts.Scope, "rules") && len(opts.ImportConfig.Rules) > 0) ||
		(utils.Contains(opts.Scope, "synonyms") && len(opts.ImportConfig.Synonyms) > 0) {
		return nil
	}
	return fmt.Errorf(
		"%s No %s found in config file",
		cs.FailureIcon(),
		utils.SliceToReadableString(opts.Scope),
	)
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

	erroredScope := opts.Scope
	opts.Scope = []string{}
	err = ask.AskMultiSelectQuestion(
		"scope (comma separated):",
		erroredScope,
		&opts.Scope,
		scopeOptions,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}
	if utils.Contains(opts.Scope, "synonyms") {
		err = ask.AskBooleanQuestion(
			"Clear and replace existing synonyms? (default: no)",
			&opts.ClearExistingSynonyms,
			false,
		)
		if err != nil {
			return err
		}
		err = ask.AskBooleanQuestion(
			"Forward synonyms to replicas? (default: no)",
			&opts.ForwardSynonymsToReplicas,
			false,
		)
		if err != nil {
			return err
		}
	}
	if utils.Contains(opts.Scope, "rules") {
		err = ask.AskBooleanQuestion(
			"Clear and replace existing rules? (default: no)",
			&opts.ClearExistingRules,
			false,
		)
		if err != nil {
			return err
		}
		err = ask.AskBooleanQuestion(
			"Forward rules to replicas? (default: no)",
			&opts.ForwardRulesToReplicas,
			false,
		)
		if err != nil {
			return err
		}
	}
	if utils.Contains(opts.Scope, "settings") {
		err = ask.AskBooleanQuestion(
			"Forward settings to replicas? (default: no)",
			&opts.ForwardSettingsToReplicas,
			false,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func readConfigFromFile(cs *iostreams.ColorScheme, filePath string) (*ImportConfigJSON, error) {
	var config *ImportConfigJSON

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s An error occurred when opening file: %w", cs.FailureIcon(), err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf(
			"%s An error occurred when reading JSON file: %w",
			cs.FailureIcon(),
			err,
		)
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, fmt.Errorf(
			"%s An error occurred when parsing JSON file: %w",
			cs.FailureIcon(),
			err,
		)
	}

	return config, nil
}
