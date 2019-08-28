package util

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

var _ = Describe("http test", func() {
	var (
		ctrl           *gomock.Controller
		roundTripper   *mhttp.MockRoundTripper
		downloader     HTTPDownloader
		targetFilePath string
		responseBody   string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		targetFilePath = "test.log"
		downloader = HTTPDownloader{
			TargetFilePath: targetFilePath,
			RoundTripper:   roundTripper,
		}
		responseBody = "fake body"
	})

	AfterEach(func() {
		os.Remove(targetFilePath)
		ctrl.Finish()
	})

	Context("DownloadFile", func() {
		It("no progress indication", func() {
			request, _ := http.NewRequest("GET", "", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)
			err := downloader.DownloadFile()
			Expect(err).To(BeNil())

			_, err = os.Stat(targetFilePath)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(targetFilePath)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(responseBody))
			return
		})

		It("with BasicAuth", func() {
			downloader = HTTPDownloader{
				TargetFilePath: targetFilePath,
				RoundTripper:   roundTripper,
				UserName:       "UserName",
				Password:       "Password",
			}

			request, _ := http.NewRequest("GET", "", nil)
			request.SetBasicAuth(downloader.UserName, downloader.Password)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Header:     http.Header{},
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)
			err := downloader.DownloadFile()
			Expect(err).To(BeNil())

			_, err = os.Stat(targetFilePath)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(targetFilePath)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(responseBody))
			return
		})
	})
})
