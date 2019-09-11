package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("center download command", func() {
	var (
		ctrl           *gomock.Controller
		roundTripper   *mhttp.MockRoundTripper
		targetFilePath string
		responseBody   string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		targetFilePath = "jenkins.war"
		responseBody = "fake response"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(targetFilePath)
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should success", func() {
			centerDownloadOption.RoundTripper = roundTripper
			centerDownloadOption.LTS = false

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

			rootCmd.SetArgs([]string{"center", "download"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			_, err = os.Stat(targetFilePath)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(targetFilePath)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(responseBody))
		})
	})
})
