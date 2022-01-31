package cmdutil

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/iostreams"
)

func TestAddJSONFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantsExport *exportFormat
		wantsError  string
	}{
		{
			name:        "no JSON flag",
			args:        []string{},
			wantsExport: nil,
		},
		{
			name:        "cannot use --jq without --json",
			args:        []string{"--jq", ".number"},
			wantsExport: nil,
			wantsError:  "cannot use `--jq` without specifying `--json`",
		},
		{
			name:        "cannot use --template without --json",
			args:        []string{"--template", "{{.number}}"},
			wantsExport: nil,
			wantsError:  "cannot use `--template` without specifying `--json`",
		},
		{
			name: "with jq filter",
			args: []string{"--json", "number", "-q.number"},
			wantsExport: &exportFormat{
				filter:   ".number",
				template: "",
			},
		},
		{
			name: "with Go template",
			args: []string{"--json", "number", "-t", "{{.number}}"},
			wantsExport: &exportFormat{
				filter:   "",
				template: "{{.number}}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
			cmd.Flags().Bool("web", false, "")
			var exporter Exporter
			AddJSONFlags(cmd, &exporter)
			cmd.SetArgs(tt.args)
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)
			_, err := cmd.ExecuteC()
			if tt.wantsError == "" {
				require.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantsError)
				return
			}
			if tt.wantsExport == nil {
				assert.Nil(t, exporter)
			} else {
				assert.Equal(t, tt.wantsExport, exporter)
			}
		})
	}
}

func Test_exportFormat_Write(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name     string
		exporter exportFormat
		args     args
		wantW    string
		wantErr  bool
	}{
		{
			name:     "regular JSON output",
			exporter: exportFormat{},
			args: args{
				data: map[string]string{"name": "hubot"},
			},
			wantW:   "{\"name\":\"hubot\"}\n",
			wantErr: false,
		},
		{
			name:     "with jq filter",
			exporter: exportFormat{filter: ".name"},
			args: args{
				data: map[string]string{"name": "hubot"},
			},
			wantW:   "hubot\n",
			wantErr: false,
		},
		{
			name:     "with Go template",
			exporter: exportFormat{template: "{{.name}}"},
			args: args{
				data: map[string]string{"name": "hubot"},
			},
			wantW:   "hubot",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			io := &iostreams.IOStreams{
				Out: w,
			}
			if err := tt.exporter.Write(io, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("exportFormat.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("exportFormat.Write() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}
