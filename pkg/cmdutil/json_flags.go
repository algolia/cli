package cmdutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/export"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/jsoncolor"
)

type JSONFlagError struct {
	error
}

func AddJSONFlags(cmd *cobra.Command, exportTarget *Exporter, defaulValue bool) {
	f := cmd.Flags()
	f.Bool("json", defaulValue, "Output JSON")
	f.StringP("jq", "q", "", "Filter JSON output using a jq `expression`")
	f.StringP("template", "t", "", "Format JSON output using a Go template")

	oldPreRun := cmd.PreRunE
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if oldPreRun != nil {
			if err := oldPreRun(c, args); err != nil {
				return err
			}
		}
		if export, err := checkJSONFlags(c, defaulValue); err == nil {
			if export == nil {
				*exportTarget = nil
			} else {
				*exportTarget = export
			}
		} else {
			return err
		}
		return nil
	}
}

func checkJSONFlags(cmd *cobra.Command, defaulValue bool) (*exportFormat, error) {
	f := cmd.Flags()
	jsonFlag := f.Lookup("json")
	jqFlag := f.Lookup("jq")
	tplFlag := f.Lookup("template")

	if jsonFlag.Changed || defaulValue {
		return &exportFormat{
			filter:   jqFlag.Value.String(),
			template: tplFlag.Value.String(),
		}, nil
	} else if jqFlag.Changed {
		return nil, errors.New("cannot use `--jq` without specifying `--json`")
	} else if tplFlag.Changed {
		return nil, errors.New("cannot use `--template` without specifying `--json`")
	}
	return nil, nil
}

type Exporter interface {
	Write(io *iostreams.IOStreams, data interface{}) error
}

type exportFormat struct {
	filter   string
	template string
}

// Write serializes data into JSON output written to w. If the object passed as data implements exportable,
// or if data is a map or slice of exportable object, ExportData() will be called on each object to obtain
// raw data for serialization.
func (e *exportFormat) Write(ios *iostreams.IOStreams, data interface{}) error {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	w := ios.Out
	if e.filter != "" {
		return export.FilterJSON(w, &buf, e.filter)
	} else if e.template != "" {
		return export.ExecuteTemplate(ios, &buf, e.template)
	} else if ios.ColorEnabled() {
		return jsoncolor.Write(w, &buf, "  ")
	}

	_, err := io.Copy(w, &buf)
	return err
}

var sliceOfEmptyInterface []interface{}
var emptyInterfaceType = reflect.TypeOf(sliceOfEmptyInterface).Elem()
