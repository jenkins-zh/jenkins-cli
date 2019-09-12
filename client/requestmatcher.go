package client

import (
	"fmt"
	"net/http"
)

// RequestMatcher to match the http request
type RequestMatcher struct {
	request *http.Request
	verbose bool

	matchOptions matchOptions
}

type matchOptions struct {
	withQuery bool
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

// Matches returns a matcher with given function
func (request *RequestMatcher) Matches(x interface{}) bool {
	target := x.(*http.Request)

	if request.verbose {
		fmt.Printf("%s=?%s , %s=?%s, %s=?%s \n", request.request.Method, target.Method, request.request.URL.Path, target.URL.Path,
			request.request.URL.Opaque, target.URL.Opaque)
	}

	match := request.request.Method == target.Method && (request.request.URL.Path == target.URL.Path ||
		request.request.URL.Path == target.URL.Opaque) //gitlab sdk did not set request path correctly

	if request.matchOptions.withQuery {
		if request.verbose {
			fmt.Printf("%s=?%s  \n", request.request.URL.RawQuery, target.URL.RawQuery)
		}
		match = match && (request.request.URL.RawQuery == target.URL.RawQuery)
	}

	return match
}

// String returns the text of current object
func (*RequestMatcher) String() string {
	return "request matcher"
}
