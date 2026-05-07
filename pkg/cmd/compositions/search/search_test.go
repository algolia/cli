package search_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compsearch "github.com/algolia/cli/pkg/cmd/compositions/search"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestSearchComposition(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		mock    string
		wantOut string
	}{
		{
			name:    "basic query returns hits",
			cli:     "my-comp shirt",
			mock:    `{"hits":[{"objectID":"obj1","name":"shirt"}],"nbHits":1,"page":0,"hitsPerPage":20,"processingTimeMS":5}`,
			wantOut: `{"hits":[{"name":"shirt","objectID":"obj1"}],"hitsPerPage":20,"nbHits":1,"page":0,"processingTimeMS":5,"results":null}`,
		},
		{
			name:    "no results",
			cli:     "my-comp nomatch",
			mock:    `{"hits":[],"nbHits":0}`,
			wantOut: `{"hits":[],"nbHits":0,"results":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/compositions/my-comp/run"),
				httpmock.StringResponse(tt.mock),
			)

			f, out := test.NewFactory(false, r, nil, "")
			cmd := compsearch.NewSearchCmd(f)
			_, err := test.Execute(cmd, tt.cli, out)
			require.NoError(t, err)

			assert.JSONEq(t, tt.wantOut, strings.TrimSpace(out.String()))
			r.Verify(t)
		})
	}
}

func TestSearchComposition_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := compsearch.NewSearchCmd(f)
	_, err := test.Execute(cmd, "my-comp", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> and a <query> argument")
}
