package export

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	indiceConfig "github.com/algolia/cli/pkg/cmd/shared/config"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
)

type ExportOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	Indices   []string
	Scope     []string
	Directory string

	SearchClient func() (*search.Client, error)
	Client       *search.Client
}

type ConfigJson struct {
	Settings *search.Settings `json:"settings,omitempty"`
	Rules    []search.Rule    `json:"rules,omitempty"`
	Synonyms []search.Synonym `json:"synonyms,omitempty"`
}

// NewExportCmd creates and returns an export command for indices config
func NewExportCmd(f *cmdutil.Factory, runF func(*ExportOptions) error) *cobra.Command {
	opts := &ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "export <index>...",
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Export the config of one or multiple indice(es)",
		Long: heredoc.Doc(`
			Export the config of one or multiple indice(s) including their settings, synonyms and rules.
		`),
		Example: heredoc.Doc(`
			# Export the config of the index 'TEST_PRODUCTS' to a .json in the current folder
			$ algolia indices config export TEST_PRODUCTS

			# Export the config of the 'TEST_PRODUCTS_1', 'TEST_PRODUCTS_2' and 'TEST_PRODUCTS_3' indices to a .json in the current folder
			$ algolia indices config export TEST_PRODUCTS_1 TEST_PRODUCTS_2 TEST_PRODUCTS_3

			# Export the synonyms and rules of the index 'TEST_PRODUCTS' to a .json in the current folder
			$ algolia indices config export TEST_PRODUCTS --scope synonyms,rules

			# Export the config of the index 'TEXT_PRODUCTS' to a .json into 'exports' folder
			$ algolia indices config export TEST_PRODUCTS --directory exports
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cs := opts.IO.ColorScheme()
			opts.Indices = args
			client, err := opts.SearchClient()
			if err != nil {
				return err
			}
			opts.Client = client

			existingIndices, err := client.ListIndices()
			if err != nil {
				return err
			}
			var existingIndicesNames []string
			for _, existingIndex := range existingIndices.Items {
				existingIndicesNames = append(existingIndicesNames, existingIndex.Name)
			}
			for _, indexToCheck := range opts.Indices {
				if !utils.Contains(existingIndicesNames, indexToCheck) {
					return fmt.Errorf("%s Indice '%s' doesn't exist", cs.FailureIcon(), indexToCheck)
				}
			}

			return runExportCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Directory, "directory", "d", "", "Directory path of the output file (directy must exist)")
	_ = cmd.MarkFlagDirname("directory")
	cmd.Flags().StringSliceVarP(&opts.Scope, "scope", "s", []string{"settings", "synonyms", "rules"}, "Scope to export (default: all)")
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "settings",
			"synonyms": "synonyms",
			"rules":    "rules",
		}, "export only"))

	return cmd
}

func runExportCmd(opts *ExportOptions) error {
	cs := opts.IO.ColorScheme()

	for _, indexName := range opts.Indices {
		indice := opts.Client.InitIndex(indexName)
		var configJson ConfigJson

		if utils.Contains(opts.Scope, "synonyms") {
			rawSynonyms, err := indiceConfig.GetSynonyms(indice)
			if err != nil {
				return fmt.Errorf("%s An error occured when retrieving synonyms: %w", cs.FailureIcon(), err)
			}
			configJson.Synonyms = rawSynonyms
		}

		if utils.Contains(opts.Scope, "rules") {
			rawRules, err := indiceConfig.GetRules(indice)
			if err != nil {
				return fmt.Errorf("%s An error occured when retrieving rules: %w", cs.FailureIcon(), err)
			}
			configJson.Rules = rawRules
		}

		if utils.Contains(opts.Scope, "settings") {
			rawSettings, err := indice.GetSettings()
			if err != nil {
				return fmt.Errorf("%s An error occured when retrieving settings: %w", cs.FailureIcon(), err)
			}
			configJson.Settings = &rawSettings
		}

		if len(configJson.Rules) == 0 && len(configJson.Synonyms) == 0 && configJson.Settings == nil {
			return fmt.Errorf("%s No config to export", cs.FailureIcon())
		}

		configJsonIndented, err := json.MarshalIndent(configJson, "", "  ")
		if err != nil {
			return fmt.Errorf("%s An error occured when creating the config json: %w", cs.FailureIcon(), err)
		}

		filePath := getConfigFileName(opts.Directory, indexName, indice.GetAppID(), strconv.FormatInt(time.Now().UTC().Unix(), 10))
		err = os.WriteFile(filePath, configJsonIndented, 0644)
		if err != nil {
			return fmt.Errorf("%s An error occured when saving the file: %w", cs.FailureIcon(), err)
		}

		fmt.Printf("%s '%s' Index config successfully exported to %s\n", cs.SuccessIcon(), indexName, filePath)
	}

	return nil
}

// Matching Algolia Dashboard file naming
// https://github.com/algolia/AlgoliaWeb/blob/develop/_client/src/routes/explorer/components/Explorer/IndexExportSettingsModal.tsx#L88
func getConfigFileName(path string, indiceName string, appId string, unixTimestamp string) string {
	rootPath := ""
	if path != "" {
		rootPath = path + "/"
	}

	return fmt.Sprintf("%sexport-%s-%s-%s.json", rootPath, indiceName, appId, unixTimestamp)
}
