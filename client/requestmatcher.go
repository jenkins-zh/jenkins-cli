package client

import (
	"fmt"
	"io/ioutil"
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
	withBody bool
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

	if request.verbose && !match {
		fmt.Printf("%s=?%s , %s=?%s, %s=?%s \n", request.request.Method, target.Method, request.request.URL.Path, target.URL.Path,
			request.request.URL.Opaque, target.URL.Opaque)
	}

	if request.matchOptions.withQuery {
		match = match && (request.request.URL.RawQuery == target.URL.RawQuery)
		if request.verbose && !match {
			fmt.Printf("query: %s=?%s  \n", request.request.URL.RawQuery, target.URL.RawQuery)
		}
	}

	if request.matchOptions.withBody {
		if request.request.Body != target.Body {
			if request.request.Body == nil || target.Body == nil {
				match = false
			} else {
				reqBody, _ := ioutil.ReadAll(request.request.Body)
				targetBody, _ := ioutil.ReadAll(target.Body)

				match = match && (string(reqBody) == string(targetBody))
				if request.verbose && !match {
					fmt.Printf("request body: %s, target body: %s \n", string(reqBody), string(targetBody))
				}
			}
		}
	}

	return match
}

// String returns the text of current object
func (*RequestMatcher) String() string {
	return "request matcher"
}
