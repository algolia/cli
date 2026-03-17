package apputil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// CreateApplicationWithRetry creates an application, retrying with a different
// region if the selected one has no available cluster.
func CreateApplicationWithRetry(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	accessToken string,
	region string,
	appName string,
) (*dashboard.Application, string, error) {
	cs := io.ColorScheme()

	for {
		if region == "" {
			var err error
			region, err = PromptRegion(io, client, accessToken)
			if err != nil {
				return nil, "", err
			}
		}

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

// CreateAndFetchApplication creates an application (with region retry) and
// fetches its full details including API keys.
func CreateAndFetchApplication(
	io *iostreams.IOStreams,
	client *dashboard.Client,
	accessToken, region, appName string,
) (*dashboard.Application, error) {
	if appName == "" {
		appName = "My First Application"
	}

	app, _, err := CreateApplicationWithRetry(io, client, accessToken, region, appName)
	if err != nil {
		return nil, err
	}

	io.StartProgressIndicatorWithLabel("Fetching API keys")
	appDetails, err := client.GetApplication(accessToken, app.ID)
	io.StopProgressIndicator()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch application details: %w", err)
	}

	return appDetails, nil
}

// ConfigureProfile creates a CLI profile from application details and
// optionally sets it as the default.
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

	if exists, existingAppID := cfg.ApplicationIDForProfile(profileName); exists && existingAppID != appDetails.ID {
		profileName = strings.ToLower(appDetails.Name + "-" + appDetails.ID)
	}

	profile := config.Profile{
		Name:          profileName,
		ApplicationID: appDetails.ID,
		APIKey:        appDetails.APIKey,
		Default:       setDefault,
	}

	if err := profile.Add(); err != nil {
		return err
	}

	if setDefault {
		if err := cfg.SetDefaultProfile(profileName); err != nil {
			fmt.Fprintf(io.ErrOut, "%s Could not set default profile: %s\n", cs.WarningIcon(), err)
		}
	}

	if io.IsStdoutTTY() {
		fmt.Fprintf(io.Out, "%s Profile %q configured for application %s.\n",
			cs.SuccessIcon(), profileName, cs.Bold(appDetails.ID))
	}

	return nil
}
