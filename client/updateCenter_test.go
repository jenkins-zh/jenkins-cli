package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("update center test", func() {
	var (
		ctrl         *gomock.Controller
		manager      *UpdateCenterManager
		roundTripper *mhttp.MockRoundTripper
		responseBody string
		donwloadFile string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		manager = &UpdateCenterManager{}
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		responseBody = "fake response"
		donwloadFile = "downloadfile.log"
	})

	AfterEach(func() {
		os.Remove(donwloadFile)
		ctrl.Finish()
	})

	Context("DownloadJenkins", func() {
		It("should success with basic cases", func() {
			manager.RoundTripper = roundTripper

			request, _ := http.NewRequest("GET", "http://mirrors.jenkins.io/war/latest/jenkins.war", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)
			err := manager.DownloadJenkins(false, donwloadFile)
			Expect(err).To(BeNil())

			_, err = os.Stat(donwloadFile)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(donwloadFile)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(responseBody))
		})
	})

	Context("Upgrade", func() {
		It("basic cases", func() {
			manager.RoundTripper = roundTripper
			manager.URL = ""

			requestCrumb, _ := http.NewRequest("GET", fmt.Sprintf("%s/crumbIssuer/api/json", ""), nil)
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

			request, _ := http.NewRequest("POST", "/updateCenter/upgrade", nil)
			request.Header.Add("CrumbRequestField", "Crumb")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			err := manager.Upgrade()
			Expect(err).To(BeNil())
		})
	})

	Context("Status", func() {
		It("should success", func() {
			manager.RoundTripper = roundTripper
			manager.URL = ""

			request, _ := http.NewRequest("GET", "/updateCenter/api/json?pretty=false&depth=1", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
			{"RestartRequiredForCompletion": true}
			`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			status, err := manager.Status()
			Expect(err).To(BeNil())
			Expect(status).NotTo(BeNil())
			Expect(status.RestartRequiredForCompletion).Should(BeTrue())
		})
	})
})
