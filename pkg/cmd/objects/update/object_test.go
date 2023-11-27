package update

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateOperation(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "invalid operation",
			operation:  "invalid",
			wantErr:    true,
			wantErrMsg: "invalid operation \"invalid\" (valid operations are Increment, Decrement, Add, AddUnique, IncrementSet and IncrementFrom)",
		},
	}

	for _, ops := range []string{"Increment", "Decrement", "Add", "AddUnique", "IncrementSet", "IncrementFrom"} {
		tests = append(tests, struct {
			name       string
			operation  string
			wantErr    bool
			wantErrMsg string
		}{
			name:      ops,
			operation: ops,
			wantErr:   false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperation(search.PartialUpdateOperation{Operation: tt.operation})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_Object_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		wantObj    Object
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "empty object",
			data:       []byte(`{}`),
			wantErr:    true,
			wantErrMsg: "objectID is required",
		},
		{
			name:       "missing objectID",
			data:       []byte(`{"foo": "bar"}`),
			wantErr:    true,
			wantErrMsg: "objectID is required",
		},
		{
			name:    "valid object",
			data:    []byte(`{"objectID": "foo"}`),
			wantErr: false,
			wantObj: Object{"objectID": "foo"},
		},
		{
			name: "nested object (not an operation)",
			data: []byte(`{
				"objectID": "foo",
				"bar": {
					"foo": "bar"
				}
			}`),
			wantErr: false,
			wantObj: Object{"objectID": "foo", "bar": map[string]interface{}{"foo": "bar"}},
		},
		{
			name: "invalid operation type",
			data: []byte(`{
				"objectID": "foo",
				"bar": {
					"operation": "invalid"
				}
			}`),
			wantErr:    true,
			wantErrMsg: "invalid operation \"invalid\" (valid operations are Increment, Decrement, Add, AddUnique, IncrementSet and IncrementFrom)",
		},
		{
			name: "valid operation",
			data: []byte(`{
				"objectID": "foo",
				"bar": {
					"operation": "Increment"
				}
			}`),
			wantErr: false,
			wantObj: Object{"objectID": "foo", "bar": search.PartialUpdateOperation{Operation: "Increment"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Object
			err := o.UnmarshalJSON(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantObj, o)
		})
	}
}
