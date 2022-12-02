//go:build gen_flags
// +build gen_flags

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
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
	searchSpecFile = "../../../api/specs/search.yml"
	pathTemplate   = "../../gen/flags.go.tpl"
	pathTemplateTs = "../../gen/flags.ts.tpl"
	pathName       = "flags.go.tpl"
	pathNameTs     = "flags.ts.tpl"
	pathOutput     = "../../cmdutil/spec_flags.go"
	pathOutputTs   = "../../gen/spec_flags.ts"
)

func main() {
	// This is the script that generates the `flags.go` file from the
	// OpenAPI spec file.

	specNames := []string{"searchParamsObject", "browseParamsObject", "indexSettings"}
	templateData, err := getTemplateData(specNames)
	if err != nil {
		panic(err)
	}

	// Load the template with a custom function map
	tmpl := template.Must(template.
		// Note that the template name MUST match the file name
		New(pathName).
		Funcs(template.FuncMap{
			"capitalize": func(s string) string {
				return strings.Title(s)
			},
		}).
		ParseFiles(pathTemplate))

	// Execute the template
	var result bytes.Buffer
	err = tmpl.Execute(&result, templateData)
	if err != nil {
		panic(err)
	}

	// Format the output of the template execution
	formatted, err := format.Source(result.Bytes())
	if err != nil {
		panic(err)
	}

	// Write the formatted source code to disk
	fmt.Printf("writing %s\n", pathOutput)
	err = ioutil.WriteFile(pathOutput, formatted, 0644)
	if err != nil {
		panic(err)
	}

	generateTs()
}

func generateTs() {
	// This is the script that generates the `flags.go` file from the
	// OpenAPI spec file.

	specNames := []string{"searchParamsObject"}
	templateData, err := getTemplateData(specNames)
	if err != nil {
		panic(err)
	}

	// Load the template with a custom function map
	tmpl := template.Must(template.
		// Note that the template name MUST match the file name
		New(pathNameTs).
		ParseFiles(pathTemplateTs))

	// Execute the template
	var result bytes.Buffer
	err = tmpl.Execute(&result, templateData.SpecFlags["searchParamsObject"])
	if err != nil {
		panic(err)
	}

	// Write the formatted source code to disk
	fmt.Printf("writing %s\n", pathOutputTs)
	err = ioutil.WriteFile(pathOutputTs, result.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

// loadProperties recursively loads the properties of the given schemaRef.
func loadProperties(schemaRef *openapi3.SchemaRef) map[string]*openapi3.Schema {
	properties := make(map[string]*openapi3.Schema)

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

// This is the function that loads the OpenAPI 3.0 spec file and
// returns the data for the template.
func getTemplateData(specNames []string) (TemplateData, error) {
	data := &TemplateData{
		SpecFlags: make(map[string]*SpecFlags),
	}
	for _, specName := range specNames {
		specParams, err := loadSpecs(searchSpecFile, specName)
		if err != nil {
			return *data, err
		}
		data.SpecFlags[specName] = getFlags(specParams)
	}
	return *data, nil
}

// getFlags returns the flags for the given spec.
func getFlags(params map[string]*openapi3.Schema) *SpecFlags {
	flags := &SpecFlags{
		Flags: make(map[string]*SpecFlag),
	}
	for name, param := range params {
		flags.Flags[name] = getFlag(param)
	}
	return flags
}

// GetGoType returns the Go type for the given OpenAPI 3.0 schema.
func GetGoType(param *openapi3.Schema) string {
	var SpecTypeGoType = map[string]string{
		"string":  "string",
		"integer": "int",
		"number":  "float64",
		"boolean": "bool",
	}
	if param.Type == "array" {
		return "[]" + GetGoType(param.Items.Value)
	}
	return SpecTypeGoType[param.Type]
}

// getFlag returns the flag for the given parameter.
func getFlag(param *openapi3.Schema) *SpecFlag {
	subType := ""
	if param.Type == "array" {
		subType = param.Items.Value.Type
	} else {
		subType = ""
	}

	var categories []string
	if param.ExtensionProps.Extensions["x-categories"] != nil {
		json.Unmarshal(param.ExtensionProps.Extensions["x-categories"].(json.RawMessage), &categories)
	}

	flag := &SpecFlag{
		Def:        param.Default,
		Type:       param.Type,
		GoType:     GetGoType(param),
		Usage:      getDescription(param),
		SubType:    subType,
		Categories: categories,
	}

	if param.OneOf != nil {
		for _, oneOf := range param.OneOf {
			flag.OneOf = append(flag.OneOf, oneOf.Value.Type)
		}
	}

	return flag
}

// getDescription returns the description for the given parameter.
func getDescription(param *openapi3.Schema) string {
	if param.Enum != nil {
		choices := make([]string, len(param.Enum))
		for i, e := range param.Enum {
			choices[i] = e.(string)
		}
		return fmt.Sprintf("%s One of: (%v).", param.Description, strings.Join(choices, ", "))
	}
	return param.Description
}
