package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("job search command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		jobSearchOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("no need roundtripper", func() {
		It("search without any options", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jobSearchCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
				cmd.Print("help")
			})
			rootCmd.SetArgs([]string{"job", "search"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("help"))
		})
	})

	Context("basic cases, need roundtrippper", func() {
		BeforeEach(func() {
			roundTripper = mhttp.NewMockRoundTripper(ctrl)
			jobSearchOption.RoundTripper = roundTripper
		})

		It("should success, search with one result item", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			keyword := "fake"

			request, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/jenkins/search/suggest?query=%s&max=10", keyword), nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"suggestions": [{"name": "fake"}]}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			rootCmd.SetArgs([]string{"job", "search", keyword})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`{
  "Suggestions": [
    {
      "Name": "fake"
    }
  ]
}
`))
		})

		It("should success, search without keyword", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := http.NewRequest("GET", "http://localhost:8080/jenkins/search/suggest?query=&max=10", nil)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"suggestions": [{"name": "fake"}]}`)),
			}
			roundTripper.EXPECT().
				RoundTrip(request).Return(response, nil)

			rootCmd.SetArgs([]string{"job", "search", "--all"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`{
  "Suggestions": [
    {
      "Name": "fake"
    }
  ]
}
`))
		})
	})
})
