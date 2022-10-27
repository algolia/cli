package config

import (
	"fmt"
	"io"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

func GetSynonyms(srcIndex *search.Index) ([]search.Synonym, error) {
	it, err := srcIndex.BrowseSynonyms()
	if err != nil {
		return nil, fmt.Errorf("cannot browse source index synonyms: %v", err)
	}

	var synonyms []search.Synonym

	for {
		synonym, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, fmt.Errorf("error while iterating source index synonyms: %v", err)
			}
		}
		synonyms = append(synonyms, synonym)
	}

	return synonyms, nil
}

func GetRules(srcIndex *search.Index) ([]search.Rule, error) {
	it, err := srcIndex.BrowseRules()
	if err != nil {
		return nil, fmt.Errorf("cannot browse source index rules: %v", err)
	}

	var rules []search.Rule

	for {
		rule, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, fmt.Errorf("error while iterating source index rules: %v", err)
			}
		}
		rules = append(rules, *rule)
	}

	return rules, nil
}

type ExportConfigJson struct {
	Settings *search.Settings `json:"settings,omitempty"`
	Rules    []search.Rule    `json:"rules,omitempty"`
	Synonyms []search.Synonym `json:"synonyms,omitempty"`
}

func GetIndiceConfig(indice *search.Index, scope []string, cs *iostreams.ColorScheme) (*ExportConfigJson, error) {
	var configJson ExportConfigJson

	if utils.Contains(scope, "synonyms") {
		rawSynonyms, err := GetSynonyms(indice)
		if err != nil {
			return nil, fmt.Errorf("%s An error occured when retrieving synonyms: %w", cs.FailureIcon(), err)
		}
		configJson.Synonyms = rawSynonyms
	}

	if utils.Contains(scope, "rules") {
		rawRules, err := GetRules(indice)
		if err != nil {
			return nil, fmt.Errorf("%s An error occured when retrieving rules: %w", cs.FailureIcon(), err)
		}
		configJson.Rules = rawRules
	}

	if utils.Contains(scope, "settings") {
		rawSettings, err := indice.GetSettings()
		if err != nil {
			return nil, fmt.Errorf("%s An error occured when retrieving settings: %w", cs.FailureIcon(), err)
		}
		configJson.Settings = &rawSettings
	}

	if len(configJson.Rules) == 0 && len(configJson.Synonyms) == 0 && configJson.Settings == nil {
		return nil, fmt.Errorf("%s No config to export", cs.FailureIcon())
	}

	return &configJson, nil
}
