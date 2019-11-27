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

var _ = Describe("plugin upgrade command", func() {
	var (
		ctrl         *gomock.Controller
		roundTripper *mhttp.MockRoundTripper
		pluginName   string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		roundTripper = mhttp.NewMockRoundTripper(ctrl)
		pluginUpgradeOption.RoundTripper = roundTripper
		rootCmd.SetArgs([]string{})
		rootOptions.Jenkins = ""
		rootOptions.ConfigFile = "test.yaml"
		pluginName = "fake"
	})

	AfterEach(func() {
		rootCmd.SetArgs([]string{})
		os.Remove(rootOptions.ConfigFile)
		rootOptions.ConfigFile = ""
		ctrl.Finish()
	})

	Context("basic cases", func() {
		It("given plugin name, should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			client.PrepareForInstallPlugin(roundTripper, "http://localhost:8080/jenkins", pluginName, "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "upgrade", pluginName})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("findUpgradeablePlugins", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			jclient := &client.PluginManager{
				JenkinsCore: client.JenkinsCore{
					RoundTripper: roundTripper,
				},
			}
			getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")

			pluginUpgradeOption.Filter = []string{"name=a"}

			result, err := pluginUpgradeOption.findUpgradeablePlugins(jclient)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(len(result)).To(Equal(1))
			Expect(result[0].ShortName).To(Equal(pluginName))

			plugins := pluginUpgradeOption.convertToArray(result)
			Expect(len(plugins)).To(Equal(1))
			Expect(plugins[0]).To(Equal(pluginName))
		})

		It("upgrade compatible plugin, should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			//response := client.PrepareShowPlugins(roundTripper, pluginName)
			client.PrepareForInstallPlugin(roundTripper, "http://localhost:8080/jenkins", pluginName, "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "upgrade", "--compatible"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})

		It("upgrade all plugin, should success", func() {
			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			request, _ := client.PrepareForOneInstalledPlugin(roundTripper, "http://localhost:8080/jenkins")
			request.SetBasicAuth("admin", "111e3a2f0231198855dceaff96f20540a9")
			//client.PrepareShowPlugins(roundTripper, pluginName)
			client.PrepareForInstallPlugin(roundTripper, "http://localhost:8080/jenkins", pluginName, "admin", "111e3a2f0231198855dceaff96f20540a9")

			rootCmd.SetArgs([]string{"plugin", "upgrade", "--all"})

			buf := new(bytes.Buffer)
			rootCmd.SetOutput(buf)
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())

			Expect(buf.String()).To(Equal(""))
		})
	})
})
