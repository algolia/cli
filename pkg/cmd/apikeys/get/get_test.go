package get

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runGetCmd(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr string
	}{
		{
			name: "get a key (success)",
			key:  "foo",
		},
		{
			name:    "get a key (error)",
			key:     "bar",
			wantErr: "API key \"bar\" does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.key == "foo" {
				name := "test"
				r.Register(
					httpmock.REST("GET", "1/keys/foo"),
					httpmock.JSONResponse(search.GetApiKeyResponse{
						Value:       "foo",
						Description: &name,
						Acl:         []search.Acl{search.ACL_SEARCH},
					}),
				)
			} else {
				r.Register(
					httpmock.REST("GET", "1/keys/bar"),
					httpmock.ErrorResponse(),
				)
			}

			f, out := test.NewFactory(false, &r, nil, "")
			cmd := NewGetCmd(f, nil)
			_, err := test.Execute(cmd, tt.key, out)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
				return
			}
		})
	}
}
