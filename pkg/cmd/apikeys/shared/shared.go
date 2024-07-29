package shared

import (
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	v4 "github.com/algolia/algoliasearch-client-go/v4/algolia/search"
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

type V4Key struct {
	ACL                    []v4.Acl `json:"acl"`
	CreatedAt              int64    `json:"createdAt"`
	Description            string   `json:"description"`
	Indexes                []string `json:"indexes"`
	MaxQueriesPerIPPerHour *int32   `json:"maxQueriesPerIPPerHour"`
	MaxHitsPerQuery        *int32   `json:"maxHitsPerQuery"`
	Referers               []string `json:"referers"`
	QueryParameters        *string  `json:"queryParameters"`
	Validity               *int32   `json:"validity"`
	Value                  string   `json:"value"`
}
