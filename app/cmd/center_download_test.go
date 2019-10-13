package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("center download command", func() {
	var (
		ctrl           *gomock.Controller
		roundTripper   *mhttp.MockRoundTripper
		targetFilePath string

		ltsResponseBody   string
		weeklyResponseBody string

		err error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		centerDownloadOption.RoundTripper = roundTripper
		targetFilePath = "jenkins.war"

		ltsResponseBody = "lts"
		weeklyResponseBody = "weekly"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		err = os.Remove(targetFilePath)
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("download the lts Jenkins", func() {
			request, _ := http.NewRequest("GET", "http://mirrors.jenkins.io/war-stable/latest/jenkins.war", nil)
			response := &http.Response{
				StatusCode: 200,
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(ltsResponseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			rootCmd.SetArgs([]string{"center", "download", "--progress=false"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			_, err = os.Stat(targetFilePath)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(targetFilePath)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(ltsResponseBody))
		}, 1)


		It("download the weekly Jenkins", func() {
			request, _ := http.NewRequest("GET", "http://mirrors.jenkins.io/war/latest/jenkins.war", nil)
			response := &http.Response{
				StatusCode: 200,
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(weeklyResponseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			rootCmd.SetArgs([]string{"center", "download", "--lts=false", "--progress=false"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			_, err = os.Stat(targetFilePath)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(targetFilePath)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(weeklyResponseBody))
		}, 1)
	})
})
