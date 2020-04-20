package cmd

import (
	"bytes"
	"github.com/golang/mock/gomock"
	. "github.com/jenkins-zh/jenkins-cli/app/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var _ = Describe("Table util test", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		config = nil
	})

	AfterEach(func() {
		config = nil
		ctrl.Finish()
	})

	Context("basic test", func() {
		It("getJenkinsNames", func() {
			config = &Config{
				JenkinsServers: []JenkinsServer{{
					Name: "a",
				}, {
					Name: "b",
				}},
			}

			names := getJenkinsNames()
			Expect(names).To(Equal([]string{"a", "b"}))

			config.JenkinsServers = []JenkinsServer{}
			names = getJenkinsNames()
			Expect(names).To(Equal([]string{}))
		})

		It("getCurrentJenkins", func() {
			config = &Config{}
			current := getCurrentJenkins()
			Expect(current).To(BeNil())

			config.Current = "test"
			config.JenkinsServers = []JenkinsServer{{
				Name: "test",
			}}
			current = getCurrentJenkins()
			Expect(current).To(Equal(&config.JenkinsServers[0]))
		})

		It("findSuiteByName", func() {
			config = &Config{}
			suite := findSuiteByName("fake")
			Expect(suite).To(BeNil())

			pluginName := "plugin-one"
			config.PluginSuites = []PluginSuite{{
				Name: pluginName,
			}}
			suite = findSuiteByName(pluginName)
			Expect(suite).NotTo(BeNil())
			Expect(suite.Name).To(Equal(pluginName))
		})

		It("getMirrors", func() {
			config = &Config{}
			mirrors := getMirrors()
			Expect(mirrors).NotTo(BeNil())
			Expect(len(mirrors)).To(Equal(1))
			Expect(mirrors[0].Name).To(Equal("default"))
		})
	})

	Context("command test", func() {
		BeforeEach(func() {
			rootOptions.Jenkins = ""
			rootOptions.ConfigFile = "test.yaml"

			data, err := generateSampleConfig()
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())
		})

		It("config command test", func() {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)

			rootCmd.SetArgs([]string{"config"})
			_, err := rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("Current Jenkins's name is"))
		})

		It("config command with description", func() {
			jenkinsDesc := "description"

			sampleConfig := getSampleConfig()
			sampleConfig.JenkinsServers[0].Description = jenkinsDesc
			data, err := yaml.Marshal(&sampleConfig)
			Expect(err).To(BeNil())
			err = ioutil.WriteFile(rootOptions.ConfigFile, data, 0664)
			Expect(err).To(BeNil())

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)

			rootCmd.SetArgs([]string{"config"})
			_, err = rootCmd.ExecuteC()
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring(jenkinsDesc))
		})
	})
})

var _ = Describe("GetConfigFromHome", func() {
	It("should success", func() {
		path, err := GetConfigFromHome()
		Expect(err).To(BeNil())
		Expect(path).To(ContainSubstring(".jenkins-cli.yaml"))
	})
})
