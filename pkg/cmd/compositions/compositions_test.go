package compositions_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/compositions"
	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/test"
)

// wantBody is the exact composition the interactive build produces from the
// scripted answers below. It is a valid composition: objectID is pre-populated
// from the positional arg, name comes from the keyed Input, and behavior selects
// the injection variant whose main source is a search source pointing at a real
// index. description, sortingStrategy, and all the optional query parameters are
// unanswered and therefore omitted.
const wantBody = `{
	"objectID": "my-comp",
	"name": "My Composition",
	"behavior": {
		"injection": {
			"main": {
				"source": {
					"search": {"index": "my-index"}
				}
			}
		}
	}
}`

// Drives `compositions upsert --interactive` through the real command tree with
// a label-keyed ScriptedPrompter on the Factory. Answers are keyed by a unique
// substring of the prompt label, so adding or reordering SDK fields does not
// break the INPUT side: unmatched prompts fall back to safe defaults (skip).
func TestCompositions_UpsertInteractive(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(httpmock.REST("PUT", "1/compositions/my-comp"), func(req *http.Request) (*http.Response, error) {
		captured, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(`{"taskID":42}`)(req)
	})
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/42"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(true, r, nil, "")
	f.Prompter = &interactive.ScriptedPrompter{
		Inputs: map[string]string{
			"name":  "My Composition",
			"index": "my-index", // behavior.injection.main.source.search.index (required)
		},
		Confirms: map[string]bool{
			// Trailing "?" pins this to the source pointer confirm
			// ("...main.source?") so it does not also match the deeper
			// "...search.params?" confirm, whose path contains ".main.source.".
			"main.source?": true,
		},
		// Both unions are keyed by their leaf "(variant)" label so the deep
		// source select does not collide with the top-level behavior select.
		Selects: map[string]string{
			"behavior (variant)": "CompositionInjectionBehavior",
			"source (variant)":   "InjectionMainSearchSource",
		},
	}

	cmd := compositions.NewCompositionsCmd(f)
	_, err := test.Execute(cmd, "upsert my-comp --interactive", out)
	require.NoError(t, err)

	assert.JSONEq(t, wantBody, string(captured))
	r.Verify(t)
}
