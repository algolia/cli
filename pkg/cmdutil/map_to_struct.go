package cmdutil

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
)

// capitalize makes the first letter of a word uppercase
func capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	firstRune := []rune(word)[0]
	rest := []rune(word)[1:]
	return string(unicode.ToUpper(firstRune)) + string(rest)
}

// MapToStruct converts a map into a struct
func MapToStruct(m map[string]any, s interface{}) error {
	val := reflect.ValueOf(s).Elem()

	for k, v := range m {
		// cmdline options are lowercase (`--query`),
		// but struct fields are capital (`Query`)
		field := val.FieldByName(capitalize(k))
		if !field.IsValid() {
			return errors.New(fmt.Sprintf("No such parameter: %s for browse\n.", k))
		}

		if !field.CanSet() {
			return errors.New(fmt.Sprintf("Can't set field: %s\n", field))
		}

		fieldValue := reflect.ValueOf(v)

		if field.Type().Kind() == reflect.Ptr &&
			fieldValue.Type().ConvertibleTo(field.Type().Elem()) {
			newValue := reflect.New(fieldValue.Type()).Elem()
			newValue.Set(fieldValue)
			field.Set(newValue.Addr())
		} else if fieldValue.Type().ConvertibleTo(field.Type()) {
			field.Set(fieldValue.Convert(field.Type()))
		} else {
			return errors.New(fmt.Sprintf("Can't convert type of %s to %s\n", fieldValue.Type(), field.Type()))
		}
	}
	return nil
}
