package authflow

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/oauth/webapp"
	"github.com/cli/browser"
)

var (
	// OAuthHost is the host of the OAuth server
	OAuthHost = "https://www.algolia.com/oauth"
	// The "GitHub CLI" OAuth app
	oauthClientID = os.Getenv("ALGOLIA_OAUTH_CLIENT_ID")
	// This value is safe to be embedded in version control
	oauthClientSecret = os.Getenv("ALGOLIA_OAUTH_CLIENT_SECRET")
)

var scopes = []string{
	"public",
	"applications:manage",
	"teams:manage",
	"keys:manage",
}

func AuthFlow(IO *iostreams.IOStreams, notice string) (string, string, error) {
	httpClient := &http.Client{}

	callbackURI := os.Getenv("ALGOLIA_OAUTH_CALLBACK_URI")

	flow, err := webapp.InitFlow()
	if err != nil {
		panic(err)
	}

	params := webapp.BrowserParams{
		ClientID:    oauthClientID,
		RedirectURI: callbackURI,
		Scopes:      scopes,
	}
	browserURL, err := flow.BrowserURL(fmt.Sprintf("%s/authorize", OAuthHost), params)
	if err != nil {
		panic(err)
	}

	go func() {
		_ = flow.StartServer(nil)
	}()

	// Note: the user's web browser must run on the same device as the running app.
	err = browser.OpenURL(browserURL)
	if err != nil {
		panic(err)
	}

	accessToken, err := flow.Wait(context.TODO(), httpClient, fmt.Sprintf("%s/token", OAuthHost), webapp.WaitOptions{
		ClientSecret: oauthClientSecret,
	})
	if err != nil {
		panic(err)
	}

	return accessToken.Token, accessToken.RefreshToken, nil
}

// RefreshToken refreshes the access token using the refresh token.
func RefreshToken(refreshToken string) (string, string, error) {
	httpClient := &http.Client{}

	accessToken, err := webapp.RefreshAccessToken(httpClient, fmt.Sprintf("%s/token", OAuthHost), webapp.RefreshOptions{
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,
		RefreshToken: refreshToken,
	})
	if err != nil {
		panic(err)
	}

	return accessToken.Token, accessToken.RefreshToken, nil
}
