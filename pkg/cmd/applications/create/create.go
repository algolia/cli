package create

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/provisionning"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

const waitTimeout = 240 * time.Second // 4 minutes

// CreateOptions are used to create a new application.
type CreateOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	ProvisionningClient func() (*provisionning.Client, error)

	Name   string
	Region string
	Plan   string

	PrintFlags *cmdutil.PrintFlags
}

// availablePlans and availableRegions are the only available plans and regions for now.
var availablePlans = []string{provisionning.PlanV8Build}
var availableRegions = []string{"EU", "UK", "USC", "USE", "USW"}

// NewCreateCmd creates and returns a create command for Applications.
func NewCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                  f.IOStreams,
		Config:              f.Config,
		ProvisionningClient: f.ProvisionningClient,
		PrintFlags:          cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new application",
		Long: heredoc.Doc(`
			Create a new application.
			
			Only free plan (V8.5 build) is available for now.
		`),

		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate plan
			if opts.Plan != provisionning.PlanV8Build {
				return fmt.Errorf("invalid plan: %s", opts.Plan)
			}

			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "The name of the application")
	cmd.Flags().StringVarP(&opts.Region, "region", "r", "", "The region of the application")
	cmd.MarkFlagRequired("region")
	cmd.RegisterFlagCompletionFunc("region", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return availableRegions, cobra.ShellCompDirectiveDefault
	})

	cmd.Flags().StringVar(&opts.Plan, "plan", provisionning.PlanV8Build, "The plan of the application. Always set to `v8.5-plg-build` for now.")
	cmd.RegisterFlagCompletionFunc("plan", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return availablePlans, cobra.ShellCompDirectiveDefault
	})

	return cmd
}

// runCreateCmd executes the create command
func runCreateCmd(opts *CreateOptions) error {
	client, err := opts.ProvisionningClient()
	if err != nil {
		return err
	}

	name := opts.Name
	if name != "" {
		name = " '" + name + "'"
	}
	opts.IO.StartProgressIndicatorWithLabel(
		fmt.Sprintf("Creating application%s in region %s (plan: %s)", name, opts.Region, opts.Plan),
	)

	params := provisionning.FreeApplicationCreationRequest{
		Region: opts.Region,
		Name:   opts.Name,
	}
	res, err := client.CreateFreeApplication(params)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.UpdateProgressIndicatorLabel("Waiting for the application to be ready...")
	err = waitForApplicationCreation(client, res.ID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Printf("%s Application%s %s successfully created in region %s [plan: %s]\n", cs.SuccessIcon(), name, cs.Bold(res.ID), cs.Bold(opts.Region), opts.Plan)
	} else {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}

		return p.Print(opts.IO, res)
	}
	return nil
}

// waitForApplicationCreation waits for the application to be ready before returning
// It polls the application status every 2 seconds until it's ready or the timeout is reached
func waitForApplicationCreation(client *provisionning.Client, applicationID string) error {
	c := time.Tick(2 * time.Second)
	timeout := time.After(waitTimeout)
	for {
		select {
		case <-c:
			status, err := client.GetApplicationCreationStatus(applicationID)
			if err != nil {
				return err
			}

			if status == "ready" {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout reached while waiting for the application to be ready")
		}
	}
}
