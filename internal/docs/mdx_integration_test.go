package docs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/internal/docs"
	"github.com/algolia/cli/pkg/cmd/root"
	clitest "github.com/algolia/cli/test"
)

func TestGenMdxTreeSupportsCurrentCommandTree(t *testing.T) {
	f, _ := clitest.NewFactory(false, nil, nil, "")
	rootCmd := root.NewRootCmd(f)

	dir := t.TempDir()
	require.NoError(t, docs.GenMdxTree(rootCmd, dir))

	rootContent := readTestFile(t, filepath.Join(dir, "index.mdx"))
	require.Contains(t, rootContent, "algolia search MOVIES --query \"toy story\"")

	logoutContent := readTestFile(t, filepath.Join(dir, "auth", "logout.mdx"))
	require.Contains(t, logoutContent, "algolia auth logout")
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content)
}
