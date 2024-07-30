package config

import (
	"fmt"
	"reflect"

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
			for _, hit := range res.(*search.SearchSynonymsResponse).Hits {
				synonyms = append(synonyms, hit)
			}
		}),
	)
	if err != nil {
		return nil, err
	}
	return synonyms, nil
}

func GetRules(client *search.APIClient, srcIndex string) ([]search.Rule, error) {
	var rules []search.Rule
	err := client.BrowseRules(
		srcIndex,
		*search.NewEmptySearchRulesParams(),
		search.WithAggregator(func(res any, _ error) {
			for _, hit := range res.(*search.SearchRulesResponse).Hits {
				rules = append(rules, hit)
			}
		}),
	)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

type ExportConfigJson struct {
	Settings *search.IndexSettings `json:"settings,omitempty"`
	Rules    []search.Rule         `json:"rules,omitempty"`
	Synonyms []search.SynonymHit   `json:"synonyms,omitempty"`
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
		configJson.Settings = SettingsResponseToIndexSettings(rawSettings)
	}

	if len(configJson.Rules) == 0 && len(configJson.Synonyms) == 0 && configJson.Settings == nil {
		return nil, fmt.Errorf("%s No config to export", cs.FailureIcon())
	}

	return &configJson, nil
}

// SettingsResponseToIndexSettings converts the SettingsResponse struct to an IndexSettings struct
func SettingsResponseToIndexSettings(r *search.SettingsResponse) *search.IndexSettings {
	settings := search.NewIndexSettings()
	settingsType := reflect.TypeOf(settings).Elem()
	settingsVal := reflect.ValueOf(settings).Elem()
	rval := reflect.ValueOf(r).Elem()

	for i := 0; i < settingsType.NumField(); i++ {
		fieldName := settingsType.Field(i).Name
		responseField := rval.FieldByName(fieldName)
		if responseField.IsValid() && responseField.CanSet() {
			settingsVal.FieldByName(fieldName).Set(responseField)
		}
	}

	return settings
}
