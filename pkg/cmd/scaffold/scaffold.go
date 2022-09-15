package scaffold

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	algoliaSearch "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

var templateChoices = []string{
	"InstantSearch.js",
	"React InstantSearch",
	"React InstantSearch Native",
	"Vue InstantSearch",
	"Angular InstantSearch",
	"InstantSearch iOS",
	"InstantSearch Android",
}

// ScaffoldOptions represents the options for the scaffold command
type ScaffoldOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*algoliaSearch.Client, error)

	Path string

	Template     string
	SearchAPIKey string
	IndexName    string
}

// NewScaffoldCmd returns a new instance of the scaffold command
func NewScaffoldCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ScaffoldOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	cmd := &cobra.Command{
		Use:   "scaffold  <directory>",
		Short: "Create a new instantsearch app in the given directory.",
		Args:  validators.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Example: heredoc.Doc(`
			# Scaffold a new React instantsearch app in the directory "my-app" for the index "products"
			$ algolia scaffold my-app -t "React InstantSearch" -i "products" -k "search-api-key"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(args[0]); !os.IsNotExist(err) {
				return fmt.Errorf("directory %s already exists", args[0])
			}

			opts.Path = args[0]

			return runScaffoldCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Template, "template", "t", "InstantSearch.js", "The template to use for the scaffold")
	_ = cmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return templateChoices, cobra.ShellCompDirectiveNoFileComp
	})
	_ = cmd.MarkFlagRequired("template")

	cmd.Flags().StringVarP(&opts.IndexName, "index-name", "i", "", "The index name to use")
	_ = cmd.RegisterFlagCompletionFunc("index-name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := f.SearchClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		indicesRes, err := client.ListIndices()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		indicesComp := make([]string, 0, len(indicesRes.Items))
		for _, index := range indicesRes.Items {
			indicesComp = append(indicesComp, fmt.Sprintf("%s\t%s records", index.Name, humanize.Comma(index.Entries)))
		}
		return indicesComp, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Flags().StringVarP(&opts.SearchAPIKey, "search-api-key", "k", "", "The search api key to use")
	_ = cmd.RegisterFlagCompletionFunc("search-api-key", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := f.SearchClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		keysRes, err := client.ListAPIKeys()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		keysComp := make([]string, 0, len(keysRes.Keys))
		for _, key := range keysRes.Keys {
			if key.Description == "Search-only API Key" {
				keysComp = append(keysComp, key.Value)
			}
		}
		return keysComp, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func runScaffoldCmd(opts *ScaffoldOptions) error {
	args, err := YarnOrNpxArgs()
	if err != nil {
		return err
	}

	args = append(args, opts.Path, "--template", opts.Template)

	appID, _ := opts.Config.Profile().GetApplicationID()
	if appID != "" {
		args = append(args, "--app-id", appID)
	}

	if opts.IndexName != "" {
		args = append(args, "--index-name", opts.IndexName)
	}

	if opts.SearchAPIKey != "" {
		args = append(args, "--api-key", opts.SearchAPIKey)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = opts.IO.Out
	cmd.Stderr = opts.IO.ErrOut
	cmd.Stdin = opts.IO.In

	return cmd.Run()
}

func YarnOrNpxArgs() ([]string, error) {
	_, err := exec.LookPath("npx")
	if err == nil {
		return []string{"npx", "create-instantsearch-app"}, nil
	}
	_, err = exec.LookPath("yarn")
	if err == nil {
		return []string{"yarn", "create", "instantsearch-app"}, nil
	}

	return []string{}, fmt.Errorf("npx or yarn is required to scaffold an instantsearch app")
}
