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

// stubSurveyAskOne replaces survey with a no-op and returns a pointer to the last prompt it received.
func stubSurveyAskOne(t *testing.T) *survey.Prompt {
	t.Helper()
	orig := prompt.SurveyAskOne
	t.Cleanup(func() {
		prompt.SurveyAskOne = orig
	})
	var last survey.Prompt
	prompt.SurveyAskOne = func(sp survey.Prompt, _ interface{}, _ ...survey.AskOpt) error {
		last = sp
		return nil
	}
	return &last
}

func TestSurveyPrompter_PrintsHelpForEveryKind(t *testing.T) {
	// Setup
	io, _, out, _ := iostreams.Test()
	stubSurveyAskOne(t)
	p := NewSurveyPrompter(io)

	// Action
	_, _ = p.Input("x", nil, WithHelp("input help"))
	_, _ = p.Confirm("x", WithHelp("confirm help"))
	_, _ = p.Select("x", []string{"a"}, WithHelp("select help"))
	_, _ = p.MultiSelect("x", []string{"a"}, WithHelp("multiselect help"))

	// Assert
	assert.Contains(t, out.String(), "input help")
	assert.Contains(t, out.String(), "confirm help")
	assert.Contains(t, out.String(), "select help")
	assert.Contains(t, out.String(), "multiselect help")
}

func TestSurveyPrompter_SelectPreselectsDefault(t *testing.T) {
	// Initial setup
	io, _, _, _ := iostreams.Test()
	last := stubSurveyAskOne(t)

	// Define testcases
	tests := []struct {
		name string
		def  string
		want any
	}{
		{
			name: "default is an option",
			def:  "b",
			want: "b",
		},
		{
			name: "default is not an option",
			def:  "z",
			want: nil,
		},
	}

	// For each testcase
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Action
			_, err := NewSurveyPrompter(io).Select("x", []string{"a", "b"}, WithDefault(tc.def))

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.want, (*last).(*survey.Select).Default)
		})
	}
}
