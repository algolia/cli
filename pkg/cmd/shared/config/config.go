package config

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

func GetSynonyms(client *search.APIClient, srcIndex string) ([]search.SynonymHit, error) {
	var synonyms []search.SynonymHit

	err := client.BrowseSynonyms(
		srcIndex,
		*search.NewEmptySearchSynonymsParams(),
		search.WithAggregator(func(res any, _ error) {
			response, _ := res.(search.SearchSynonymsResponse)
			synonyms = append(synonyms, response.Hits...)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve synonyms from source index: %s: %v", srcIndex, err)
	}
	return synonyms, nil
}

func GetRules(client *search.APIClient, srcIndex string) ([]search.Rule, error) {
	var rules []search.Rule

	err := client.BrowseRules(
		srcIndex,
		*search.NewEmptySearchRulesParams(),
		search.WithAggregator(func(res any, _ error) {
			response, _ := res.(search.SearchRulesResponse)
			rules = append(rules, response.Hits...)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve rules from source index: %s: %v", srcIndex, err)
	}
	return rules, nil
}

type ExportConfigJson struct {
	Settings *search.SettingsResponse `json:"settings,omitempty"`
	Rules    []search.Rule            `json:"rules,omitempty"`
	Synonyms []search.SynonymHit      `json:"synonyms,omitempty"`
}

func GetIndexConfig(
	client *search.APIClient,
	index string,
	scope []string,
	cs *iostreams.ColorScheme,
) (*ExportConfigJson, error) {
	var configJson ExportConfigJson

	if utils.Contains(scope, "synonyms") {
		rawSynonyms, err := GetSynonyms(client, index)
		if err != nil {
			return nil, fmt.Errorf(
				"%s An error occurred when retrieving synonyms: %w",
				cs.FailureIcon(),
				err,
			)
		}
		configJson.Synonyms = rawSynonyms
	}

	if utils.Contains(scope, "rules") {
		rawRules, err := GetRules(client, index)
		if err != nil {
			return nil, fmt.Errorf(
				"%s An error occurred when retrieving rules: %w",
				cs.FailureIcon(),
				err,
			)
		}
		configJson.Rules = rawRules
	}

	if utils.Contains(scope, "settings") {
		rawSettings, err := client.GetSettings(client.NewApiGetSettingsRequest(index))
		if err != nil {
			return nil, fmt.Errorf(
				"%s An error occurred when retrieving settings: %w",
				cs.FailureIcon(),
				err,
			)
		}
		configJson.Settings = rawSettings
	}

	if len(configJson.Rules) == 0 && len(configJson.Synonyms) == 0 && configJson.Settings == nil {
		return nil, fmt.Errorf("%s No config to export", cs.FailureIcon())
	}

	return &configJson, nil
}
