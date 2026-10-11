package sortingstrategy_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmd/compositions/sortingstrategy"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/interactive"
	"github.com/algolia/cli/test"
)

func TestSortingStrategy_File(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("POST", "1/compositions/my-comp/sortingStrategy"), httpmock.StringResponse(`{"taskID":7}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/7"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := sortingstrategy.NewSortingStrategyCmd(f)
	out.InBuf.WriteString(`{"Price (asc)":"products_price_asc"}`)
	_, err := test.Execute(cmd, "my-comp --file -", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":7}`, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestSortingStrategy_Interactive(t *testing.T) {
	r := &httpmock.Registry{}
	var captured []byte
	r.Register(httpmock.REST("POST", "1/compositions/my-comp/sortingStrategy"), func(req *http.Request) (*http.Response, error) {
		captured, _ = io.ReadAll(req.Body)
		return httpmock.StringResponse(`{"taskID":8}`)(req)
	})
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/8"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(true, r, nil, "")
	// One entry: label "Price (asc)" -> index "products_price_asc". Keys match
	// the engine's map prompt labels (count/key/value), as in TestBuild_StringMap.
	f.Prompter = &interactive.ScriptedPrompter{Inputs: map[string]string{
		"entries":         "1",
		"key[0]":          "Price (asc)",
		`["Price (asc)"]`: "products_price_asc",
	}}

	cmd := sortingstrategy.NewSortingStrategyCmd(f)
	_, err := test.Execute(cmd, "my-comp --interactive", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"Price (asc)":"products_price_asc"}`, string(captured))
	r.Verify(t)
}

func TestSortingStrategy_InteractiveAndFileConflict(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(true, r, nil, "")
	cmd := sortingstrategy.NewSortingStrategyCmd(f)
	_, err := test.Execute(cmd, "my-comp --file strategy.json --interactive", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one of `--file` or `--interactive`")
}

func TestSortingStrategy_InteractiveNoTTY(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := sortingstrategy.NewSortingStrategyCmd(f)
	_, err := test.Execute(cmd, "my-comp --interactive", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a terminal")
}

func TestSortingStrategy_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := sortingstrategy.NewSortingStrategyCmd(f)
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> argument")
}
