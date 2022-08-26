package httpmock

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Matcher func(req *http.Request) bool
type Responder func(req *http.Request) (*http.Response, error)

type Stub struct {
	matched   bool
	Matcher   Matcher
	Responder Responder
}

func REST(method, p string) Matcher {
	return func(req *http.Request) bool {
		if !strings.EqualFold(req.Method, method) {
			return false
		}
		if req.URL.Path != "/"+p {
			return false
		}
		return true
	}
}

func StringResponse(body string) Responder {
	return func(req *http.Request) (*http.Response, error) {
		return httpResponse(200, req, bytes.NewBufferString(body)), nil
	}
}

func JSONResponse(body interface{}) Responder {
	return func(req *http.Request) (*http.Response, error) {
		b, _ := json.Marshal(body)
		return httpResponse(200, req, bytes.NewBuffer(b)), nil
	}
}

func httpResponse(status int, req *http.Request, body io.Reader) *http.Response {
	return &http.Response{
		StatusCode: status,
		Request:    req,
		Body:       ioutil.NopCloser(body),
		Header:     http.Header{},
	}
}
