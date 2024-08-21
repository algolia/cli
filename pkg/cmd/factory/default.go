package factory

import (
	"fmt"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/api/provisionning"
	"github.com/algolia/cli/internal/authflow"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

func New(appVersion string, cfg config.IConfig) *cmdutil.Factory {
	f := &cmdutil.Factory{
		Config:         cfg,
		ExecutableName: "gh",
	}
	f.IOStreams = ioStreams()
	f.SearchClient = searchClient(f, appVersion)
	f.CrawlerClient = crawlerClient(f)
	f.ProvisionningClient = provisionningClient(f)

	return f
}

func ioStreams() *iostreams.IOStreams {
	io := iostreams.System()
	return io
}

func searchClient(f *cmdutil.Factory, appVersion string) func() (*search.Client, error) {
	return func() (*search.Client, error) {
		appID, err := f.Config.Profile().GetApplicationID()
		if err != nil {
			return nil, err
		}
		APIKey, err := f.Config.Profile().GetAPIKey()
		if err != nil {
			return nil, err
		}

		clientCfg := search.Configuration{
			AppID:          appID,
			APIKey:         APIKey,
			ExtraUserAgent: fmt.Sprintf("Algolia CLI (%s)", appVersion),
			Hosts:          f.Config.Profile().GetSearchHosts(),
		}
		return search.NewClientWithConfig(clientCfg), nil
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

func provisionningClient(f *cmdutil.Factory) func() (*provisionning.Client, error) {
	return func() (*provisionning.Client, error) {
		token, error := f.Config.Auth().Token()
		if error != nil {
			return nil, error
		}
		client := provisionning.NewClient(token)
		// Test the client by listing the applications
		_, err := client.ListApplications()
		if err != nil {
			// If the client returns a 401, we should refresh the token
			if strings.Contains(err.Error(), "Unauthorized") {
				refreshToken, err := f.Config.Auth().RefreshToken()
				if err != nil {
					return nil, fmt.Errorf("No valid refresh token found. Please login again")
				}
				token, refreshToken, err := authflow.RefreshToken(refreshToken)
				if err != nil {
					return nil, fmt.Errorf("Failed to refresh the token: %s", err)
				}
				err = f.Config.Auth().Login(token, refreshToken)
				if err != nil {
					return nil, err
				}
				client = provisionning.NewClient(token)
			} else {
				return nil, err
			}
		}
		return client, nil
	}
}
