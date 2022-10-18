package open

import (
	"fmt"
	"sort"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/printers"
)

type OpenUrl struct {
	Default   string
	WithAppId string
}

var openUrlMap = map[string]OpenUrl{
	"api":       {Default: "https://www.algolia.com/doc/api-reference/rest-api/"},
	"dashboard": {Default: "https://www.algolia.com/dashboard", WithAppId: "https://www.algolia.com/apps/%s/dashboard"},
	"codex":     {Default: "https://www.algolia.com/developers/code-exchange/"},
	"devhub":    {Default: "https://www.algolia.com/developers/"},
	"docs":      {Default: "https://algolia.com/doc/"},
	"cli-docs":  {Default: "https://algolia.com/doc/tools/cli/get-started/overview/"},
	"cli-repo":  {Default: "https://github.com/algolia/cli"},
	"status":    {Default: "https://status.algolia.com/", WithAppId: "https://www.algolia.com/apps/%s/monitoring/status"},
}

func openNames() []string {
	keys := make([]string, 0, len(openUrlMap))
	for k := range openUrlMap {
		keys = append(keys, k)
	}

	return keys
}

func getNameUrlMap(applicationID string) map[string]string {
	nameUrlMap := make(map[string]string)
	for _, openName := range openNames() {
		url := openUrlMap[openName].Default
		if applicationID != "" && openUrlMap[openName].WithAppId != "" {
			url = fmt.Sprintf(openUrlMap[openName].WithAppId, applicationID)
		}
		nameUrlMap[openName] = url
	}

	return nameUrlMap
}

// OpenOptions represents the options for the open command
type OpenOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams

	List     bool
	Shortcut string
}

func NewOpenCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &OpenOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}
	cmd := &cobra.Command{
		Use:       "open <shortcut>",
		ValidArgs: openNames(),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if opts.List {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return openNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Short: "Quickly open Algolia pages",
		Long:  `The open command provices shortcuts to quickly let you open pages to Algolia within your browser. 'algolia open --list' for a list of supported shortcuts.`,
		Example: heredoc.Doc(`
			# Display the list of supported shortcuts
			$ algolia open --list

			# Open the dashboard for the current application
			$ algolia open dashboard
			
			# Open the API reference
			$ algolia open api
			
			# Open the documentation
			$ algolia open docs

			# Open the CLI documentation
			$ algolia open cli-docs

			# Open the status page
			$ algolia open status
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Shortcut = args[0]
			}
			return runOpenCmd(opts)
		},
	}

	cmd.Flags().BoolP("list", "l", false, "List all supported shortcuts")

	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runOpenCmd(opts *OpenOptions) error {
	profile := opts.config.Profile()
	applicationID, _ := profile.GetApplicationID()
	nameUrlMap := getNameUrlMap(applicationID)

	if opts.List || opts.Shortcut == "" {
		fmt.Println("open quickly opens Algolia pages. To use, run 'algolia open <shortcut>'.")
		fmt.Println("open supports the following shortcuts:")
		fmt.Println()

		shortcuts := openNames()
		sort.Strings(shortcuts)

		table := printers.NewTablePrinter(opts.IO)
		if table.IsTTY() {
			table.AddField("SHORTCUT", nil, nil)
			table.AddField("URL", nil, nil)
			table.EndRow()
		}

		for shortcutName, url := range nameUrlMap {
			table.AddField(shortcutName, nil, nil)
			table.AddField(url, nil, nil)
			table.EndRow()
		}

		return table.Render()
	}

	var err error
	if url, ok := nameUrlMap[opts.Shortcut]; ok {
		err = open.Browser(url)

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported open command, given: %s", opts.Shortcut)
	}

	return nil
}
