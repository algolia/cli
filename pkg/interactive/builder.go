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

// SuggestFunc returns completion candidates for a field, where scope holds the string values entered so far keyed by JSON name.
type SuggestFunc func(fieldPath, jsonName string, scope map[string]string) []string

// Builder prompts for every field of a struct via reflection.
type Builder struct {
	Prompter          Prompter
	MaxDepth          int         // 0 uses DefaultMaxDepth
	ParamBagThreshold int         // 0 uses DefaultParamBagThreshold
	Suggest           SuggestFunc // nil disables dynamic completion

	scope map[string]string // string values entered so far, keyed by JSON name
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
	b.scope = map[string]string{}
	return b.buildValue(rv.Elem().Type().Name(), rv.Elem(), 0, nil, leafMeta{})
}

// leafMeta carries a field's generated help and example, and is empty when the SDK exposes none.
type leafMeta struct {
	help     string
	example  string
	jsonName string
}

func (m leafMeta) promptOptions() []PromptOption {
	return []PromptOption{WithHelp(m.help), WithDefault(m.example)}
}

// resolveFieldMeta returns zero help and example when the model has no FieldDescription or FieldExample methods.
func resolveFieldMeta(owner reflect.Value, jsonName string) leafMeta {
	return leafMeta{
		help:     callStringField(owner, "FieldDescription", jsonName),
		example:  callStringField(owner, "FieldExample", jsonName),
		jsonName: jsonName,
	}
}

func callStringField(owner reflect.Value, method, arg string) string {
	m := findMethod(owner, method)
	if !m.IsValid() || !isStringToString(m) {
		return ""
	}
	return m.Call([]reflect.Value{reflect.ValueOf(arg)})[0].String()
}

// findMethod falls back to the address of owner so pointer-receiver methods are found.
func findMethod(owner reflect.Value, name string) reflect.Value {
	m := owner.MethodByName(name)
	if !m.IsValid() && owner.CanAddr() {
		m = owner.Addr().MethodByName(name)
	}
	return m
}

func isStringToString(m reflect.Value) bool {
	mt := m.Type()
	return mt.NumIn() == 1 && mt.NumOut() == 1 &&
		mt.In(0) == reflect.TypeFor[string]() && mt.Out(0) == reflect.TypeFor[string]()
}

func (b *Builder) buildValue(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
	if depth > b.maxDepth() || slices.Contains(stack, rv.Type()) {
		return b.promptRawJSON(label, rv, meta)
	}

	switch rv.Kind() {
	case reflect.Pointer:
		return b.buildPointer(label, rv, depth, stack, meta)
	case reflect.Struct:
		return b.buildStruct(label, rv, depth, stack, meta)
	case reflect.Slice:
		return b.buildSlice(label, rv, depth, stack, meta)
	case reflect.Map:
		return b.buildMap(label, rv, depth, stack, meta)
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		_, err := b.assignScalar(label, rv, false, meta)
		return err
	case reflect.Interface:
		return b.promptRawJSON(label, rv, meta)
	default:
		return fmt.Errorf("interactive: unsupported kind %s for %s", rv.Kind(), label)
	}
}

func (b *Builder) inputInt(label string, bits int, opts ...PromptOption) (int64, error) {
	s, err := b.Prompter.Input(label+" (integer)", countValidator(), opts...)
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

// enumSkip is the Select option that leaves an optional enum unset.
const enumSkip = "(skip)"

// enumValues reports ok=false unless Values() returns a slice of v's own type, so callers fall back to IsValid.
func enumValues(v reflect.Value) ([]reflect.Value, bool) {
	m := findMethod(v, "Values")
	if !m.IsValid() {
		return nil, false
	}
	mt := m.Type()
	if mt.NumIn() != 0 || mt.NumOut() != 1 || mt.Out(0).Kind() != reflect.Slice || mt.Out(0).Elem() != v.Type() {
		return nil, false
	}
	list := m.Call(nil)[0]
	out := make([]reflect.Value, list.Len())
	for i := range out {
		out[i] = list.Index(i)
	}
	return out, len(out) > 0
}

// assignEnum writes the picked enum value into v, and an optional enum also offers enumSkip.
func (b *Builder) assignEnum(label string, v reflect.Value, optional bool, values []reflect.Value, meta leafMeta) (bool, error) {
	labels := make([]string, len(values))
	for i, ev := range values {
		labels[i] = fmt.Sprint(ev.Interface())
	}
	options := labels
	if optional {
		options = append([]string{enumSkip}, labels...)
	}
	pick, err := b.Prompter.Select(label, options, meta.promptOptions()...)
	if err != nil {
		return false, err
	}
	if optional && pick == enumSkip {
		return false, nil
	}
	for i, l := range labels {
		if l == pick {
			v.Set(values[i])
			return true, nil
		}
	}
	return false, fmt.Errorf("%q is not a valid %s", pick, v.Type().Name())
}

// suggestFor returns nil when there is no Suggest callback and no example, otherwise Tab fills the SDK example even without a Suggest callback.
func (b *Builder) suggestFor(path string, meta leafMeta) func(string) []string {
	if b.Suggest == nil && meta.example == "" {
		return nil
	}
	return func(string) []string {
		var out []string
		if b.Suggest != nil {
			out = b.Suggest(path, meta.jsonName, b.scope)
		}
		if meta.example != "" {
			out = append(out, meta.example)
		}
		return out
	}
}

// assignScalar returns whether v was set, and an empty answer skips the field only when optional.
func (b *Builder) assignScalar(label string, v reflect.Value, optional bool, meta leafMeta) (bool, error) {
	// SDK enums expose their allowed set via Values(), so a Select replaces free text and the value cannot be mistyped.
	if values, ok := enumValues(v); ok {
		return b.assignEnum(label, v, optional, values, meta)
	}
	switch v.Kind() {
	case reflect.String:
		return b.assignString(label, v, optional, meta)
	case reflect.Bool:
		return b.assignBool(label, v, optional, meta)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return b.assignInt(label, v, optional, meta)
	case reflect.Float32, reflect.Float64:
		return b.assignFloat(label, v, optional, meta)
	default:
		return false, fmt.Errorf("interactive: unsupported scalar kind %s for %s", v.Kind(), label)
	}
}

// stringValidator defers to the SDK's IsValid for enum-like strings, so a hand-maintained list never drifts.
func stringValidator(t reflect.Type, label string, optional bool) func(string) error {
	switch {
	case isEnumType(t):
		return enumValidator(t, label, optional)
	case !optional:
		return requiredString(label)
	default:
		return nil
	}
}

func (b *Builder) assignString(label string, v reflect.Value, optional bool, meta leafMeta) (bool, error) {
	promptOpts := meta.promptOptions()
	if fn := b.suggestFor(label, meta); fn != nil {
		promptOpts = append(promptOpts, WithSuggest(fn))
	}
	s, err := b.Prompter.Input(label, stringValidator(v.Type(), label, optional), promptOpts...)
	if err != nil || (optional && s == "") {
		return false, err
	}
	v.SetString(s)
	if meta.jsonName != "" {
		b.scope[meta.jsonName] = s
	}
	return true, nil
}

// assignBool asks a Confirm for a required bool, but takes free-text true/false for an optional *bool so an empty answer can mean skip.
func (b *Builder) assignBool(label string, v reflect.Value, optional bool, meta leafMeta) (bool, error) {
	if !optional {
		val, err := b.Prompter.Confirm(label, meta.promptOptions()...)
		if err != nil {
			return false, err
		}
		v.SetBool(val)
		return true, nil
	}
	s, err := b.Prompter.Input(label+" (true/false, empty to skip)", boolValidator(), meta.promptOptions()...)
	if err != nil || s == "" {
		return false, err
	}
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("invalid boolean for %s: %w", label, err)
	}
	v.SetBool(val)
	return true, nil
}

// scalarHint is the suffix that tells the user what to type and whether an empty answer skips.
func scalarHint(kind string, optional bool) string {
	if optional {
		return " (" + kind + ", empty to skip)"
	}
	return " (" + kind + ")"
}

// assignInt treats an empty answer as skip, which is only reachable when optional because intValidator rejects empty otherwise.
func (b *Builder) assignInt(label string, v reflect.Value, optional bool, meta leafMeta) (bool, error) {
	bits := v.Type().Bits()
	s, err := b.Prompter.Input(label+scalarHint("integer", optional), intValidator(label, bits, optional), meta.promptOptions()...)
	if err != nil || s == "" {
		return false, err
	}
	n, err := strconv.ParseInt(s, 10, bits)
	if err != nil {
		return false, fmt.Errorf("invalid integer for %s: %w", label, err)
	}
	v.SetInt(n)
	return true, nil
}

// assignFloat parses at the field's precision so a float32 overflow fails here instead of becoming +Inf, which encoding/json cannot marshal.
func (b *Builder) assignFloat(label string, v reflect.Value, optional bool, meta leafMeta) (bool, error) {
	bits := v.Type().Bits()
	s, err := b.Prompter.Input(label+scalarHint("number", optional), floatValidator(label, bits, optional), meta.promptOptions()...)
	if err != nil || s == "" {
		return false, err
	}
	f, err := strconv.ParseFloat(s, bits)
	if err != nil {
		return false, fmt.Errorf("invalid number for %s: %w", label, err)
	}
	v.SetFloat(f)
	return true, nil
}

func (b *Builder) buildPointer(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
	elem := rv.Type().Elem()
	switch elem.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		// Optional scalar: build into a fresh element and only assign the pointer
		// if the user provided a value.
		ptr := reflect.New(elem)
		set, err := b.assignScalar(label, ptr.Elem(), true, meta)
		if err != nil || !set {
			return err
		}
		rv.Set(ptr)
		return nil
	}

	// Pointer to struct, slice, map, or other composite.
	set, err := b.Prompter.Confirm("Set "+label+"?", meta.promptOptions()...)
	if err != nil || !set {
		return err
	}
	// The Confirm above already showed the description, so the prompts inside do not repeat it.
	meta.help = ""
	ptr := reflect.New(elem)
	if err := b.buildValue(label, ptr.Elem(), depth+1, stack, meta); err != nil {
		return err
	}
	rv.Set(ptr)
	return nil
}

func (b *Builder) promptRawJSON(label string, rv reflect.Value, meta leafMeta) error {
	// Hand-written SDK examples are not always valid JSON, and a non-JSON default would block skipping.
	if !json.Valid([]byte(meta.example)) {
		meta.example = ""
	}
	raw, err := b.Prompter.Input(label+" (raw JSON, empty to skip)", jsonValidator(), meta.promptOptions()...)
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

func (b *Builder) buildStruct(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
	t := rv.Type()
	if shouldSkipType(t) {
		return nil
	}
	if isUnionType(t) {
		return b.buildUnion(label, rv, depth, stack, meta)
	}
	if isParamBag(t, b.paramBagThreshold()) {
		return b.buildParamBag(label, rv, depth, stack, meta)
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
		jsonName := jsonFieldName(jsonTag)
		fieldLabel := label + "." + jsonName
		if err := b.buildValue(fieldLabel, rv.Field(i), depth+1, stack, resolveFieldMeta(rv, jsonName)); err != nil {
			return err
		}
	}
	return nil
}

type unionVariant struct {
	fieldIdx int
	name     string
}

// unionVariants lists the pointer fields of a oneOf wrapper, named by their type.
func unionVariants(t reflect.Type) []unionVariant {
	var variants []unionVariant
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() || f.Type.Kind() != reflect.Pointer {
			continue
		}
		name := f.Type.Elem().Name()
		if name == "" {
			name = f.Name
		}
		variants = append(variants, unionVariant{i, name})
	}
	return variants
}

// variantLabel appends the variant's description when the owner exposes one through VariantDescription.
func variantLabel(owner reflect.Value, name string) string {
	desc := callStringField(owner, "VariantDescription", name)
	if desc == "" {
		return name
	}
	return name + ": " + desc
}

func (b *Builder) buildUnion(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
	variants := unionVariants(rv.Type())
	if len(variants) == 0 {
		return nil
	}
	labels := make([]string, len(variants))
	for i, v := range variants {
		labels[i] = variantLabel(rv, v.name)
	}
	pick, err := b.Prompter.Select(label+" (variant)", labels, WithHelp(meta.help))
	if err != nil {
		return err
	}
	idx := slices.Index(labels, pick)
	if idx < 0 {
		return nil
	}
	return b.buildVariant(label, rv, variants[idx], depth, stack)
}

func (b *Builder) buildVariant(label string, rv reflect.Value, v unionVariant, depth int, stack []reflect.Type) error {
	field := rv.Field(v.fieldIdx)
	ptr := reflect.New(field.Type().Elem())
	stack = append(stack, rv.Type())
	if err := b.buildValue(label+"."+v.name, ptr.Elem(), depth+1, stack, leafMeta{}); err != nil {
		return err
	}
	field.Set(ptr)
	return nil
}

func (b *Builder) buildParamBag(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
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
	picked, err := b.Prompter.MultiSelect(label, options, WithHelp(meta.help))
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
		if err := b.buildValue(fieldLabel, rv.Field(e.fieldIdx), depth+1, stack, resolveFieldMeta(rv, e.jsonName)); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildSlice(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
	count, err := b.inputInt("how many "+label, 32, WithHelp(meta.help))
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
		if err := b.buildValue(itemLabel, out.Index(i), depth+1, stack, leafMeta{}); err != nil {
			return err
		}
	}
	rv.Set(out)
	return nil
}

func (b *Builder) buildMap(label string, rv reflect.Value, depth int, stack []reflect.Type, meta leafMeta) error {
	t := rv.Type()
	if t.Key().Kind() != reflect.String {
		return b.promptRawJSON(label, rv, meta)
	}
	count, err := b.inputInt("how many "+label+" entries", 32, WithHelp(meta.help))
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
		if out.MapIndex(reflect.ValueOf(key)).IsValid() {
			overwrite, err := b.Prompter.Confirm(fmt.Sprintf("%q is already set; overwrite?", key))
			if err != nil {
				return err
			}
			if !overwrite {
				continue
			}
		}
		valPtr := reflect.New(valType)
		if err := b.buildValue(fmt.Sprintf("%s[%q]", label, key), valPtr.Elem(), depth+1, stack, leafMeta{}); err != nil {
			return err
		}
		out.SetMapIndex(reflect.ValueOf(key), valPtr.Elem())
	}
	rv.Set(out)
	return nil
}
