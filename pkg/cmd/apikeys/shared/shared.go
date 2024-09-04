package shared

import (
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

// JSONKey is the same as search.Key without omitting values
type JSONKey struct {
	ACL                    []search.Acl `json:"acl"`
	CreatedAt              int64        `json:"createdAt"`
	Description            string       `json:"description"`
	Indexes                []string     `json:"indexes"`
	MaxQueriesPerIPPerHour *int32       `json:"maxQueriesPerIPPerHour"`
	MaxHitsPerQuery        *int32       `json:"maxHitsPerQuery"`
	Referers               []string     `json:"referers"`
	QueryParameters        *string      `json:"queryParameters"`
	Validity               *int32       `json:"validity"`
	Value                  string       `json:"value"`
}
