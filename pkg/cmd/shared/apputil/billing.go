package apputil

import (
	"fmt"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/iostreams"
	pkgopen "github.com/algolia/cli/pkg/open"
	"github.com/algolia/cli/pkg/prompt"
)

func OfferBilling(
	io *iostreams.IOStreams,
	browser func(string) error,
	dashboardURL string,
	plan dashboard.Plan,
) error {
	cs := io.ColorScheme()
	url := dashboardURL + "/account/billing/details"

	fmt.Fprintf(
		io.Out,
		"\nThe %s plan requires a payment method on file before a paid application can be provisioned.\nThe CLI can't collect card details.\n",
		cs.Bold(plan.Name),
	)

	if io.CanPrompt() && io.IsStdoutTTY() {
		if browser == nil {
			browser = pkgopen.Browser
		}
		open := true
		if err := prompt.Confirm("Open the billing page to add a payment method?", &open); err != nil {
			return err
		}
		if !open {
			return nil
		}
		fmt.Fprintf(io.Out, "Opening %s\n", cs.Bold(url))
		return browser(url)
	}

	fmt.Fprintf(
		io.Out,
		"Add a payment method here, then re-run with --plan %s:\n%s\n",
		plan.ID,
		url,
	)
	return fmt.Errorf("the %q plan requires a payment method; none is on file", plan.Name)
}
