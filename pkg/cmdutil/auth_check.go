package cmdutil

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/config"
)

func DisableAuthCheck(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}

	cmd.Annotations["skipAuthCheck"] = "true"
}

func CheckAuth(cfg config.Config) bool {
	_, err := cfg.Application.GetID()
	if err != nil {
		return false
	}
	_, err = cfg.Application.GetAdminAPIKey()
	return err == nil
}

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help", "powershell", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
		return false
	}

	if cmd.Parent() != nil && cmd.Parent().Name() == "completion" {
		return false
	}

	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}
