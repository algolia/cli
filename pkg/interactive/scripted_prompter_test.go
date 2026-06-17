package interactive

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScriptedPrompter_ValidatorError(t *testing.T) {
	// A non-nil validator that rejects the scripted answer surfaces as an error
	// rather than looping (the fake does not re-prompt).
	p := &ScriptedPrompter{Inputs: map[string]string{"n": "abc"}}
	_, err := p.Input("n", func(s string) error {
		if s == "abc" {
			return errors.New("bad")
		}
		return nil
	})
	require.Error(t, err)
}

func TestScriptedPrompter_Matched(t *testing.T) {
	p := &ScriptedPrompter{
		Inputs:       map[string]string{"name": "widget"},
		Confirms:     map[string]bool{"address": true},
		Selects:      map[string]string{"variant": "string"},
		MultiSelects: map[string][]int{"paramBag": {0, 2}},
	}

	in, err := p.Input("scalars.name", nil)
	require.NoError(t, err)
	assert.Equal(t, "widget", in)

	c, err := p.Confirm("Set person.address?")
	require.NoError(t, err)
	assert.True(t, c)

	sel, err := p.Select("union (variant)", []string{"string", "int32"})
	require.NoError(t, err)
	assert.Equal(t, "string", sel)

	ms, err := p.MultiSelect("paramBag", []string{"a", "b", "c"})
	require.NoError(t, err)
	assert.Equal(t, []int{0, 2}, ms)
}

func TestScriptedPrompter_Defaults(t *testing.T) {
	// An empty script: every unmatched prompt returns a safe default, so unknown
	// or newly added fields are simply skipped.
	p := &ScriptedPrompter{}

	in, _ := p.Input("anything", nil)
	assert.Equal(t, "", in)
	c, _ := p.Confirm("anything")
	assert.False(t, c)
	sel, _ := p.Select("anything", []string{"first", "second"})
	assert.Equal(t, "first", sel) // first option
	ms, _ := p.MultiSelect("anything", []string{"a"})
	assert.Nil(t, ms)
}

func TestScriptedPrompter_LongestKeyWins(t *testing.T) {
	// Both keys are substrings of the element label; the more specific one wins.
	p := &ScriptedPrompter{Inputs: map[string]string{"tags": "2", "tags[0]": "a"}}

	got, _ := p.Input("person.tags[0]", nil)
	assert.Equal(t, "a", got)
	got, _ = p.Input("how many person.tags (integer)", nil)
	assert.Equal(t, "2", got)
}
