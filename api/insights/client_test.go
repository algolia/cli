package insights

import (
	"reflect"
	"testing"

	algoliaInsights "github.com/algolia/algoliasearch-client-go/v4/algolia/insights"
	"github.com/stretchr/testify/require"
)

func TestNewClientSetsRequestedRegion(t *testing.T) {
	client, err := NewClient("test-app-id", "test-api-key", algoliaInsights.DE)
	require.NoError(t, err)

	cfg := client.GetConfiguration()
	require.Equal(t, algoliaInsights.DE, cfg.Region)
	require.NotEmpty(t, cfg.Hosts)

	host := reflect.ValueOf(cfg.Hosts[0]).FieldByName("host").String()
	require.Equal(t, "insights.de.algolia.io", host)
}
