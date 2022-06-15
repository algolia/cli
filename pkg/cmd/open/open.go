package open

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/printers"
)

var nameURLmap = map[string]string{
	"api":       "https://www.algolia.com/doc/api-reference/rest-api/",
	"dashboard": "https://www.algolia.com/apps%s/dashboard",
	"codex":     "https://www.algolia.com/developers/code-exchange/",
	"devhub":    "https://www.algolia.com/developers/",
	"docs":      "https://algolia.com/doc/",
}

func openNames() []string {
	keys := make([]string, 0, len(nameURLmap))
	for k := range nameURLmap {
		keys = append(keys, k)
	}

	return keys
}

func getLongestShortcut(shortcuts []string) int {
	longest := 0
	for _, s := range shortcuts {
		if len(s) > longest {
			longest = len(s)
		}
	}

	return longest
}

func padName(name string, length int) string {
	difference := length - len(name)

	var b strings.Builder

	fmt.Fprint(&b, name)

	for i := 0; i < difference; i++ {
		fmt.Fprint(&b, " ")
	}

	return b.String()
}

// OpenOptions represents the options for the open command
type OpenOptions struct {
	config *config.Config
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
	var applicationID string
	app := opts.config.Application
	if app.ID == "" {
		applicationID = ""
	} else {
		applicationID = "/" + app.ID
	}

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

		for _, shortcut := range shortcuts {
			url := nameURLmap[shortcut]
			if strings.Contains(url, "%s") {
				url = fmt.Sprintf(url, applicationID)
			}

			table.AddField(shortcut, nil, nil)
			table.AddField(url, nil, nil)
			table.EndRow()
		}

		return table.Render()
	}

	var err error
	if url, ok := nameURLmap[opts.Shortcut]; ok {
		if strings.Contains(url, "%s") {
			err = open.Browser(fmt.Sprintf(url, applicationID))
		} else {
			err = open.Browser(url)
		}

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported open command, given: %s", opts.Shortcut)
	}

	return nil
}
