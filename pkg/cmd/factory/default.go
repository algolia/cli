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
	f.AgentStudioClient = agentStudioClient(f)

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

func agentStudioClient(f *cmdutil.Factory) func() (*agentstudio.Client, error) {
	return func() (*agentstudio.Client, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		apiKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		return agentstudio.NewClient(appID, apiKey), nil
	}
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
