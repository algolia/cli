package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikey"
	"github.com/algolia/cli/pkg/cmd/application"
	"github.com/algolia/cli/pkg/cmd/index"
	"github.com/algolia/cli/pkg/cmd/objects"
	"github.com/algolia/cli/pkg/cmd/open"
	"github.com/algolia/cli/pkg/cmd/rule"
	"github.com/algolia/cli/pkg/cmd/settings"
	"github.com/algolia/cli/pkg/cmd/synonym"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/version"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "algolia <command> <subcommand> [flags]",
		Short: "Algolia CLI",
		Long:  "The official command-line tool to interact with Algolia.",

		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ algolia objects browse
			$ algolia apikey create --acl search
			$ algolia rule list
		`),
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.SetVersionTemplate(version.Template)

	cmd.PersistentFlags().StringVarP(&f.Config.Application.Name, "application", "a", "default", "The application to use")
	cmd.RegisterFlagCompletionFunc("application", cmdutil.ConfiguredApplicationsCompletionFunc(f))

	cmd.PersistentFlags().StringVarP(&f.Config.Application.ID, "application-id", "", "", "The application ID")
	cmd.PersistentFlags().StringVarP(&f.Config.Application.AdminAPIKey, "admin-api-key", "", "", "The admin API key")

	cmd.Flags().BoolP("version", "v", false, "Get the version of the Algolia CLI")

	// CLI related commands
	cmd.AddCommand(application.NewApplicationCmd(f))

	// Convenience commands
	cmd.AddCommand(open.NewOpenCmd(f))

	// API related commands
	cmd.AddCommand(index.NewIndexCmd(f))
	cmd.AddCommand(objects.NewObjectsCmd(f))
	cmd.AddCommand(apikey.NewAPIKeyCmd(f))
	cmd.AddCommand(settings.NewSettingsCmd(f))
	cmd.AddCommand(rule.NewRuleCmd(f))
	cmd.AddCommand(synonym.NewSynonymCmd(f))

	return cmd
}
