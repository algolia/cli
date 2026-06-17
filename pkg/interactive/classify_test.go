package interactive

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unionSample struct {
	AsString *string
	AsInt    *int32
}

type plainSample struct {
	Name string `json:"name"`
	Age  *int32 `json:"age,omitempty"`
}

type paramBagSample struct {
	A *string `json:"a,omitempty"`
	B *string `json:"b,omitempty"`
	C *string `json:"c,omitempty"`
}

func TestIsUnionType(t *testing.T) {
	assert.True(t, isUnionType(reflect.TypeOf(unionSample{})))
	assert.False(t, isUnionType(reflect.TypeOf(plainSample{})))
}

func TestIsParamBag(t *testing.T) {
	// threshold of 2 makes the 3-field struct a param bag; the plain struct has
	// a required field so it never qualifies.
	assert.True(t, isParamBag(reflect.TypeOf(paramBagSample{}), 2))
	assert.False(t, isParamBag(reflect.TypeOf(plainSample{}), 2))
}

func TestJSONFieldName(t *testing.T) {
	assert.Equal(t, "name", jsonFieldName("name,omitempty"))
	assert.Equal(t, "age", jsonFieldName("age"))
}

func TestIsRequired(t *testing.T) {
	tp := reflect.TypeOf(plainSample{})
	nameField, _ := tp.FieldByName("Name")
	ageField, _ := tp.FieldByName("Age")
	assert.True(t, isRequired(nameField))
	assert.False(t, isRequired(ageField))
}
