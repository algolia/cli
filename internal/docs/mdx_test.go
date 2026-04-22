package docs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestGetCommandsRecursivelyIncludesNestedCommands(t *testing.T) {
	root := &cobra.Command{Use: "algolia", Short: "Algolia CLI"}
	one := &cobra.Command{Use: "one", Short: "Level one"}
	two := &cobra.Command{Use: "two", Short: "Level two"}
	three := &cobra.Command{Use: "three", Short: "Level three"}
	four := &cobra.Command{
		Use:   "four",
		Short: "Level four",
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}

	three.AddCommand(four)
	two.AddCommand(three)
	one.AddCommand(two)
	root.AddCommand(one)

	commands := getCommands(root)
	require.Len(t, commands, 1)
	require.Len(t, commands[0].SubCommands, 1)
	require.Len(t, commands[0].SubCommands[0].SubCommands, 1)
	require.Len(t, commands[0].SubCommands[0].SubCommands[0].SubCommands, 1)
	require.Equal(
		t,
		"algolia one two three four",
		commands[0].SubCommands[0].SubCommands[0].SubCommands[0].Name,
	)
}

func TestGenMdxTreeWritesNestedCommandPages(t *testing.T) {
	root := &cobra.Command{Use: "algolia", Short: "Algolia CLI"}
	events := &cobra.Command{Use: "events", Short: "Manage events"}
	sources := &cobra.Command{Use: "sources", Short: "Manage event sources"}
	list := &cobra.Command{
		Use:     "list",
		Short:   "List event sources",
		Example: "# List sources\n$ algolia events sources list",
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}
	list.Flags().StringP("format", "F", "json", "Output format")

	sources.AddCommand(list)
	events.AddCommand(sources)
	root.AddCommand(events)

	dir := t.TempDir()
	require.NoError(t, GenMdxTree(root, dir))

	rootContent := readTestFile(t, filepath.Join(dir, "index.mdx"))
	require.Contains(t, rootContent, "slug: tools/cli/commands")
	require.Contains(t, rootContent, "[`algolia events`](/tools/cli/commands/events)")

	eventsContent := readTestFile(t, filepath.Join(dir, "events", "index.mdx"))
	require.Contains(t, eventsContent, "slug: tools/cli/commands/events")
	require.Contains(t, eventsContent, "[`algolia events sources`](/tools/cli/commands/events/sources)")

	listContent := readTestFile(t, filepath.Join(dir, "events", "sources", "list", "index.mdx"))
	require.Contains(t, listContent, "slug: tools/cli/commands/events/sources/list")
	require.Contains(t, listContent, "`algolia events sources list [flags]`")
	require.Contains(t, listContent, "`-F`, `--format`")
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content)
}
