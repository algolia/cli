package generate

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

func TestNewGenerateCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		errMsg    string
		wantsOpts GenerateOptions
	}{
		{
			name:     "missing required flags",
			cli:      "",
			tty:      true,
			wantsErr: true,
			errMsg:   "required flag(s) \"datasource\", \"prompt\", \"query\" not set",
		},
		{
			name:     "missing datasource",
			cli:      "--query \"hello\" --prompt prompt-123",
			tty:      true,
			wantsErr: true,
			errMsg:   "required flag(s) \"datasource\" not set",
		},
		{
			name:     "missing prompt",
			cli:      "--query \"hello\" --datasource ds-123",
			tty:      true,
			wantsErr: true,
			errMsg:   "required flag(s) \"prompt\" not set",
		},
		{
			name:     "missing query",
			cli:      "--datasource ds-123 --prompt prompt-123",
			tty:      true,
			wantsErr: true,
			errMsg:   "required flag(s) \"query\" not set",
		},
		{
			name:     "valid with minimum flags",
			cli:      "--query \"hello\" --datasource ds-123 --prompt prompt-123",
			tty:      true,
			wantsErr: false,
			wantsOpts: GenerateOptions{
				Query:        "hello",
				DataSourceID: "ds-123",
				PromptID:     "prompt-123",
				LogRegion:    "us",
				NbHits:       4,
			},
		},
		{
			name:     "valid with all flags",
			cli:      "--query \"hello\" --datasource ds-123 --prompt prompt-123 --region de --id resp-123 --hits 10 --filters \"brand:apple\" --object-ids id1,id2 --attributes name,price --conversation-id conv-123 --save --use-cache",
			tty:      true,
			wantsErr: false,
			wantsOpts: GenerateOptions{
				Query:                "hello",
				DataSourceID:         "ds-123",
				PromptID:             "prompt-123",
				LogRegion:            "de",
				ObjectID:             "resp-123",
				NbHits:               10,
				AdditionalFilters:    "brand:apple",
				WithObjectIDs:        []string{"id1", "id2"},
				AttributesToRetrieve: []string{"name", "price"},
				ConversationID:       "conv-123",
				Save:                 true,
				UseCache:             true,
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

			var opts *GenerateOptions
			cmd := NewGenerateCmd(f, func(o *GenerateOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.Query, opts.Query)
			assert.Equal(t, tt.wantsOpts.DataSourceID, opts.DataSourceID)
			assert.Equal(t, tt.wantsOpts.PromptID, opts.PromptID)
			assert.Equal(t, tt.wantsOpts.LogRegion, opts.LogRegion)
			assert.Equal(t, tt.wantsOpts.ObjectID, opts.ObjectID)
			assert.Equal(t, tt.wantsOpts.NbHits, opts.NbHits)
			assert.Equal(t, tt.wantsOpts.AdditionalFilters, opts.AdditionalFilters)
			assert.Equal(t, tt.wantsOpts.WithObjectIDs, opts.WithObjectIDs)
			assert.Equal(t, tt.wantsOpts.AttributesToRetrieve, opts.AttributesToRetrieve)
			assert.Equal(t, tt.wantsOpts.ConversationID, opts.ConversationID)
			assert.Equal(t, tt.wantsOpts.Save, opts.Save)
			assert.Equal(t, tt.wantsOpts.UseCache, opts.UseCache)
		})
	}
}

func Test_runGenerateCmd(t *testing.T) {
	tests := []struct {
		name      string
		opts      GenerateOptions
		isTTY     bool
		httpStubs func(*httpmock.Registry)
		wantOut   string
	}{
		{
			name: "generates response (tty)",
			opts: GenerateOptions{
				Query:        "hello",
				DataSourceID: "ds-123",
				PromptID:     "prompt-123",
				LogRegion:    "us",
				NbHits:       4,
			},
			isTTY: true,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "generate/response"),
					httpmock.JSONResponse(genai.GenerateResponseOutput{
						ObjectID: "resp-123",
						Response: "Hello! How can I help you today?",
					}),
				)
			},
			wantOut: "✓ Response generated with ID: resp-123\n\nHello! How can I help you today?\n",
		},
		{
			name: "generates response (non-tty)",
			opts: GenerateOptions{
				Query:        "hello",
				DataSourceID: "ds-123",
				PromptID:     "prompt-123",
				LogRegion:    "us",
				NbHits:       4,
			},
			isTTY: false,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "generate/response"),
					httpmock.JSONResponse(genai.GenerateResponseOutput{
						ObjectID: "resp-123",
						Response: "Hello! How can I help you today?",
					}),
				)
			},
			wantOut: "Hello! How can I help you today?",
		},
		{
			name: "generates response with all options",
			opts: GenerateOptions{
				Query:                "hello",
				DataSourceID:         "ds-123",
				PromptID:             "prompt-123",
				LogRegion:            "de",
				ObjectID:             "resp-123",
				NbHits:               10,
				AdditionalFilters:    "brand:apple",
				WithObjectIDs:        []string{"id1", "id2"},
				AttributesToRetrieve: []string{"name", "price"},
				ConversationID:       "conv-123",
				Save:                 true,
				UseCache:             true,
			},
			isTTY: true,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "generate/response"),
					httpmock.JSONResponse(genai.GenerateResponseOutput{
						ObjectID: "resp-123",
						Response: "Here is information about the apple products you requested.",
					}),
				)
			},
			wantOut: "✓ Response generated with ID: resp-123\n\nHere is information about the apple products you requested.\n",
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

			err := runGenerateCmd(&tt.opts)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
