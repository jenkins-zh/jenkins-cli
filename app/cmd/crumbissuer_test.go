package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("crumb command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		jenkinsRoot  string
		username     string
		token        string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		crumbIssuerOptions.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		jenkinsRoot = "http://localhost:8080/jenkins"
		username = "admin"
		token = "111e3a2f0231198855dceaff96f20540a9"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("should success", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForGetIssuer(roundTripper, jenkinsRoot, username, token)

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			rootCmd.SetArgs([]string{"crumb"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("CrumbRequestField=Crumb\n"))
		})
	})
})
