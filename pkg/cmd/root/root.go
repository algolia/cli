package root

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikey"
	"github.com/algolia/cli/pkg/cmd/application"
	"github.com/algolia/cli/pkg/cmd/indices"
	"github.com/algolia/cli/pkg/cmd/rule"
	"github.com/algolia/cli/pkg/cmd/settings"
	"github.com/algolia/cli/pkg/cmd/synonym"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/version"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "algolia",
		Short: "Algolia CLI",
		Long:  "The official command-line tool to interact with Algolia.",

		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ algolia indices list
			$ algolia apikeys create --acl search
			$ algolia rules export TEST_index > rules.json
			$ algolia rules import TEST_index -F rules.json
			$ algolia settings set TEST_index "attributesForFaceting": ["category"]
		`),
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.SetVersionTemplate(version.Template)

	cmd.PersistentFlags().StringVarP(&f.Config.App.Name, "application", "a", "default", "The application to use")
	cmd.RegisterFlagCompletionFunc("application", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		profiles := f.Config.Applications()
		completions := make([]string, 0, len(profiles))

		// We want to show the profile name and the Application ID as the description.
		// https://github.com/spf13/cobra/blob/master/shell_completions.md#descriptions-for-completions
		for profileName, AppID := range profiles {
			completions = append(completions, fmt.Sprintf("%s\t%s", profileName, AppID))
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Flags().BoolP("version", "v", false, "Get the version of the Algolia CLI")

	// Child commands
	cmd.AddCommand(application.NewApplicationCmd(f))

	cmd.AddCommand(indices.NewIndicesCmd(f))
	cmd.AddCommand(apikey.NewAPIKeyCmd(f))
	cmd.AddCommand(settings.NewSettingsCmd(f))
	cmd.AddCommand(rule.NewRuleCmd(f))
	cmd.AddCommand(synonym.NewSynonymCmd(f))

	return cmd
}
