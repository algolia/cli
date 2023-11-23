package analyze

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_computeObjectStats(t *testing.T) {
	scenarios := []struct {
		Name   string
		Input  map[string]interface{}
		Output *Stats
	}{
		{
			Name: "simple",
			Input: map[string]interface{}{
				"foo": "bar", // string
				"baz": 1,     // int
				"qux": true,  // bool
				"quux": []interface{}{ // array
					"corge",
					"grault",
					"garply",
				},
				"waldo": map[string]interface{}{ // object
					"fred": "plugh",
				},
			},
			Output: &Stats{
				Attributes: map[string]*AttributeStats{
					"foo": {
						Count: 1,
						Types: map[AttributeType]float64{
							String: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"baz": {
						Count: 1,
						Types: map[AttributeType]float64{
							Numeric: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"qux": {
						Count: 1,
						Types: map[AttributeType]float64{
							Boolean: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"quux": {
						Count: 1,
						Types: map[AttributeType]float64{
							Array: 1,
						},
						StringValues: map[string]int{
							"corge":  1,
							"grault": 1,
							"garply": 1,
						},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"waldo": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"waldo.fred": {
						Count: 1,
						Types: map[AttributeType]float64{
							String: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
				},
			},
		},
		{
			Name: "nested",
			Input: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": map[string]interface{}{
							"qux": map[string]interface{}{
								"quux": map[string]interface{}{
									"waldo": map[string]interface{}{
										"fred": "plugh",
									},
								},
							},
						},
					},
				},
			},
			Output: &Stats{
				Attributes: map[string]*AttributeStats{
					"foo": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"foo.bar": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"foo.bar.baz": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"foo.bar.baz.qux": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"foo.bar.baz.qux.quux": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"foo.bar.baz.qux.quux.waldo": {
						Count: 1,
						Types: map[AttributeType]float64{
							Object: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
					"foo.bar.baz.qux.quux.waldo.fred": {
						Count: 1,
						Types: map[AttributeType]float64{
							String: 1,
						},
						StringValues:  map[string]int{},
						NumericValues: map[float64]int{},
						BooleanValues: map[bool]int{},
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			stats := Stats{
				Attributes: make(map[string]*AttributeStats),
			}
			out := computeObjectStats(&stats, "", s.Input)
			require.NotNil(t, out)
			assert.Equal(t, s.Output, out)
		})
	}
}

func Test_getType(t *testing.T) {
	scenarios := []struct {
		Name   string
		Input  interface{}
		Output AttributeType
	}{
		{
			Name:   "string",
			Input:  "foo",
			Output: String,
		},
		{
			Name:   "int",
			Input:  1,
			Output: Numeric,
		},
		{
			Name:   "float64",
			Input:  1.0,
			Output: Numeric,
		},
		{
			Name:   "bool",
			Input:  true,
			Output: Boolean,
		},
		{
			Name:   "array",
			Input:  []interface{}{},
			Output: Array,
		},
		{
			Name:   "object",
			Input:  map[string]interface{}{},
			Output: Object,
		},
		{
			Name:   "null",
			Input:  nil,
			Output: Null,
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			assert.Equal(t, s.Output, getType(s.Input))
		})
	}
}

func Test_inSettings(t *testing.T) {
	scenarios := []struct {
		Name      string
		Attribute string
		Settings  map[string]interface{}
		Output    []string
	}{
		{
			Name:      "found",
			Attribute: "foo.bar",
			Settings: map[string]interface{}{
				"searchableAttributes": []interface{}{
					"foo.bar",
				},
			},
			Output: []string{"searchableAttributes"},
		},
		{
			Name:      "not found",
			Attribute: "foo.bar",
			Settings:  map[string]interface{}{},
			Output:    nil,
		},
		{
			Name:      "found (wrapped)",
			Attribute: "foo.bar",
			Settings: map[string]interface{}{
				"searchableAttributes": []interface{}{
					"ordered(foo.bar)",
				},
			},
			Output: []string{"searchableAttributes"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			assert.Equal(t, s.Output, inSettings(s.Settings, s.Attribute))
		})
	}
}
