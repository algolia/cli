package docs

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const mdxDocsRootSlug = "doc/tools/cli/commands"

//go:embed mdx.tpl
var mdxTemplate string

type mdxPage struct {
	Command
	Slug     string
	SubPages []mdxPage
}

func GenMdxTree(cmd *cobra.Command, dir string) error {
	tpl, err := template.New("mdx.tpl").Funcs(template.FuncMap{
		"getExamples": func(cmd Command) []Example {
			return cmd.ExamplesList()
		},
	}).Parse(mdxTemplate)
	if err != nil {
		return err
	}

	return writeMdxPageTree(tpl, dir, newMdxPage(describeCommand(cmd)))
}

func newMdxPage(cmd Command) mdxPage {
	page := mdxPage{
		Command: cmd,
		Slug:    buildMdxSlug(commandPathParts(cmd)),
	}

	for _, subCommand := range cmd.SubCommands {
		page.SubPages = append(page.SubPages, newMdxPage(subCommand))
	}

	return page
}

func writeMdxPageTree(tpl *template.Template, dir string, page mdxPage) error {
	filename := filepath.Join(dir, commandOutputDir(page.Command), "index.mdx")
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tpl.Execute(file, page); err != nil {
		return err
	}

	for _, subPage := range page.SubPages {
		if err := writeMdxPageTree(tpl, dir, subPage); err != nil {
			return err
		}
	}

	return nil
}

func commandPathParts(cmd Command) []string {
	parts := strings.Fields(cmd.Name)
	if len(parts) <= 1 {
		return nil
	}

	return parts[1:]
}

func commandOutputDir(cmd Command) string {
	return filepath.Join(commandPathParts(cmd)...)
}

func buildMdxSlug(parts []string) string {
	if len(parts) == 0 {
		return mdxDocsRootSlug
	}

	return mdxDocsRootSlug + "/" + strings.Join(parts, "/")
}
