package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func GenMdxTree(cmd *cobra.Command, dir string) error {
	tpl, err := template.New("mdx.tpl").Funcs(template.FuncMap{
		"getExamples": func(cmd Command) []Example {
			return cmd.ExamplesList()
		},
	}).ParseFiles("internal/docs/mdx.tpl")
	if err != nil {
		return err
	}

	commands := getCommands(cmd)

	for _, c := range commands {
		c.Slug = strings.ReplaceAll(c.Name, " ", "-")
		filename := filepath.Join(dir, fmt.Sprintf("%s.mdx", c.Slug))
		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		err = tpl.Execute(file, c)
		if err != nil {
			return err
		}
	}

	return nil
}
