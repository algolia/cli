package importRecords

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func runCommand(http *httpmock.Registry, in string, cli string) (*test.CmdOut, error) {
	io, stdin, stdout, stderr := iostreams.Test()
	stdin.WriteString(in)

	client := search.NewClientWithConfig(search.Configuration{
		Requester: http,
	})

	factory := &cmdutil.Factory{
		IOStreams: io,
		SearchClient: func() (search.ClientInterface, error) {
			return client, nil
		},
	}

	cmd := NewImportCmd(factory)

	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	if stdin != nil {
		cmd.SetIn(stdin)
	} else {
		cmd.SetIn(&bytes.Buffer{})
	}
	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)

	_, err = cmd.ExecuteC()
	return &test.CmdOut{
		OutBuf: stdout,
		ErrBuf: stderr,
	}, err
}

func Test_runExportCmd(t *testing.T) {

	tmpFile := filepath.Join(t.TempDir(), "objects.json")
	err := ioutil.WriteFile(tmpFile, []byte("{\"objectID\":\"foo\"}"), 0600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cli     string
		stdin   string
		wantOut string
	}{
		{
			name:    "from stdin",
			cli:     "foo -F -",
			stdin:   `{"objectID": "foo"}`,
			wantOut: "",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("foo -F '%s'", tmpFile),
			stdin:   `{"objectID": "foo"}`,
			wantOut: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("POST", "1/indexes/foo/batch"), httpmock.JSONResponse(search.BatchRes{}))
			defer r.Verify(t)

			out, err := runCommand(&r, tt.stdin, tt.cli)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
