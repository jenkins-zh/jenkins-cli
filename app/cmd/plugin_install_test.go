package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/mock/mhttp"
)

var _ = Describe("plugin install command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		rootURL      string
		username     string
		token        string

		pluginName string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginInstallOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"

		rootURL = "http://localhost:8080/jenkins"
		username = "admin"
		token = "111e3a2f0231198855dceaff96f20540a9"

		pluginName = "fake"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("install one plugin", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForInstallPlugin(roundTripper, rootURL, pluginName, username, token)

			rootCmd.SetArgs([]string{"plugin", "install", pluginName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("unknow suite", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			rootCmd.SetArgs([]string{"plugin", "install", pluginName, "--suite", "fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("error: cannot found suite fake"))
		})
	})

	Context("convertToArray", func() {
		It("empty plugins", func() {
			plugins := make([]client.AvailablePlugin, 0)
			result := convertToArray(plugins)
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(0))
		})

		It("all installed plugins", func() {
			plugins := []client.AvailablePlugin{{
				Installed: true,
			}}
			result := convertToArray(plugins)
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(0))
		})

		It("one no installed plugin", func() {
			plugins := []client.AvailablePlugin{{
				Installed: false,
				Name:      "fake",
			}}
			result := convertToArray(plugins)
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(1))
		})
	})
})
