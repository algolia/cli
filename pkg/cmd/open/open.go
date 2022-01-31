package open

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/open"
)

var nameURLmap = map[string]string{
	"api":       "https://www.algolia.com/doc/api-reference/rest-api/",
	"apiref":    "https://www.algolia.com/doc/api-reference/api-parameters/",
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

type openCmd struct {
	cfg *config.Config
	cmd *cobra.Command
}

func newOpenCmd(config *config.Config) *openCmd {
	oc := &openCmd{
		cfg: config,
	}
	oc.cmd = &cobra.Command{
		Use:       "open",
		ValidArgs: openNames(),
		Short:     "Quickly open Algolia pages",
		Long: `The open command provices shortcuts to quickly let you open pages to Algolia with
in your browser. A full list of support shortcuts can be seen with 'algolia open --list'`,
		Example: `algolia open --list
  algolia open api
  algolia open docs
  algolia open dashboard/webhooks
  algolia open dashboard/billing`,
		RunE: oc.runOpenCmd,
	}

	oc.cmd.Flags().Bool("list", false, "List all supported shortcuts")

	return oc
}

func (oc *openCmd) runOpenCmd(cmd *cobra.Command, args []string) error {
	list, err := cmd.Flags().GetBool("list")
	if err != nil {
		return err
	}

	applicationID, err := oc.cfg.App.GetID()
	if err != nil {
		applicationID = ""
	} else {
		applicationID = "/" + applicationID
	}

	if list || len(args) == 0 {
		fmt.Println("open quickly opens Algolia pages. To use, run 'algolia open <shortcut>'.")
		fmt.Println("open supports the following shortcuts:")
		fmt.Println()

		shortcuts := openNames()
		sort.Strings(shortcuts)

		longest := getLongestShortcut(shortcuts)

		fmt.Printf("%s%s\n", padName("shortcut", longest), "    url")
		fmt.Printf("%s%s\n", padName("--------", longest), "    ---------")

		for _, shortcut := range shortcuts {

			url := nameURLmap[shortcut]
			if strings.Contains(url, "%s") {
				url = fmt.Sprintf(url, applicationID)
			}

			paddedName := padName(shortcut, longest)
			fmt.Printf("%s => %s\n", paddedName, url)
		}

		return nil
	}

	if url, ok := nameURLmap[args[0]]; ok {
		if strings.Contains(url, "%s") {
			err = open.Browser(fmt.Sprintf(url, applicationID))
		} else {
			err = open.Browser(url)
		}

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported open command, given: %s", args[0])
	}

	return nil
}
