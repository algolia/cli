package interactive

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// promptCall is one prompt the builder issued.
type promptCall struct {
	kind  string
	label string
	cfg   PromptConfig
}

// recordingPrompter records every prompt and answers with the first option or an empty value.
type recordingPrompter struct {
	calls      []promptCall
	confirmYes bool
}

func (r *recordingPrompter) record(kind, label string, opts []PromptOption) {
	r.calls = append(r.calls, promptCall{kind, label, ApplyPromptOptions(opts)})
}

func (r *recordingPrompter) Input(label string, _ func(string) error, opts ...PromptOption) (string, error) {
	r.record("input", label, opts)
	return "", nil
}

func (r *recordingPrompter) Confirm(label string, opts ...PromptOption) (bool, error) {
	r.record("confirm", label, opts)
	return r.confirmYes, nil
}

func (r *recordingPrompter) Select(label string, options []string, opts ...PromptOption) (string, error) {
	r.record("select", label, opts)
	return options[0], nil
}

func (r *recordingPrompter) MultiSelect(label string, _ []string, opts ...PromptOption) ([]int, error) {
	r.record("multiselect", label, opts)
	return nil, nil
}

func (r *recordingPrompter) findCall(kind, label string) (promptCall, bool) {
	for _, c := range r.calls {
		if c.kind == kind && c.label == label {
			return c, true
		}
	}
	return promptCall{}, false
}

// wrongValuesEnum has a Values() list of the wrong element type.
type wrongValuesEnum string

func (wrongValuesEnum) IsValid() bool    { return true }
func (wrongValuesEnum) Values() []string { return []string{"a"} }

// wrongAccessor has a FieldDescription whose parameter is a named string type.
type wrongAccessor struct {
	Name string `json:"name"`
}

func (wrongAccessor) FieldDescription(_ wrongValuesEnum) string { return "never used" }

// describedUnion mimics a generated oneOf wrapper that exposes per-variant descriptions.
type describedUnion struct {
	AsString *string `json:"-"`
	AsNumber *int32  `json:"-"`
}

func (describedUnion) VariantDescription(variant string) string {
	if variant == "string" {
		return "free text"
	}
	return ""
}

// metaHolder mimics a generated SDK model with one field per prompt path that carries help.
type metaHolder struct {
	Tone   tone            `json:"tone"`
	Alive  bool            `json:"alive"`
	Choice describedUnion  `json:"choice"`
	Pick   *describedUnion `json:"pick,omitempty"`
	Bag    paramBag        `json:"bag"`
	Any    any             `json:"any"`
	Bad    any             `json:"bad"`
	Keyed  map[int]string  `json:"keyed"`
	Odd    wrongValuesEnum `json:"odd"`
	Inner  wrongAccessor   `json:"inner"`
}

var metaExamples = map[string]string{
	"tone":  "cool",
	"any":   `"x"`,
	"bad":   "not json",
	"keyed": `"x"`,
}

func (metaHolder) FieldDescription(field string) string { return "describes " + field }

func (metaHolder) FieldExample(field string) string { return metaExamples[field] }

func TestBuild_FieldMetaReachesPrompts(t *testing.T) {
	// Initial setup
	rec := &recordingPrompter{confirmYes: true}
	b := &Builder{Prompter: rec, ParamBagThreshold: 2}
	require.NoError(t, b.Build(&metaHolder{}))

	// Define testcases
	tests := []struct {
		name  string
		kind  string
		label string
		help  string
		def   string
	}{
		{
			name:  "enum select preselects the example",
			kind:  "select",
			label: "metaHolder.tone",
			help:  "describes tone",
			def:   "cool",
		},
		{
			name:  "bool confirm",
			kind:  "confirm",
			label: "metaHolder.alive",
			help:  "describes alive",
		},
		{
			name:  "union variant select",
			kind:  "select",
			label: "metaHolder.choice (variant)",
			help:  "describes choice",
		},
		{
			name:  "param bag multiselect",
			kind:  "multiselect",
			label: "metaHolder.bag",
			help:  "describes bag",
		},
		{
			name:  "pointer field shows help on the confirm",
			kind:  "confirm",
			label: "Set metaHolder.pick?",
			help:  "describes pick",
		},
		{
			name:  "pointer field does not repeat help on the inner prompt",
			kind:  "select",
			label: "metaHolder.pick (variant)",
		},
		{
			name:  "raw JSON keeps a valid JSON example",
			kind:  "input",
			label: "metaHolder.any (raw JSON, empty to skip)",
			help:  "describes any",
			def:   `"x"`,
		},
		{
			name:  "raw JSON drops an example that is not JSON",
			kind:  "input",
			label: "metaHolder.bad (raw JSON, empty to skip)",
			help:  "describes bad",
		},
		{
			name:  "non-string-keyed map falls back to raw JSON",
			kind:  "input",
			label: "metaHolder.keyed (raw JSON, empty to skip)",
			help:  "describes keyed",
			def:   `"x"`,
		},
		{
			name:  "Values with the wrong element type falls back to a text input",
			kind:  "input",
			label: "metaHolder.odd",
			help:  "describes odd",
		},
		{
			name:  "FieldDescription with a named string parameter yields no help",
			kind:  "input",
			label: "metaHolder.inner.name",
		},
	}

	// For each testcase
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Action
			call, found := rec.findCall(tc.kind, tc.label)

			// Assert
			require.True(t, found, "no %s prompt labelled %q", tc.kind, tc.label)
			assert.Equal(t, tc.help, call.cfg.Help)
			assert.Equal(t, tc.def, call.cfg.Default)
		})
	}
}

func TestBuild_UnionPickMapsDescribedLabelBackToVariant(t *testing.T) {
	// Initial setup
	str := "hello"
	num := int32(7)

	// Define testcases
	tests := []struct {
		name string
		pick string
		in   string
		want describedUnion
	}{
		{
			name: "variant with a description",
			pick: "string: free text",
			in:   "hello",
			want: describedUnion{AsString: &str},
		},
		{
			name: "variant without a description",
			pick: "int32",
			in:   "7",
			want: describedUnion{AsNumber: &num},
		},
	}

	// For each testcase
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			b := &Builder{Prompter: &ScriptedPrompter{
				Selects: map[string]string{"variant": tc.pick},
				Inputs:  map[string]string{"describedUnion.": tc.in},
			}}
			var v describedUnion

			// Action
			require.NoError(t, b.Build(&v))

			// Assert
			assert.Equal(t, tc.want, v)
		})
	}
}

func TestBuild_SDKMetadataAndSuggestReachPrompt(t *testing.T) {
	// Setup
	rec := &recordingPrompter{}
	b := &Builder{
		Prompter: rec,
		Suggest: func(_, _ string, _ map[string]string) []string {
			return []string{"idx_alpha"}
		},
	}

	// Action
	require.NoError(t, b.Build(&composition.Composition{}))

	// Assert
	call, found := rec.findCall("input", "Composition.name")
	require.True(t, found, "no input prompt labelled Composition.name")
	assert.NotEmpty(t, call.cfg.Help)
	assert.NotEmpty(t, call.cfg.Default)
	require.NotNil(t, call.cfg.Suggest)
	assert.Equal(t, []string{"idx_alpha", call.cfg.Default}, call.cfg.Suggest(""))

	// Action (no Suggest)
	bare := &recordingPrompter{}
	require.NoError(t, (&Builder{Prompter: bare}).Build(&composition.Composition{}))

	// Assert (example is the only candidate)
	bareCall, found := bare.findCall("input", "Composition.name")
	require.True(t, found)
	require.NotNil(t, bareCall.cfg.Suggest)
	assert.Equal(t, []string{bareCall.cfg.Default}, bareCall.cfg.Suggest(""))
}
