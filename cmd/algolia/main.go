package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/errs"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"

	"github.com/algolia/algolia-cli/pkg/cmd/factory"
	"github.com/algolia/algolia-cli/pkg/cmd/root"
	"github.com/algolia/algolia-cli/pkg/cmdutil"
	"github.com/algolia/algolia-cli/pkg/config"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	hasDebug := os.Getenv("DEBUG") != ""

	cfg := config.Config{}
	cfg.InitConfig()
	cmdFactory := factory.New(&cfg)

	stderr := cmdFactory.IOStreams.ErrOut

	if !cmdFactory.IOStreams.ColorEnabled() {
		surveyCore.DisableColor = true
	} else {
		// override survey's poor choice of color
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
			case "white":
				if cmdFactory.IOStreams.ColorSupport256() {
					return fmt.Sprintf("\x1b[%d;5;%dm", 38, 242)
				}
				return ansi.ColorCode("default")
			default:
				return ansi.ColorCode(style)
			}
		}
	}

	rootCmd := root.NewRootCmd(cmdFactory)
	authError := errors.New("authError")

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		if err == cmdutil.SilentError {
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

		// Handle Algolia client specific errors
		var algoliaErr errs.AlgoliaErr
		if errors.As(err, &algoliaErr) && algoliaErr.Status == 401 {
			fmt.Fprintln(stderr, "Try authenticating with: algolia login")
		} else if msg := algoliaErr.Message; msg != "" {
			fmt.Fprintln(stderr, msg)
		}

		return exitError
	}

	return exitOK
}

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
