package cmdutil

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
)

func TestCompositionIDs(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/compositions"),
		httpmock.StringResponse(`{"items":[{"objectID":"c1"},{"objectID":"c2"}],"nbPages":1}`),
	)

	clientF := func() (*composition.APIClient, error) {
		return composition.NewClientWithConfig(composition.CompositionConfiguration{
			Configuration: transport.Configuration{AppID: "default", ApiKey: "default", Requester: r},
		})
	}

	complete := CompositionIDs(clientF)

	ids, directive := complete(nil, nil, "")
	assert.Equal(t, []string{"c1", "c2"}, ids)
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)

	// only the first positional argument is completed.
	ids, _ = complete(nil, []string{"c1"}, "")
	require.Nil(t, ids)
}

func TestCompositionAndRuleIDs(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", "1/compositions"),
		httpmock.StringResponse(`{"items":[{"objectID":"c1"},{"objectID":"c2"}],"nbPages":1}`))
	r.Register(httpmock.REST("POST", "1/compositions/c1/rules/search"),
		httpmock.StringResponse(`{"hits":[{"objectID":"r1"},{"objectID":"r2"}],"nbHits":2,"page":0,"nbPages":1}`))

	clientF := func() (*composition.APIClient, error) {
		return composition.NewClientWithConfig(composition.CompositionConfiguration{
			Configuration: transport.Configuration{AppID: "default", ApiKey: "default", Requester: r},
		})
	}

	complete := CompositionAndRuleIDs(clientF)

	// first arg completes composition ids.
	ids, _ := complete(nil, nil, "")
	assert.Equal(t, []string{"c1", "c2"}, ids)

	// second arg completes the chosen composition's rule ids.
	ids, _ = complete(nil, []string{"c1"}, "")
	assert.Equal(t, []string{"r1", "r2"}, ids)
}
