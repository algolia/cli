package cmdutil

import (
	"encoding/json"

	"github.com/algolia/cli/pkg/utils"
)

// JSONValue is a flag.Value that marshals a JSON object into a string.
type JSONValue struct {
	Value interface{}
	types []string
}

// NewJSONValue creates a new JSONVar.
func NewJSONVar(types ...string) *JSONValue {
	return &JSONValue{
		types: types,
	}
}

// String returns the string representation of the JSON object.
func (j *JSONValue) String() string {
	b, err := json.Marshal(j.Value)
	if err != nil {
		return "failed to marshal object"
	}
	return string(b)
}

// Set parses the JSON string into the value.
func (j *JSONValue) Set(s string) error {
	if err := json.Unmarshal([]byte(s), &j.Value); err != nil {
		if utils.Contains(j.types, "string") {
			j.Value = s
			return nil
		}
		return err
	}
	return nil
}

// Type returns the type of the value.
func (j *JSONValue) Type() string {
	return "json"
}
