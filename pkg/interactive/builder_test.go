package interactive

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scalars struct {
	Name    string  `json:"name"`
	Count   int32   `json:"count"`
	Ratio   float64 `json:"ratio"`
	Enabled bool    `json:"enabled"`
}

func TestBuild_Scalars(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs:   map[string]string{"name": "widget", "count": "7", "ratio": "1.5"},
		Confirms: map[string]bool{"enabled": true},
	}}

	var v scalars
	require.NoError(t, b.Build(&v))

	assert.Equal(t, "widget", v.Name)
	assert.Equal(t, int32(7), v.Count)
	assert.Equal(t, 1.5, v.Ratio)
	assert.True(t, v.Enabled)
}

type optionals struct {
	Note *string `json:"note,omitempty"`
	Max  *int32  `json:"max,omitempty"`
}

func TestBuild_OptionalSkipped(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{}}

	var v optionals
	require.NoError(t, b.Build(&v))

	assert.Nil(t, v.Note)
	assert.Nil(t, v.Max)
}

func TestBuild_OptionalSet(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs: map[string]string{"note": "hi", "max": "9"},
	}}

	var v optionals
	require.NoError(t, b.Build(&v))

	require.NotNil(t, v.Note)
	assert.Equal(t, "hi", *v.Note)
	require.NotNil(t, v.Max)
	assert.Equal(t, int32(9), *v.Max)
}

type optionalScalars struct {
	Flag  *bool    `json:"flag,omitempty"`
	Score *float64 `json:"score,omitempty"`
}

func TestBuild_OptionalBoolFloatSet(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs: map[string]string{"flag": "true", "score": "2.5"},
	}}

	var v optionalScalars
	require.NoError(t, b.Build(&v))

	require.NotNil(t, v.Flag)
	assert.True(t, *v.Flag)
	require.NotNil(t, v.Score)
	assert.Equal(t, 2.5, *v.Score)
}

func TestBuild_OptionalBoolFloatSkipped(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{}}

	var v optionalScalars
	require.NoError(t, b.Build(&v))

	assert.Nil(t, v.Flag)
	assert.Nil(t, v.Score)
}

func TestBuild_RequiredEmptyRejected(t *testing.T) {
	// A required string with no scripted answer resolves to "" and is rejected
	// by the validator (on a real terminal survey would re-prompt).
	b := &Builder{Prompter: &ScriptedPrompter{}}

	var v scalars
	err := b.Build(&v)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is required")
}

func TestBuild_TypeMismatchRejected(t *testing.T) {
	// A non-numeric answer for an integer field is rejected by the validator.
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs:   map[string]string{"name": "x", "count": "abc", "ratio": "1"},
		Confirms: map[string]bool{"enabled": true},
	}}

	var v scalars
	err := b.Build(&v)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "whole number")
}

// shade mimics an SDK enum: a named string type with an IsValid() bool method.
type shade string

func (s shade) IsValid() bool { return s == "red" || s == "green" || s == "blue" }

type enumHolder struct {
	Shade shade `json:"shade"`
}

func TestBuild_Enum(t *testing.T) {
	// Enums are validated via their IsValid() method (no registry); a valid
	// free-text answer is accepted.
	b := &Builder{Prompter: &ScriptedPrompter{Inputs: map[string]string{"shade": "green"}}}

	var v enumHolder
	require.NoError(t, b.Build(&v))
	assert.Equal(t, shade("green"), v.Shade)
}

func TestBuild_EnumRejectsInvalid(t *testing.T) {
	// An out-of-set value fails IsValid; the validator surfaces the error.
	b := &Builder{Prompter: &ScriptedPrompter{Inputs: map[string]string{"shade": "purple"}}}

	var v enumHolder
	err := b.Build(&v)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid shade")
}

type address struct {
	City string `json:"city"`
}

type person struct {
	ID      string   `json:"id"`
	Address *address `json:"address,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

func TestBuild_PrePopulatedPreserved(t *testing.T) {
	// ID is pre-set; builder must not re-prompt it. Only address (confirm:no)
	// and tags (count:0) are walked.
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs: map[string]string{"tags": "0"},
	}}

	v := person{ID: "keep-me"}
	require.NoError(t, b.Build(&v))

	assert.Equal(t, "keep-me", v.ID)
	assert.Nil(t, v.Address)
	assert.Empty(t, v.Tags)
}

func TestBuild_NestedPointerStruct(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs:   map[string]string{"id": "id1", "city": "Berlin", "tags": "0"},
		Confirms: map[string]bool{"address": true},
	}}

	var v person
	require.NoError(t, b.Build(&v))

	assert.Equal(t, "id1", v.ID)
	require.NotNil(t, v.Address)
	assert.Equal(t, "Berlin", v.Address.City)
}

func TestBuild_Slice(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs: map[string]string{"id": "id1", "tags": "2", "tags[0]": "a", "tags[1]": "b"},
	}}

	var v person
	// address pointer: confirm defaults to false in the fake, so it stays nil.
	require.NoError(t, b.Build(&v))
	assert.Equal(t, []string{"a", "b"}, v.Tags)
}

type union struct {
	AsString *string `json:"-"`
	AsNumber *int32  `json:"-"`
}

func TestBuild_Union(t *testing.T) {
	b := &Builder{Prompter: &ScriptedPrompter{
		Selects: map[string]string{"variant": "string"},
		Inputs:  map[string]string{"union.string": "via-string"},
	}}

	var v union
	require.NoError(t, b.Build(&v))
	require.NotNil(t, v.AsString)
	assert.Equal(t, "via-string", *v.AsString)
	assert.Nil(t, v.AsNumber)
}

type stringMapHolder struct {
	Labels map[string]string `json:"labels,omitempty"`
}

func TestBuild_StringMap(t *testing.T) {
	// Labels is a non-pointer map reached directly (no "Set?" confirm): the
	// engine prompts a count then that many key/value pairs.
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs: map[string]string{
			"entries": "2",
			"key[0]":  "k1",
			"key[1]":  "k2",
			`["k1"]`:  "v1",
			`["k2"]`:  "v2",
		},
	}}

	var v stringMapHolder
	require.NoError(t, b.Build(&v))
	assert.Equal(t, map[string]string{"k1": "v1", "k2": "v2"}, v.Labels)
}

type paramBag struct {
	Alpha *string `json:"alpha,omitempty"`
	Beta  *string `json:"beta,omitempty"`
	Gamma *string `json:"gamma,omitempty"`
}

func TestBuild_ParamBag(t *testing.T) {
	// threshold 2 -> 3 optional fields qualifies. Select indexes 0 and 2.
	b := &Builder{
		Prompter: &ScriptedPrompter{
			MultiSelects: map[string][]int{"paramBag": {0, 2}},
			Inputs:       map[string]string{"alpha": "a-val", "gamma": "g-val"},
		},
		ParamBagThreshold: 2,
	}

	var v paramBag
	require.NoError(t, b.Build(&v))
	require.NotNil(t, v.Alpha)
	assert.Equal(t, "a-val", *v.Alpha)
	assert.Nil(t, v.Beta)
	require.NotNil(t, v.Gamma)
	assert.Equal(t, "g-val", *v.Gamma)
}

type errPrompter struct{ err error }

func (e errPrompter) Input(string, func(string) error) (string, error) { return "", e.err }
func (e errPrompter) Confirm(string) (bool, error)                     { return false, e.err }
func (e errPrompter) Select(string, []string) (string, error)          { return "", e.err }
func (e errPrompter) MultiSelect(string, []string) ([]int, error)      { return nil, e.err }

func TestBuild_PrompterErrorPropagates(t *testing.T) {
	sentinel := errors.New("prompt failed")
	b := &Builder{Prompter: errPrompter{err: sentinel}}

	var v scalars
	err := b.Build(&v)
	require.Error(t, err)
	assert.ErrorIs(t, err, sentinel)
}

type recursive struct {
	Name string     `json:"name"`
	Self *recursive `json:"self,omitempty"`
}

func TestBuild_CycleGuardFallsBackToRawJSON(t *testing.T) {
	// The Self pointer recurses into the same type. Confirming "yes" enters the
	// pointer, and at that point the type is already on the recursion stack, so
	// the cycle guard fires and degrades to a raw-JSON prompt instead of looping
	// forever. The empty raw-JSON input leaves the nested value zero. The test
	// completing at all proves the guard stopped the recursion; we additionally
	// check the recursion is bounded one level deep (Self.Self stays nil).
	b := &Builder{Prompter: &ScriptedPrompter{
		Inputs:   map[string]string{"name": "top"},
		Confirms: map[string]bool{"self": true},
	}}

	var v recursive
	require.NoError(t, b.Build(&v))
	assert.Equal(t, "top", v.Name)
	// Confirm:yes allocates Self; the cycle guard then skips its contents, so
	// Self is a non-nil empty leaf and recursion did not continue.
	require.NotNil(t, v.Self)
	assert.Equal(t, "", v.Self.Name)
	assert.Nil(t, v.Self.Self)
}

func TestBuild_MaxDepthFallsBackToRawJSON(t *testing.T) {
	// With MaxDepth 1, descending into the nested address struct (depth 2)
	// exceeds the limit and degrades to a raw-JSON prompt. Empty input leaves
	// the address contents zero. The build still completes without error.
	b := &Builder{
		Prompter: &ScriptedPrompter{
			Inputs:   map[string]string{"id": "id1"},
			Confirms: map[string]bool{"address": true},
		},
		MaxDepth: 1,
	}

	var v person
	require.NoError(t, b.Build(&v))
	assert.Equal(t, "id1", v.ID)
}
