package models_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agentstudio/providers/models"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestModels_All(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/providers/models"),
		httpmock.StringResponse(`{"openai":["gpt-4o"],"anthropic":["claude-sonnet-5"]}`),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := models.NewModelsCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"gpt-4o"`)
	r.Verify(t)
}

func TestModels_ForProvider(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "agent-studio/1/providers/my-provider/models"),
		httpmock.StringResponse(`["gpt-4o","gpt-4o-mini"]`),
	)

	f, out := test.NewFactory(false, r, nil, "")
	cmd := models.NewModelsCmd(f)
	_, err := test.Execute(cmd, "--provider-id my-provider", out)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(out.String()), `"gpt-4o-mini"`)
	r.Verify(t)
}
