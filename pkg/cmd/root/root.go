//go:generate go run ../../gen/gen_flags.go

package root

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikey"
	"github.com/algolia/cli/pkg/cmd/application"
	"github.com/algolia/cli/pkg/cmd/art"
	"github.com/algolia/cli/pkg/cmd/factory"
	"github.com/algolia/cli/pkg/cmd/index"
	"github.com/algolia/cli/pkg/cmd/objects"
	"github.com/algolia/cli/pkg/cmd/open"
	"github.com/algolia/cli/pkg/cmd/rules"
	"github.com/algolia/cli/pkg/cmd/search"
	"github.com/algolia/cli/pkg/cmd/settings"
	"github.com/algolia/cli/pkg/cmd/synonyms"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/pkg/version"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "algolia <command> <subcommand> [flags]",
		Version: version.Version,
		Short:   "Algolia CLI",
		Long:    "The official command-line tool to interact with Algolia.",

		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ algolia search MY_INDEX --query "foo"
			$ algolia objects browse MY_INDEX
			$ algolia apikey create --acl search
			$ algolia rule import MY_INDEX -f rules.json
		`),
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.SetVersionTemplate(version.Template)
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		return rootUsageFunc(f.IOStreams.Out, cmd)
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(f, cmd, args)
	})
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.PersistentFlags().StringVarP(&f.Config.Application.Name, "application", "a", "", "The application to use")
	_ = cmd.RegisterFlagCompletionFunc("application", cmdutil.ConfiguredApplicationsCompletionFunc(f))

	cmd.PersistentFlags().StringVarP(&f.Config.Application.ID, "application-id", "", "", "The application ID")
	cmd.PersistentFlags().StringVarP(&f.Config.Application.AdminAPIKey, "admin-api-key", "", "", "The admin API key")

	cmd.Flags().BoolP("version", "v", false, "Get the version of the Algolia CLI")

	// CLI related commands
	cmd.AddCommand(application.NewApplicationCmd(f))

	// Convenience commands
	cmd.AddCommand(open.NewOpenCmd(f))

	// API related commands
	cmd.AddCommand(search.NewSearchCmd(f))
	cmd.AddCommand(index.NewIndexCmd(f))
	cmd.AddCommand(objects.NewObjectsCmd(f))
	cmd.AddCommand(apikey.NewAPIKeyCmd(f))
	cmd.AddCommand(settings.NewSettingsCmd(f))
	cmd.AddCommand(rules.NewRulesCmd(f))
	cmd.AddCommand(synonyms.NewSynonymsCmd(f))

	// ??? related commands
	cmd.AddCommand(art.NewArtCmd(f))

	return cmd
}

func Execute() exitCode {
	hasDebug := os.Getenv("DEBUG") != ""
	hasTelemetry := os.Getenv("ALGOLIA_CLI_TELEMETRY") != "0"

	// Set up the command factory.
	cfg := config.Config{}
	cfg.InitConfig()
	cmdFactory := factory.New(&cfg)
	stderr := cmdFactory.IOStreams.ErrOut

	// Set up the root command.
	rootCmd := NewRootCmd(cmdFactory)

	// Pre-command auth check and telemetry setup.
	authError := errors.New("authError")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmdutil.IsAuthCheckEnabled(cmd) {
			if err := cmdutil.CheckAuth(cfg); err != nil {
				fmt.Fprintf(stderr, "Authentication error: %s\n", err)
				fmt.Fprintln(stderr, "Please run `algolia application add` to configure your first application.")
				return authError
			}
		}

		if !cmdutil.ShouldTrackUsage(cmd) {
			return nil
		}

		// Initialize telemetry context.
		appID, err := cfg.Application.GetID()
		if err != nil {
			appID = ""
		}
		telemetryMetadata := telemetry.GetEventMetadata(cmd.Context())
		telemetryMetadata.SetCobraCommandContext(cmd)
		telemetryMetadata.SetAppID(appID)
		telemetryMetadata.SetConfiguredApplicationsNb(len(cfg.ConfiguredApplications()))

		ctx := cmd.Context()
		telemetryClient := telemetry.GetTelemetryClient(ctx)

		// Identify the user.
		err = telemetryClient.Identify(ctx)
		if err != nil && hasDebug {
			fmt.Fprintf(stderr, "Failed to identify user: %s\n", err)
			return err
		}

		// Send telemetry.
		err = telemetryClient.Track(ctx, "Command Invoked")
		if err != nil && hasDebug {
			fmt.Fprintf(stderr, "Error tracking telemetry: %s\n", err)
		}

		return nil
	}

	// Command context is used to pass information to the telemetry client.
	ctx, err := createContext(rootCmd, stderr, hasDebug, hasTelemetry)
	if err != nil {
		printError(stderr, err, rootCmd, hasDebug)
		return exitError
	}

	// Run the command.
	cmd, err := rootCmd.ExecuteContextC(ctx)

	if cmdutil.ShouldTrackUsage(cmd) {
		// Post-command telemetry
		ctx = cmd.Context()
		telemetryClient := telemetry.GetTelemetryClient(ctx)
		telemetryErr := telemetryClient.Track(ctx, "Command Finished")
		if telemetryErr != nil && hasDebug {
			fmt.Fprintf(stderr, "Error tracking telemetry: %s\n", err)
		}
		telemetryErr = telemetryClient.Close() // flush telemetry events
		if telemetryErr != nil && hasDebug {
			fmt.Fprintf(stderr, "Error closing telemetry client: %s\n", err)
		}
	}

	// Handle eventual errors.
	if err != nil {
		if err == cmdutil.ErrSilent {
			return exitError
		} else if cmdutil.IsUserCancellation(err) {
			if errors.Is(err, terminal.InterruptErr) {
				// ensure the next shell prompt will start on its own line
				fmt.Fprint(stderr, "\n")
			}
			return exitCancel
		} else if errors.Is(err, authError) {
			return exitAuth
		}

		printError(stderr, err, cmd, hasDebug)
		return exitError
	}

	return exitOK
}

// createContext creates a context with telemetry.
func createContext(cmd *cobra.Command, stderr io.Writer, hasDebug bool, hasTelemetry bool) (context.Context, error) {
	ctx := context.Background()
	telemetryMetadata := telemetry.NewEventMetadata()
	updatedCtx := telemetry.WithEventMetadata(ctx, telemetryMetadata)

	var telemetryClient telemetry.TelemetryClient
	var err error
	if hasTelemetry {
		telemetryClient, err = telemetry.NewAnalyticsTelemetryClient(hasDebug)
		// Fail silently if telemetry is not available unless in debug mode.
		if err != nil && hasDebug {
			fmt.Fprintf(stderr, "Error creating telemetry client: %s\n", err)
			return nil, err
		}
	} else {
		telemetryClient = &telemetry.NoOpTelemetryClient{}
	}
	contextWithTelemetry := telemetry.WithTelemetryClient(updatedCtx, telemetryClient)
	return contextWithTelemetry, nil
}

// printError prints an error to the stderr, with additional information if applicable.
func printError(out io.Writer, err error, cmd *cobra.Command, debug bool) {
	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		fmt.Fprintf(out, "error connecting to %s\n", dnsError.Name)
		if debug {
			fmt.Fprintln(out, dnsError)
		}
		fmt.Fprintln(out, "check your internet connection or https://status.algolia.com")
		return
	}

	fmt.Fprintln(out, err)

	var flagError *cmdutil.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, cmd.UsageString())
	}
}
