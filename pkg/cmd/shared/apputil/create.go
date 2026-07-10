package apputil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/keychain"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/telemetry"
)

// CreateApplicationWithRetry creates an application, retrying with a different
// region if the selected one has no available cluster.
//
// The optional tracker (nil-safe) records which step the flow is in, so the
// telemetry of the calling flow can tell where the user stopped.
func CreateApplicationWithRetry(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	accessToken string,
	region string,
	appName string,
	tracker *telemetry.FlowTracker,
) (*dashboard.Application, string, error) {
	cs := io.ColorScheme()

	for {
		if region == "" {
			if !io.CanPrompt() {
				return nil, "", fmt.Errorf(
					"no region specified; pass --region (e.g. EU, UK, USC, USE, USW)",
				)
			}
			tracker.SetStep(telemetry.StepRegion)
			var err error
			region, err = PromptRegion(io, client, accessToken)
			if err != nil {
				return nil, "", err
			}
		}

		tracker.SetStep(telemetry.StepAPICall)
		io.StartProgressIndicatorWithLabel("Creating application")
		app, err := client.CreateApplication(accessToken, region, appName)
		io.StopProgressIndicator()

		if err == nil {
			fmt.Fprintf(io.Out, "%s Application %s created in region %q\n",
				cs.SuccessIcon(), cs.Bold(app.ID), region)
			return app, region, nil
		}

		var clusterErr *dashboard.ErrClusterUnavailable
		if errors.As(err, &clusterErr) {
			fmt.Fprintf(io.Out, "%s No cluster available in region %q. Please select another region.\n",
				cs.WarningIcon(), region)
			region = ""

			if !io.CanPrompt() {
				return nil, "", fmt.Errorf("no cluster available in region %q — try a different --region", clusterErr.Region)
			}
			continue
		}

		return nil, "", fmt.Errorf("application creation failed: %w", err)
	}
}

// PromptRegion fetches regions from the API and prompts the user to select one.
func PromptRegion(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	accessToken string,
) (string, error) {
	io.StartProgressIndicatorWithLabel("Fetching regions")
	regions, err := client.ListRegions(accessToken)
	io.StopProgressIndicator()
	if err != nil {
		return "", fmt.Errorf("failed to fetch regions: %w", err)
	}

	if len(regions) == 0 {
		return "", fmt.Errorf("no regions available")
	}

	regionOptions := make([]string, len(regions))
	for i, r := range regions {
		if r.Name != "" {
			regionOptions[i] = fmt.Sprintf("%s (%s)", r.Code, r.Name)
		} else {
			regionOptions[i] = r.Code
		}
	}

	var selected int
	err = prompt.SurveyAskOne(
		&survey.Select{
			Message: "Region:",
			Options: regionOptions,
		},
		&selected,
	)
	if err != nil {
		return "", err
	}

	return regions[selected].Code, nil
}

// PromptName prompts the user for the application name.
func PromptName() (string, error) {
	defaultName := "My First Application"

	var name string
	err := prompt.SurveyAskOne(
		&survey.Input{
			Message: "Name:",
			Default: defaultName,
		},
		&name,
	)
	if err != nil {
		return "", err
	}

	if name == "" {
		name = defaultName
	}
	return name, nil
}

// CreateAndFetchApplication creates an application (with region retry) and
// generates an API key for it. It returns the region the application was
// actually created in, which may differ from the requested one.
func CreateAndFetchApplication(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	accessToken, region, appName string,
	tracker *telemetry.FlowTracker,
) (*dashboard.Application, string, error) {
	app, createdRegion, err := CreateApplicationWithRetry(io, client, accessToken, region, appName, tracker)
	if err != nil {
		return nil, "", err
	}

	tracker.SetStep(telemetry.StepAPIKey)
	if err := EnsureAPIKey(io, client, accessToken, app); err != nil {
		return nil, "", err
	}

	return app, createdRegion, nil
}

// EnsureAPIKey generates a write API key for the application.
// Callers should skip this if the local profile already has a key.
func EnsureAPIKey(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	accessToken string,
	app *dashboard.Application,
) error {
	cs := io.ColorScheme()
	io.StartProgressIndicatorWithLabel("Generating API key")
	created, err := client.CreateAPIKey(accessToken, app.ID, dashboard.WriteACL, "Algolia CLI")
	io.StopProgressIndicator()
	if err != nil {
		return fmt.Errorf("failed to generate API key: %w", err)
	}

	app.APIKey = created.Value
	app.APIKeyUUID = created.UUID
	fmt.Fprintf(io.Out, "%s API key generated for application %s\n",
		cs.SuccessIcon(), cs.Bold(app.ID))
	return nil
}

// ConfigureProfile persists the application credentials in the new model
// (state.toml + OS keychain) and optionally makes it the current application.
func ConfigureProfile(
	io *iostreams.IOStreams,
	cfg config.IConfig,
	appDetails *dashboard.Application,
	profileName string,
	setDefault bool,
) error {
	cs := io.ColorScheme()

	if profileName == "" {
		profileName = appDetails.Name
	}
	if profileName == "" {
		profileName = appDetails.ID
	}
	profileName = strings.ToLower(profileName)

	// Another application already carries this alias: derive a unique one.
	if otherID, ok := cfg.ApplicationIDByAlias(profileName); ok && otherID != appDetails.ID {
		profileName = strings.ToLower(appDetails.Name + "-" + appDetails.ID)
	}

	if err := cfg.SaveApplication(
		appDetails.ID, profileName, appDetails.APIKeyUUID, appDetails.APIKey, setDefault,
	); err != nil {
		return err
	}

	if io.IsStdoutTTY() {
		fmt.Fprintf(io.Out, "%s Application %s configured (alias %q).\n",
			cs.SuccessIcon(), cs.Bold(appDetails.ID), profileName)
	}

	return nil
}

// ReuseExistingAPIKey looks for an API key already stored for the application
// (keychain first, then legacy config.toml profiles). If found, it sets
// app.APIKey and returns true so callers skip creating a new key. A key reused
// from a legacy profile has no known UUID, so api_key_uuid is left as-is.
func ReuseExistingAPIKey(cfg config.IConfig, app *dashboard.Application) bool {
	if secrets, err := keychain.LoadAppSecrets(app.ID); err == nil && secrets != nil &&
		secrets.APIKey != "" {
		app.APIKey = secrets.APIKey
		return true
	}

	for _, p := range cfg.ConfiguredProfiles() {
		if p.ApplicationID == app.ID && p.APIKey != "" {
			app.APIKey = p.APIKey
			return true
		}
	}

	return false
}
