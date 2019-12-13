package cmd

import (
	"io/ioutil"
	"os"

	"github.com/jenkins-zh/jenkins-cli/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("casc command check", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		rootURL      string
		user         string
		password     string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		cascOptions.RoundTripper = roundTripper
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		rootURL = "http://localhost:8080/jenkins"
		user = "admin"
		password = "111e3a2f0231198855dceaff96f20540a9"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("without casc plugin", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			req, _ := client.PrepareForOneInstalledPlugin(roundTripper, rootURL)
			req.SetBasicAuth(user, password)

			err = cascOptions.Check()
			Expect(err).To(HaveOccurred())
		})

		It("with casc plugin", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			req, _ := client.PrepareForOneInstalledPluginWithPluginName(roundTripper, rootURL,
				"configuration-as-code")
			req.SetBasicAuth(user, password)

			err = cascOptions.Check()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
