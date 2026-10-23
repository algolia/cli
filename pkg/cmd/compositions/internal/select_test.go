package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/interactive"
)

func TestResolveCompositionID_FromArg(t *testing.T) {
	// an explicit arg is returned as-is, no client or prompt needed.
	id, err := compinternal.ResolveCompositionID([]string{"c1"}, true, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, "c1", id)
}

func TestResolveCompositionID_NoArgNoTTY(t *testing.T) {
	_, err := compinternal.ResolveCompositionID(nil, false, nil, nil)
	assert.Error(t, err)
}

func TestResolveCompositionID_Selector(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("GET", "1/compositions"),
		httpmock.StringResponse(`{"items":[{"objectID":"alpha"},{"objectID":"beta"}],"nbPages":1}`))

	prompter := &interactive.ScriptedPrompter{Selects: map[string]string{"Composition": "beta"}}

	id, err := compinternal.ResolveCompositionID(nil, true, prompter, compositionClientFor(r))
	require.NoError(t, err)
	assert.Equal(t, "beta", id)
}
