package util

import (
	"bytes"
	"fmt"
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

		It("with error request", func() {
			downloader = HTTPDownloader{
				URL: "fake url",
			}
			err := downloader.DownloadFile()
			Expect(err).To(HaveOccurred())
		})

		It("with error response", func() {
			downloader = HTTPDownloader{
				RoundTripper: roundTripper,
			}

			request, _ := http.NewRequest("GET", "", nil)
			response := &http.Response{}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, fmt.Errorf("fake error"))
			err := downloader.DownloadFile()
			Expect(err).To(HaveOccurred())
		})

		It("status code isn't 200", func() {
			downloader = HTTPDownloader{
				RoundTripper: roundTripper,
				Debug:        true,
			}

			request, _ := http.NewRequest("GET", "", nil)
			response := &http.Response{
				StatusCode: 400,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)
			err := downloader.DownloadFile()
			Expect(err).To(HaveOccurred())

			const debugFile = "debug-download.html"

			_, err = os.Stat(debugFile)
			Expect(err).To(BeNil())

			content, readErr := ioutil.ReadFile(debugFile)
			Expect(readErr).To(BeNil())
			Expect(string(content)).To(Equal(responseBody))

			defer os.Remove(debugFile)
		})

		It("showProgress", func() {
			downloader = HTTPDownloader{
				RoundTripper:   roundTripper,
				ShowProgress:   true,
				TargetFilePath: targetFilePath,
			}

			request, _ := http.NewRequest("GET", "", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)
			err := downloader.DownloadFile()
			Expect(err).To(BeNil())
		})
	})
})
