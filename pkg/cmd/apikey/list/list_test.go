package list

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func runCommand(isTTY bool, cli string, key string) (*test.CmdOut, error) {
	io, _, stdout, stderr := iostreams.Test()
	io.SetStdoutTTY(isTTY)
	io.SetStdinTTY(isTTY)
	io.SetStderrTTY(isTTY)

	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/keys"),
		httpmock.JSONResponse(search.ListAPIKeysRes{
			Keys: []search.Key{
				{
					Value:                  "foo",
					Description:            "test",
					ACL:                    []string{"*"},
					Validity:               0,
					MaxHitsPerQuery:        0,
					MaxQueriesPerIPPerHour: 0,
					Referers:               []string{},
					CreatedAt:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		}),
	)

	client := search.NewClientWithConfig(search.Configuration{
		Requester: &r,
	})

	factory := &cmdutil.Factory{
		IOStreams: io,
		SearchClient: func() (*search.Client, error) {
			return client, nil
		},
	}

	cmd := NewListCmd(factory, nil)

	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	cmd.SetIn(&bytes.Buffer{})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(ioutil.Discard)

	_, err = cmd.ExecuteC()
	return &test.CmdOut{
		OutBuf: stdout,
		ErrBuf: stderr,
	}, err
}

func Test_runListCmd(t *testing.T) {
	tests := []struct {
		name    string
		isTTY   bool
		wantOut string
	}{
		{
			name:    "list",
			isTTY:   false,
			wantOut: "\ttest\t[*]\t[]\tNever expire\t0\t0\t[]\ta long while ago\n",
		},
		{
			name:  "list_tty",
			isTTY: true,
			wantOut: `KEY  DESCRIPTION  ACL  INDICES  VALIDITY  MAX H...  MAX Q...  REFERERS  CREAT...
     test         [*]  []       Never...  0         0         []        a lon...
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := runCommand(tt.isTTY, "", "")
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
