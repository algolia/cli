package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/version"
	"github.com/algolia/cli/test"
)

func TestBetaAgentsPreRunE_skipWhenReleaseBuild(t *testing.T) {
	prev := version.Distribution
	t.Cleanup(func() { version.Distribution = prev })
	version.Distribution = ""

	f, bio := test.NewFactory(false, nil, nil, "")
	require.NoError(t, betaAgentsPreRunE(f)(nil, nil))
	assert.Empty(t, bio.ErrBuf.String())
}

func TestBetaAgentsPreRunE_warnWhenDistributionSet(t *testing.T) {
	prev := version.Distribution
	t.Cleanup(func() { version.Distribution = prev })
	version.Distribution = "beta"

	f, bio := test.NewFactory(false, nil, nil, "")
	require.NoError(t, betaAgentsPreRunE(f)(nil, nil))
	got := bio.ErrBuf.String()
	assert.Contains(t, got, "beta CLI")
	assert.Contains(t, got, "agents")
}
