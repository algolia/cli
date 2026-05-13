package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/compositions/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestGetComposition(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		mock    string
		wantOut string
	}{
		{
			name:    "composition with injection behavior",
			cli:     "my-comp",
			mock:    `{"objectID":"my-comp","name":"My Comp","behavior":{"injection":{"main":{}}}}`,
			wantOut: `{"behavior":{"injection":{"main":{}}},"name":"My Comp","objectID":"my-comp"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/compositions/my-comp"),
				httpmock.StringResponse(tt.mock),
			)

			f, out := test.NewFactory(false, r, nil, "")
			cmd := get.NewGetCmd(f)
			_, err := test.Execute(cmd, tt.cli, out)
			require.NoError(t, err)

			assert.JSONEq(t, tt.wantOut, strings.TrimSpace(out.String()))
			r.Verify(t)
		})
	}
}

func TestGetComposition_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> argument")
}
