//go:generate go run ../../gen/gen_flags.go

package root

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/safeexec"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/internal/update"
	"github.com/algolia/cli/pkg/cmd/apikeys"
	"github.com/algolia/cli/pkg/cmd/art"
	"github.com/algolia/cli/pkg/cmd/dictionary"
	"github.com/algolia/cli/pkg/cmd/factory"
	"github.com/algolia/cli/pkg/cmd/indices"
	"github.com/algolia/cli/pkg/cmd/objects"
	"github.com/algolia/cli/pkg/cmd/open"
	"github.com/algolia/cli/pkg/cmd/profile"
	"github.com/algolia/cli/pkg/cmd/rules"
	"github.com/algolia/cli/pkg/cmd/search"
	"github.com/algolia/cli/pkg/cmd/settings"
	"github.com/algolia/cli/pkg/cmd/synonyms"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/telemetry"
	"github.com/algolia/cli/pkg/utils"
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
			$ algolia apikeys create --acl search
			$ algolia rules import MY_INDEX -f rules.json
		`),
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cmd.SetVersionTemplate(version.Template)
	cmd.SetUsageFunc(rootUsageFunc(f.IOStreams, cmd))
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(f, cmd, args)
	})
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.PersistentFlags().StringVarP(&f.Config.Profile().Name, "profile", "p", "", "The profile to use")
	_ = cmd.RegisterFlagCompletionFunc("profile", cmdutil.ConfiguredProfilesCompletionFunc(f))

	cmd.PersistentFlags().StringVarP(&f.Config.Profile().ApplicationID, "application-id", "", "", "The application ID")
	cmd.PersistentFlags().StringVarP(&f.Config.Profile().AdminAPIKey, "admin-api-key", "", "", "The admin API key")

	cmd.Flags().BoolP("version", "v", false, "Get the version of the Algolia CLI")

	// CLI related commands
	cmd.AddCommand(profile.NewProfileCmd(f))

	// Convenience commands
	cmd.AddCommand(open.NewOpenCmd(f))

	// API related commands
	cmd.AddCommand(search.NewSearchCmd(f))
	cmd.AddCommand(indices.NewIndicesCmd(f))
	cmd.AddCommand(objects.NewObjectsCmd(f))
	cmd.AddCommand(apikeys.NewAPIKeysCmd(f))
	cmd.AddCommand(settings.NewSettingsCmd(f))
	cmd.AddCommand(rules.NewRulesCmd(f))
	cmd.AddCommand(synonyms.NewSynonymsCmd(f))
	cmd.AddCommand(dictionary.NewDictionaryCmd(f))

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
	cmdFactory := factory.New(version.Version, &cfg)
	stderr := cmdFactory.IOStreams.ErrOut

	// Set up the update notifier.
	updateMessageChan := make(chan *update.ReleaseInfo)
	go func() {
		rel, err := checkForUpdate(cfg, version.Version)
		if err != nil && hasDebug {
			fmt.Fprintf(stderr, "Error checking for update: %s\n", err)
		}
		updateMessageChan <- rel
	}()

	// Set up the root command.
	rootCmd := NewRootCmd(cmdFactory)

	// Pre-command auth check and telemetry setup.
	authError := errors.New("authError")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmdutil.IsAuthCheckEnabled(cmd) {
			if err := cmdutil.CheckAuth(cfg); err != nil {
				fmt.Fprintf(stderr, "Authentication error: %s\n", err)
				fmt.Fprintln(stderr, "Please run `algolia profile add` to configure your first profile.")
				return authError
			}
		}

		if !cmdutil.ShouldTrackUsage(cmd) {
			return nil
		}

		// Initialize telemetry context.
		appID, err := cfg.Profile().GetApplicationID()
		if err != nil {
			appID = ""
		}
		telemetryMetadata := telemetry.GetEventMetadata(cmd.Context())
		telemetryMetadata.SetCobraCommandContext(cmd)
		telemetryMetadata.SetAppID(appID)
		telemetryMetadata.SetConfiguredApplicationsNb(len(cfg.ConfiguredProfiles()))

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

		go telemetryClient.Close() // flush telemetry events

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

	// If there is an update available, notify the user.
	newRelease := <-updateMessageChan
	if newRelease != nil {
		cs := cmdFactory.IOStreams.ColorScheme()
		isHomebrew := isUnderHomebrew(cmdFactory.Executable())
		fmt.Fprintf(stderr, "\n\n%s %s → %s\n",
			cs.Yellow("A new release of the Algolia CLI is available:"),
			cs.Cyan(strings.TrimPrefix(version.Version, "v")),
			cs.Cyan(strings.TrimPrefix(newRelease.Version, "v")))
		if isHomebrew {
			fmt.Fprintf(stderr, "To upgrade, run: %s\n", "brew update && brew upgrade algolia")
		}
		fmt.Fprintf(stderr, "%s\n\n",
			cs.Yellow(newRelease.URL))
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

func shouldCheckForUpdate() bool {
	if os.Getenv("ALGOLIA_NO_UPDATE_NOTIFIER") != "" {
		return false
	}
	return !utils.IsCI() && utils.IsTerminal(os.Stdout) && utils.IsTerminal(os.Stderr)
}

func checkForUpdate(cfg config.Config, currentVersion string) (*update.ReleaseInfo, error) {
	if !shouldCheckForUpdate() {
		return nil, nil
	}
	stateFilePath := filepath.Join(cfg.GetConfigFolder(os.Getenv("XDG_CONFIG_HOME")), "state.yml")
	client := http.Client{}
	return update.CheckForUpdate(&client, stateFilePath, currentVersion)
}

// Check whether the gh binary was found under the Homebrew prefix
func isUnderHomebrew(ghBinary string) bool {
	brewExe, err := safeexec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(ghBinary, brewBinPrefix)
}
