package upsert_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmd/compositions/rules/upsert"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

// validRuleBody is a well-formed rule body.
// consequence.behavior must include injection or multifeed - SDK MarshalJSON panics on empty CompositionBehavior.
const validRuleBody = `{"objectID":"rule-1","conditions":[{"anchoring":"is","pattern":"shirt"}],"consequence":{"behavior":{"injection":{"promote":[{"objectIDs":["obj1"],"position":0}]}}}}`

func TestUpsertRule(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("PUT", "1/compositions/my-comp/rules/rule-1"), httpmock.StringResponse(`{"taskID":55}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/55"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	out.InBuf.WriteString(validRuleBody)
	_, err := test.Execute(cmd, "my-comp rule-1 --file -", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":55}`, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestUpsertRule_WaitsForPublished(t *testing.T) {
	// Verifies polling continues through successive notPublished states.
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("PUT", "1/compositions/my-comp/rules/rule-1"), httpmock.StringResponse(`{"taskID":66}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/66"), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/66"), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/66"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	out.InBuf.WriteString(validRuleBody)
	_, err := test.Execute(cmd, "my-comp rule-1 --file -", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":66}`, strings.TrimSpace(out.String()))
	r.Verify(t)

	taskPolls := 0
	for _, req := range r.Requests {
		if strings.Contains(req.URL.Path, "/task/") {
			taskPolls++
		}
	}
	assert.Equal(t, 3, taskPolls, "expected 3 task status polls (2x notPublished + 1x published)")
}

func TestUpsertRule_MissingFile(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	_, err := test.Execute(cmd, "my-comp rule-1", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file")
}

func TestUpsertRule_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	out.InBuf.WriteString(`not-json`)
	_, err := test.Execute(cmd, "my-comp rule-1 --file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing rule JSON")
}

func TestUpsertRule_MissingArgs(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> and a <rule-id> argument")
}
