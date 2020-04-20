package cmd

import (
	"os"

	. "github.com/jenkins-zh/jenkins-cli/app/config"
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

		config = &Config{
			Current: "fake",
			JenkinsServers: []JenkinsServer{JenkinsServer{
				Name:     "fake",
				URL:      rootURL,
				UserName: user,
				Token:    password,
			}},
		}
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		config = nil
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("without casc plugin", func() {
			req, _ := client.PrepareForOneInstalledPlugin(roundTripper, rootURL)
			req.SetBasicAuth(user, password)

			err := cascOptions.Check()
			Expect(err).To(HaveOccurred())
		})

		It("with casc plugin", func() {
			req, _ := client.PrepareForOneInstalledPluginWithPluginName(roundTripper, rootURL,
				"configuration-as-code")
			req.SetBasicAuth(user, password)

			err := cascOptions.Check()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
