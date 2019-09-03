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

			request, _ := client.PrepareForEmptyInstalledPluginList(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal("number name version update\n"))
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

			Expect(buf.String()).To(Equal(`number name version update
0      fake 1.0     true
`))
		})

		It("one plugin output with json format", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "list", "fake", "--output", "json", "--filter", "hasUpdate", "--filter", "name=fake"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(pluginsJSON()))
		})
	})
})

func pluginsJSON() string {
	return `[
  {
    "Active": false,
    "Enabled": false,
    "Bundled": false,
    "Downgradable": false,
    "Deleted": false,
    "Enable": false,
    "ShortName": "fake",
    "LongName": "",
    "Version": "1.0",
    "URL": "",
    "HasUpdate": true,
    "Pinned": false,
    "RequiredCoreVesion": "",
    "MinimumJavaVersion": "",
    "SupportDynamicLoad": "",
    "BackVersion": ""
  }
]`
}
