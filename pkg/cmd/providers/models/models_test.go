package models_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/providers/models"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

const (
	providerID   = "22222222-2222-4222-8222-222222222222"
	catalogPath  = "agent-studio/1/providers/models"
	providerPath = "agent-studio/1/providers/" + providerID + "/models"
	catalogBody  = `{"openai":["gpt-4o","gpt-4o-mini"],"anthropic":["claude-3-5-sonnet"]}`
	providerBody = `["gpt-4o","gpt-4o-mini"]`
)

func TestModels_Catalog(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", catalogPath), httpmock.StringResponse(catalogBody))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := models.NewModelsCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.JSONEq(t, catalogBody, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestModels_ForProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", providerPath), httpmock.StringResponse(providerBody))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := models.NewModelsCmd(f)
	_, err := test.Execute(cmd, providerID, out)
	require.NoError(t, err)

	assert.JSONEq(t, providerBody, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestModels_TooManyArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := models.NewModelsCmd(f)
	_, err := test.Execute(cmd, providerID+" extra", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "optional <provider-id>")
	assert.Empty(t, r.Requests)
}
