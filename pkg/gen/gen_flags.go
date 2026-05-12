//go:build gen_flags
// +build gen_flags

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

type TemplateData struct {
	SpecFlags map[string]*SpecFlags
}

type SpecFlags struct {
	Flags map[string]*SpecFlag
}

type SpecFlag struct {
	Def        interface{}
	Type       string
	GoType     string
	SubType    string
	OneOf      []string
	Usage      string
	Categories []string
}

const (
	pathTemplate = "../../gen/flags.go.tpl"
	pathName     = "flags.go.tpl"
)

// SpecConfig describes a single OpenAPI document and which schemas to flatten
// into Cobra flag bindings.
type SpecConfig struct {
	File       string   // path to the spec, relative to pkg/gen
	Schemas    []string // top-level component schema names
	DocLinks   bool     // append Algolia search-API doc links to descriptions
	OutputPath string   // generated Go file (relative to pkg/gen)
}

var specs = []SpecConfig{
	{
		File: "../../../api/specs/search.yml",
		Schemas: []string{
			"searchParamsObject",
			"browseParamsObject",
			"indexSettings",
			"deleteByParams",
		},
		DocLinks:   true,
		OutputPath: "../../cmdutil/spec_flags.go",
	},
	{
		File: "../../../api/specs/agent-studio.json",
		Schemas: []string{
			"AgentConfigCreate",
			"AgentCompletionRequest",
		},
		DocLinks:   false,
		OutputPath: "../../cmdutil/agent_studio_flags.go",
	},
}

func main() {
	tmpl := template.Must(template.
		New(pathName).
		Funcs(template.FuncMap{
			"capitalize": func(s string) string {
				return strings.Title(s)
			},
		}).
		ParseFiles(pathTemplate))

	for _, spec := range specs {
		data, err := getTemplateData(spec)
		if err != nil {
			panic(fmt.Errorf("loading %s: %w", spec.File, err))
		}

		var result bytes.Buffer
		if err := tmpl.Execute(&result, data); err != nil {
			panic(err)
		}

		formatted, err := format.Source(result.Bytes())
		if err != nil {
			panic(err)
		}

		fmt.Printf("writing %s\n", spec.OutputPath)
		if err := os.WriteFile(spec.OutputPath, formatted, 0o600); err != nil {
			panic(err)
		}
	}
}

// loadProperties recursively loads the properties of the given schemaRef.
func loadProperties(schemaRef *openapi3.SchemaRef) map[string]*openapi3.Schema {
	properties := make(map[string]*openapi3.Schema)

	// Load the direct properties of the current  (ex: `deleteByParams`)
	if schemaRef.Value.Properties != nil {
		for name, property := range schemaRef.Value.Properties {
			properties[name] = property.Value
		}
	}

	// Load the properties of the allOf schemas (ex: `searchParamsObject`)
	for _, schema := range schemaRef.Value.AllOf {
		if schema.Value.Properties != nil {
			for name, param := range schema.Value.Properties {
				properties[name] = param.Value
			}
		} else {
			for name, param := range loadProperties(schema) {
				properties[name] = param
			}
		}
	}

	return properties
}

// loadSpecs loads the parameters from a OpenAPI 3.0 schema.
func loadSpecs(specFile, specName string) (map[string]*openapi3.Schema, error) {
	doc, err := openapi3.NewLoader().LoadFromFile(specFile)
	if err != nil {
		return nil, err
	}

	schemaRef, ok := doc.Components.Schemas[specName]
	if !ok {
		return nil, fmt.Errorf("schema %s not found", specName)
	}

	return loadProperties(schemaRef), nil
}

// getTemplateData loads all flags for a single SpecConfig.
func getTemplateData(spec SpecConfig) (TemplateData, error) {
	data := TemplateData{SpecFlags: make(map[string]*SpecFlags)}
	for _, name := range spec.Schemas {
		params, err := loadSpecs(spec.File, name)
		if err != nil {
			return data, err
		}
		data.SpecFlags[name] = getFlags(params, spec.DocLinks)
	}
	return data, nil
}

// getFlags returns the flags for the given spec.
func getFlags(params map[string]*openapi3.Schema, withDocLinks bool) *SpecFlags {
	flags := &SpecFlags{Flags: make(map[string]*SpecFlag)}
	for name, param := range params {
		flags.Flags[name] = getFlag(name, param, withDocLinks)
	}
	return flags
}

// unwrapNullable handles OpenAPI 3.1 `anyOf: [T, {type: null}]` nullability.
// When a schema has no Type set but its AnyOf contains exactly one non-null
// branch and one explicit null branch, return the non-null branch. Otherwise
// return the input unchanged.
func unwrapNullable(s *openapi3.Schema) *openapi3.Schema {
	if s == nil || s.Type != nil || len(s.AnyOf) == 0 {
		return s
	}
	var nonNull *openapi3.Schema
	for _, branch := range s.AnyOf {
		if branch == nil || branch.Value == nil {
			continue
		}
		if branch.Value.Type != nil && branch.Value.Type.Is("null") {
			continue
		}
		if nonNull != nil {
			return s // more than one non-null branch — leave as-is
		}
		nonNull = branch.Value
	}
	if nonNull == nil {
		return s
	}
	return nonNull
}

// schemaTypeString returns the primary OpenAPI type for a schema, accounting
// for the v0.135 *Types representation. Returns "" if no concrete type is set.
func schemaTypeString(s *openapi3.Schema) string {
	if s == nil || s.Type == nil {
		return ""
	}
	for _, t := range *s.Type {
		if t != "" && t != "null" {
			return t
		}
	}
	return ""
}

// GetGoType returns the Go type for the given OpenAPI 3.0/3.1 schema.
func GetGoType(param *openapi3.Schema) string {
	param = unwrapNullable(param)
	specTypeGoType := map[string]string{
		"string":  "string",
		"integer": "int",
		"number":  "float64",
		"boolean": "bool",
	}
	t := schemaTypeString(param)
	if t == "array" && param.Items != nil && param.Items.Value != nil {
		return "[]" + GetGoType(param.Items.Value)
	}
	return specTypeGoType[t]
}

// getFlag returns the flag for the given parameter.
func getFlag(name string, param *openapi3.Schema, withDocLinks bool) *SpecFlag {
	param = unwrapNullable(param)
	t := schemaTypeString(param)

	subType := ""
	if t == "array" && param.Items != nil && param.Items.Value != nil {
		subType = schemaTypeString(unwrapNullable(param.Items.Value))
	}

	// Arrays of unions / objects can't be modeled as typed slices; fall through
	// to the JSONVar branch in the template by clearing the type.
	if t == "array" && subType != "string" && subType != "integer" && subType != "number" {
		t = ""
	}

	var categories []string
	if raw, ok := param.Extensions["x-categories"]; ok {
		switch v := raw.(type) {
		case []any:
			for _, item := range v {
				if s, ok := item.(string); ok {
					categories = append(categories, s)
				}
			}
		case []string:
			categories = append(categories, v...)
		}
	}

	flag := &SpecFlag{
		Def:        param.Default,
		Type:       t,
		GoType:     GetGoType(param),
		Usage:      getDescription(name, param, withDocLinks),
		SubType:    subType,
		Categories: categories,
	}

	if param.OneOf != nil {
		for _, oneOf := range param.OneOf {
			flag.OneOf = append(flag.OneOf, schemaTypeString(unwrapNullable(oneOf.Value)))
		}
	}

	return flag
}

// shortDescription returns the first sentence of the parameter description.
func shortDescription(description string) string {
	// Handle sentences ending with a colon
	s := strings.Split(description, ":\n")
	// Handle sentences ending with a period
	s = strings.Split(s[0], ".\n")
	s[0] = replaceMarkdownLinks(s[0])
	s[0] = strings.ReplaceAll(s[0], "`", "")

	if !strings.HasSuffix(s[0], ".") {
		s[0] += "."
	}

	return strings.TrimSpace(s[0])
}

func replaceMarkdownLinks(text string) string {
	re := regexp.MustCompile("\\[([^\\[\\]]*)\\]\\([^\\(\\)]*\\)")
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		linkText := match[1]
		text = strings.Replace(text, match[0], linkText, 1)
	}

	return text
}

// getDescription returns the short description for the given parameter.
// It's the first sentence of the parameter description followed by possible values if it's an enum,
// followed by a link to the API param reference page when withDocLinks is true.
func getDescription(name string, param *openapi3.Schema, withDocLinks bool) string {
	// These params don't have an API reference page
	if name == "semanticSearch" || name == "cursor" || name == "reRankingApplyFilter" {
		withDocLinks = false
	}

	description := shortDescription(param.Description)

	// Add choices if param is an enum
	if param.Enum != nil {
		choices := make([]string, len(param.Enum))
		for i, e := range param.Enum {
			choices[i] = fmt.Sprintf("%v", e)
		}
		description = fmt.Sprintf("%s One of: %v.", description, strings.Join(choices, ", "))
	}

	// Add link to the API param reference page
	if withDocLinks {
		link := fmt.Sprintf("https://www.algolia.com/doc/api-reference/api-parameters/%s/", name)
		description = fmt.Sprintf("%s\nSee: %s", description, link)
	}
	return description
}
