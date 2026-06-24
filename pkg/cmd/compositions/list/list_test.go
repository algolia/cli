package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/compositions/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestListCompositions(t *testing.T) {
	tests := []struct {
		name    string
		mock    string
		wantOut string
	}{
		{
			name:    "single composition",
			mock:    `{"items":[{"objectID":"my-comp","name":"My Comp","behavior":{"injection":{"main":{}}}}],"nbPages":1,"page":0,"hitsPerPage":20,"nbHits":1}`,
			wantOut: `{"hitsPerPage":20,"items":[{"behavior":{"injection":{"main":{}}},"name":"My Comp","objectID":"my-comp"}],"nbHits":1,"nbPages":1,"page":0}`,
		},
		{
			name:    "empty list",
			mock:    `{"items":[],"nbPages":0,"page":0,"hitsPerPage":20,"nbHits":0}`,
			wantOut: `{"hitsPerPage":20,"items":[],"nbHits":0,"nbPages":0,"page":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/compositions"),
				httpmock.StringResponse(tt.mock),
			)

			f, out := test.NewFactory(false, r, nil, "")
			cmd := list.NewListCmd(f)
			_, err := test.Execute(cmd, "", out)
			require.NoError(t, err)

			assert.JSONEq(t, tt.wantOut, strings.TrimSpace(out.String()))
			r.Verify(t)
		})
	}
}
