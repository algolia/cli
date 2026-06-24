package list_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/compositions/rules/list"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

// ruleJSON is a valid composition rule with consequence.behavior.injection populated
// (v4.38.0+ requires a valid oneOf value in CompositionBehavior).
const ruleJSON = `{"objectID":"rule-1","conditions":[{"anchoring":"is","pattern":"shirt"}],"consequence":{"behavior":{"injection":{"main":{}}}}}`

func TestListRules(t *testing.T) {
	tests := []struct {
		name    string
		mock    string
		wantOut string
	}{
		{
			name:    "single rule",
			mock:    `{"hits":[` + ruleJSON + `],"nbHits":1,"page":0,"nbPages":1,"hitsPerPage":20}`,
			wantOut: `{"hits":[{"conditions":[{"anchoring":"is","pattern":"shirt"}],"consequence":{"behavior":{"injection":{"main":{}}}},"objectID":"rule-1"}],"nbHits":1,"nbPages":1,"page":0}`,
		},
		{
			name:    "empty list",
			mock:    `{"hits":[],"nbHits":0,"page":0,"nbPages":0,"hitsPerPage":20}`,
			wantOut: `{"hits":[],"nbHits":0,"nbPages":0,"page":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpmock.Registry{}
			r.Register(
				httpmock.REST("POST", "1/compositions/my-comp/rules/search"),
				httpmock.StringResponse(tt.mock),
			)

			f, out := test.NewFactory(false, r, nil, "")
			cmd := list.NewListCmd(f)
			_, err := test.Execute(cmd, "my-comp", out)
			require.NoError(t, err)

			assert.JSONEq(t, tt.wantOut, strings.TrimSpace(out.String()))
			r.Verify(t)
		})
	}
}

func TestListRules_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := list.NewListCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> argument")
}
