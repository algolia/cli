package shared

import (
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// JSONKey is the same as search.Key without omitting values
type JSONKey struct {
	ACL                    []string              `json:"acl"`
	CreatedAt              time.Time             `json:"createdAt"`
	Description            string                `json:"description"`
	Indexes                []string              `json:"indexes"`
	MaxQueriesPerIPPerHour int                   `json:"maxQueriesPerIPPerHour"`
	MaxHitsPerQuery        int                   `json:"maxHitsPerQuery"`
	Referers               []string              `json:"referers"`
	QueryParameters        search.KeyQueryParams `json:"queryParameters"`
	Validity               time.Duration         `json:"validity"`
	Value                  string                `json:"value"`
}
