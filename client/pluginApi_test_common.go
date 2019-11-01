package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

// PrepareShowTrend only for test
func PrepareShowTrend(roundTripper *mhttp.MockRoundTripper, keyword string) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://plugins.jenkins.io/api/plugin/%s", keyword), nil)
	response := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Request:    request,
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{"stats": {"installations":[{"total":1512},{"total":3472},{"total":4385},{"total":3981}]}}
		`)),
	}
	roundTripper.EXPECT().
		RoundTrip(request).Return(response, nil)
}
