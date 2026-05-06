package factory

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/call"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

func New(appVersion string, cfg config.IConfig) *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config:         cfg,
		ExecutableName: "gh",
	}
	f.IOStreams = ioStreams(f)
	f.SearchClient = searchClient(f, appVersion)
	f.CrawlerClient = crawlerClient(f)
	f.AgentStudioClient = agentStudioClient(f, appVersion)

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

func crawlerClient(f *cmdutil.Factory) func() (*crawler.Client, error) {
	return func() (*crawler.Client, error) {
		userID, err := f.Config.Profile().GetCrawlerUserID()
		if err != nil {
			return nil, err
		}
		APIKey, err := f.Config.Profile().GetCrawlerAPIKey()
		if err != nil {
			return nil, err
		}

		return crawler.NewClient(userID, APIKey), nil
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
