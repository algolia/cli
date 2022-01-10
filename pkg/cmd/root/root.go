package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algolia-cli/pkg/cmd/apikeys"
	"github.com/algolia/algolia-cli/pkg/cmd/indices"
	"github.com/algolia/algolia-cli/pkg/cmd/login"
	"github.com/algolia/algolia-cli/pkg/cmd/settings"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/algolia/algolia-cli/pkg/version"
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
			$ algolia instantsearch my-app
		`),
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.SetVersionTemplate(version.Template)

	cmd.PersistentFlags().StringVarP(&f.Config.Profile.ProfileName, "profile", "p", "default", "The profile name to read from for config")

	cmd.Flags().BoolP("version", "v", false, "Get the version of the Algolia CLI")

	// Child commands
	cmd.AddCommand(login.NewLoginCmd(f))
	cmd.AddCommand(indices.NewIndicesCmd(f))
	cmd.AddCommand(apikeys.NewAPIKeysCmd(f))
	cmd.AddCommand(settings.NewSettingsCmd(f))

	return cmd
}
