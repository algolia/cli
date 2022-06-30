package cmdutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JSONValue(t *testing.T) {
	tests := []struct {
		in    string
		types []string
		out   interface{}
		err   bool
	}{
		{
			in: "string",
			types: []string{
				"string",
			},
			out: "string",
			err: false,
		},
		{
			in:    "string",
			types: []string{},
			out:   nil,
			err:   true,
		},
		{
			in:    `["array"]`,
			types: []string{},
			out:   []interface{}{"array"},
			err:   false,
		},
		{
			in:    "true",
			types: []string{},
			out:   true,
			err:   false,
		},
	}

	for _, test := range tests {
		j := NewJSONVar(test.types...)
		err := j.Set(test.in)

		if err != nil && !test.err {
			t.Errorf("unexpected error: %v", err)
		}
		if err == nil && test.err {
			t.Errorf("expected error, got none")
		}

		assert.Equal(t, test.out, j.Value)
	}
}
