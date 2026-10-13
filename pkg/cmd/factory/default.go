package factory

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/call"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/pkg/iostreams"
)

func New(appVersion string, cfg config.IConfig) *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config:         cfg,
		ExecutableName: "gh",
	}
	f.IOStreams = ioStreams(f)
	f.Prompter = interactive.NewSurveyPrompter(f.IOStreams)
	f.SearchClient = searchClient(f, appVersion)
	f.CrawlerClient = crawlerClient(f)
	f.CompositionClient = compositionClient(f, appVersion)
	f.AgentStudioClient = agentStudioClient(f, appVersion)
	f.AgentStudioAPIClient = agentStudioAPIClient(f, appVersion)

	return f
}

func ioStreams(_ *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()
	return io
}

func searchClient(f *cmdutil.Factory, appVersion string) func() (*search.APIClient, error) {
	return func() (*search.APIClient, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		userAgent, err := getUserAgentInfo(appID, apiKey, appVersion)
		if err != nil {
			return nil, err
		}
		if userAgent == "" {
			return nil, fmt.Errorf("user agent must not be empty")
		}

		clientConf := search.SearchConfiguration{
			Configuration: transport.Configuration{
				AppID:                           appID,
				ApiKey:                          apiKey,
				UserAgent:                       userAgent,
				ExposeIntermediateNetworkErrors: true,
			},
		}

		// Read custom hosts from flags, environment, or profile, or use default ones
		hosts := GetStatefulHosts(f.Config.Profile().GetSearchHosts())
		if len(hosts) > 0 {
			clientConf.Configuration.Hosts = hosts
		}

		return search.NewClientWithConfig(clientConf)
	}
}

func agentStudioClient(f *cmdutil.Factory, appVersion string) func() (*agentstudio.Client, error) {
	return func() (*agentstudio.Client, error) {
		profile := f.Config.Profile()
		appID, err := profile.GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := profile.GetAPIKey()
		if err != nil {
			return nil, err
		}

		baseURL, err := resolveAgentStudioBaseURL(
			profile.GetAgentStudioURL(),
			agentstudio.DefaultBaseURL,
			appID,
		)
		if err != nil {
			return nil, err
		}

		userID := "cli"
		if profile.Name != "" {
			userID = "cli-" + profile.Name
		}

		return agentstudio.NewClient(agentstudio.Config{
			BaseURL:       baseURL,
			ApplicationID: appID,
			APIKey:        apiKey,
			UserID:        userID,
			UserAgent:     fmt.Sprintf("algolia-cli/%s agentstudio", appVersion),
		})
	}
}

// Timeouts for SDK-backed Agent Studio calls. The SDK defaults (25s read, 2s
// connect) are too tight for conversation exports and for slow connections.
const (
	sdkRequestTimeout = 2 * time.Minute
	sdkConnectTimeout = 10 * time.Second
)

// agentStudioAPIClient builds the official Agent Studio SDK client used for
// the standard CRUD surface. It targets the same backend as the local client:
// the base URL is resolved with resolveAgentStudioBaseURL (profile / env
// override, then the build-time default, then the cluster proxy
// https://<appID>.algolia.net/agent-studio). The SDK hard-codes the
// /agent-studio/1 path, so agentstudio.SDKRequester maps it onto the resolved
// base path. The X-Algolia-User-ID header preserves the same cli /
// cli-<profile> attribution the local client sends.
func agentStudioAPIClient(f *cmdutil.Factory, appVersion string) func() (*agentStudio.APIClient, error) {
	return func() (*agentStudio.APIClient, error) {
		profile := f.Config.Profile()
		appID, err := profile.GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := profile.GetAPIKey()
		if err != nil {
			return nil, err
		}

		baseURL, err := resolveAgentStudioBaseURL(
			profile.GetAgentStudioURL(),
			agentstudio.DefaultBaseURL,
			appID,
		)
		if err != nil {
			return nil, err
		}
		scheme, host, err := agentstudio.SplitBaseURL(baseURL)
		if err != nil {
			return nil, err
		}

		client, err := agentStudio.NewClientWithConfig(agentStudio.AgentStudioConfiguration{
			Configuration: transport.Configuration{
				AppID:         appID,
				ApiKey:        apiKey,
				DefaultHeader: make(map[string]string),
				UserAgent:     fmt.Sprintf("algolia-cli/%s agentstudio", appVersion),
				Requester:     agentstudio.NewSDKRequester(baseURL),
				Hosts: []transport.StatefulHost{
					transport.NewStatefulHost(scheme, host, call.IsReadWrite),
				},
				ReadTimeout:    sdkRequestTimeout,
				WriteTimeout:   sdkRequestTimeout,
				ConnectTimeout: sdkConnectTimeout,
			},
		})
		if err != nil {
			return nil, err
		}

		userID := "cli"
		if profile.Name != "" {
			userID = "cli-" + profile.Name
		}
		client.AddDefaultHeader("X-Algolia-User-ID", userID)

		return client, nil
	}
}

// fetchCrawlerUserID retrieves the crawler user ID from the dashboard API
// using the stored OAuth token. The new model never stores it locally, so it
// is fetched when needed. Overridable in tests.
var fetchCrawlerUserID = func() (string, error) {
	client := dashboard.NewClient(auth.OAuthClientID())
	token, err := auth.GetValidToken(client)
	if err != nil {
		return "", err
	}
	user, err := client.GetCrawlerUser(token)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func crawlerClient(f *cmdutil.Factory) func() (*crawler.Client, error) {
	return func() (*crawler.Client, error) {
		APIKey, err := f.Config.Profile().GetCrawlerAPIKey()
		if err != nil {
			return nil, err
		}

		userID, err := f.Config.Profile().GetCrawlerUserID()
		if err != nil {
			userID, err = fetchCrawlerUserID()
			if err != nil {
				return nil, fmt.Errorf(
					"the crawler user ID is not configured and could not be fetched (run `algolia auth login`): %w",
					err,
				)
			}
		}

		return crawler.NewClient(userID, APIKey), nil
	}
}

func compositionClient(f *cmdutil.Factory, appVersion string) func() (*composition.APIClient, error) {
	return func() (*composition.APIClient, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		userAgent, err := getUserAgentInfo(appID, apiKey, appVersion)
		if err != nil {
			return nil, err
		}

		clientConf := composition.CompositionConfiguration{
			Configuration: transport.Configuration{
				AppID:                           appID,
				ApiKey:                          apiKey,
				UserAgent:                       userAgent,
				ExposeIntermediateNetworkErrors: true,
			},
		}

		return composition.NewClientWithConfig(clientConf)
	}
}

// resolveAgentStudioBaseURL picks the Agent Studio base URL from, in order:
//   - profileOverride (env var ALGOLIA_AGENT_STUDIO_URL or the profile's
//     agent_studio_url field — both surfaced via Profile.GetAgentStudioURL),
//   - buildDefault (the package-level agentstudio.DefaultBaseURL set via
//     ldflags by `task build` from $ALGOLIA_AGENT_STUDIO_URL),
//   - the cluster-proxy fallback https://<appID>.algolia.net/agent-studio.
//
// Extracted from agentStudioClient so the priority chain is exercised in
// isolation by tests without needing a config mock.
func resolveAgentStudioBaseURL(profileOverride, buildDefault, appID string) (string, error) {
	override := profileOverride
	if override == "" {
		override = buildDefault
	}
	return agentstudio.ResolveHost(agentstudio.HostOptions{
		Override:      override,
		ApplicationID: appID,
	})
}

// getUserAgentInfo returns the standard user agent info plus Algolia CLI
func getUserAgentInfo(appID string, apiKey string, appVersion string) (string, error) {
	client, err := search.NewClient(appID, apiKey)
	if err != nil {
		return "", err
	}
	return client.GetConfiguration().UserAgent + fmt.Sprintf("; Algolia CLI (%s)", appVersion), nil
}

// GetStatefulHosts reads the hosts information from the profile and turns into the right structure
func GetStatefulHosts(hosts []string) []transport.StatefulHost {
	var out []transport.StatefulHost
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}

		// Bare hostnames (no scheme) need a scheme prefix for url.Parse to
		// correctly place the value in the Host field instead of Path.
		if !strings.Contains(host, "://") {
			host = "https://" + host
		}

		parsedURL, err := url.Parse(host)
		if err != nil || parsedURL.Host == "" {
			continue
		}

		statefulHost := transport.NewStatefulHost(
			parsedURL.Scheme,
			parsedURL.Host,
			call.IsReadWrite,
		)
		out = append(out, statefulHost)
	}
	return out
}
