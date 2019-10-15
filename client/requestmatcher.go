package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// RequestMatcher to match the http request
type RequestMatcher struct {
	request *http.Request
	target *http.Request
	verbose bool

	matchOptions matchOptions
}

type matchOptions struct {
	withQuery bool
	withBody  bool
}

// NewRequestMatcher create a request matcher will match request method and request path
func NewRequestMatcher(request *http.Request) *RequestMatcher {
	return &RequestMatcher{request: request}
}

// NewVerboseRequestMatcher create a verbose request matcher will match request method and request path
func NewVerboseRequestMatcher(request *http.Request) *RequestMatcher {
	return &RequestMatcher{request: request, verbose: true}
}

// WithQuery returns a matcher with query
func (request *RequestMatcher) WithQuery() *RequestMatcher {
	request.matchOptions.withQuery = true
	return request
}

// WithBody returns a matcher with body
func (request *RequestMatcher) WithBody() *RequestMatcher {
	request.matchOptions.withBody = true
	return request
}

// Matches returns a matcher with given function
func (request *RequestMatcher) Matches(x interface{}) bool {
	target := x.(*http.Request)
	request.target = target

	match := request.request.Method == target.Method && (request.request.URL.Path == target.URL.Path ||
		request.request.URL.Path == target.URL.Opaque)

	if match {
		match = matchHeader(request.request.Header, request.target.Header)
	}

	if request.matchOptions.withQuery && match {
		match = request.request.URL.RawQuery == target.URL.RawQuery
	}

	if request.matchOptions.withBody && match {
		reqBody, _ := getStrFromReader(request.request)
		targetBody, _ := getStrFromReader(target)
		match = reqBody == targetBody
	}

	return match
}

func matchHeader(left, right http.Header) bool {
	if len(left) != len(right) {
		return false
	}

	for k, v := range left {
		if k == "Content-Type" { // it's hard to compare file upload cases
			continue
		}
		if tv, ok := right[k]; !ok || !reflect.DeepEqual(v, tv) {
			return false
		}
	}
	return true
}

func getStrFromReader(request *http.Request) (text string, err error) {
	reader := request.Body
	if reader == nil {
		return
	}

	if data, err := ioutil.ReadAll(reader); err == nil {
		text = string(data)

		// it could be read twice
		payload := strings.NewReader(text)
		request.Body = ioutil.NopCloser(payload)
	}
	return
}

// String returns the text of current object
func (request *RequestMatcher) String() string {
	target := request.target
	return fmt.Sprintf("%v", target)
}
