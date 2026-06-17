package interactive

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
)

// DefaultMaxDepth bounds recursion to defeat self-referential unions. When the
// guard fires the builder degrades to a raw-JSON prompt.
const DefaultMaxDepth = 20

// DefaultParamBagThreshold is the optional-field count above which a struct is
// treated as a parameter bag (the user multi-selects which fields to set).
const DefaultParamBagThreshold = 15

// Builder prompts for every field of a struct via reflection.
type Builder struct {
	Prompter          Prompter
	MaxDepth          int // 0 uses DefaultMaxDepth
	ParamBagThreshold int // 0 uses DefaultParamBagThreshold
}

func (b *Builder) maxDepth() int {
	if b.MaxDepth <= 0 {
		return DefaultMaxDepth
	}
	return b.MaxDepth
}

func (b *Builder) paramBagThreshold() int {
	if b.ParamBagThreshold <= 0 {
		return DefaultParamBagThreshold
	}
	return b.ParamBagThreshold
}

// Build reflectively prompts for every field of the struct pointed to by v.
// v must be a non-nil pointer to a struct. Fields already non-zero on v are
// preserved (lets the caller pre-populate identifiers).
func (b *Builder) Build(v any) error {
	if b.Prompter == nil {
		return errors.New("interactive: Builder.Prompter must be set")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return errors.New("interactive: Build needs a non-nil pointer to a struct")
	}
	return b.buildValue(rv.Elem().Type().Name(), rv.Elem(), 0, nil)
}

func (b *Builder) buildValue(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	if depth > b.maxDepth() || slices.Contains(stack, rv.Type()) {
		return b.promptRawJSON(label, rv)
	}

	switch rv.Kind() {
	case reflect.Pointer:
		return b.buildPointer(label, rv, depth, stack)
	case reflect.Struct:
		return b.buildStruct(label, rv, depth, stack)
	case reflect.Slice:
		return b.buildSlice(label, rv, depth, stack)
	case reflect.Map:
		return b.buildMap(label, rv, depth, stack)
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		_, err := b.assignScalar(label, rv, false)
		return err
	case reflect.Interface:
		return b.promptRawJSON(label, rv)
	default:
		return fmt.Errorf("interactive: unsupported kind %s for %s", rv.Kind(), label)
	}
}

func (b *Builder) inputInt(label string, bits int) (int64, error) {
	s, err := b.Prompter.Input(label+" (integer)", countValidator())
	if err != nil {
		return 0, err
	}
	if s == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(s, 10, bits)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s: %w", label, err)
	}
	return n, nil
}

// isEnumType reports whether t is a named string type exposing a value-receiver
// IsValid() bool method, the convention every Algolia SDK enum is generated
// with. Such a field is validated against the SDK's own allowed-value check
// rather than a hand-maintained list.
func isEnumType(t reflect.Type) bool {
	if t.Kind() != reflect.String || t.Name() == "string" {
		return false
	}
	m, ok := t.MethodByName("IsValid")
	if !ok {
		return false
	}
	ft := m.Func.Type() // the receiver counts as the first in-param
	return ft.NumIn() == 1 && ft.NumOut() == 1 && ft.Out(0).Kind() == reflect.Bool
}

// enumValidator validates that s is an allowed value of enum type t by calling
// the type's IsValid() method reflectively. Empty is allowed only when optional.
func enumValidator(t reflect.Type, label string, optional bool) func(string) error {
	return func(s string) error {
		if s == "" {
			if optional {
				return nil
			}
			return fmt.Errorf("%s is required", label)
		}
		cand := reflect.New(t).Elem()
		cand.SetString(s)
		if !cand.MethodByName("IsValid").Call(nil)[0].Bool() {
			return fmt.Errorf("%q is not a valid %s", s, t.Name())
		}
		return nil
	}
}

// assignScalar prompts for a scalar value and writes it into v, which must be
// settable. optional reports whether an empty answer means "skip" (leaving v
// untouched) rather than writing the zero value, and selects the prompt hints.
// It returns whether v was set. Bool keeps a deliberate asymmetry: a required
// bool is a yes/no Confirm, while an optional *bool is free-text true/false with
// an empty answer meaning skip.
func (b *Builder) assignScalar(label string, v reflect.Value, optional bool) (bool, error) {
	switch v.Kind() {
	case reflect.String:
		// SDK enums (named string types with an IsValid() bool method) are
		// validated against the SDK's own allowed-value check; plain strings use
		// the required/optional rule.
		var validate func(string) error
		switch {
		case isEnumType(v.Type()):
			validate = enumValidator(v.Type(), label, optional)
		case !optional:
			validate = requiredString(label)
		}
		s, err := b.Prompter.Input(label, validate)
		if err != nil {
			return false, err
		}
		if optional && s == "" {
			return false, nil
		}
		v.SetString(s)
		return true, nil
	case reflect.Bool:
		if optional {
			s, err := b.Prompter.Input(label+" (true/false, empty to skip)", boolValidator())
			if err != nil || s == "" {
				return false, err
			}
			val, perr := strconv.ParseBool(s)
			if perr != nil {
				return false, fmt.Errorf("invalid boolean for %s: %w", label, perr)
			}
			v.SetBool(val)
			return true, nil
		}
		val, err := b.Prompter.Confirm(label)
		if err != nil {
			return false, err
		}
		v.SetBool(val)
		return true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		hint := " (integer)"
		if optional {
			hint = " (integer, empty to skip)"
		}
		s, err := b.Prompter.Input(label+hint, intValidator(label, v.Type().Bits(), optional))
		if err != nil {
			return false, err
		}
		if s == "" {
			// Only reachable when optional: intValidator rejects empty for
			// required fields, so the prompter re-prompts or errors before here.
			return false, nil
		}
		n, perr := strconv.ParseInt(s, 10, v.Type().Bits())
		if perr != nil {
			return false, fmt.Errorf("invalid integer for %s: %w", label, perr)
		}
		v.SetInt(n)
		return true, nil
	case reflect.Float32, reflect.Float64:
		hint := " (number)"
		if optional {
			hint = " (number, empty to skip)"
		}
		s, err := b.Prompter.Input(label+hint, floatValidator(label, v.Type().Bits(), optional))
		if err != nil {
			return false, err
		}
		if s == "" {
			// Only reachable when optional (see the integer case above).
			return false, nil
		}
		// Parse at the field's precision so a float32 overflow is caught here
		// rather than later becoming +Inf, which encoding/json cannot marshal.
		f, perr := strconv.ParseFloat(s, v.Type().Bits())
		if perr != nil {
			return false, fmt.Errorf("invalid number for %s: %w", label, perr)
		}
		v.SetFloat(f)
		return true, nil
	default:
		return false, fmt.Errorf("interactive: unsupported scalar kind %s for %s", v.Kind(), label)
	}
}

func (b *Builder) buildPointer(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	elem := rv.Type().Elem()
	switch elem.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		// Optional scalar: build into a fresh element and only assign the pointer
		// if the user provided a value.
		ptr := reflect.New(elem)
		set, err := b.assignScalar(label, ptr.Elem(), true)
		if err != nil || !set {
			return err
		}
		rv.Set(ptr)
		return nil
	}

	// Pointer to struct, slice, map, or other composite.
	set, err := b.Prompter.Confirm("Set " + label + "?")
	if err != nil || !set {
		return err
	}
	ptr := reflect.New(elem)
	if err := b.buildValue(label, ptr.Elem(), depth+1, stack); err != nil {
		return err
	}
	rv.Set(ptr)
	return nil
}

func (b *Builder) promptRawJSON(label string, rv reflect.Value) error {
	raw, err := b.Prompter.Input(label+" (raw JSON, empty to skip)", jsonValidator())
	if err != nil || raw == "" {
		return err
	}
	if !rv.CanAddr() {
		return fmt.Errorf("interactive: cannot unmarshal raw JSON into unaddressable %s", label)
	}
	if err := json.Unmarshal([]byte(raw), rv.Addr().Interface()); err != nil {
		return fmt.Errorf("invalid JSON for %s: %w", label, err)
	}
	return nil
}

func (b *Builder) buildStruct(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	t := rv.Type()
	if shouldSkipType(t) {
		return nil
	}
	if isUnionType(t) {
		return b.buildUnion(label, rv, depth, stack)
	}
	if isParamBag(t, b.paramBagThreshold()) {
		return b.buildParamBag(label, rv, depth, stack)
	}
	stack = append(stack, t)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		// Preserve fields the caller pre-populated.
		if !rv.Field(i).IsZero() {
			continue
		}
		fieldLabel := label + "." + jsonFieldName(jsonTag)
		if err := b.buildValue(fieldLabel, rv.Field(i), depth+1, stack); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildUnion(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	t := rv.Type()
	type variant struct {
		fieldIdx int
		name     string
	}
	var variants []variant
	var labels []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() || f.Type.Kind() != reflect.Pointer {
			continue
		}
		name := f.Type.Elem().Name()
		if name == "" {
			name = f.Name
		}
		variants = append(variants, variant{i, name})
		labels = append(labels, name)
	}
	if len(variants) == 0 {
		return nil
	}
	pick, err := b.Prompter.Select(label+" (variant)", labels)
	if err != nil {
		return err
	}
	for _, v := range variants {
		if v.name != pick {
			continue
		}
		field := rv.Field(v.fieldIdx)
		ptr := reflect.New(field.Type().Elem())
		stack = append(stack, t)
		if err := b.buildValue(label+"."+v.name, ptr.Elem(), depth+1, stack); err != nil {
			return err
		}
		field.Set(ptr)
		return nil
	}
	return nil
}

func (b *Builder) buildParamBag(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	t := rv.Type()
	type entry struct {
		fieldIdx int
		jsonName string
		typeName string
	}
	var entries []entry
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		ft := f.Type
		for ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		typeName := ft.Name()
		if typeName == "" {
			typeName = ft.Kind().String()
		}
		entries = append(entries, entry{i, jsonFieldName(jsonTag), typeName})
	}
	if len(entries) == 0 {
		return nil
	}
	options := make([]string, len(entries))
	for i, e := range entries {
		options[i] = fmt.Sprintf("%s (%s)", e.jsonName, e.typeName)
	}
	picked, err := b.Prompter.MultiSelect(label, options)
	if err != nil {
		return err
	}
	stack = append(stack, t)
	for _, idx := range picked {
		if idx < 0 || idx >= len(entries) {
			continue // defend against a Prompter returning out-of-range indexes
		}
		e := entries[idx]
		fieldLabel := label + "." + e.jsonName
		if err := b.buildValue(fieldLabel, rv.Field(e.fieldIdx), depth+1, stack); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildSlice(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	count, err := b.inputInt("how many "+label, 32)
	if err != nil {
		return err
	}
	if count <= 0 {
		return nil
	}
	t := rv.Type()
	out := reflect.MakeSlice(t, int(count), int(count))
	for i := 0; i < int(count); i++ {
		itemLabel := fmt.Sprintf("%s[%d]", label, i)
		if err := b.buildValue(itemLabel, out.Index(i), depth+1, stack); err != nil {
			return err
		}
	}
	rv.Set(out)
	return nil
}

func (b *Builder) buildMap(label string, rv reflect.Value, depth int, stack []reflect.Type) error {
	t := rv.Type()
	if t.Key().Kind() != reflect.String {
		return b.promptRawJSON(label, rv)
	}
	count, err := b.inputInt("how many "+label+" entries", 32)
	if err != nil {
		return err
	}
	if count <= 0 {
		// Leave the field at its zero value (nil map) so it is omitted, matching
		// buildSlice. A confirmed-but-empty map would otherwise serialize as {}.
		return nil
	}
	valType := t.Elem()
	out := reflect.MakeMapWithSize(t, int(count))
	for i := 0; i < int(count); i++ {
		keyLabel := fmt.Sprintf("%s key[%d]", label, i)
		key, err := b.Prompter.Input(keyLabel, requiredString(keyLabel))
		if err != nil {
			return err
		}
		valPtr := reflect.New(valType)
		if err := b.buildValue(fmt.Sprintf("%s[%q]", label, key), valPtr.Elem(), depth+1, stack); err != nil {
			return err
		}
		out.SetMapIndex(reflect.ValueOf(key), valPtr.Elem())
	}
	rv.Set(out)
	return nil
}
