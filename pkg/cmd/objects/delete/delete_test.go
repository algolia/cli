package delete

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewDeleteCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts DeleteOptions
	}{
		{
			name:     "no --confirm without tty",
			cli:      "foo --object-ids 1",
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "--confirm without tty",
			cli:      "foo --object-ids 1 --confirm",
			tty:      true,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Index:     "foo",
				ObjectIDs: []string{
					"1",
				},
			},
		},
		{
			name:     "no --confirm with tty",
			cli:      "foo --object-ids 1",
			tty:      true,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: true,
				Index:     "foo",
				ObjectIDs: []string{
					"1",
				},
			},
		},
		{
			name:     "multiple --object-ids",
			cli:      "foo --object-ids 1,2 --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Index:     "foo",
				ObjectIDs: []string{
					"1",
					"2",
				},
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

			var opts *DeleteOptions
			cmd := NewDeleteCmd(f, func(o *DeleteOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.Index, opts.Index)
			assert.Equal(t, tt.wantsOpts.ObjectIDs, opts.ObjectIDs)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runDeleteCmd(t *testing.T) {
	tests := []struct {
		name             string
		cli              string
		indice           string
		objectIDs        []string
		nbHits           int32
		exhaustiveNbHits bool
		isTTY            bool
		wantOut          string
	}{
		{
			name:   "single object-id, no TTY",
			cli:    "foo --object-ids 1 --confirm",
			indice: "foo",
			objectIDs: []string{
				"1",
			},
			isTTY:   false,
			wantOut: "",
		},
		{
			name:   "single object-id, TTY",
			cli:    "foo --object-ids 1 --confirm",
			indice: "foo",
			objectIDs: []string{
				"1",
			},
			isTTY:   true,
			wantOut: "✓ Successfully deleted exactly 1 object from foo\n",
		},
		{
			name:   "multiple object-ids, TTY",
			cli:    "foo --object-ids 1,2 --confirm",
			indice: "foo",
			objectIDs: []string{
				"1",
				"2",
			},
			isTTY:   true,
			wantOut: "✓ Successfully deleted exactly 2 objects from foo\n",
		},
		{
			name:      "filters, TTY",
			cli:       "foo --filters 'foo:bar' --confirm",
			indice:    "foo",
			objectIDs: []string{},
			nbHits:    2,
			isTTY:     true,
			wantOut:   "✓ Successfully deleted approximately 2 objects from foo\n",
		},
		{
			name:             "filters, TTY",
			cli:              "foo --filters 'foo:bar' --confirm",
			indice:           "foo",
			objectIDs:        []string{},
			nbHits:           2,
			exhaustiveNbHits: true,
			isTTY:            true,
			wantOut:          "✓ Successfully deleted exactly 2 objects from foo\n",
		},
		{
			name:      "filters and object-ids, TTY",
			cli:       "foo --filters 'foo:bar' --object-ids 1,2 --confirm",
			indice:    "foo",
			objectIDs: []string{"1", "2"},
			nbHits:    2,
			isTTY:     true,
			wantOut:   "✓ Successfully deleted approximately 4 objects from foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			for _, id := range tt.objectIDs {
				// Checking that the object exists
				r.Register(
					httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/%s", tt.indice, id)),
					httpmock.JSONResponse(search.GetObjectsResponse{}),
				)
			}
			if tt.nbHits > 0 {
				// Searching for the objects to delete (if filters are used)
				r.Register(
					httpmock.REST("POST", fmt.Sprintf("1/indexes/%s/query", tt.indice)),
					httpmock.JSONResponse(search.BrowseResponse{
						NbHits:           &tt.nbHits,
						ExhaustiveNbHits: &tt.exhaustiveNbHits,
					}),
				)
				// Deleting the objects
				r.Register(
					httpmock.REST("POST", fmt.Sprintf("1/indexes/%s/deleteByQuery", tt.indice)),
					httpmock.JSONResponse(search.DeletedAtResponse{}),
				)
			}
			r.Register(
				httpmock.REST("POST", fmt.Sprintf("1/indexes/%s/batch", tt.indice)),
				httpmock.JSONResponse(search.BatchResponse{}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewDeleteCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
