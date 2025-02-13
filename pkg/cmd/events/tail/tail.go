package tail

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"

	algoliaInsights "github.com/algolia/algoliasearch-client-go/v4/algolia/insights"
	"github.com/algolia/cli/api/insights"
)

const (
	// DefaultRegion is the default region to use.
	DefaultRegion = algoliaInsights.US

	// Interval is the interval between each request to fetch events.
	Interval = 3 * time.Second
)

// TailOptions contains all the options for the `events tail` command.
type TailOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	Region     string
	PrintFlags *cmdutil.PrintFlags
}

// NewTailCmd returns a new command for tailing events.
func NewTailCmd(f *cmdutil.Factory, runF func(*TailOptions) error) *cobra.Command {
	opts := &TailOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:  "tail",
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"acls": "analytics",
		},
		Short: "Tail events",
		Long: heredoc.Doc(`
			Tail events from your Algolia application.

			By default, this command will tail events for the United States region.
			If your Analytics data is stored in a different region than the United States (e.g. Germany/Europe), you can specify the region using the --region (-r) flag.
		`),
		Example: heredoc.Doc(`
			# Tail events
			$ algolia events tail

			# Tail events for a specific region matching the Analytics region of your application
			$ algolia events tail -r de

			# Tail events and output them as JSON
			$ algolia events tail --output json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runTailCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.Region, "region", "r", string(DefaultRegion), "Region where your analytics data is stored and processed.")
	_ = cmd.RegisterFlagCompletionFunc("region", cmdutil.StringCompletionFunc(map[string]string{
		string(algoliaInsights.US): "United States",
		string(algoliaInsights.DE): "Germany (Europe)",
	}))

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runTailCmd(opts *TailOptions) error {
	appID, err := opts.Config.Profile().GetApplicationID()
	if err != nil {
		return err
	}
	apiKey, err := opts.Config.Profile().GetAPIKey()
	if err != nil {
		return err
	}

	// This is the CLIs custom augmented API client
	client, err := insights.NewClient(appID, apiKey, algoliaInsights.Region(opts.Region))
	if err != nil {
		return err
	}

	var p printers.Printer
	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err = opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
	} else if opts.IO.IsStdoutTTY() {
		fmt.Fprint(opts.IO.Out, "\nWaiting for events... Press Ctrl+C to stop.\n")
	}

	c := time.Tick(Interval)
	for t := range c {
		utc := t.UTC()
		events, err := client.GetEvents(utc.Add(-1*time.Second), utc, 1000)
		if err != nil {
			if strings.Contains(err.Error(), "The log processing region does not match") {
				cs := opts.IO.ColorScheme()
				errDetails := heredoc.Docf(`
					%s The Analytics storage region of your application does not match the region you specified (%s).
					Please specify the correct region using the --region (-r) flag.
					You can view the Analytics storage region of your application in the Algolia dashboard: https://www.algolia.com/infra/analytics
				`, cs.FailureIcon(), opts.Region)
				return errors.New(errDetails)
			}
		}

		for _, event := range events.Events {
			if p != nil {
				if err := p.Print(opts.IO, event); err != nil {
					return err
				}
			} else {
				if err := printEvent(opts.IO, event); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func printEvent(io *iostreams.IOStreams, event insights.EventWrapper) error {
	cs := io.ColorScheme()

	timeLayout := "2006-01-02 15:04:05"
	formatedTime := event.Event.Timestamp.Format(timeLayout)
	formatedTime = cs.Gray(formatedTime)

	colorizedStatus := cs.Green(fmt.Sprint(event.Status))
	if event.Status > 200 {
		colorizedStatus = cs.Red(fmt.Sprint(event.Status))
	}

	_, err := fmt.Fprintf(
		io.Out,
		"%s [%s] %s %s [%s] %s\n",
		cs.Bold(formatedTime),
		colorizedStatus,
		event.Event.EventType,
		cs.Bold(event.Event.Index),
		event.Event.EventName,
		event.Event.UserToken,
	)
	return err
}
