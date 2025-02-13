package configexport

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	indexConfig "github.com/algolia/cli/pkg/cmd/shared/config"
	"github.com/algolia/cli/pkg/cmd/shared/handler"
	config "github.com/algolia/cli/pkg/cmd/shared/handler/indices"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/utils"

	"github.com/algolia/cli/pkg/validators"
)

// NewExportCmd creates and returns an export command for index config
func NewExportCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &config.ExportOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4SearchClient,
	}

	cmd := &cobra.Command{
		Use:               "export <index> [--scope <scope>...] [--directory]",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.V4IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "settings",
		},
		Short: "Export an index configuration (settings, synonyms, rules) to a file",
		Long: heredoc.Doc(`
			Export an index configuration (settings, synonyms, rules) to a file.
		`),
		Example: heredoc.Doc(`
			# Export the config of the index 'MOVIES' to a .json in the current folder
			$ algolia index config export MOVIES

			# Export the synonyms and rules of the index 'MOVIES' to a .json in the current folder
			$ algolia index config export MOVIES --scope synonyms,rules

			# Export the config of the index 'MOVIES' to a .json into 'exports' folder
			$ algolia index config export MOVIES --directory exports
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			client, err := opts.SearchClient()
			if err != nil {
				return err
			}
			listResponse, err := client.ListIndices(client.NewApiListIndicesRequest())
			if err != nil {
				return err
			}
			var indices []string
			for _, i := range listResponse.Items {
				indices = append(indices, i.Name)
			}
			opts.Indices = indices
			exportConfigHandler := &handler.IndexConfigExportHandler{
				Opts: opts,
			}

			err = handler.HandleFlags(exportConfigHandler, opts.IO.CanPrompt())
			if err != nil {
				return err
			}

			return runExportCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.Directory, "directory", "d", "", "Directory path of the output file (default: current folder)")
	_ = cmd.MarkFlagDirname("directory")
	cmd.Flags().
		StringSliceVarP(&opts.Scope, "scope", "s", []string{"settings", "synonyms", "rules"}, "Scope to export (default: all)")
	_ = cmd.RegisterFlagCompletionFunc("scope",
		cmdutil.StringSliceCompletionFunc(map[string]string{
			"settings": "settings",
			"synonyms": "synonyms",
			"rules":    "rules",
		}, "export only"))

	return cmd
}

func runExportCmd(opts *config.ExportOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	configJSON, err := indexConfig.GetIndexConfig(client, opts.Index, opts.Scope, cs)
	if err != nil {
		return err
	}

	configJSONIndented, err := json.MarshalIndent(configJSON, "", "  ")
	if err != nil {
		return fmt.Errorf(
			"%s An error occurred when creating the config json: %w",
			cs.FailureIcon(),
			err,
		)
	}

	filePath := config.GetConfigFileName(
		opts.Directory,
		opts.Index,
		client.GetConfiguration().AppID,
	)
	// Gosec wants permissions of 0600 or less, but I don't want to change it
	err = os.WriteFile(filePath, configJSONIndented, 0o644) // nolint:gosec
	if err != nil {
		return fmt.Errorf("%s An error occurred when saving the file: %w", cs.FailureIcon(), err)
	}
	currentDir, _ := os.Getwd()
	rootPath := "."
	if opts.Directory != "" {
		rootPath = currentDir
	}
	fmt.Printf(
		"%s '%s' Index config (%s) successfully exported to %s\n",
		cs.SuccessIcon(),
		opts.Index,
		utils.SliceToReadableString(opts.Scope),
		fmt.Sprintf("%s/%s", rootPath, filePath),
	)

	return nil
}
