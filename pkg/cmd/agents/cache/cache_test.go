package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/cache"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestInvalidateCache(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", "agent-studio/1/agents/my-agent/cache"), httpmock.StringResponse(``))

	f, out := test.NewFactory(false, r, nil, "")
	cmd := cache.NewInvalidateCacheCmd(f)
	_, err := test.Execute(cmd, "my-agent", out)
	require.NoError(t, err)
	r.Verify(t)
}

func TestInvalidateCache_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := cache.NewInvalidateCacheCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires an <agent-id> argument")
}
