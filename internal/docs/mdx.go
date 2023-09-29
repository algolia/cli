package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func GenMdxFile(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		if c.Hidden || c.Name() == "help" {
			continue
		}
		command := newCommand(c)
		command.Slug = strings.ReplaceAll(command.Name, " ", "-")
		filename := filepath.Join(dir, fmt.Sprintf("%s.mdx", command.Slug))
		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		tpl, err := template.ParseFiles("internal/docs/mdx.tpl")
		if err != nil {
			return err
		}

		if c.HasAvailableSubCommands() {
			for _, s := range c.Commands() {
				sub := newCommand(s)
				if s.HasAvailableSubCommands() {
					for _, sus := range s.Commands() {
						sub.SubCommands = append(sub.SubCommands, newCommand(sus))
					}
				}
				command.SubCommands = append(command.SubCommands, sub)
			}
		}
		err = tpl.Execute(file, command)
		if err != nil {
			return err
		}
	}

	return nil
}
