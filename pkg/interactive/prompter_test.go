package interactive

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

func TestSurveyPrompter_UsesPromptVars(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	p := NewSurveyPrompter(io)

	origAskOne := prompt.SurveyAskOne
	t.Cleanup(func() {
		prompt.SurveyAskOne = origAskOne
	})

	prompt.SurveyAskOne = func(sp survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
		switch sp.(type) {
		case *survey.Input:
			*(response.(*string)) = "typed"
		case *survey.Confirm:
			*(response.(*bool)) = true
		case *survey.Select:
			*(response.(*string)) = "picked"
		case *survey.MultiSelect:
			*(response.(*[]int)) = []int{1}
		}
		return nil
	}

	in, err := p.Input("x", nil)
	require.NoError(t, err)
	assert.Equal(t, "typed", in)

	c, err := p.Confirm("x")
	require.NoError(t, err)
	assert.True(t, c)

	sel, err := p.Select("x", []string{"a", "b"})
	require.NoError(t, err)
	assert.Equal(t, "picked", sel)

	multi, err := p.MultiSelect("x", []string{"a", "b"})
	require.NoError(t, err)
	assert.Equal(t, []int{1}, multi)
}
