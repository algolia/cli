package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func GenYamlTree(cmd *cobra.Command, dir string) error {
	commands := getCommands(cmd)

	for _, c := range commands {
		commandPath := strings.ReplaceAll(c.Name, " ", "_")
		filename := filepath.Join(dir, fmt.Sprintf("%s.yml", commandPath))
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write a comment to the file, ignoring vale errors
		_, err = f.WriteString("# <!-- vale off -->\n")
		if err != nil {
			return err
		}

		encoder := yaml.NewEncoder(f)
		defer encoder.Close()

		err = encoder.Encode(c)
		if err != nil {
			return err
		}
	}

	return nil
}
