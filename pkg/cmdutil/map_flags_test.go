package cmdutil

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func Test_FlagValuesMap(t *testing.T) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flagSet.String("string", "", "")
	flagSet.Int("int", 0, "")
	flagSet.Bool("bool", false, "")
	flagSet.Float64("float64", 0.0, "")
	flagSet.StringSlice("stringSlice", []string{}, "")
	flagSet.IntSlice("intSlice", []int{}, "")
	flagSet.BoolSlice("boolSlice", []bool{}, "")
	flagSet.Float64Slice("float64Slice", []float64{}, "")
	flagSet.Var(&JSONValue{}, "json", "")

	_ = flagSet.Set("string", "string")
	_ = flagSet.Set("int", "1")
	_ = flagSet.Set("bool", "true")
	_ = flagSet.Set("float64", "1.0")
	_ = flagSet.Set("stringSlice", "string")
	_ = flagSet.Set("intSlice", "1")
	_ = flagSet.Set("boolSlice", "true")
	_ = flagSet.Set("float64Slice", "1.0")
	_ = flagSet.Set("json", `["json"]`)

	flagValuesMap, err := FlagValuesMap(flagSet)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, "string", flagValuesMap["string"])
	assert.Equal(t, 1, flagValuesMap["int"])
	assert.Equal(t, true, flagValuesMap["bool"])
	assert.Equal(t, 1.0, flagValuesMap["float64"])
	assert.Equal(t, []string{"string"}, flagValuesMap["stringSlice"])
	assert.Equal(t, []int{1}, flagValuesMap["intSlice"])
	assert.Equal(t, []bool{true}, flagValuesMap["boolSlice"])
	assert.Equal(t, []float64{1.0}, flagValuesMap["float64Slice"])
	assert.Equal(t, []interface{}{"json"}, flagValuesMap["json"])
}
