package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func TestNewUpgradeCmd(t *testing.T) {
	f, _ := test.NewFactory(false, nil, nil, "")
	cmd := NewUpgradeCmd(f)

	assert.Equal(t, "upgrade", cmd.Name())
	assert.Equal(t, "true", cmd.Annotations["skipAuthCheck"])

	require.NotNil(t, cmd.Flags().Lookup("plan"))
	require.NotNil(t, cmd.Flags().Lookup("dry-run"))

	acceptTerms := cmd.Flags().Lookup("accept-terms")
	require.NotNil(t, acceptTerms)
	assert.Equal(t, "y", acceptTerms.Shorthand)
}
