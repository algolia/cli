package browse

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func runCommand(http *httpmock.Registry, hits []map[string]interface{}, cli string) (*test.CmdOut, error) {
	io, _, stdout, stderr := iostreams.Test()

	client := search.NewClientWithConfig(search.Configuration{
		Requester: http,
	})

	factory := &cmdutil.Factory{
		IOStreams: io,
		SearchClient: func() (*search.Client, error) {
			return client, nil
		},
	}

	cmd := NewBrowseCmd(factory)

	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	cmd.SetIn(&bytes.Buffer{})
	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)

	_, err = cmd.ExecuteC()
	return &test.CmdOut{
		OutBuf: stdout,
		ErrBuf: stderr,
	}, err
}

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		hits    []map[string]interface{}
		wantOut string
	}{
		{
			name:    "single object",
			cli:     "foo",
			hits:    []map[string]interface{}{{"objectID": "foo"}},
			wantOut: "{\"objectID\":\"foo\"}\n",
		},
		{
			name:    "multiple objects",
			cli:     "foo",
			hits:    []map[string]interface{}{{"objectID": "foo"}, {"objectID": "bar"}},
			wantOut: "{\"objectID\":\"foo\"}\n{\"objectID\":\"bar\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("POST", "1/indexes/foo/browse"), httpmock.JSONResponse(search.QueryRes{
				Hits: tt.hits,
			}))
			defer r.Verify(t)

			out, err := runCommand(&r, tt.hits, tt.cli)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
