package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

// RequestMatcher to match the http request
type RequestMatcher struct {
	request *http.Request
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

	match := request.request.Method == target.Method && (request.request.URL.Path == target.URL.Path ||
		request.request.URL.Path == target.URL.Opaque)

	if !match {
		match = reflect.DeepEqual(request.request.Header, target.Header)
	}

	if request.matchOptions.withQuery && !match {
		match = request.request.URL.RawQuery == target.URL.RawQuery
	}

	reqBody, _ := getStrFromReader(request.request.Body)
	targetBody, _ := getStrFromReader(target.Body)
	if request.matchOptions.withBody && !match {
		match = reqBody == targetBody
	}

	if !match {
		fmt.Printf("%s=?%s , %s=?%s, %s=?%s \n", request.request.Method, target.Method, request.request.URL.Path, target.URL.Path,
			request.request.URL.Opaque, target.URL.Opaque)
		if request.matchOptions.withQuery {
			fmt.Printf("query: %s=?%s \n", request.request.URL.RawQuery, target.URL.RawQuery)
		} else if request.matchOptions.withBody {
			fmt.Printf("request body: %s, target body: %s \n", reqBody, targetBody)
		}
	}

	return match
}

func getStrFromReader(reader io.ReadCloser) (text string, err error) {
	if reader == nil {
		return
	}

	if data, err := ioutil.ReadAll(reader); err == nil {
		text = string(data)
	}
	return
}

// String returns the text of current object
func (*RequestMatcher) String() string {
	return "request matcher"
}
