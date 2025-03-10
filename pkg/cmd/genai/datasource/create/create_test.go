package create

import (
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewCreateCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		errMsg    string
		wantsOpts CreateOptions
	}{
		{
			name:     "required source flag",
			cli:      "MyDataSource",
			tty:      true,
			wantsErr: true,
			errMsg:   "required flag(s) \"source\" not set",
		},
		{
			name:     "valid with minimum flags",
			cli:      "MyDataSource --source products",
			tty:      true,
			wantsErr: false,
			wantsOpts: CreateOptions{
				Name:   "MyDataSource",
				Source: "products",
			},
		},
		{
			name:     "valid with all flags",
			cli:      "MyDataSource --source products --filters \"category:phones\" --id custom-id",
			tty:      true,
			wantsErr: false,
			wantsOpts: CreateOptions{
				Name:     "MyDataSource",
				Source:   "products",
				Filters:  "category:phones",
				ObjectID: "custom-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()
			if tt.tty {
				io.SetStdinTTY(tt.tty)
				io.SetStdoutTTY(tt.tty)
			}

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			var opts *CreateOptions
			cmd := NewCreateCmd(f, func(o *CreateOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.Name, opts.Name)
			assert.Equal(t, tt.wantsOpts.Source, opts.Source)
			assert.Equal(t, tt.wantsOpts.Filters, opts.Filters)
			assert.Equal(t, tt.wantsOpts.ObjectID, opts.ObjectID)
		})
	}
}

func Test_runCreateCmd(t *testing.T) {
	tests := []struct {
		name      string
		opts      CreateOptions
		isTTY     bool
		httpStubs func(*httpmock.Registry)
		wantOut   string
	}{
		{
			name: "creates with minimum options (tty)",
			opts: CreateOptions{
				Name:   "MyDataSource",
				Source: "products",
			},
			isTTY: true,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "create/data_source"),
					httpmock.JSONResponse(genai.DataSourceResponse{
						ObjectID: "ds-123",
					}),
				)
			},
			wantOut: "✓ Data source MyDataSource created with ID: ds-123\n",
		},
		{
			name: "creates with all options (tty)",
			opts: CreateOptions{
				Name:     "MyDataSource",
				Source:   "products",
				Filters:  "category:phones",
				ObjectID: "custom-id",
			},
			isTTY: true,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "create/data_source"),
					httpmock.JSONResponse(genai.DataSourceResponse{
						ObjectID: "custom-id",
					}),
				)
			},
			wantOut: "✓ Data source MyDataSource created with ID: custom-id\n",
		},
		{
			name: "creates (non-tty)",
			opts: CreateOptions{
				Name:   "MyDataSource",
				Source: "products",
			},
			isTTY: false,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "create/data_source"),
					httpmock.JSONResponse(genai.DataSourceResponse{
						ObjectID: "ds-123",
					}),
				)
			},
			wantOut: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.httpStubs != nil {
				tt.httpStubs(&r)
			}
			defer r.Verify(t)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			tt.opts.IO = f.IOStreams
			tt.opts.Config = f.Config
			tt.opts.GenAIClient = f.GenAIClient

			err := runCreateCmd(&tt.opts)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
