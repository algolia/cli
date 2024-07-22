package authflow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/oauth/webapp"
	"github.com/cli/browser"
)

var (
	// The "GitHub CLI" OAuth app
	oauthClientID = "Q4d7MJnfy4-QGEvf5gZ9IIhwSifR_9N2DviotpiA58s"
	// This value is safe to be embedded in version control
	oauthClientSecret = "xwfSbqQWJYTPUauQmQ72dOEULPY0Ia7L6c0vOrsn_7I"
)

func AuthFlow(oauthHost string, IO *iostreams.IOStreams, notice string, additionalScopes []string) (string, string, error) {
	httpClient := &http.Client{}

	minimumScopes := []string{"public"}
	scopes := append(minimumScopes, additionalScopes...)

	callbackURI := "http://localhost:3456/callback"

	flow, err := webapp.InitFlow()
	if err != nil {
		panic(err)
	}

	params := webapp.BrowserParams{
		ClientID:    oauthClientID,
		RedirectURI: callbackURI,
		Scopes:      scopes,
	}
	browserURL, err := flow.BrowserURL("https://www.algolia.com/oauth/authorize", params)
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

	accessToken, err := flow.Wait(context.TODO(), httpClient, "https://github.com/login/oauth/access_token", webapp.WaitOptions{
		ClientSecret: oauthClientSecret,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access token: %s\n", accessToken.Token)

	return accessToken.Token, accessToken.RefreshToken, nil
}
