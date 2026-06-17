package interactive

import (
	"reflect"
	"strings"
)

// isUnionType reports whether t looks like an OpenAPI oneOf wrapper: a struct
// whose every exported field is a pointer with no JSON name (either untagged,
// as the Algolia SDK emits, or tagged json:"-").
func isUnionType(t reflect.Type) bool {
	if t.Kind() != reflect.Struct || t.NumField() == 0 {
		return false
	}
	exported := 0
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		exported++
		jsonTag := f.Tag.Get("json")
		if f.Type.Kind() != reflect.Pointer || (jsonTag != "" && jsonTag != "-") {
			return false
		}
	}
	return exported > 0
}

// isParamBag reports whether t is a large optional-only parameter object: more
// than threshold optional exported fields and zero required fields.
func isParamBag(t reflect.Type, threshold int) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	exported := 0
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		exported++
		if isRequired(f) {
			return false
		}
	}
	return exported > threshold
}

// isRequired reports whether a struct field is required: a non-pointer with a
// json tag that does not contain "omitempty".
func isRequired(f reflect.StructField) bool {
	if f.Type.Kind() == reflect.Pointer {
		return false
	}
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return false
	}
	return !strings.Contains(tag, "omitempty")
}

// shouldSkipType filters out the SDK's internal /utils helper types.
func shouldSkipType(t reflect.Type) bool {
	return strings.Contains(t.PkgPath(), "/utils")
}

// jsonFieldName returns the field name from a json struct tag, dropping options
// like ",omitempty".
func jsonFieldName(tag string) string {
	parts := strings.Split(tag, ",")
	if parts[0] != "" {
		return parts[0]
	}
	return tag
}
