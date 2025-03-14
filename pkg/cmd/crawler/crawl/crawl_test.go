package crawl

import (
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewCrawlCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts CrawlOptions
	}{
		{
			name:     "single URL, no save, without tty",
			cli:      "my-crawler --urls http://example.com",
			tty:      false,
			wantsErr: false,
			wantsOpts: CrawlOptions{
				ID:   "my-crawler",
				URLs: []string{"http://example.com"},
				Save: false,
			},
		},
		{
			name:     "single URL, no save, with tty",
			cli:      "my-crawler --urls http://example.com",
			tty:      true,
			wantsErr: false,
			wantsOpts: CrawlOptions{
				ID:   "my-crawler",
				URLs: []string{"http://example.com"},
				Save: false,
			},
		},
		{
			name:     "single URL, save, without tty",
			cli:      "my-crawler --urls http://example.com --save",
			tty:      false,
			wantsErr: false,
			wantsOpts: CrawlOptions{
				ID:   "my-crawler",
				URLs: []string{"http://example.com"},
				Save: true,
			},
		},
		{
			name:     "multiple URLs, no save, without tty",
			cli:      "my-crawler --urls http://example.com,http://example.com/doc",
			tty:      false,
			wantsErr: false,
			wantsOpts: CrawlOptions{
				ID:   "my-crawler",
				URLs: []string{"http://example.com", "http://example.com/doc"},
				Save: false,
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

			var opts *CrawlOptions
			cmd := NewCrawlCmd(f, func(o *CrawlOptions) error {
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

			assert.Equal(t, tt.wantsOpts.ID, opts.ID)
			assert.Equal(t, tt.wantsOpts.URLs, opts.URLs)
		})
	}
}

func Test_runCrawlCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		id      string
		urls    []string
		isTTY   bool
		wantErr string
		wantOut string
	}{
		{
			name:    "no TTY",
			cli:     "my-crawler --urls http://example.com",
			id:      "my-crawler",
			urls:    []string{"http://example.com"},
			isTTY:   false,
			wantOut: "",
		},
		{
			name:    "TTY",
			cli:     "my-crawler --urls http://example.com",
			id:      "my-crawler",
			urls:    []string{"http://example.com"},
			isTTY:   true,
			wantOut: "✓ Successfully requested crawl for 1 URL on crawler my-crawler\n",
		},
		{
			name:    "TTY, multiple URLs",
			cli:     "my-crawler --urls http://example.com,http://example.com/doc",
			id:      "my-crawler",
			urls:    []string{"http://example.com", "http://example.com/doc"},
			isTTY:   true,
			wantOut: "✓ Successfully requested crawl for 2 URLs on crawler my-crawler\n",
		},
		{
			name:    "TTY, error (message+code)",
			cli:     "my-crawler --urls http://example.com",
			id:      "my-crawler",
			urls:    []string{"http://example.com"},
			isTTY:   true,
			wantErr: "X Crawler API error: [not-found] Crawler not-found not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.wantErr == "" {
				r.Register(
					httpmock.REST("POST", "api/1/crawlers/"+tt.id+"/urls/crawl"),
					httpmock.JSONResponse(crawler.TaskIDResponse{TaskID: "taskID"}),
				)
			} else {
				r.Register(httpmock.REST("POST", "api/1/crawlers/"+tt.id+"/urls/crawl"), httpmock.ErrorResponseWithBody(crawler.ErrResponse{Err: crawler.Err{Code: "not-found", Message: "Crawler not-found not found"}}))
			}
			defer r.Verify(t)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewCrawlCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
