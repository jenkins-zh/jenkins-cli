package cmd

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("backup command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		backupOption.RoundTripper = roundTripper
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		ctrl.Finish()
	})

	Context("basic cases", func() {
		BeforeEach(func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("backup successfully", func() {
			request, _ := http.NewRequest(http.MethodGet, "/thinBackup/backupManual", nil)
			response := &http.Response{
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       nil,
			}
			roundTripper.EXPECT().RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)
			rootCmd.SetArgs([]string{"backup"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})

		It("not found resources", func() {
			request, _ := http.NewRequest(http.MethodGet, "/thinBackup/backupManual", nil)
			response := &http.Response{
				StatusCode: 404,
				Proto:      "HTTP/1.1",
				Request:    request,
				Body:       nil,
			}
			roundTripper.EXPECT().RoundTrip(client.NewRequestMatcher(request)).Return(response, nil)
			rootCmd.SetArgs([]string{"backup"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(Equal("not found resources"))
		})
	})
})
