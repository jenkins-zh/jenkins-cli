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

var _ = Describe("plugin list command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginListOption.RoundTripper = roundTripper
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

	Context("basic cases", func() {
		It("no plugin in the list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForEmptyInstalledPluginList(roundTripper, "http://localhost:8080/jenkins", 1)
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("ShortName Version HasUpdate\n"))
		})

		It("one plugin in the list", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list", "fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`ShortName Version HasUpdate
fake      1.0     true
`))
		})

		It("one plugin in the list without headers", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list", "fake", "--no-headers"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(`fake 1.0 true
`))
		})

		It("one plugin output with json format", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list", "fake", "--output", "json", "--filter", "HasUpdate=true",
				"--filter", "ShortName=fake", "--filter", "Enable=true", "--filter", "Active=true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(pluginsJSON()))
		})

		It("one plugin output with yaml format", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list", "fake", "--output", "yaml", "--filter", "HasUpdate=true",
				"--filter", "ShortName=fake", "--filter", "Enable=true", "--filter", "Active=true"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(pluginYaml()))
		})

		It("one plugin output with not support format", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list", "fake", "--output", "fake"})

			_, err = rootCmd.ExecuteC()
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(ContainSubstring("not support format fake"))
		})
	})
})

func pluginsJSON() string {
	return `[
  {
    "Active": true,
    "Enabled": false,
    "Bundled": false,
    "Downgradable": false,
    "Deleted": false,
    "Enable": true,
    "ShortName": "fake",
    "LongName": "",
    "Version": "1.0",
    "URL": "",
    "HasUpdate": true,
    "Pinned": false,
    "RequiredCoreVesion": "",
    "MinimumJavaVersion": "",
    "SupportDynamicLoad": "",
    "BackVersion": "",
    "Dependencies": null
  }
]`
}

func pluginYaml() string {
	return `- plugin:
    active: true
    enabled: false
    bundled: false
    downgradable: false
    deleted: false
  enable: true
  shortname: fake
  longname: ""
  version: "1.0"
  url: ""
  hasupdate: true
  pinned: false
  requiredcorevesion: ""
  minimumjavaversion: ""
  supportdynamicload: ""
  backversion: ""
  dependencies: []
`
}
