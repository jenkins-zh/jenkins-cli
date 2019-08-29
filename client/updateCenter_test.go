package client

import (
	"bytes"
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
})
