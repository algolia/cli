package get_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/compositions/rules/get"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestGetRule(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		mock    string
		wantOut string
	}{
		{
			name:    "rule with injection consequence",
			cli:     "my-comp rule-1",
			mock:    `{"objectID":"rule-1","conditions":[{"anchoring":"is","pattern":"shirt"}],"consequence":{"behavior":{"injection":{"main":{}}}}}`,
			wantOut: `{"conditions":[{"anchoring":"is","pattern":"shirt"}],"consequence":{"behavior":{"injection":{"main":{}}}},"objectID":"rule-1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("GET", "1/compositions/my-comp/rules/rule-1"),
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

func TestGetRule_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := get.NewGetCmd(f)
	_, err := test.Execute(cmd, "my-comp", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> and a <rule-id> argument")
}
