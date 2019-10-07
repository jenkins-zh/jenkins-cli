package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("common test", func() {
	var (
		ctrl         *gomock.Controller
		jenkinsCore  JenkinsCore
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		jenkinsCore = JenkinsCore{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		jenkinsCore.RoundTripper = roundTripper
		jenkinsCore.URL = "http://localhost"
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Request", func() {
		var (
			method  string
			api     string
			headers map[string]string
			payload io.Reader
		)

		BeforeEach(func() {
			method = "GET"
			api = "/fake"
		})

		It("normal case for get request", func() {
			request, _ := http.NewRequest(method, fmt.Sprintf("%s%s", jenkinsCore.URL, api), payload)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			statusCode, data, err := jenkinsCore.Request(method, api, headers, payload)
			Expect(err).To(BeNil())
			Expect(statusCode).To(Equal(200))
			Expect(string(data)).To(Equal(""))
		})

		It("normal case for post request", func() {
			method = "POST"
			request, _ := http.NewRequest(method, fmt.Sprintf("%s%s", jenkinsCore.URL, api), payload)
			request.Header.Add("CrumbRequestField", "Crumb")
			request.Header.Add("Fake", "fake")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), payload)
			responseCrumb := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{"crumbRequestField":"CrumbRequestField","crumb":"Crumb"}
				`)),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			headers = make(map[string]string, 1)
			headers["fake"] = "fake"
			statusCode, data, err := jenkinsCore.Request(method, api, headers, payload)
			Expect(err).To(BeNil())
			Expect(statusCode).To(Equal(200))
			Expect(string(data)).To(Equal(""))
		})
	})

	Context("GetCrumb", func() {
		It("without crumb setting", func() {
			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 404,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			crumb, err := jenkinsCore.GetCrumb()
			Expect(crumb).To(BeNil())
			Expect(err).To(BeNil())
		})

		It("with crumb setting", func() {
			RequestCrumb(roundTripper, jenkinsCore.URL)

			crumb, err := jenkinsCore.GetCrumb()
			Expect(err).To(BeNil())
			Expect(crumb).NotTo(BeNil())
			Expect(crumb.CrumbRequestField).To(Equal("CrumbRequestField"))
			Expect(crumb.Crumb).To(Equal("Crumb"))
		})

		It("with error from server", func() {
			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)

			_, err := jenkinsCore.GetCrumb()
			Expect(err).To(HaveOccurred())
		})

		It("with 404 error from server", func() {
			err := jenkinsCore.ErrorHandle(404, []byte{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Not found resources"))
		})

		It("with 403 error from server", func() {
			err := jenkinsCore.ErrorHandle(403, []byte{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("The current user no permission"))
		})

		It("with CrumbHandle error from server", func() {
			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", jenkinsCore.URL, "/crumbIssuer/api/json"), nil)
			responseCrumb := &http.Response{
				StatusCode: 500,
				Proto:      "HTTP/1.1",
				Request:    requestCrumb,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(requestCrumb).Return(responseCrumb, nil)
			err := jenkinsCore.CrumbHandle(requestCrumb)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unexpected status code: 500"))
		})

		It("test GetClient", func() {
			jenkinsCore.RoundTripper = nil
			jenkinsCore.Proxy = "kljasdsll"
			jenkinsCore.ProxyAuth = "kljaslkdjkslad"
			jclient := jenkinsCore.GetClient()
			Expect(jclient).NotTo(BeNil())
		})
	})
})
