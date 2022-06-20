package cmdutil

import (
	"github.com/spf13/pflag"

	"github.com/algolia/cli/pkg/utils"
)

// FlagValuesMap returns a map of flag values for the given FlagSet.
func FlagValuesMap(flags *pflag.FlagSet, only ...string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	flags.Visit(func(flag *pflag.Flag) {
		// Skip if we only want to load a subset of flags.
		if only != nil && !utils.Contains(only, flag.Name) {
			return
		}
		switch flag.Value.Type() {
		case "string":
			val, err := flags.GetString(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "int":
			val, err := flags.GetInt(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "bool":
			val, err := flags.GetBool(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "float64":
			val, err := flags.GetFloat64(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "stringSlice":
			val, err := flags.GetStringSlice(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "intSlice":
			val, err := flags.GetIntSlice(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "boolSlice":
			val, err := flags.GetBoolSlice(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "float64Slice":
			val, err := flags.GetFloat64Slice(flag.Name)
			if err == nil {
				values[flag.Name] = val
			}
		case "json":
			values[flag.Name] = flag.Value.(*JSONValue).Value
		default:
			panic("unsupported flag type")
		}
	})
	return values, nil
}
