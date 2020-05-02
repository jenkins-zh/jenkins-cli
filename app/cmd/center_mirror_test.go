package cmd

import (
	"github.com/jenkins-zh/jenkins-cli/client"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("center mirror command", func() {
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
		centerUpgradeOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
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
		It("enable mirror site", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForChangeUpdateCenterSite(roundTripper, rootURL, user, password,
				"default", "https://updates.jenkins-zh.cn/update-center.json")
			client.PrepareForSetMirrorCertificate(roundTripper, rootURL, user, password, true)

			rootCmd.SetArgs([]string{"center", "mirror"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})

		It("disable mirror site", func() {
			data, err := GenerateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForChangeUpdateCenterSite(roundTripper, rootURL, user, password,
				"default", "https://updates.jenkins.io/update-center.json")
			client.PrepareForSetMirrorCertificate(roundTripper, rootURL, user, password, false)

			rootCmd.SetArgs([]string{"center", "mirror", "--enable=false"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
		})
	})
})
