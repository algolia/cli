package set

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runSetCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		wantOut string
	}{
		{
			name:    "without forwardToReplicas",
			cli:     "foo --advancedSyntax",
			wantOut: "✓ Set settings on foo\n",
		},
		{
			name:    "with forwardToReplicas",
			cli:     "foo --advancedSyntax --forward-to-replicas",
			wantOut: "✓ Set settings on foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("PUT", "1/indexes/foo/settings"), httpmock.JSONResponse(search.UpdateTaskRes{}))
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, "")
			cmd := NewSetCmd(f)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
